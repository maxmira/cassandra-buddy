package nodetool

import (
	"regexp"
	"strings"
)

var datacenterPat = regexp.MustCompile(`^\s*Datacenter: (.+)$`)
var nodePat = regexp.MustCompile(`^\s*([UD][NLJM])\s+([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)\s+([0-9\.]+ [PGMK]?B)\s+([0-9]+)\s+([0-9\?\.\%]+)\s+([a-zA-Z0-9\-]+)\s+(.+)$`)

type Node struct {
	State   string
	Address string
	Load    string
	Tokens  string
	Owns    string
	HostID  string
	Rack    string
}
type Datacenter struct {
	Name  string
	Nodes []Node
}

type Status struct {
	Datacenters []Datacenter
}

func NewStatus(d []byte) *Status {
	datacenters := make([]Datacenter, 0)

	for _, line := range strings.Split(string(d), "\n") {

		if dcParts := datacenterPat.FindAllStringSubmatch(line, 2); dcParts != nil {
			if len(dcParts[0]) != 2 {
				continue
			}
			datacenters = append(datacenters, Datacenter{Name: dcParts[0][1], Nodes: make([]Node, 0)})
			continue
		}
		if len(datacenters) < 1 {
			continue //without a DC we can't do much else
		}
		if nodeParts := nodePat.FindAllStringSubmatch(line, 7); nodeParts != nil {
			if len(nodeParts[0]) != 8 {
				continue
			}
			datacenters[len(datacenters)-1].Nodes = append(datacenters[len(datacenters)-1].Nodes, Node{State: nodeParts[0][1], Address: nodeParts[0][2], Load: nodeParts[0][3], Tokens: nodeParts[0][4], Owns: nodeParts[0][5], HostID: nodeParts[0][6], Rack: nodeParts[0][7]})
			continue
		}

	}

	return &Status{Datacenters: datacenters}
}
