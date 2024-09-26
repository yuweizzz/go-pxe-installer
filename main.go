package main

import (
	"log"
	"net"
	"time"
	"os"
	"fmt"
	"io"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
	tftp "github.com/pin/tftp/v3"
)

func handler(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	// this function will just print the received DHCPv4 message, without replying
	log.Print(m.Summary())

	// from coredhcp
	if m.OpCode != dhcpv4.OpcodeBootRequest {
		log.Printf("unsupported opcode %d. Only BootRequest (%d) is supported", m.OpCode, dhcpv4.OpcodeBootRequest)
		return
	}

	reply, err := dhcpv4.NewReplyFromRequest(m)
	if err != nil {
		log.Printf("failed to build reply: %v", err)
		return
	}

	switch mt := m.MessageType(); mt {
	case dhcpv4.MessageTypeDiscover:
		reply.YourIPAddr = net.IPv4(10,0,2,100)
		reply.ServerIPAddr = net.IPv4(10,0,2,5)
		reply.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeOffer))
		reply.UpdateOption(dhcpv4.OptServerIdentifier(net.IPv4(10,0,2,5)))
		reply.UpdateOption(dhcpv4.OptSubnetMask(net.IPv4Mask(255, 255, 255, 0)))
		reply.UpdateOption(dhcpv4.OptBroadcastAddress(net.IPv4(10,0,2,255)))
		reply.UpdateOption(dhcpv4.OptDNS(net.IPv4(10,0,2,1)))
		reply.UpdateOption(dhcpv4.OptRouter(net.IPv4(10,0,2,1)))
		hours, _ := time.ParseDuration("1h")
		reply.UpdateOption(dhcpv4.OptIPAddressLeaseTime(hours))
		hours, _ = time.ParseDuration("3h")
		reply.UpdateOption(dhcpv4.OptRebindingTimeValue(hours))
		log.Println(m.ClientArch())
		if arch := m.ClientArch(); len(arch) > 0 {
			switch arch[0] {
			case 7:
				// EFI_X86_64
				reply.UpdateOption(dhcpv4.OptBootFileName("syslinux.efi"))
			default:
				reply.UpdateOption(dhcpv4.OptBootFileName("pxelinux.0"))
			}
		}
	case dhcpv4.MessageTypeRequest:
		reply.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeAck))
		reply.YourIPAddr = net.IPv4(10,0,2,100)
		reply.ServerIPAddr = net.IPv4(10,0,2,5)
		reply.UpdateOption(dhcpv4.OptServerIdentifier(net.IPv4(10,0,2,5)))
		reply.UpdateOption(dhcpv4.OptSubnetMask(net.IPv4Mask(255, 255, 255, 0)))
		reply.UpdateOption(dhcpv4.OptBroadcastAddress(net.IPv4(10,0,2,255)))
		reply.UpdateOption(dhcpv4.OptDNS(net.IPv4(10,0,2,1)))
		reply.UpdateOption(dhcpv4.OptRouter(net.IPv4(10,0,2,1)))
		hours, _ := time.ParseDuration("1h")
		reply.UpdateOption(dhcpv4.OptIPAddressLeaseTime(hours))
		hours, _ = time.ParseDuration("3h")
		reply.UpdateOption(dhcpv4.OptRebindingTimeValue(hours))
		if arch := m.ClientArch(); len(arch) > 0 {
			switch arch[0] {
			case 7:
				// EFI_X86_64
				reply.UpdateOption(dhcpv4.OptBootFileName("syslinux.efi"))
			default:
				reply.UpdateOption(dhcpv4.OptBootFileName("pxelinux.0"))
			}
		}
	default:
		log.Printf("Unhandled message type: %v", mt)
		return
	}

	if _, err := conn.WriteTo(reply.ToBytes(), peer); err != nil {
		log.Printf("Cannot reply to client: %v", err)
	}
	log.Print(reply.Summary())
}

// readHandler is called when client starts file download from server
func readHandler(filename string, rf io.ReaderFrom) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	n, err := rf.ReadFrom(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	fmt.Printf("%d bytes sent\n", n)
	return nil
}


func runtftp() {
	// use nil in place of handler to disable read or write operations
	s := tftp.NewServer(readHandler, nil)
	s.SetTimeout(5 * time.Second) // optional
	err := s.ListenAndServe(":69") // blocks until s.Shutdown() is called
	if err != nil {
		fmt.Fprintf(os.Stdout, "server: %v\n", err)
		os.Exit(1)
	}
}


func main() {
	laddr := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 67,
	}
	server, err := server4.NewServer("enp0s3", &laddr, handler)
	if err != nil {
		log.Fatal(err)
	}

	// This never returns. If you want to do other stuff, dump it into a
	// goroutine.
	go server.Serve()
	runtftp()
}
