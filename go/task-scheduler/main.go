package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Task represents a scheduled command
type Task struct {
	ID        int
	Command   string
	ExecuteAt time.Time
	Timer     *time.Timer // To allow cancellation
}

var (
	tasks     = make(map[int]*Task)
	nextTaskID = 1
	mu        sync.Mutex // Mutex to protect tasks map and nextTaskID
)

func scheduleTask(delaySeconds int, command string) {
	mu.Lock()
	defer mu.Unlock()

	taskID := nextTaskID
	nextTaskID++

	executeAt := time.Now().Add(time.Duration(delaySeconds) * time.Second)

	task := &Task{
		ID:        taskID,
		Command:   command,
		ExecuteAt: executeAt,
	}

	task.Timer = time.AfterFunc(time.Duration(delaySeconds)*time.Second, func() {
		log.Printf("Executing task %d: %s\n", task.ID, task.Command)
		cmd := exec.Command("sh", "-c", task.Command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Printf("Task %d execution failed: %v\n", task.ID, err)
		}

		mu.Lock()
		delete(tasks, task.ID) // Remove task after execution
		mu.Unlock()
	})

	tasks[taskID] = task
	fmt.Printf("Task %d scheduled to run in %d seconds (at %s): %s\n", taskID, delaySeconds, executeAt.Format("15:04:05"), command)
}

func listTasks() {
	mu.Lock()
	defer mu.Unlock()

	if len(tasks) == 0 {
		fmt.Println("No tasks scheduled.")
		return
	}

	fmt.Println("Scheduled Tasks:")
	for _, task := range tasks {
		fmt.Printf("  ID: %d, Command: %s, Executes At: %s (in %s)\n",
			task.ID, task.Command, task.ExecuteAt.Format("15:04:05"), time.Until(task.ExecuteAt).Round(time.Second))
	}
}

func cancelTask(taskID int) {
	mu.Lock()
	defer mu.Unlock()

	task, ok := tasks[taskID]
	if !ok {
		fmt.Printf("Task %d not found.\n", taskID)
		return
	}

	task.Timer.Stop() // Stop the timer
	delete(tasks, taskID)
	fmt.Printf("Task %d cancelled.\n", taskID)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s schedule <delay_in_seconds> <command>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s list\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s cancel <task_id>\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	command := args[0]

	switch command {
	case "schedule":
		if len(args) < 3 {
			fmt.Println("Usage: schedule <delay_in_seconds> <command>")
			os.Exit(1)
		}
		delay, err := strconv.Atoi(args[1])
		if err != nil || delay < 0 {
			fmt.Println("Invalid delay. Must be a non-negative integer.")
			os.Exit(1)
		}
		cmdToRun := strings.Join(args[2:], " ")
		scheduleTask(delay, cmdToRun)
	case "list":
		listTasks()
	case "cancel":
		if len(args) < 2 {
			fmt.Println("Usage: cancel <task_id>")
			os.Exit(1)
		}
		taskID, err := strconv.Atoi(args[1])
		if err != nil || taskID <= 0 {
			fmt.Println("Invalid task ID. Must be a positive integer.")
			os.Exit(1)
		}
		cancelTask(taskID)
	default:
		flag.Usage()
		os.Exit(1)
	}

	// Keep the main goroutine alive for scheduled tasks to run
	select {}
}