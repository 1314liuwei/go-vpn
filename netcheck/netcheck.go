package netcheck

import (
	"log"
	"net"
	"time"

	"github.com/gogf/gf/v2/util/gconv"

	"github.com/gogf/gf/v2/frame/g"
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
	log.Println(addr.String())
	stunAddr := &net.UDPAddr{IP: addr.IP, Port: 3478}

	udp, err := net.ListenUDP("udp4", nil)
	if err != nil {
		return "", err
	}
	defer udp.Close()

	var (
		stepRes []map[string]map[string]interface{}
	)

	buff := make([]byte, 1024)
	// Step1:
	{
		req := buildRequestHeader(0)
		_, err = udp.WriteToUDP(req, stunAddr)
		if err != nil {
			log.Fatal(err)
			return "", err
		}
		log.Println("step1 send success!")

		err := udp.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			return "", err
		}
		read, _, err := udp.ReadFromUDP(buff)
		if err != nil {
			netErr, ok := err.(*net.OpError)
			if ok && netErr.Timeout() {
				return "UDP is not allowed", nil
			} else {
				return "", err
			}
		}
		attributes, err := parseResponseAttributes(buff[:read])
		if err != nil {
			log.Fatal(err)
			return "", err
		}
		g.Dump(attributes)
		stepRes = append(stepRes, attributes)
	}

	// Step2:
	{
		req := buildChangePortAndIPRequest(true, true)
		_, err = udp.WriteToUDP(req, stunAddr)
		if err != nil {
			log.Fatal(err)
			return "", err
		}
		log.Println("step2 send success!")

		err := udp.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			return "", err
		}
		read, _, err := udp.ReadFromUDP(buff)
		if err != nil {
			if err != nil {
				netErr, ok := err.(*net.OpError)
				if ok && netErr.Timeout() {
					stepRes = append(stepRes, nil)
					goto step3
				} else {
					return "", err
				}
			}
		}

		attributes, err := parseResponseAttributes(buff[:read])
		if err != nil {
			log.Fatal(err)
			return "", err
		}
		g.Dump(attributes)
		stepRes = append(stepRes, attributes)
	}

step3:
	// Step3:
	{
		req := buildChangePortAndIPRequest(false, true)
		_, err = udp.WriteToUDP(req, stunAddr)
		if err != nil {
			log.Fatal(err)
			return "", err
		}
		log.Println("step3 send success!")

		err := udp.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			return "", err
		}
		read, _, err := udp.ReadFromUDP(buff)
		if err != nil {
			if err != nil {
				netErr, ok := err.(*net.OpError)
				if ok && netErr.Timeout() {
					stepRes = append(stepRes, nil)
					goto step4
				} else {
					return "", err
				}
			}
		}

		attributes, err := parseResponseAttributes(buff[:read])
		if err != nil {
			log.Fatal(err)
			return "", err
		}
		g.Dump(attributes)
		stepRes = append(stepRes, attributes)
	}

step4:
	// Step4:
	{
		if len(stepRes) > 1 {
			step1 := stepRes[0]
			changeAddr, ok := step1["CHANGED_ADDRESS"]["IP"]
			if ok {
				stunAddr.IP = net.ParseIP(gconv.String(changeAddr))
			}
		}
		req := buildRequestHeader(0)
		_, err := udp.WriteToUDP(req, stunAddr)
		if err != nil {
			log.Fatal(err)
			return "", err
		}
		log.Println("step4 send success!")

		err = udp.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			return "", err
		}
		read, _, err := udp.ReadFromUDP(buff)
		if err != nil {
			if err != nil {
				netErr, ok := err.(*net.OpError)
				if ok && netErr.Timeout() {
					stepRes = append(stepRes, nil)
					goto end
				} else {
					return "", err
				}
			}
		}

		attributes, err := parseResponseAttributes(buff[:read])
		if err != nil {
			log.Fatal(err)
			return "", err
		}
		stepRes = append(stepRes, attributes)
	}

end:
	g.Dump(stepRes)
	return "", nil
}
