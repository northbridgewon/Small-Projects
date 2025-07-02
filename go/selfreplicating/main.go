package selfreplicating

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// main is the entry point of the program. It orchestrates the replication lifecycle.
func main() {
	log.Println("Starting replication cycle...")

	// Stage 1: Locate and Read Self
	originalPath, binaryContent, err := readSelf()
	if err != nil {
		log.Fatalf("Fatal: Could not read own binary: %v", err)
	}
	log.Printf("Successfully read own binary from: %s", originalPath)

	// Stage 2: Replicate Self
	newPath, err := replicateSelf(binaryContent)
	if err != nil {
		log.Fatalf("Fatal: Could not replicate self: %v", err)
	}
	log.Printf("Successfully replicated binary to: %s", newPath)

	// Stage 3: Execute New Generation
	err = launchDetached(newPath)
	if err != nil {
		log.Fatalf("Fatal: Could not launch new generation: %v", err)
	}
	log.Printf("Successfully launched new generation process from: %s", newPath)

	// Stage 4: Terminate and Cleanup
	err = selfDestruct(originalPath)
	if err != nil {
		// This is a non-fatal error, as the main goal was achieved.
		// The OS will clean up the temp dir, but the original binary remains.
		log.Printf("Warning: Failed to self-destruct original binary: %v", err)
	}
	log.Printf("Self-destruct initiated for: %s. Parent process exiting.", originalPath)
}

// readSelf locates the current running executable, resolves any symlinks,
// and reads its entire content into a byte slice.
func readSelf() (string, []byte, error) {
	// os.Executable() is the most reliable way to find the path to the current binary.
	rawPath, err := os.Executable()
	if err != nil {
		return "", nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	// Resolve any symlinks to get the canonical path. This is crucial because
	// we need to read the actual file, not the link.
	canonicalPath, err := filepath.EvalSymlinks(rawPath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to evaluate symlinks for %s: %w", rawPath, err)
	}

	// os.ReadFile is the simplest way to read an entire file into memory.
	binaryContent, err := os.ReadFile(canonicalPath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read own binary content from %s: %w", canonicalPath, err)
	}

	return canonicalPath, binaryContent, nil
}

// replicateSelf creates a new unique temporary directory and writes the binary content
// to a new file within it, ensuring it is executable.
func replicateSelf(content []byte) (string, error) {
	// os.MkdirTemp creates a unique directory in the OS's default temp location (e.g., /tmp).
	// This is secure and avoids race conditions.
	tempDir, err := os.MkdirTemp("", "replicant-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Defer cleanup in case of subsequent errors. If this function returns an error,
	// the created directory will be removed.
	defer func() {
		if err != nil {
			os.RemoveAll(tempDir)
		}
	}()

	// Define the path for the new executable.
	newBinaryPath := filepath.Join(tempDir, "replica")

	// Use the atomic write pattern: write to a temp file first, then rename.
	// This prevents a corrupted binary if the write is interrupted.
	tempFilePath := newBinaryPath + ".tmp"

	// Write the binary content with executable permissions.
	// 0755 is rwxr-xr-x: owner can read/write/execute, group/others can read/execute.
	// The final permissions will be affected by the system's umask.
	err = os.WriteFile(tempFilePath, content, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to write temporary binary to %s: %w", tempFilePath, err)
	}

	// Atomically rename the temp file to its final name. This is an atomic operation
	// on Linux if the source and destination are on the same filesystem.
	err = os.Rename(tempFilePath, newBinaryPath)
	if err != nil {
		return "", fmt.Errorf("failed to atomically rename binary to %s: %w", newBinaryPath, err)
	}

	return newBinaryPath, nil
}

// launchDetached starts the new executable as a fully detached process.
func launchDetached(path string) error {
	// Create the command to execute the new binary.
	cmd := exec.Command(path)

	// For true detachment on Linux, we set the process group ID.
	// This creates a new session, making the child process independent of the parent's
	// terminal and process group, effectively daemonizing it.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// cmd.Start() is non-blocking. It launches the process and returns immediately.
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start new process: %w", err)
	}

	// cmd.Process.Release() is crucial. It detaches the parent from the child,
	// allowing the parent to exit without killing the child. The child process
	// will be reparented by 'init' (PID 1).
	err = cmd.Process.Release()
	if err != nil {
		// If release fails, we should try to kill the process we just started.
		cmd.Process.Kill()
		return fmt.Errorf("failed to release new process: %w", err)
	}

	return nil
}

// selfDestruct removes the original executable file from the filesystem.
func selfDestruct(path string) error {
	// On Linux, a running executable can be unlinked (removed). The OS keeps the
	// file data in memory for the running process and cleans it up on exit.
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("failed to remove original binary %s: %w", path, err)
	}
	return nil
}
