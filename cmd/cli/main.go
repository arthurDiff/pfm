package main

import (
	"github.com/arthurDiff/pfm/internal/stun"
	"github.com/spf13/cobra"
)

// https://github.com/spf13/cobra
var rootCmd = &cobra.Command{
	Use:   "pfm",
	Short: "Port forward designated port and make it accesible",
	Long: `Port forward given port so your friend can access your destination.
	pfm --platform linux --port --stun stun:stun1.l.google.com:3478`,
	Run: portForwardMe,
}

var (
	// for firewall config
	platform string
	port     uint16
	// for stun server
	stunAddr string
)

func portForwardMe(cmd *cobra.Command, args []string) {
	stunClient := stun.NewClient(stunAddr)
	defer stunClient.Close()

	// TODO
	// CONFIG FIREWALL TO PROVIDED OPEN PORT (Defer reset firewall setting)
	// IN FOR LOOP
	// - SHOULD GET CURRENTLY CONFIGURED DNS IP
	// - GET ACTUAL IP
	// -- IF MISMATCH UPDATE CONFIG
	// MAYBE PANIC NOTIFIER
}

func init() {
	rootCmd.Flags().StringVar(&platform, "platform", "", "OS Platform (supports linux or windows)")
	rootCmd.MarkFlagRequired("platform")

	rootCmd.Flags().Uint16VarP(&port, "port", "p", 0, "Port to forward request to")
	rootCmd.MarkFlagRequired("port")

	rootCmd.Flags().StringVarP(&stunAddr, "stun", "s", "stun:stun1.l.google.com:3478", "STUN address to use")
}

func main() {
	rootCmd.Execute()
}
