package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// State represents the current state of a Raft node.
type State int

const (
	Follower  State = iota
	Candidate
	Leader
)

// LogEntry represents an entry in the Raft log.
type LogEntry struct {
	Term        int
	CommandType string // e.g., "PUT_FILE", "DELETE_FILE"
	CommandData []byte // JSON encoded command data
}

// FileStateMachine represents the application state that Raft replicates.
type FileStateMachine struct {
	mu      sync.Mutex
	baseDir string // Directory where files are stored
}

// NewFileStateMachine creates a new file state machine.
func NewFileStateMachine(baseDir string) *FileStateMachine {
	return &FileStateMachine{
		baseDir: baseDir,
	}
}

// Apply applies a log entry to the state machine.
func (fsm *FileStateMachine) Apply(entry LogEntry) error {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	log.Printf("Applying command %s to state machine: %s", entry.CommandType, string(entry.CommandData))

	switch entry.CommandType {
	case "PUT_FILE":
		var cmd PutFileCommand
		if err := json.Unmarshal(entry.CommandData, &cmd); err != nil {
			return fmt.Errorf("failed to unmarshal PutFileCommand: %w", err)
		}
		filePath := filepath.Join(fsm.baseDir, cmd.Filename)
		return os.WriteFile(filePath, cmd.Content, 0644)
	case "DELETE_FILE":
		var cmd DeleteFileCommand
		if err := json.Unmarshal(entry.CommandData, &cmd); err != nil {
			return fmt.Errorf("failed to unmarshal DeleteFileCommand: %w", err)
		}
		filePath := filepath.Join(fsm.baseDir, cmd.Filename)
		return os.Remove(filePath)
	default:
		return fmt.Errorf("unknown command type: %s", entry.CommandType)
	}
}

// PutFileCommand represents a command to put a file.
type PutFileCommand struct {
	Filename string
	Content  []byte
}

// DeleteFileCommand represents a command to delete a file.
type DeleteFileCommand struct {
	Filename string
}

// Raft represents a single Raft node.
type Raft struct {
	mu        sync.Mutex          // Mutex to protect shared state
	id        int                 // Unique ID of this Raft node
	peers     []string            // Network addresses of other Raft nodes
	state     State               // Current state of the Raft node
	currentTerm int                 // Current term number
	votedFor  int                 // Candidate ID that received vote in current term
	leaderId  int                 // Current leader's ID

	log         []LogEntry          // The Raft log
	commitIndex int                 // Index of highest log entry known to be committed
	lastApplied int                 // Index of highest log entry applied to state machine

	// For leaders
	nextIndex  []int // For each server, index of the next log entry to send to that server
	matchIndex []int // For each server, index of highest log entry known to be replicated on server

	// Election timeout
	electionTimeout time.Duration
	lastHeartbeat   time.Time

	// RPC client connections to peers
	peerClients []*rpc.Client

	// State machine for applying committed commands
	stateMachine *FileStateMachine
}

// RequestVoteArgs is the arguments for a RequestVote RPC.
type RequestVoteArgs struct {
	Term        int // candidate's term
	CandidateId int // candidate requesting vote
	LastLogIndex int // index of candidate's last log entry
	LastLogTerm  int // term of candidate's last log entry
}

// RequestVoteReply is the reply for a RequestVote RPC.
type RequestVoteReply struct {
	Term        int  // currentTerm, for candidate to update itself
	VoteGranted bool // true means candidate received vote
}

// AppendEntriesArgs is the arguments for an AppendEntries RPC (heartbeat or log entries).
type AppendEntriesArgs struct {
	Term         int        // leader's term
	LeaderId     int        // so follower can redirect clients
	PrevLogIndex int        // index of log entry immediately preceding new ones
	PrevLogTerm  int        // term of PrevLogIndex entry
	Entries      []LogEntry // log entries to store (empty for heartbeat; may send more than one for efficiency)
	LeaderCommit int        // leader's commitIndex
}

// AppendEntriesReply is the reply for an AppendEntries RPC.
type AppendEntriesReply struct {
	Term    int  // currentTerm, for leader to update itself
	Success bool // true if follower contained entry matching prevLogIndex and prevLogTerm
	XTerm   int  // term in the conflicting entry (if any)
	XIndex  int  // index of first entry with that term (if any)
	XLen    int  // log length (if no conflict)
}

// NewRaft creates a new Raft node.
func NewRaft(id int, peers []string, baseDir string) *Raft {
	rf := &Raft{
		id:        id,
		peers:     peers,
		state:     Follower,
		currentTerm: 0,
		votedFor:  -1,
		leaderId:  -1,
		log:         make([]LogEntry, 1), // Log is 1-indexed, so 0th entry is dummy
		commitIndex: 0,
		lastApplied: 0,
		nextIndex:   make([]int, len(peers)),
		matchIndex:  make([]int, len(peers)),
		electionTimeout: time.Duration(150+rand.Intn(150)) * time.Millisecond,
		lastHeartbeat:   time.Now(),
		stateMachine: NewFileStateMachine(baseDir),
	}
	return rf
}

// RequestVote RPC handler.
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) error {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	log.Printf("Node %d (Term %d, State %s) received RequestVote from %d (Term %d, LastLogIndex %d, LastLogTerm %d)",
		rf.id, rf.currentTerm, rf.state, args.CandidateId, args.Term, args.LastLogIndex, args.LastLogTerm)

	reply.Term = rf.currentTerm
	reply.VoteGranted = false

	// 1. Reply false if args.Term < currentTerm
	if args.Term < rf.currentTerm {
		return nil
	}

	// If args.Term > currentTerm, convert to follower
	if args.Term > rf.currentTerm {
		rf.becomeFollower(args.Term)
	}

	lastLogIndex := len(rf.log) - 1
	lastLogTerm := rf.log[lastLogIndex].Term

	// 2. If votedFor is null or candidateId, and candidate's log is at least as up-to-date as receiver's log, grant vote
	if (rf.votedFor == -1 || rf.votedFor == args.CandidateId) &&
		(args.LastLogTerm > lastLogTerm || (args.LastLogTerm == lastLogTerm && args.LastLogIndex >= lastLogIndex)) {
		rf.votedFor = args.CandidateId
		reply.VoteGranted = true
		log.Printf("Node %d (Term %d, State %s) granted vote to %d", rf.id, rf.currentTerm, rf.state, args.CandidateId)
		rf.lastHeartbeat = time.Now() // Reset election timer on granting vote
	}
	return nil
}

// AppendEntries RPC handler.
func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) error {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	log.Printf("Node %d (Term %d, State %s) received AppendEntries from %d (Term %d, PrevLogIndex %d, PrevLogTerm %d, Entries %d, LeaderCommit %d)",
		rf.id, rf.currentTerm, rf.state, args.LeaderId, args.Term, args.PrevLogIndex, args.PrevLogTerm, len(args.Entries), args.LeaderCommit)

	reply.Term = rf.currentTerm
	reply.Success = false

	// 1. Reply false if args.Term < currentTerm
	if args.Term < rf.currentTerm {
		return nil
	}

	// If leader's term is greater or equal, become follower and reset election timer
	if args.Term >= rf.currentTerm {
		rf.becomeFollower(args.Term)
		rf.leaderId = args.LeaderId
		rf.lastHeartbeat = time.Now()
	}

	// 2. Reply false if log doesn't contain an entry at PrevLogIndex whose term matches PrevLogTerm
	if args.PrevLogIndex >= len(rf.log) || rf.log[args.PrevLogIndex].Term != args.PrevLogTerm {
		reply.XTerm = -1
		reply.XIndex = -1
		reply.XLen = len(rf.log)
		if args.PrevLogIndex < len(rf.log) {
			reply.XTerm = rf.log[args.PrevLogIndex].Term
			// Find the first index for XTerm
			for i := 1; i <= args.PrevLogIndex; i++ {
				if rf.log[i].Term == reply.XTerm {
					reply.XIndex = i
					break
				}
			}
		}
		return nil
	}

	// 3. If an existing entry conflicts with a new one (same index but different terms),
	// delete the existing entry and all that follow it
	i := 0
	for ; i < len(args.Entries); i++ {
		logIndex := args.PrevLogIndex + 1 + i
		if logIndex < len(rf.log) {
			if rf.log[logIndex].Term != args.Entries[i].Term {
				rf.log = rf.log[:logIndex] // Delete conflicting entry and all that follow
				break
			}
		} else {
			break // Reached end of current log, append new entries
		}
	}

	// Append any new entries not already in the log
	rf.log = append(rf.log, args.Entries[i:]...)
	reply.Success = true

	// 4. If leaderCommit > commitIndex, set commitIndex = min(leaderCommit, index of last new entry)
	if args.LeaderCommit > rf.commitIndex {
		lastNewEntryIndex := len(rf.log) - 1
		rf.commitIndex = min(args.LeaderCommit, lastNewEntryIndex)
		rf.applyCommittedEntries()
	}

	return nil
}

// applyCommittedEntries applies committed log entries to the state machine.
func (rf *Raft) applyCommittedEntries() {
	for rf.lastApplied < rf.commitIndex {
		rf.lastApplied++
		entry := rf.log[rf.lastApplied]
		if err := rf.stateMachine.Apply(entry); err != nil {
			log.Printf("Error applying log entry %d to state machine: %v", rf.lastApplied, err)
		}
	}
}

// min returns the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// becomeFollower transitions the node to Follower state.
func (rf *Raft) becomeFollower(term int) {
	rf.state = Follower
	rf.currentTerm = term
	rf.votedFor = -1
	rf.leaderId = -1
	log.Printf("Node %d transitioned to Follower (Term %d)", rf.id, rf.currentTerm)
}

// becomeCandidate transitions the node to Candidate state.
func (rf *Raft) becomeCandidate() {
	rf.state = Candidate
	rf.currentTerm++
	rf.votedFor = rf.id
	rf.leaderId = -1
	rf.lastHeartbeat = time.Now() // Reset election timer
	log.Printf("Node %d transitioned to Candidate (Term %d)", rf.id, rf.currentTerm)
}

// becomeLeader transitions the node to Leader state.
func (rf *Raft) becomeLeader() {
	rf.state = Leader
	rf.leaderId = rf.id
	log.Printf("Node %d transitioned to Leader (Term %d)", rf.id, rf.currentTerm)
	// Initialize nextIndex and matchIndex for all followers
	lastLogIndex := len(rf.log) - 1
	for i := 0; i < len(rf.peers); i++ {
		rf.nextIndex[i] = lastLogIndex + 1
		rf.matchIndex[i] = 0
	}
	// Send initial heartbeats
	rf.sendHeartbeats()
	go rf.startHeartbeatTimer()
}

// startHeartbeatTimer sends heartbeats periodically.
func (rf *Raft) startHeartbeatTimer() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		rf.mu.Lock()
		if rf.state != Leader {
			rf.mu.Unlock()
			return
		}
		rf.sendHeartbeats()
		rf.mu.Unlock()
	}
}

// sendHeartbeats sends AppendEntries RPCs to all peers.
func (rf *Raft) sendHeartbeats() {
	for i, peer := range rf.peers {
		if i == rf.id { // Don't send to self
			continue
		}

		// For each follower, send AppendEntries RPC with log entries starting from nextIndex[i]
		prevLogIndex := rf.nextIndex[i] - 1
		prevLogTerm := rf.log[prevLogIndex].Term

		args := &AppendEntriesArgs{
			Term:         rf.currentTerm,
			LeaderId:     rf.id,
			PrevLogIndex: prevLogIndex,
			PrevLogTerm:  prevLogTerm,
			Entries:      rf.log[rf.nextIndex[i]:],
			LeaderCommit: rf.commitIndex,
		}

		go func(peerAddr string, peerId int) {
			reply := &AppendEntriesReply{}
			log.Printf("Node %d (Leader) sending AppendEntries to %s (Term %d, PrevLogIndex %d, Entries %d)", rf.id, peerAddr, args.Term, args.PrevLogIndex, len(args.Entries))
			err := rf.call(peerAddr, "Raft.AppendEntries", args, reply)
			if err != nil {
				log.Printf("Node %d (Leader) failed to send AppendEntries to %s: %v", rf.id, peerAddr, err)
				return
			}

			rf.mu.Lock()
			defer rf.mu.Unlock()

			if rf.state != Leader || rf.currentTerm != args.Term {
				return // Leader changed or term changed
			}

			if reply.Term > rf.currentTerm {
				log.Printf("Node %d (Leader) discovered higher term from %s: %d. Becoming Follower.", rf.id, peerAddr, reply.Term)
				rf.becomeFollower(reply.Term)
				return
			}

			if reply.Success {
				rf.nextIndex[peerId] = args.PrevLogIndex + len(args.Entries) + 1
				rf.matchIndex[peerId] = args.PrevLogIndex + len(args.Entries)
				log.Printf("Node %d (Leader) AppendEntries to %s successful. nextIndex: %d, matchIndex: %d", rf.id, peerAddr, rf.nextIndex[peerId], rf.matchIndex[peerId])

				// Update commitIndex if a majority of followers have replicated the entry
				N := rf.commitIndex
				for N < len(rf.log) {
					count := 0
					for _, matchIdx := range rf.matchIndex {
						if matchIdx >= N && rf.log[N].Term == rf.currentTerm {
							count++
						}
					}
					if count > len(rf.peers)/2 {
						N++
					} else {
						break
					}
				}
				if N > rf.commitIndex {
					rf.commitIndex = N
					rf.applyCommittedEntries()
				}
			} else {
				// Decrement nextIndex and retry AppendEntries
				log.Printf("Node %d (Leader) AppendEntries to %s failed. Decrementing nextIndex from %d", rf.id, peerAddr, rf.nextIndex[peerId])
				if reply.XTerm != -1 {
					// Conflict by term
					found := false
					for i := len(rf.log) - 1; i >= 1; i-- {
						if rf.log[i].Term == reply.XTerm {
							rf.nextIndex[peerId] = i + 1
							found = true
							break
						}
					}
					if !found {
						rf.nextIndex[peerId] = reply.XIndex // Or reply.XLen
					}
				} else if reply.XIndex != -1 {
					// Conflict by index
					rf.nextIndex[peerId] = reply.XIndex
				} else {
					rf.nextIndex[peerId] = reply.XLen // Conflict by log length
				}
				if rf.nextIndex[peerId] < 1 {
					rf.nextIndex[peerId] = 1
				}
			}
		}(peer, i)
	}
}

// call RPC method on a peer.
func (rf *Raft) call(addr string, serviceMethod string, args interface{}, reply interface{}) error {
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer client.Close()

	err = client.Call(serviceMethod, args, reply)
	return err
}

// runRaft runs the main loop of a Raft node.
func (rf *Raft) runRaft() {
	// Register RPC methods
	err := rpc.Register(rf)
	if err != nil {
		log.Fatalf("Failed to register RPC: %v", err)
	}

	// Start RPC server
	l, err := net.Listen("tcp", rf.peers[rf.id])
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", rf.peers[rf.id], err)
	}
	go rpc.Accept(l)

	// Main Raft loop
	for {
		rf.mu.Lock()
		state := rf.state
		lastHeartbeat := rf.lastHeartbeat
		electionTimeout := rf.electionTimeout
		rf.mu.Unlock()

		switch state {
		case Follower:
			if time.Since(lastHeartbeat) > electionTimeout {
				log.Printf("Node %d (Follower) election timeout. Becoming Candidate.", rf.id)
				rf.mu.Lock()
				rf.becomeCandidate()
				rf.mu.Unlock()
				go rf.startElection()
			}
		case Candidate:
			if time.Since(lastHeartbeat) > electionTimeout {
				log.Printf("Node %d (Candidate) election timeout. Starting new election.", rf.id)
				rf.mu.Lock()
				rf.becomeCandidate()
				rf.mu.Unlock()
				go rf.startElection()
			}
		case Leader:
			// Leader loop is handled by startHeartbeatTimer goroutine
		}
		time.Sleep(10 * time.Millisecond) // Small delay to prevent busy-waiting
	}
}

// startElection initiates an election.
func (rf *Raft) startElection() {
	rf.mu.Lock()
	currentTerm := rf.currentTerm
	candidateId := rf.id
	lastLogIndex := len(rf.log) - 1
	lastLogTerm := rf.log[lastLogIndex].Term
	rf.mu.Unlock()

	votesReceived := 1 // Vote for self
	var votesMu sync.Mutex

	args := &RequestVoteArgs{
		Term:        currentTerm,
		CandidateId: candidateId,
		LastLogIndex: lastLogIndex,
		LastLogTerm:  lastLogTerm,
	}

	for i, peer := range rf.peers {
		if i == rf.id {
			continue // Don't send to self
		}
		go func(peerAddr string) {
			reply := &RequestVoteReply{}
			log.Printf("Node %d (Candidate) sending RequestVote to %s (Term %d)", rf.id, peerAddr, args.Term)
			err := rf.call(peerAddr, "Raft.RequestVote", args, reply)
			if err != nil {
				log.Printf("Node %d (Candidate) failed to send RequestVote to %s: %v", rf.id, peerAddr, err)
				return
			}

			rf.mu.Lock()
			defer rf.mu.Unlock()

			if rf.state != Candidate || rf.currentTerm != currentTerm {
				return // Already changed state or term
			}

			if reply.Term > rf.currentTerm {
				log.Printf("Node %d (Candidate) discovered higher term from %s: %d. Becoming Follower.", rf.id, peerAddr, reply.Term)
				rf.becomeFollower(reply.Term)
				return
			}

			if reply.VoteGranted {
				votesMu.Lock()
				votesReceived++
				votesMu.Unlock()
				log.Printf("Node %d (Candidate) received vote from %s. Total votes: %d", rf.id, peerAddr, votesReceived)
				if votesReceived > len(rf.peers)/2 && rf.state == Candidate {
					log.Printf("Node %d (Candidate) received majority votes. Becoming Leader.", rf.id)
					rf.becomeLeader()
				}
			}
		}(peer)
	}
}

// Propose a command to the Raft cluster.
func (rf *Raft) Propose(commandType string, commandData []byte) error {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if rf.state != Leader {
		return fmt.Errorf("not leader")
	}

	newEntry := LogEntry{
		Term:    rf.currentTerm,
		CommandType: commandType,
		CommandData: commandData,
	}
	rf.log = append(rf.log, newEntry)
	log.Printf("Node %d (Leader) proposed new entry: %+v", rf.id, newEntry)

	// Immediately try to send AppendEntries to all followers
	rf.sendHeartbeats()

	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <node_id>", os.Args[0])
	}

	nodeID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Invalid node ID: %v", err)
	}

	// Example usage: 3 nodes
	peers := []string{"127.0.0.1:8001", "127.0.0.1:8002", "127.0.0.1:8003"}

	if nodeID < 0 || nodeID >= len(peers) {
		log.Fatalf("Node ID %d is out of bounds for %d peers", nodeID, len(peers))
	}

	// Create a base directory for this node's files
	nodeBaseDir := fmt.Sprintf("node_data_%d", nodeID)
	if err := os.MkdirAll(nodeBaseDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory %s: %v", nodeBaseDir, err)
	}

	rf := NewRaft(nodeID, peers, nodeBaseDir)
	rf.runRaft()

	select {} // Block forever to keep the goroutine alive
}