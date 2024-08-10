package main

import (
	"fmt"
	"strconv"
	"strings"
	"syscall/js"
)

func main() {
	c := make(chan struct{}, 0)
	fmt.Println("WebAssembly module is ready")
	js.Global().Set("nextGeneration", js.FuncOf(nextGenerationJS))
	<-c
}

func nextGenerationJS(this js.Value, args []js.Value) interface{} {
	// Wrap the function in a closure to handle panics
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	fmt.Println("nextGenerationJS called")
	if len(args) != 1 {
		fmt.Println("Invalid number of arguments")
		return js.Null()
	}

	// Convert JS Map to Go map
	jsGrid := args[0]
	grid := make(map[string]bool)
	keys := js.Global().Get("Array").Call("from", jsGrid.Call("keys"))
	for i := 0; i < keys.Length(); i++ {
		key := keys.Index(i).String()
		value := jsGrid.Call("get", key).Bool()
		grid[key] = value
	}

	fmt.Printf("Input grid size: %d\n", len(grid))

	// Compute next generation
	newGrid := nextGeneration(grid)

	fmt.Printf("Output grid size: %d\n", len(newGrid))

	// Convert Go map back to JS Map
	result := js.Global().Get("Map").New()
	for key, value := range newGrid {
		result.Call("set", key, value)
	}

	return result
}

func nextGeneration(grid map[string]bool) map[string]bool {
	newGrid := make(map[string]bool)
	neighbors := make(map[string]int)

	// Count neighbors for all cells
	for key := range grid {
		x, y := parseKey(key)
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if dx == 0 && dy == 0 {
					continue
				}
				neighborKey := fmt.Sprintf("%d,%d", x+dx, y+dy)
				neighbors[neighborKey]++
			}
		}
	}

	fmt.Printf("Neighbors counted: %d\n", len(neighbors))

	// Apply rules
	for key, count := range neighbors {
		if grid[key] {
			if count == 2 || count == 3 {
				newGrid[key] = true
			}
		} else {
			if count == 3 {
				newGrid[key] = true
			}
		}
	}

	fmt.Printf("New grid cells: %d\n", len(newGrid))

	return newGrid
}

func parseKey(key string) (int, int) {
	parts := strings.Split(key, ",")
	x, _ := strconv.Atoi(parts[0])
	y, _ := strconv.Atoi(parts[1])
	return x, y
}
