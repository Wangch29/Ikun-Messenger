package config

type Config struct {
	ID    int    `mapstructure:"id"`
	Nodes []Node `mapstructure:"nodes"`
}

type Node struct {
	ID         int    `mapstructure:"id"`
	RaftAddr   string `mapstructure:"raft_addr"`
	KVAddr     string `mapstructure:"kv_addr"`
	IMHttpAddr string `mapstructure:"im_http_addr"`
	IMGrpcAddr string `mapstructure:"im_grpc_addr"`
}

// global config
var Global Config
