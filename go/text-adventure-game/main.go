package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Room represents a location in the game world.
type Room struct {
	Name        string
	Description string
	Exits       map[string]*Room // e.g., "north": &room2
	Items       []*Item
}

// Item represents an object in the game.
type Item struct {
	Name        string
	Description string
}

// Player represents the player character.
type Player struct {
	CurrentRoom *Room
	Inventory   []*Item
}

// Game represents the overall game state.
type Game struct {
	Player *Player
	Rooms  map[string]*Room // All rooms in the game, keyed by name
}

func main() {
	fmt.Println("Welcome to the Go Text Adventure Game!")

	// --- Game Initialization (Placeholder) ---
	// Create rooms
	startRoom := &Room{
		Name:        "Starting Room",
		Description: "You are in a dimly lit room. There is a dusty old book on a table.",
		Exits:       make(map[string]*Room),
	}
	anotherRoom := &Room{
		Name:        "Another Room",
		Description: "This room is brighter, with a window overlooking a garden.",
		Exits:       make(map[string]*Room),
	}

	// Link rooms
	startRoom.Exits["north"] = anotherRoom
	anotherRoom.Exits["south"] = startRoom

	// Create items
	dustyBook := &Item{
		Name:        "dusty book",
		Description: "A very old book, its pages are brittle.",
	}
	startRoom.Items = append(startRoom.Items, dustyBook)

	// Initialize player
	player := &Player{
		CurrentRoom: startRoom,
	}

	// Initialize game
	game := &Game{
		Player: player,
		Rooms: map[string]*Room{
			"Starting Room": startRoom,
			"Another Room":  anotherRoom,
		},
	}
	// --- End Game Initialization ---

	reader := bufio.NewReader(os.Stdin)

	// Game Loop
	for {
		fmt.Println("\n----------------------------------------")
		fmt.Printf("Location: %s\n", game.Player.CurrentRoom.Name)
		fmt.Println(game.Player.CurrentRoom.Description)

		if len(game.Player.CurrentRoom.Items) > 0 {
			fmt.Print("You see: ")
			for i, item := range game.Player.CurrentRoom.Items {
				fmt.Print(item.Name)
				if i < len(game.Player.CurrentRoom.Items)-1 {
					fmt.Print(", ")
				}
			}
			fmt.Println(".")
		}

		fmt.Print("Exits: ")
		exits := []string{}
		for exit := range game.Player.CurrentRoom.Exits {
			exits = append(exits, exit)
		}
		fmt.Println(strings.Join(exits, ", ") + ".")

		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))

		// Basic Command Parsing
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		args := []string{}
		if len(parts) > 1 {
			args = parts[1:]
		}

		switch command {
		case "go":
			if len(args) > 0 {
				direction := args[0]
				if nextRoom, ok := game.Player.CurrentRoom.Exits[direction]; ok {
					game.Player.CurrentRoom = nextRoom
				} else {
					fmt.Println("You can't go that way.")
				}
			} else {
				fmt.Println("Go where?")
			}
		case "look":
			// Already displayed at the start of the loop, but can be expanded
			// to look at specific items or directions.
			fmt.Println(game.Player.CurrentRoom.Description)
			if len(game.Player.CurrentRoom.Items) > 0 {
				fmt.Print("You see: ")
				for i, item := range game.Player.CurrentRoom.Items {
					fmt.Print(item.Name)
					if i < len(game.Player.CurrentRoom.Items)-1 {
						fmt.Print(", ")
					}
				}
				fmt.Println(".")
			}
		case "quit", "exit":
			fmt.Println("Thanks for playing!")
			return
		default:
			fmt.Println("I don't understand that command.")
		}
	}
}