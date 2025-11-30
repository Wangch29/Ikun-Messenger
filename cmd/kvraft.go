package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Wangch29/IkunMessenger/kvraft"
	"github.com/Wangch29/IkunMessenger/raft"
	"github.com/spf13/cobra"
)

var (
	kvraftMe int
)

// raft node addresses
var raftPeers = []string{
	"127.0.0.1:25000",
	"127.0.0.1:25001",
	"127.0.0.1:25002",
}

// kv node addresses
var kvPeers = []string{
	"127.0.0.1:6000",
	"127.0.0.1:6001",
	"127.0.0.1:6002",
}

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
	if kvraftMe < 0 || kvraftMe >= len(raftPeers) {
		log.Fatalf("Invalid node ID: %d", kvraftMe)
	}

	// 1. Start Raft
	applyCh := make(chan raft.ApplyMsg)
	rf := raft.Make(raftPeers, kvraftMe, raft.NewMemoryStorage(), applyCh)

	go func() {
		parts := strings.Split(raftPeers[kvraftMe], ":")
		if err := rf.StartServer(":" + parts[1]); err != nil {
			log.Fatalf("Raft server failed: %v", err)
		}
	}()

	// 2. Start KV Server
	kv := kvraft.NewKVServer(kvraftMe, rf, applyCh, 1000)
	go func() {
		parts := strings.Split(kvPeers[kvraftMe], ":")
		if err := kv.StartKVServer(":" + parts[1]); err != nil {
			log.Fatalf("KV server failed: %v", err)
		}
	}()

	log.Printf("KVRaft Node %d started. Raft: %s, KV: %s", kvraftMe, raftPeers[kvraftMe], kvPeers[kvraftMe])

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
