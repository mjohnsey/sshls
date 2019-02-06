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
)

type Hosts []*sshconfig.SSHHost

func (s Hosts) Len() int      { return len(s) }
func (s Hosts) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByHost struct{ Hosts }

func (s ByHost) Less(i, j int) bool { return s.Hosts[i].Host[0] < s.Hosts[j].Host[0] }

var hosts Hosts

func (hosts Hosts) PrettyPrintStrings() []string {
	result := make([]string, len(hosts))
	for i, host := range hosts {
		result[i] = fmt.Sprintf("%s - %s\n", host.Host, host.HostName)
	}
	return result
}

func Run(cmd *cobra.Command, args []string) {
	if hosts == nil {
		log.Fatal("Could not get the hosts!")
	}
	printHosts(hosts)
}

func printHosts(hosts Hosts) {
	for _, host := range Hosts.PrettyPrintStrings(hosts) {
		fmt.Printf(host)
	}
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "sshls",
	Short: "This lists all your ssh hosts",
	Run: func(cmd *cobra.Command, args []string) {
		Run(cmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Cobra command run on initialization
func init() {
	cobra.OnInitialize(getHosts)
}

// This will get the ssh config file and set the objects
func getHosts() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Problem getting the home dir!"))
	}
	// Grab the ssh config and parse it
	sshConfigFile := filepath.Join(home, ".ssh", "config")
	theHosts, err := sshconfig.ParseSSHConfig(sshConfigFile)

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
	sort.Sort(ByHost{whitelist})
	// set the hosts
	hosts = whitelist
}
