/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "babydump",
	Short: "A simple Linux packet sniffer using raw sockets",
	Long: `A low-level Linux packet sniffer using raw sockets,
achieving tcpdump/Wireshark-like parsing without 
relying on libpcap.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		print(args)
		fmt.Println("Hello, World!")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.babydump.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// var Verbose bool
	// rootCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	// print(Verbose)

	rootCmd.Flags().Bool("list-interfaces", false, "print the list of network interfaces available on the system")

	var NInterface string
	rootCmd.Flags().StringVarP(&NInterface, "interface", "i", "", "the network interface to caputure on (required)")
	rootCmd.MarkFlagsOneRequired("list-interfaces", "interface")

}
