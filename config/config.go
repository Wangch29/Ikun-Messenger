package config

import "fmt"

type Config struct {
	ID    int    `mapstructure:"id"`
	Nodes []Node `mapstructure:"nodes"`
}

type Node struct {
	ID         int    `mapstructure:"id"`
	Host       string `mapstructure:"host"`
	RaftPort   int    `mapstructure:"raft_port"`
	KVPort     int    `mapstructure:"kv_port"`
	IMHttpPort int    `mapstructure:"im_http_port"`
	IMGrpcPort int    `mapstructure:"im_grpc_port"`
}

// global config
var Global Config

func (n Node) RaftAddr() string {
	return fmt.Sprintf("%s:%d", n.Host, n.RaftPort)
}

func (n Node) KVAddr() string {
	return fmt.Sprintf("%s:%d", n.Host, n.KVPort)
}

func (n Node) IMHttpAddr() string {
	return fmt.Sprintf("%s:%d", n.Host, n.IMHttpPort)
}

func (n Node) IMGrpcAddr() string {
	return fmt.Sprintf("%s:%d", n.Host, n.IMGrpcPort)
}
