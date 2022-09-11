package netcheck

import (
	"context"
	"log"
	"net"
	"net/netip"
	"time"

	"github.com/gogf/gf/v2/util/gconv"
)

const (
	UDPBlocked = 127

	SetConnTimeout = 3 * time.Second // 连接超时时间
	SetTryoutTimes = 5               // 超时重传次数
)

var (
	STUNServers = []string{
		"stun.qq.com",
		"stun.bige0.com",
	}
)

type FilteringBehavior int

const (
	UnknownFiltering FilteringBehavior = iota
	EndpointIndependentFiltering
	AddressAndPortDependentFiltering
	AddressDependentFiltering
)

type MappingBehavior int

const (
	UnknownMapping MappingBehavior = iota
	NoNAT
	EndpointIndependentMapping
	AddressAndPortDependentMapping
	AddressDependentMapping
)

func NatTypeTest(ctx context.Context) (MappingBehavior, FilteringBehavior, error) {
	addr, err := net.ResolveIPAddr("ip", STUNServers[0])
	if err != nil {
		return UnknownMapping, UnknownFiltering, err
	}
	stunAddr := &net.UDPAddr{IP: addr.IP, Port: 3478}

	udp, err := net.ListenUDP("udp4", nil)
	if err != nil {
		return UnknownMapping, UnknownFiltering, err
	}
	defer udp.Close()

	// 探测 Mapping 行为
	mBehavior, err := probeMappingBehavior(ctx, udp, stunAddr)
	if err != nil {
		return UnknownMapping, UnknownFiltering, nil
	}
	// 探测 Filtering 行为
	fBehavior, err := probeFilterBehavior(ctx, udp, stunAddr)
	if err != nil {
		return UnknownMapping, UnknownFiltering, nil
	}

	return mBehavior, fBehavior, nil
}

func probeMappingBehavior(ctx context.Context, conn *net.UDPConn, stunAddr *net.UDPAddr) (MappingBehavior, error) {
	log.Println("Start probe mapping behavior...")
	defer log.Println("End probe mapping behavior!")

	// Step1: 探测主机是否位于 NAT 后面
	req := buildRequestHeader(0)
	attributes1, err := probeSendAndReceive(ctx, conn, stunAddr, req)
	if err != nil {
		netErr, ok := err.(*net.OpError)
		if ok && netErr.Timeout() {
			return UDPBlocked, nil
		} else {
			return UnknownMapping, err
		}
	}

	laddr, err := netip.ParseAddrPort(conn.LocalAddr().String())
	if err != nil {
		return UnknownMapping, err
	}
	rAddr, err := netip.ParseAddr(attributes1[AttributesName[MappedAddress]].IP.To4().String())
	if err != nil {
		return UnknownMapping, err
	}
	raddr := netip.AddrPortFrom(rAddr, gconv.Uint16(attributes1[AttributesName[MappedAddress]].Port))
	if isLocalAddrEqualRemoteAddr(laddr, raddr) {
		return NoNAT, nil
	}

	// Step: 探测 Endpoint Independent Mapping
	// 使用 STUN 服务器的另一个地址进行检测
	stunAddr.IP = attributes1[AttributesName[ChangedAddress]].IP
	req = buildRequestHeader(0)
	attributes2, err := probeSendAndReceive(ctx, conn, stunAddr, req)
	if err != nil {
		netErr, ok := err.(*net.OpError)
		if ok && netErr.Timeout() {
			return UDPBlocked, nil
		} else {
			return UnknownMapping, err
		}
	}

	// 如果第一次的地址和这一次的地址一致，则代表是 EndpointIndependentMapping 类型；否则需要进一步判断
	if attributes1[AttributesName[MappedAddress]].IP.To4().String() == attributes2[AttributesName[MappedAddress]].IP.To4().String() &&
		attributes1[AttributesName[MappedAddress]].Port == attributes2[AttributesName[MappedAddress]].Port {
		return EndpointIndependentMapping, nil
	}

	// Step3: 探测 AddressDependentMapping 和 AddressAndPortDependentMapping
	stunAddr.Port = attributes1[AttributesName[ChangedAddress]].Port
	req = buildRequestHeader(0)
	attributes3, err := probeSendAndReceive(ctx, conn, stunAddr, req)
	if err != nil {
		netErr, ok := err.(*net.OpError)
		if ok && netErr.Timeout() {
			return UDPBlocked, nil
		} else {
			return UnknownMapping, err
		}
	}

	// 如果第二次的地址和这一次的地址一致，则代表是 AddressDependentMapping 类型；否则是 AddressAndPortDependentMapping 类型
	if attributes2[AttributesName[MappedAddress]].IP.To4().String() == attributes3[AttributesName[MappedAddress]].IP.To4().String() &&
		attributes2[AttributesName[MappedAddress]].Port == attributes3[AttributesName[MappedAddress]].Port {
		return AddressDependentMapping, nil
	} else {
		return AddressAndPortDependentMapping, nil
	}
}

func isLocalAddrEqualRemoteAddr(laddr, raddr netip.AddrPort) bool {
	if laddr.Addr() != netip.IPv4Unspecified() {
		return laddr.Addr() == raddr.Addr() && laddr.Port() == raddr.Port()
	}

	if laddr.Port() != raddr.Port() {
		return false
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false
	}
	for _, addr := range addrs {
		if addr.String() == raddr.Addr().String() {
			return true
		}
	}
	return false
}

func probeFilterBehavior(ctx context.Context, conn *net.UDPConn, stunAddr *net.UDPAddr) (FilteringBehavior, error) {
	log.Println("Start probe filtering behavior...")
	defer log.Println("End probe filtering behavior!")

	// Step1: 探测 Endpoint-Independent Filtering
	req := buildChangePortAndIPRequest(true, true)
	_, err := probeSendAndReceive(ctx, conn, stunAddr, req)
	if err != nil {
		netErr, ok := err.(*net.OpError)
		if ok && netErr.Timeout() {
			goto step2
		} else {
			return UDPBlocked, err
		}
	} else {
		return EndpointIndependentFiltering, err
	}

step2:
	// Step2: 探测 Address-Dependent Filtering 和 Address and Port-Dependent Filtering
	req = buildChangePortAndIPRequest(false, true)
	_, err = probeSendAndReceive(ctx, conn, stunAddr, req)
	if err != nil {
		netErr, ok := err.(*net.OpError)
		if ok && netErr.Timeout() {
			return AddressAndPortDependentFiltering, nil
		} else {
			return UnknownFiltering, err
		}
	}
	return AddressDependentFiltering, nil
}

func probeSendAndReceive(ctx context.Context, conn *net.UDPConn, stunAddr *net.UDPAddr, req []byte) (map[string]AddressAttribute, error) {
	buff := make([]byte, 1024)

	for i := 0; i < SetTryoutTimes; i++ {
		log.Printf("The number of attempts is %d", i)
		err := conn.SetWriteDeadline(time.Now().Add(SetConnTimeout))
		if err != nil {
			if i != SetTryoutTimes-1 {
				continue
			}
			return nil, err
		}

		_, err = conn.WriteToUDP(req, stunAddr)
		if err != nil {
			if i != SetTryoutTimes-1 {
				continue
			} else {
				return nil, err
			}
		}

		err = conn.SetReadDeadline(time.Now().Add(SetConnTimeout))
		if err != nil {
			if i != SetTryoutTimes-1 {
				continue
			}
			return nil, err
		}

		read, _, err := conn.ReadFromUDP(buff)
		if err != nil {
			if i != SetTryoutTimes-1 {
				continue
			} else {
				return nil, err
			}
		}
		attributes, err := parseResponseAttributes(buff[:read])
		if err != nil {
			if i != SetTryoutTimes-1 {
				continue
			}
			return nil, err
		} else {
			return attributes, nil
		}
	}
	return nil, nil
}
