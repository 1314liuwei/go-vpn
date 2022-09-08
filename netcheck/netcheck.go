package netcheck

import (
	"log"
	"net"
	"net/netip"
	"time"

	"github.com/gogf/gf/v2/util/gconv"
)

const (
	UDPBlocked = 127
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

func NatTypeTest() (string, error) {
	addr, err := net.ResolveIPAddr("ip", STUNServers[0])
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	stunAddr := &net.UDPAddr{IP: addr.IP, Port: 3478}

	udp, err := net.ListenUDP("udp4", nil)
	if err != nil {
		return "", err
	}
	defer udp.Close()
	log.Println(udp.LocalAddr().String())

	mBehavior, err := probeMappingBehavior(udp, stunAddr)
	if err != nil {
		return "", err
	}
	log.Println("MappingBehavior: ", mBehavior)

	fBehavior, err := probeFilterBehavior(udp, stunAddr)
	if err != nil {
		return "", err
	}
	log.Println("FilterBehavior: ", fBehavior)
	return "", nil
}

func probeMappingBehavior(conn *net.UDPConn, stunAddr *net.UDPAddr) (MappingBehavior, error) {
	buff := make([]byte, 1024)

	// Step1: 探测主机是否位于 NAT 后面
	req := buildRequestHeader(0)
	_, err := conn.WriteToUDP(req, stunAddr)
	if err != nil {
		log.Fatal(err)
		return UnknownMapping, err
	}

	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return UnknownMapping, err
	}
	read, _, err := conn.ReadFromUDP(buff)

	if err != nil {
		netErr, ok := err.(*net.OpError)
		if ok && netErr.Timeout() {
			return UnknownMapping, nil
		} else {
			return UDPBlocked, err
		}
	}
	attributes1, err := parseResponseAttributes(buff[:read])
	if err != nil {
		log.Fatal(err)
		return UnknownMapping, err
	}

	laddr, err := netip.ParseAddrPort(conn.LocalAddr().String())
	if err != nil {
		return 0, err
	}
	rAddr, err := netip.ParseAddr(attributes1[AttributesName[MappedAddress]].IP.To4().String())
	if err != nil {
		return 0, err
	}
	raddr := netip.AddrPortFrom(rAddr, gconv.Uint16(attributes1[AttributesName[MappedAddress]].Port))
	if isLocalAddrEqualRemoteAddr(laddr, raddr) {
		return NoNAT, nil
	}

	// Step: 探测 Endpoint Independent Mapping
	// 使用 STUN 服务器的另一个地址进行检测
	stunAddr.IP = attributes1[AttributesName[ChangedAddress]].IP
	req = buildRequestHeader(0)
	_, err = conn.WriteToUDP(req, stunAddr)
	if err != nil {
		return UnknownMapping, err
	}

	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return UnknownMapping, err
	}
	read, _, err = conn.ReadFromUDP(buff)
	if err != nil {
		netErr, ok := err.(*net.OpError)
		if ok && netErr.Timeout() {
			return UnknownMapping, nil
		} else {
			return UDPBlocked, err
		}
	}
	attributes2, err := parseResponseAttributes(buff[:read])
	if err != nil {
		log.Fatal(err)
		return UnknownMapping, err
	}

	// 如果第一次的地址和这一次的地址一致，则代表是 EndpointIndependentMapping 类型；否则需要进一步判断
	if attributes1[AttributesName[MappedAddress]].IP.To4().String() == attributes2[AttributesName[MappedAddress]].IP.To4().String() &&
		attributes1[AttributesName[MappedAddress]].Port == attributes2[AttributesName[MappedAddress]].Port {
		return EndpointIndependentMapping, nil
	}

	// Step3: 探测 AddressDependentMapping 和 AddressAndPortDependentMapping
	stunAddr.Port = attributes1[AttributesName[ChangedAddress]].Port
	req = buildRequestHeader(0)
	_, err = conn.WriteToUDP(req, stunAddr)
	if err != nil {
		return UnknownMapping, err
	}

	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return UnknownMapping, err
	}
	read, _, err = conn.ReadFromUDP(buff)
	if err != nil {
		netErr, ok := err.(*net.OpError)
		if ok && netErr.Timeout() {
			return UnknownMapping, nil
		} else {
			return UDPBlocked, err
		}
	}
	attributes3, err := parseResponseAttributes(buff[:read])
	if err != nil {
		return UnknownMapping, err
	}

	// 如果第二次的地址和这一次的地址一致，则代表是 AddressDependentMapping 类型；否则是 AddressAndPortDependentMapping 类型
	if attributes2[AttributesName[MappedAddress]].IP.To4().String() == attributes3[AttributesName[MappedAddress]].IP.To4().String() &&
		attributes2[AttributesName[MappedAddress]].Port == attributes3[AttributesName[MappedAddress]].Port {
		return AddressDependentMapping, nil
	} else {
		return AddressAndPortDependentMapping, nil
	}
	return UnknownMapping, nil
}

func probeFilterBehavior(conn *net.UDPConn, stunAddr *net.UDPAddr) (FilteringBehavior, error) {
	buff := make([]byte, 1024)

	// Step1: 探测主机是否位于 NAT 后面
	req := buildChangePortAndIPRequest(true, true)
	_, err := conn.WriteToUDP(req, stunAddr)
	if err != nil {
		return UnknownFiltering, err
	}

	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return UnknownFiltering, err
	}
	read, _, err := conn.ReadFromUDP(buff)
	if err != nil {
		netErr, ok := err.(*net.OpError)
		if ok && netErr.Timeout() {
			goto step2
		} else {
			return UDPBlocked, err
		}
	}
	_, err = parseResponseAttributes(buff[:read])
	if err == nil {
		return EndpointIndependentFiltering, err
	}

step2:
	req = buildChangePortAndIPRequest(false, true)
	_, err = conn.WriteToUDP(req, stunAddr)
	if err != nil {
		return UnknownFiltering, err
	}

	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return UnknownFiltering, err
	}
	read, _, err = conn.ReadFromUDP(buff)
	if err != nil {
		netErr, ok := err.(*net.OpError)
		if ok && netErr.Timeout() {
			return AddressAndPortDependentFiltering, nil
		} else {
			return UDPBlocked, err
		}
	}
	_, err = parseResponseAttributes(buff[:read])
	if err == nil {
		return AddressDependentFiltering, nil
	}
	return UnknownFiltering, nil
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
