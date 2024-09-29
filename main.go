package main

import (
	"embed"
	"io"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
	tftp "github.com/pin/tftp/v3"
)

//go:embed tftpboot
var tftpRoot embed.FS

func handler(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
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
		reply.YourIPAddr = net.IPv4(10, 0, 2, 100)
		reply.ServerIPAddr = net.IPv4(10, 0, 2, 5)
		reply.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeOffer))
		reply.UpdateOption(dhcpv4.OptServerIdentifier(net.IPv4(10, 0, 2, 5)))
		reply.UpdateOption(dhcpv4.OptSubnetMask(net.IPv4Mask(255, 255, 255, 0)))
		reply.UpdateOption(dhcpv4.OptBroadcastAddress(net.IPv4(10, 0, 2, 255)))
		reply.UpdateOption(dhcpv4.OptDNS(net.IPv4(10, 0, 2, 1)))
		reply.UpdateOption(dhcpv4.OptRouter(net.IPv4(10, 0, 2, 1)))
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
	case dhcpv4.MessageTypeRequest:
		reply.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeAck))
		reply.YourIPAddr = net.IPv4(10, 0, 2, 100)
		reply.ServerIPAddr = net.IPv4(10, 0, 2, 5)
		reply.UpdateOption(dhcpv4.OptServerIdentifier(net.IPv4(10, 0, 2, 5)))
		reply.UpdateOption(dhcpv4.OptSubnetMask(net.IPv4Mask(255, 255, 255, 0)))
		reply.UpdateOption(dhcpv4.OptBroadcastAddress(net.IPv4(10, 0, 2, 255)))
		reply.UpdateOption(dhcpv4.OptDNS(net.IPv4(10, 0, 2, 1)))
		reply.UpdateOption(dhcpv4.OptRouter(net.IPv4(10, 0, 2, 1)))
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
		Info("Unhandled message type: ", mt)
		return
	}

	if _, err := conn.WriteTo(reply.ToBytes(), peer); err != nil {
		Error("Cannot reply to client: ", err)
	}
	Debug(reply.Summary())
}

func overWrite(filename string) string {
	Debug("Raw filename: ", filename)
	if filepath.IsAbs(filename) {
		filename = strings.Replace(filename, "/", "", 1)
	}
	Debug("overWrited filename: ", filename)
	return filename
}

// readHandler is called when client starts file download from server
func readHandler(filename string, rf io.ReaderFrom) error {
	// enter root filesystem
	root, _ := fs.Sub(tftpRoot, "tftpboot")
	// use relative path to access file
	filename = overWrite(filename)
	file, err := root.Open(filename)
	if err != nil {
		Error(err)
		return err
	}
	n, err := rf.ReadFrom(file)
	if err != nil {
		Error(err)
		return err
	}
	Info("sent ", n, " bytes")
	return nil
}

func runtftp() {
	// use nil in place of handler to disable read or write operations
	s := tftp.NewServer(readHandler, nil)
	s.SetTimeout(5 * time.Second)  // optional
	err := s.ListenAndServe(":69") // blocks until s.Shutdown() is called
	if err != nil {
		Panic("server: ", err)
	}
}

func main() {
	Initial("debug", os.Stdout)
	laddr := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 67,
	}
	server, err := server4.NewServer("enp0s3", &laddr, handler)
	if err != nil {
		Fatal(err)
	}

	// This never returns. If you want to do other stuff, dump it into a
	// goroutine.
	go server.Serve()
	runtftp()
}
