package main

import (
	"log"
	"net"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
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
	default:
		log.Printf("Unhandled message type: %v", mt)
		return
	}

	if _, err := conn.WriteTo(reply.ToBytes(), peer); err != nil {
		log.Printf("Cannot reply to client: %v", err)
	}
	log.Print(reply.Summary())
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
	server.Serve()
}
