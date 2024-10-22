package main

import (
	"fmt"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
)

func DHCPHandler(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	// print the received DHCPv4 message
	Debug(m.Summary())

	// modify from coredhcp
	if m.OpCode != dhcpv4.OpcodeBootRequest {
		Error("Unsupported opcode: ", m.OpCode, ". Only BootRequest (", dhcpv4.OpcodeBootRequest, ") is supported")
		return
	}

	reply, err := dhcpv4.NewReplyFromRequest(m)
	if err != nil {
		Error("Failed to create new reply from dhcp qequest: ", err)
		return
	}

	switch mt := m.MessageType(); mt {
	case dhcpv4.MessageTypeDiscover:
		reply.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeOffer))
	case dhcpv4.MessageTypeRequest:
		reply.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeAck))
	default:
		Info("Unhandled message type: ", mt)
		return
	}

	if arch := m.ClientArch(); len(arch) > 0 {
		switch arch[0] {
		case 7:
			// EFI_X86_64
			reply.UpdateOption(dhcpv4.OptBootFileName("syslinux.efi"))
		default:
			reply.UpdateOption(dhcpv4.OptBootFileName("pxelinux.0"))
		}
	}

	reply.YourIPAddr = net.IPv4(0, 0, 0, 0)
	reply.ServerIPAddr = net.ParseIP(serverIPAddr)
	reply.SetBroadcast()
	// 60, 66, 67
	reply.UpdateOption(dhcpv4.OptClassIdentifier("PXEClient"))
	// next server
	reply.UpdateOption(dhcpv4.OptTFTPServerName(serverIPAddr))

	if _, err := conn.WriteTo(reply.ToBytes(), peer); err != nil {
		Error("Cannot reply to client: ", err)
	}
	Debug(reply.Summary())
}

func Rundhcp(iface string, port int) {
	address := fmt.Sprintf(":%d", port)
	laddr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		Fatal(err)
	}
	server, err := server4.NewServer(iface, laddr, DHCPHandler)
	if err != nil {
		Fatal(err)
	}
	server.Serve()
}
