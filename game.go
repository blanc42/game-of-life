package main

// Cell represents the state of a cell in the grid
type Cell int

const (
	Dead Cell = iota
	Alive
)

// Grid represents the game board
type Grid [][]Cell

// NextGeneration computes the next generation of the grid
func NextGeneration(grid Grid) Grid {
	newGrid := make(Grid, len(grid))
	for i := range grid {
		newGrid[i] = make([]Cell, len(grid[i]))
		for j := range grid[i] {
			newGrid[i][j] = updateCell(grid, i, j)
		}
	}
	return newGrid
}

func updateCell(grid Grid, i, j int) Cell {
	neighbors := countNeighbors(grid, i, j)
	switch grid[i][j] {
	case Alive:
		if neighbors < 2 || neighbors > 3 {
			return Dead
		}
	case Dead:
		if neighbors == 3 {
			return Alive
		}
	}
	return grid[i][j]
}

func countNeighbors(grid Grid, i, j int) int {
	count := 0
	for di := -1; di <= 1; di++ {
		for dj := -1; dj <= 1; dj++ {
			if di == 0 && dj == 0 {
				continue
			}
			ni, nj := i+di, j+dj
			if ni >= 0 && ni < len(grid) && nj >= 0 && nj < len(grid[ni]) {
				if grid[ni][nj] == Alive {
					count++
				}
			}
		}
	}
	return count
}
