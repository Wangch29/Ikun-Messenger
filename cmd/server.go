package cmd

import (
	"log/slog"
	"net"
	"strconv"
	"strings"

	"github.com/Wangch29/ikun-messenger/api/impb"
	"github.com/Wangch29/ikun-messenger/im"
	"github.com/Wangch29/ikun-messenger/kvraft"
	"github.com/Wangch29/ikun-messenger/raft"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var serverMe int

// TODO: For local development only.
var imPorts = []string{"8080", "8081", "8082"}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the ikun messenger server",
	Long:  "Start the ikun messenger server",
	Run:   runServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVarP(&serverMe, "me", "m", 0, "Node ID (0, 1, 2)")
}

func runServer(cmd *cobra.Command, args []string) {
	if serverMe < 0 || serverMe >= len(imPorts) {
		slog.Error("Invalid node ID: %d", serverMe)
		return
	}

	applyCh := make(chan raft.ApplyMsg)
	rf := raft.Make(raftPeers, serverMe, raft.NewMemoryStorage(), applyCh)

	// Start Raft.
	go func() {
		_, port, found := strings.Cut(raftPeers[serverMe], ":")
		if !found {
			slog.Error("Invalid address format (missing port)", "addr", raftPeers[serverMe])
			return
		}
		if err := rf.StartServer(":" + port); err != nil {
			slog.Error("Failed to start raft server", "err", err)
			return
		}
	}()

	// Start KV server.
	kv := kvraft.NewKVServer(serverMe, rf, applyCh, 1000)
	go func() {
		_, port, found := strings.Cut(kvPeers[serverMe], ":")
		if !found {
			slog.Error("Invalid address format (missing port)", "addr", kvPeers[serverMe])
			return
		}
		if err := kv.StartKVServer(":" + port); err != nil {
			slog.Error("Failed to start kv server", "err", err)
			return
		}
	}()

	// create clerk.
	ck := kvraft.MakeClerk(kvPeers, int64(serverMe))

	nodeAddr := "127.0.0.1:" + imPorts[serverMe]
	imServer := im.NewIMServer(serverMe, ck, nodeAddr)

	// Start IM gRPC server.
	grpcPort := 17000 + serverMe
	go func() {
		lis, err := net.Listen("tcp", ":"+strconv.Itoa(grpcPort))
		if err != nil {
			slog.Error("Failed to listen", "err", err)
			return
		}
		slog.Info("IM gRPC Server listening", ":", grpcPort)
		s := grpc.NewServer()
		impb.RegisterIMServiceServer(s, imServer)
		if err := s.Serve(lis); err != nil {
			slog.Error("Failed to serve", "err", err)
		}
	}()

	slog.Info("Node started",
		"id", serverMe,
		"raft_addr", raftPeers[serverMe],
		"kv_addr", kvPeers[serverMe],
		"im_addr", nodeAddr,
	)

	// start IM server.
	if err := imServer.Start(nodeAddr); err != nil {
		slog.Error("Failed to start im server", "err", err)
		return
	}
}
