package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
)

type DHCPHandler struct {
	DHCPAddr net.IP
	TFTPAddr net.IP
}

type DHCPServer struct {
	Handler *DHCPHandler
	Iface   string
	Port    int
}

func (s *DHCPServer) Run() {
	laddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		Fatal(err)
	}
	server, err := server4.NewServer(s.Iface, laddr, s.Handler.Update)
	if err != nil {
		Fatal(err)
	}
	server.Serve()
}

func (h *DHCPHandler) Update(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	// print the received DHCPv4 message
	Debug(m.Summary())

	// 1. DHCP OpCode supported
	// 2. DHCP pxe request
	// 3. DHCP message type
	// 4. update DHCP reply

	// 1
	// modify from coredhcp
	if m.OpCode != dhcpv4.OpcodeBootRequest {
		Error("Unsupported opcode: ", m.OpCode, ". Only BootRequest (", dhcpv4.OpcodeBootRequest, ") is supported")
		return
	}

	// 2
	classIdentifier := m.ClassIdentifier()
	if !strings.HasPrefix(classIdentifier, "PXEClient") {
		Info("Wrong ClassIdentifier, not need to reply")
		return
	}

	// 3
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

	// 4
	reply.SetBroadcast()
	reply.YourIPAddr = net.IPv4(0, 0, 0, 0)
	reply.ServerIPAddr = h.DHCPAddr
	// Option 60
	reply.UpdateOption(dhcpv4.OptClassIdentifier("PXEClient"))
	// Option 66, Next server IP address
	reply.UpdateOption(dhcpv4.OptTFTPServerName(h.TFTPAddr.String()))
	// Option 67
	if arch := m.ClientArch(); len(arch) > 0 {
		switch arch[0] {
		case 7:
			// EFI_X86_64
			reply.UpdateOption(dhcpv4.OptBootFileName("syslinux.efi"))
		default:
			reply.UpdateOption(dhcpv4.OptBootFileName("pxelinux.0"))
		}
	}

	if _, err := conn.WriteTo(reply.ToBytes(), peer); err != nil {
		Error("Cannot reply to client: ", err)
	}
	Debug(reply.Summary())
}
