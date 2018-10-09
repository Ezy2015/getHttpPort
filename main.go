package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	snapshotLen int32 = 1024
	promiscuous bool  = false
	err         error
	timeout     time.Duration = 30 * time.Second
	handle      *pcap.Handle
)

func portFormat(port string) string {
	var p string
	if strings.Contains(port, "(") {
		p = port[:strings.Index(port, "(")]
	} else {
		p = port
	}
	return p
}

func printPacketInfo(packet gopacket.Packet, cfg *Config, logger *log.Logger) error {
	var flag = false
	var srcPort string

	applicationLayer := packet.ApplicationLayer()
	if applicationLayer != nil {
		if len(string(applicationLayer.Payload())) > 5 && string(applicationLayer.Payload())[:5] == "HTTP/" {
			ipLayer := packet.Layer(layers.LayerTypeIPv4)
			if ipLayer != nil {
				ip, _ := ipLayer.(*layers.IPv4) // 找到IP
				var tmp = false
				for _, v := range cfg.DestNets {
					_, net, _ := net.ParseCIDR(v)
					if net.Contains(ip.SrcIP) {
						tmp = true
					}
				}

				if tmp {
					tcpLayer := packet.Layer(layers.LayerTypeTCP)
					if tcpLayer != nil {
						tcp, _ := tcpLayer.(*layers.TCP)
						srcPort = portFormat(tcp.SrcPort.String())
						for _, k := range cfg.Ports {
							if k == srcPort {
								flag = true
							}
						}
						if !flag {
							logger.Println(ip.SrcIP, srcPort)
							cfg.Ports = append(cfg.Ports, srcPort)
							// appendToFile(cfg.LogName, ip.SrcIP.String()+" "+string(srcPort)+"\n")
						}
					}
				}
			}
		}
	}
	return nil
}

func main() {
	cfg := getConf()
	device := cfg.Nic
	logFile, err := os.Create(cfg.LogName)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error")
	}
	logger := log.New(logFile, "[Info]", log.LstdFlags)

	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		printPacketInfo(packet, cfg, logger)
	}
}
