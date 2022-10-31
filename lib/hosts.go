package lib

import (
	"encoding/json"
	"fmt"

	"github.com/mikkeloscar/sshconfig"
)

type Hosts []*sshconfig.SSHHost

func (s Hosts) Len() int      { return len(s) }
func (s Hosts) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByHost struct{ Hosts }

func (s ByHost) Less(i, j int) bool { return s.Hosts[i].Host[0] < s.Hosts[j].Host[0] }

func (hosts Hosts) PrettyPrintStrings() []string {
	result := make([]string, len(hosts))
	for i, host := range hosts {
		result[i] = fmt.Sprintf("%s - %s (%s)\n", host.Host[0], host.HostName, host.User)
	}
	return result
}

func (hosts Hosts) AsJsonString() *string {
	result := make(map[string]map[string]string)
	for _, host := range hosts {
		result[host.Host[0]] = make(map[string]string)
		result[host.Host[0]]["hostname"] = host.HostName
		result[host.Host[0]]["user"] = host.User
	}
	jsonString, _ := json.Marshal(result)
	asStr := string(jsonString)
	return &asStr
}
