package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
)

const (
	IPXE_SCRIPT          string = "ipxe.script"
	X86_64_UEFI_BOOTFILE string = "ipxe-x86_64.efi"
	X86_64_BIOS_BOOTFILE string = "ipxe-x86_64.pxe"
	ARM64_UEFI_BOOTFILE  string = "ipxe-arm64.efi"
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
	// 2. If DHCP PXE request
	// 3. Update DHCP reply

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
		// Disable PXE discovery
		// DHCP Option 43, Suboption 6: PXE_DISCOVERY_CONTROL
		// bit 0 = If set, disable broadcast discovery.
		// bit 1 = If set, disable multicast discovery.
		// bit 2 = If set, only use/accept servers in PXE_BOOT_SERVERS.
		// bit 3 = If set, and a boot file name is present in the initial DHCP or ProxyDHCP offer packet,
		// download the boot file (do not prompt/menu/discover).
		// bit 4-7 = Must be 0.
		vendorSpecificInformation := dhcpv4.Options{
			6: []byte{8},
		}
		reply.UpdateOption(dhcpv4.OptGeneric(dhcpv4.OptionVendorSpecificInformation, vendorSpecificInformation.ToBytes()))
	case dhcpv4.MessageTypeRequest:
		reply.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeAck))
	default:
		Info("Unhandled message type: ", mt)
		return
	}

	// use Broadcast in OFFER and ACK
	reply.SetBroadcast()
	reply.YourIPAddr = net.IPv4(0, 0, 0, 0)
	// Next server IP address
	reply.ServerIPAddr = h.TFTPAddr
	// Option 54, DHCP Server Identifier
	reply.UpdateOption(dhcpv4.OptServerIdentifier(h.DHCPAddr))
	reply.ServerHostName = h.DHCPAddr.String()
	// Option 60, Vendor Class Identifier, always "PXEClient"
	reply.UpdateOption(dhcpv4.OptClassIdentifier("PXEClient"))
	// Option 66, TFTP Server Name
	reply.UpdateOption(dhcpv4.OptTFTPServerName(h.TFTPAddr.String()))
	// Option 67, Boot file Name
	if arch := m.ClientArch(); len(arch) > 0 {
		switch arch[0] {
		case 7:
			// EFI_X86_64
			reply.BootFileName = X86_64_UEFI_BOOTFILE
			reply.UpdateOption(dhcpv4.OptBootFileName(X86_64_UEFI_BOOTFILE))
		case 11:
			// EFI_ARM64
			reply.BootFileName = ARM64_UEFI_BOOTFILE
			reply.UpdateOption(dhcpv4.OptBootFileName(ARM64_UEFI_BOOTFILE))
		default:
			// BIOS X86_64
			reply.BootFileName = X86_64_BIOS_BOOTFILE
			reply.UpdateOption(dhcpv4.OptBootFileName(X86_64_BIOS_BOOTFILE))
		}
	}

	if userClass := m.UserClass(); len(userClass) > 0 && userClass[0] == "iPXE" {
		// iPXE Breaking the infinite loop
		reply.BootFileName = IPXE_SCRIPT
		reply.UpdateOption(dhcpv4.OptBootFileName(IPXE_SCRIPT))
	}

	if _, err := conn.WriteTo(reply.ToBytes(), peer); err != nil {
		Error("Cannot reply to client: ", err)
	}
	Debug(reply.Summary())
}
