package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/mikkeloscar/sshconfig"
	"github.com/pkg/errors"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	lib "github.com/mjohnsey/sshls/lib"
)

var hosts lib.Hosts
var jsonFormat bool
var sshFilePath string

func Run(cmd *cobra.Command, args []string) {
	if hosts == nil {
		log.Fatal("Could not get the hosts!")
	}
	printHosts(hosts)
}

func printHosts(hosts lib.Hosts) {
	if jsonFormat {
		fmt.Println(*hosts.AsJsonString())
	} else {
		for _, host := range lib.Hosts.PrettyPrintStrings(hosts) {
			fmt.Printf(host)
		}
	}
}

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sshls",
	Short: "This lists all your ssh hosts",
	Run: func(cmd *cobra.Command, args []string) {
		getHosts(&sshFilePath)
		Run(cmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Cobra command run on initialization
func init() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Problem getting the home dir!"))
	}
	// Grab the ssh config and parse it
	sshConfigFile := filepath.Join(home, ".ssh", "config")

	rootCmd.PersistentFlags().BoolVar(&jsonFormat, "json", false, "format to json")
	rootCmd.PersistentFlags().StringVar(&sshFilePath, "configFile", sshConfigFile, "path to ssh config (defaults to ~/.ssh/config")
}

// This will get the ssh config file and set the objects
func getHosts(sshConfigFile *string) {
	// Grab the ssh config and parse it
	theHosts, err := sshconfig.ParseSSHConfig(*sshConfigFile)
	if err != nil {
		log.Fatalln("Could not open the SSH config file")
	}

	whitelist := make([]*sshconfig.SSHHost, 0)

	blacklist := map[string]bool{
		"*": false,
	}

	for _, host := range theHosts {
		if _, ok := blacklist[host.Host[0]]; !ok {
			whitelist = append(whitelist, host)
		}
	}

	// sort by host name
	sort.Sort(lib.ByHost{whitelist})
	// set the hosts
	hosts = whitelist
}
