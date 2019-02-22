package main

import (
	"flag"
	"fmt"
	"log"
	"math"
)

// import "../util`"

// type Cell

func PowerLevel(x int, y int, serialNumber int) int {
	rackID := x + 10
	power := rackID * y
	power = power + serialNumber
	power = power * rackID
	power = (power / 100) % 10
	power = power - 5
	return power
}

func calcPowerLevel(serialNumber int) func(int, int) int {
	return func(x int, y int) int {
		return PowerLevel(x, y, serialNumber)
	}
}

func calcGridValue(grid [300][300]int, x int, y int) int {
	return grid[x+0][y] + grid[x+0][y+1] + grid[x+0][y+2] +
		grid[x+1][y] + grid[x+1][y+1] + grid[x+1][y+2] +
		grid[x+2][y] + grid[x+2][y+1] + grid[x+2][y+2]
}

func solution1(serialNumber int) int {
	//coordinates are grid[x][y]
	var grid = [300][300]int{}

	powerLevel := calcPowerLevel(serialNumber)

	//fill in the furthest right columns and bottom 2 rows
	for i := 0; i < 300; i++ {
		grid[i][298] = powerLevel(i, 298)
		grid[i][299] = powerLevel(i, 299)
		grid[298][i] = powerLevel(298, i)
		grid[299][i] = powerLevel(299, i)
	}

	bestValue := math.MinInt64
	var bestX, bestY int

	// starting in the bottom right corner move right -> left and bottom -> top
	for x := 297; x > -1; x-- {
		for y := 297; y > -1; y-- {
			grid[x][y] = powerLevel(x, y)
			val := calcGridValue(grid, x, y)
			if val > bestValue {
				bestValue = val
				bestX = x
				bestY = y
			}

		}
	}

	fmt.Printf("Best Coordinates %d, %d ", bestX, bestY)
	return bestValue
}

func solution2(serialNumber int) int {

	//Start by pop
	var grid = [300][300]int{}
	powerLevel := calcPowerLevel(serialNumber)
	for x := 0; x < 300; x++ {
		for y := 0; y < 300; y++ {
			grid[x][y] = powerLevel(x, y)
		}
	}

}

func main() {
	input := flag.Int("input", -1, "input to program")
	flag.Parse()

	if *input == -1 {
		log.Fatal("Input was not set")
	}

	// file, scanner := util.GetScanner(fname)

	// file, scanner := getScanner(fname)
	// defer file.Close()
	// initDataStructures(scanner)

	fmt.Println("Solution 1:", solution1(*input))
	fmt.Println("Solution 2", solution2(*input))

}
