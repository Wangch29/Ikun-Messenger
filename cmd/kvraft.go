package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Wangch29/ikun-messenger/config"
	"github.com/Wangch29/ikun-messenger/kvraft"
	"github.com/Wangch29/ikun-messenger/raft"
	"github.com/spf13/cobra"
)

var (
	kvraftMe int
)

var kvraftCmd = &cobra.Command{
	Use:   "kvraft",
	Short: "Start KVRaft test node",
	Long:  `Start a KVRaft node for testing the distributed KV storage.`,
	Run:   runKVRaft,
}

func init() {
	rootCmd.AddCommand(kvraftCmd)
	kvraftCmd.Flags().IntVarP(&kvraftMe, "me", "m", 0, "Node ID (0, 1, 2)")
}

func runKVRaft(cmd *cobra.Command, args []string) {
	if len(config.Global.Nodes) == 0 {
		log.Fatalf("No nodes found in config")
	}

	if kvraftMe < 0 || kvraftMe >= len(config.Global.Nodes) {
		log.Fatalf("Invalid node ID: %d", kvraftMe)
	}

	// Construct peer lists from config
	var raftPeers []string
	var kvPeers []string
	for _, node := range config.Global.Nodes {
		raftPeers = append(raftPeers, node.RaftAddr)
		kvPeers = append(kvPeers, node.KVAddr)
	}

	myConfig := config.Global.Nodes[kvraftMe]

	// 1. Start Raft
	applyCh := make(chan raft.ApplyMsg)
	rf := raft.Make(raftPeers, kvraftMe, raft.NewMemoryStorage(), applyCh)

	go func() {
		_, port, found := strings.Cut(myConfig.RaftAddr, ":")
		if !found {
			log.Fatalf("Invalid raft address format: %s", myConfig.RaftAddr)
		}
		if err := rf.StartServer(":" + port); err != nil {
			log.Fatalf("Raft server failed: %v", err)
		}
	}()

	// 2. Start KV Server
	kv := kvraft.NewKVServer(kvraftMe, rf, applyCh, 1000)
	go func() {
		_, port, found := strings.Cut(myConfig.KVAddr, ":")
		if !found {
			log.Fatalf("Invalid kv address format: %s", myConfig.KVAddr)
		}
		if err := kv.StartKVServer(":" + port); err != nil {
			log.Fatalf("KV server failed: %v", err)
		}
	}()

	log.Printf("KVRaft Node %d started. Raft: %s, KV: %s", kvraftMe, myConfig.RaftAddr, myConfig.KVAddr)

	// 3. Create Clerk
	ck := kvraft.MakeClerk(kvPeers, int64(kvraftMe))

	// 4. REPL
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		parts := strings.Split(input, " ")

		switch parts[0] {
		case "put":
			if len(parts) != 3 {
				fmt.Println("Usage: put <key> <value>")
				continue
			}
			ck.Put(parts[1], parts[2])
			fmt.Println("Put OK")

		case "get":
			if len(parts) != 2 {
				fmt.Println("Usage: get <key>")
				continue
			}
			value := ck.Get(parts[1])
			fmt.Printf("Value: %s\n", value)

		case "exit":
			os.Exit(0)

		default:
			fmt.Println("Unknown command. Available: put, get, exit")
		}
	}
}
