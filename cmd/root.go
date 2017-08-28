// Copyright © 2017 Michael Johnsey <mjohnsey@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"
	"github.com/pkg/errors"
	"log"
	"path/filepath"
	"github.com/mikkeloscar/sshconfig"
	"sort"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

type Hosts []*sshconfig.SSHHost

func (s Hosts) Len() int      { return len(s) }
func (s Hosts) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByHost struct{ Hosts }

func (s ByHost) Less(i, j int) bool { return s.Hosts[i].Host[0] < s.Hosts[j].Host[0] }

var Interactive bool
var PingOnly bool
var hosts Hosts

func (hosts Hosts) PrettyPrintStrings() []string {
	result := make([]string, len(hosts))
	for i, host := range hosts {
		result[i] = fmt.Sprintf("%s - %s\n", host.Host, host.HostName)
	}
	return result
}

func main(cmd *cobra.Command, args []string, isInteractive bool, onlyPing bool) {
	if hosts == nil {
		log.Fatal("Could not get the hosts!")
	}
	if isInteractive {
		fmt.Println("Interactive!")
		if onlyPing{
			fmt.Println("I will only ping not open an ssh session!")
		}
	} else {
		printHosts(hosts)
	}
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
		main(cmd, args, Interactive, PingOnly)
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
	RootCmd.Flags().BoolVarP(&Interactive, "interactive", "i", false, "Run this in interactive mode")
	// This flag gets ignored unless you run in interactive mode (probably a better way of doing it but don't feel like figuring it out)
	RootCmd.Flags().BoolVarP(&PingOnly, "ping", "p", false, "Only run with ping when you select a host")
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
	// sort by host name
	sort.Sort(ByHost{theHosts})
	// set the hosts
	hosts = theHosts
}
