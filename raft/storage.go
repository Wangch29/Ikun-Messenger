package raft

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Storage interface {
	Save(raftState []byte, snapshot []byte)
	ReadRaftState() []byte
	ReadSnapshot() []byte
	RaftStateSize() int
}

type MemoryStorage struct {
	mu        sync.Mutex
	raftState []byte
	snapshot  []byte
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (ms *MemoryStorage) Save(raftState []byte, snapshot []byte) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if raftState != nil {
		ms.raftState = raftState
	}
	if snapshot != nil {
		ms.snapshot = snapshot
	}
}

func (ms *MemoryStorage) ReadRaftState() []byte {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	return ms.raftState
}

func (ms *MemoryStorage) ReadSnapshot() []byte {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	return ms.snapshot
}

func (ms *MemoryStorage) RaftStateSize() int {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	return len(ms.raftState)
}

type FileStorage struct {
	mu        sync.Mutex
	stateFile string
	snapFile  string
}

func NewFileStorage(dir string, nodeId int) *FileStorage {
	os.MkdirAll(dir, 0755)
	return &FileStorage{
		stateFile: filepath.Join(dir, fmt.Sprintf("raft-%d-state.bin", nodeId)),
		snapFile:  filepath.Join(dir, fmt.Sprintf("raft-%d-snap.bin", nodeId)),
	}
}

func (fs *FileStorage) Save(raftState []byte, snapshot []byte) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	if raftState != nil {
		_ = os.WriteFile(fs.stateFile, raftState, 0644)
	}
	if snapshot != nil {
		_ = os.WriteFile(fs.snapFile, snapshot, 0644)
	}
}

func (fs *FileStorage) ReadRaftState() []byte {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	data, _ := os.ReadFile(fs.stateFile)
	return data
}

func (fs *FileStorage) ReadSnapshot() []byte {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	data, _ := os.ReadFile(fs.snapFile)
	return data
}

func (fs *FileStorage) RaftStateSize() int {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	info, err := os.Stat(fs.stateFile)
	if err != nil {
		return 0
	}
	return int(info.Size())
}
