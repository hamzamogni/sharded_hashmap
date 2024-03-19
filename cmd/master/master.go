package main

import (
	"hash/fnv"
	"os"
	"strings"
)

type Node struct {
	HOST string
	PORT string
}

func (s *Node) String() string {
	return s.HOST + ":" + s.PORT
}

type Master struct {
	Nodes []Node
}

func NewMaster() *Master {
	m := &Master{}

	shards_env := os.Getenv("shards")
	shards := strings.Split(shards_env, ",")
	for _, shard := range shards {
		host_port := strings.Split(shard, ":")
		m.Nodes = append(m.Nodes, Node{host_port[0], host_port[1]})
	}

	return m
}

func (m *Master) GetNodeString() string {
	var shards []string
	for _, shard := range m.Nodes {
		shards = append(shards, shard.String())
	}
	return strings.Join(shards, ",")
}

// GetNode returns the shard node for a given key
func (m *Master) GetNode(key string) Node {
    hash := m.hashKey(key)
    return m.Nodes[hash%uint32(len(m.Nodes))]
}

func (m *Master) hashKey(key string) uint32 {
    h := fnv.New32a()
    h.Write([]byte(key))
    return h.Sum32()
}

