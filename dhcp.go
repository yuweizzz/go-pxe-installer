package main

import (
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

func DHCPHandler(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	// this function will just print the received DHCPv4 message, without replying
	Debug(m.Summary())

	// from coredhcp
	if m.OpCode != dhcpv4.OpcodeBootRequest {
		Error("unsupported opcode: ", m.OpCode, ". Only BootRequest (", dhcpv4.OpcodeBootRequest, ") is supported")
		return
	}

	reply, err := dhcpv4.NewReplyFromRequest(m)
	if err != nil {
		Error("failed to build reply: ", err)
		return
	}

	switch mt := m.MessageType(); mt {
	case dhcpv4.MessageTypeDiscover:
		reply.YourIPAddr = net.IPv4(0, 0, 0, 0)
		reply.ServerIPAddr = net.IPv4(10, 0, 2, 5)
		reply.ServerHostName = "10.0.2.5"
		// next server
		reply.UpdateOption(dhcpv4.OptTFTPServerName("10.0.2.5"))
		reply.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeOffer))
		reply.UpdateOption(dhcpv4.OptServerIdentifier(net.IPv4(10, 0, 2, 5)))
		reply.UpdateOption(dhcpv4.OptClassIdentifier("PXEClient"))
		//reply.UpdateOption(dhcpv4.OptSubnetMask(net.IPv4Mask(255, 255, 255, 0)))
		//reply.UpdateOption(dhcpv4.OptBroadcastAddress(net.IPv4(10, 0, 2, 255)))
		//reply.UpdateOption(dhcpv4.OptDNS(net.IPv4(10, 0, 2, 1)))
		//reply.UpdateOption(dhcpv4.OptRouter(net.IPv4(10, 0, 2, 1)))
		//hours, _ := time.ParseDuration("1h")
		//reply.UpdateOption(dhcpv4.OptIPAddressLeaseTime(hours))
		//hours, _ = time.ParseDuration("3h")
		//reply.UpdateOption(dhcpv4.OptRebindingTimeValue(hours))
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
		reply.YourIPAddr = net.IPv4(0, 0, 0, 0)
		reply.ServerIPAddr = net.IPv4(10, 0, 2, 5)
		reply.ServerHostName = "10.0.2.5"
		reply.UpdateOption(dhcpv4.OptTFTPServerName("10.0.2.5"))
		reply.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeOffer))
		reply.UpdateOption(dhcpv4.OptServerIdentifier(net.IPv4(10, 0, 2, 5)))
		reply.UpdateOption(dhcpv4.OptClassIdentifier("PXEClient"))
		//reply.UpdateOption(dhcpv4.OptSubnetMask(net.IPv4Mask(255, 255, 255, 0)))
		//reply.UpdateOption(dhcpv4.OptBroadcastAddress(net.IPv4(10, 0, 2, 255)))
		//reply.UpdateOption(dhcpv4.OptDNS(net.IPv4(10, 0, 2, 1)))
		//reply.UpdateOption(dhcpv4.OptRouter(net.IPv4(10, 0, 2, 1)))
		//hours, _ := time.ParseDuration("1h")
		//reply.UpdateOption(dhcpv4.OptIPAddressLeaseTime(hours))
		//hours, _ = time.ParseDuration("3h")
		//reply.UpdateOption(dhcpv4.OptRebindingTimeValue(hours))
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
		Info("Unhandled message type: ", mt)
		return
	}

	if _, err := conn.WriteTo(reply.ToBytes(), peer); err != nil {
		Error("Cannot reply to client: ", err)
	}
	Debug(reply.Summary())
}
