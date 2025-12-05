package main

import (
	"bufio"
	"log"
	"os"

	"github.com/artrctx/pfm/internal/firewall"
	"github.com/artrctx/pfm/internal/upnp"
	"github.com/spf13/cobra"
)

// https://dev.to/pradumnasaraf/how-to-publish-a-golang-package-i12
var rootCmd = &cobra.Command{
	Use:   "pfm",
	Short: "Port forward designated port and make it accesible",
	Long: `Port forward given port so your friend can access your destination. 
	(ONLY SUPPORTS LINUX FOR NOW)
	pfm --port 25565 --protocol tcp --firewall iptables`,
	Run: portForwardMe,
}

var (
	// port to open
	port uint16
	// protocol tcp || udp
	protocol string
	// for firewall config
	firewallProvider string
)

var reader = bufio.NewReader(os.Stdin)

func getKeyPress(input chan<- rune) {
	char, _, err := reader.ReadRune()
	if err != nil {
		log.Panic(err)
	}
	input <- char
}

func portForwardMe(cmd *cobra.Command, args []string) {
	provider, err := firewall.GetProvider(firewallProvider)
	if err != nil {
		log.Panicln(err)
	}
	protocol, err := firewall.GetProtocol(protocol)
	if err != nil {
		log.Panicln(err)
	}

	log.Printf("Initializing firewall client for %v\n", provider)
	fw, err := firewall.New(provider)
	if err != nil {
		log.Panicln(err)
	}
	log.Printf("Adding firewall ruleset for protocol %v | port %v\n", protocol, port)
	ruleset, err := fw.AllowPort(port, protocol)
	if err != nil {
		log.Panicln(err)
	}
	defer ruleset.Close()

	log.Println("Initializing UPNP")
	upnpClient, err := upnp.NewClient()
	if err != nil {
		log.Panicln(err)
	}

	externAddr, err := upnpClient.GetExternalIPAddress()
	if err != nil {
		log.Panicln(err)
	}

	log.Printf("Start UPnP Port Mapping")

	localAddr := upnpClient.LocalAddr()
	// One Hour Lease
	if err := upnpClient.AddPortMapping("", port, string(protocol), port, localAddr.String(), true, "portforwardme", 0); err != nil {
		log.Panicln(err)
	}
	defer upnpClient.DeletePortMapping("", port, string(protocol))

	log.Printf("Port Forwarding Mapped Successfully!\nExtern: %v:%v | Local: %v:%v\nPress q to quit", externAddr, port, localAddr, port)
	input := make(chan rune, 1)
	go getKeyPress(input)

	for {
		i := <-input
		if i == 'q' {
			break
		}
		log.Printf("Pressed %v but only supports q to quit", i)
		go getKeyPress(input)
	}
	log.Println("Stopping Execution...")
}

func init() {
	rootCmd.Flags().Uint16VarP(&port, "port", "p", 0, "Port to forward request to")
	rootCmd.MarkFlagRequired("port")

	rootCmd.Flags().StringVar(&protocol, "protocol", "tcp", "Protocol to listen to (supports tcp or udp)")

	rootCmd.Flags().StringVarP(&firewallProvider, "firewall", "f", "iptables", "Firewall provider name (currently supports iptables only)")
}

func main() {
	rootCmd.Execute()
}
