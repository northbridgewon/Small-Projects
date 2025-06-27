package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FileMetadata represents information about a shared file.
type FileMetadata struct {
	Filename string
	Filesize int64
}

// NodeInfo represents information about a peer node, including its shared files.
type NodeInfo struct {
	Address     string
	SharedFiles []FileMetadata
}

// FileRequest represents a request for a specific file.
type FileRequest struct {
	Filename string
}

// FileChunk represents a chunk of a file being transferred.
type FileChunk struct {
	Filename string
	Offset   int64
	Data     []byte
	EOF      bool
}

// P2PService is the RPC service for peer-to-peer communication.
type P2PService struct {
	mu          sync.RWMutex
	sharedFiles map[string]string // map[filename]filepath
}

// NewP2PService creates a new P2PService instance.
func NewP2PService(sharedFiles map[string]string) *P2PService {
	return &P2PService{
		sharedFiles: sharedFiles,
	}
}

// Announce is an RPC method for a node to announce its presence and shared files.
func (s *P2PService) Announce(nodeInfo NodeInfo, reply *string) error {
	log.Printf("Received Announce from %s. Shared files: %v\n", nodeInfo.Address, nodeInfo.SharedFiles)
	// In a real system, you'd update a global peer list here.
	*reply = "ACK"
	return nil
}

// ListSharedFiles is an RPC method to list files shared by this node.
func (s *P2PService) ListSharedFiles(args string, reply *[]FileMetadata) error {
	log.Printf("Received request to list shared files from %s\n", args)
	s.mu.RLock()
	defer s.mu.RUnlock()

	var files []FileMetadata
	for filename, filePath := range s.sharedFiles {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Printf("Warning: Could not get file info for %s: %v\n", filePath, err)
			continue
		}
		files = append(files, FileMetadata{Filename: filename, Filesize: fileInfo.Size()})
	}
	*reply = files
	return nil
}

// RequestFile is an RPC method to request a file from a peer.
func (s *P2PService) RequestFile(req FileRequest, stream *FileChunk) error {
	log.Printf("Received request for file: %s\n", req.Filename)

	s.mu.RLock()
	filePath, ok := s.sharedFiles[req.Filename]
	s.mu.RUnlock()

	if !ok {
		return fmt.Errorf("file %q not found on this node", req.Filename)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", filePath, err)
	}
	defer file.Close()

	// For simplicity, we'll send the whole file in one go for now.
	// In a real system, you'd stream chunks.
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %w", filePath, err)
	}

	stream.Filename = req.Filename
	stream.Offset = 0
	stream.Data = data
	stream.EOF = true

	log.Printf("Sending file %s (%d bytes) to requester.\n", req.Filename, len(data))
	return nil
}

// Node represents a peer in the P2P network
type Node struct {
	Address     string
	SharedFiles map[string]string // map[filename]filepath
	mu          sync.RWMutex
}

// NewNode creates a new P2P node
func NewNode(address string) *Node {
	return &Node{
		Address:     address,
		SharedFiles: make(map[string]string),
	}
}

// IndexFiles scans a directory and adds files to the shared files map
func (n *Node) IndexFiles(dir string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %q: %w", dir, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(dir, file.Name())
			n.SharedFiles[file.Name()] = filePath
			log.Printf("Indexed file: %s (%s)\n", file.Name(), filePath)
		}
	}
	return nil
}

// listLocalFiles prints the files shared by the local node
func (n *Node) listLocalFiles() {
	n.mu.RLock()
	defer n.mu.RUnlock()

	if len(n.SharedFiles) == 0 {
		fmt.Println("No files are being shared by this node.")
		return
	}

	fmt.Println("Files shared by this node:")
	for filename, filepath := range n.SharedFiles {
		fmt.Printf("  - %s (Path: %s)\n", filename, filepath)
	}
}

// announceToPeer sends an Announce RPC call to a peer
func (n *Node) announceToPeer(peerAddress string) error {
	client, err := rpc.Dial("tcp", peerAddress)
	if err != nil {
		return fmt.Errorf("failed to dial peer %q: %w", peerAddress, err)
	}
	defer client.Close()

	var fileMetadata []FileMetadata
	n.mu.RLock()
	for filename, filePath := range n.SharedFiles {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Printf("Warning: Could not get file info for %s: %v\n", filePath, err)
			continue
		}
		fileMetadata = append(fileMetadata, FileMetadata{Filename: filename, Filesize: fileInfo.Size()})
	}
	n.mu.RUnlock()

	nodeInfo := NodeInfo{
		Address:     n.Address,
		SharedFiles: fileMetadata,
	}

	var reply string
	err = client.Call("P2PService.Announce", nodeInfo, &reply)
	if err != nil {
		return fmt.Errorf("failed to announce to peer %q: %w", peerAddress, err)
	}
	log.Printf("Announced to %s: %s\n", peerAddress, reply)
	return nil
}

// requestFileFromPeer requests a file from a peer and saves it locally
func (n *Node) requestFileFromPeer(peerAddress, filename, saveDir string) error {
	client, err := rpc.Dial("tcp", peerAddress)
	if err != nil {
		return fmt.Errorf("failed to dial peer %q: %w", peerAddress, err)
	}
	defer client.Close()

	req := FileRequest{Filename: filename}
	var chunk FileChunk

	log.Printf("Requesting file %q from %q\n", filename, peerAddress)
	err = client.Call("P2PService.RequestFile", req, &chunk)
	if err != nil {
		return fmt.Errorf("failed to request file %q from %q: %w", filename, peerAddress, err)
	}

	if !chunk.EOF {
		return fmt.Errorf("expected full file, got partial chunk for %q", filename)
	}

	savePath := filepath.Join(saveDir, filename)
	err = ioutil.WriteFile(savePath, chunk.Data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save file %q: %w", savePath, err)
	}

	log.Printf("Successfully downloaded and saved %q to %q\n", filename, savePath)
	return nil
}

// discoverPeers scans a range of ports on a host and lists shared files from discovered peers.
func (n *Node) discoverPeers(scanHost string, startPort, endPort int) {
	log.Printf("Discovering peers on %s from port %d to %d...\n", scanHost, startPort, endPort)
	var discoveredPeers []string

	for p := startPort; p <= endPort; p++ {
		peerAddr := fmt.Sprintf("%s:%d", scanHost, p)
		conn, err := net.DialTimeout("tcp", peerAddr, 500*time.Millisecond) // Short timeout
		if err == nil {
			conn.Close()
			discoveredPeers = append(discoveredPeers, peerAddr)
		}
	}

	if len(discoveredPeers) == 0 {
		fmt.Printf("No peers found on %s in port range %d-%d.\n", scanHost, startPort, endPort)
		return
	}

	fmt.Println("\nDiscovered Peers and their Shared Files:")
	for _, peerAddr := range discoveredPeers {
		fmt.Printf("  Peer: %s\n", peerAddr)
		client, err := rpc.Dial("tcp", peerAddr)
		if err != nil {
			log.Printf("    Error connecting to peer %s: %v\n", peerAddr, err)
			continue
		}
		defer client.Close()

		var sharedFiles []FileMetadata
		err = client.Call("P2PService.ListSharedFiles", n.Address, &sharedFiles)
		if err != nil {
			log.Printf("    Error listing files from %s: %v\n", peerAddr, err)
			continue
		}

		if len(sharedFiles) == 0 {
			fmt.Println("      No files shared.")
		} else {
			for _, file := range sharedFiles {
				fmt.Printf("      - %s (%d bytes)\n", file.Filename, file.Filesize)
			}
		}
	}
}

// autoDownload scans for peers, lists their files, and downloads them.
func (n *Node) autoDownload(scanHost string, startPort, endPort int, saveDir string) {
	log.Printf("Initiating auto-download from peers on %s from port %d to %d...\n", scanHost, startPort, endPort)
	var discoveredPeers []string

	for p := startPort; p <= endPort; p++ {
		peerAddr := fmt.Sprintf("%s:%d", scanHost, p)
		conn, err := net.DialTimeout("tcp", peerAddr, 500*time.Millisecond)
		if err == nil {
			conn.Close()
			discoveredPeers = append(discoveredPeers, peerAddr)
		}
	}

	if len(discoveredPeers) == 0 {
		fmt.Printf("No peers found on %s in port range %d-%d for auto-download.\n", scanHost, startPort, endPort)
		return
	}

	fmt.Println("\nDiscovered Peers for Auto-Download:")
	for _, peerAddr := range discoveredPeers {
		fmt.Printf("  Peer: %s\n", peerAddr)
		client, err := rpc.Dial("tcp", peerAddr)
		if err != nil {
			log.Printf("    Error connecting to peer %s: %v\n", peerAddr, err)
			continue
		}
		defer client.Close()

		var sharedFiles []FileMetadata
		err = client.Call("P2PService.ListSharedFiles", n.Address, &sharedFiles)
		if err != nil {
			log.Printf("    Error listing files from %s: %v\n", peerAddr, err)
			continue
		}

		if len(sharedFiles) == 0 {
			fmt.Println("      No files shared.")
		} else {
			for _, file := range sharedFiles {
				fmt.Printf("      Attempting to download: %s\n", file.Filename)
				err := n.requestFileFromPeer(peerAddr, file.Filename, saveDir)
				if err != nil {
					log.Printf("        Error downloading %s from %s: %v\n", file.Filename, peerAddr, err)
				}
			}
		}
	}
}

func main() {
	port := flag.Int("port", 8080, "Port for the node to listen on")
	shareDir := flag.String("share-dir", ".", "Directory to share files from")
	command := flag.String("cmd", "start", "Command to execute (start, list-files, request-file, discover, auto-download)")
	peer := flag.String("peer", "", "Address of a peer to connect to (e.g., localhost:8081)")
	requestFile := flag.String("file", "", "Filename to request from a peer")
	saveDir := flag.String("save-dir", ".", "Directory to save requested files")

	scanHost := flag.String("scan-host", "localhost", "Host to scan for peers")
	scanStartPort := flag.Int("scan-start-port", 8080, "Starting port for peer scan")
	scanEndPort := flag.Int("scan-end-port", 8090, "Ending port for peer scan")

	flag.Parse()

	address := fmt.Sprintf("localhost:%d", *port)
	node := NewNode(address)

	// Index files from the shared directory
	err := node.IndexFiles(*shareDir)
	if err != nil {
		log.Fatalf("Error indexing files: %v", err)
	}

	switch *command {
	case "start":
		log.Printf("Node starting on %s, sharing files from %s\n", node.Address, *shareDir)

		// Register RPC service
		rpcService := NewP2PService(node.SharedFiles)
		rpc.Register(rpcService)

		// Start RPC listener
		listener, err := net.Listen("tcp", node.Address)
		if err != nil {
			log.Fatalf("Error listening on %s: %v\n", node.Address, err)
		}
		defer listener.Close()

		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					log.Printf("Error accepting connection: %v", err)
					continue
				}
				rpc.ServeConn(conn)
			}
		}()

		// Announce to a peer if specified
		if *peer != "" {
			log.Printf("Announcing to peer: %s\n", *peer)
			err := node.announceToPeer(*peer)
			if err != nil {
				log.Printf("Error announcing to peer %q: %v\n", *peer, err)
			}
		}

		select {} // Keep the main goroutine alive
	case "list-files":
		node.listLocalFiles()
	case "request-file":
		if *peer == "" || *requestFile == "" {
			fmt.Println("Usage: -cmd request-file -peer <peer_address> -file <filename> [-save-dir <directory>]")
			os.Exit(1)
		}
		err := node.requestFileFromPeer(*peer, *requestFile, *saveDir)
		if err != nil {
			log.Fatalf("Error requesting file: %v", err)
		}
	case "discover":
		node.discoverPeers(*scanHost, *scanStartPort, *scanEndPort)
	case "auto-download":
		node.autoDownload(*scanHost, *scanStartPort, *scanEndPort, *saveDir)
	default:
		fmt.Printf("Unknown command: %s\n", *command)
		flag.Usage()
		os.Exit(1)
	}
}
