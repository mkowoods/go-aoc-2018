package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
)

type Panel struct {
	id     int
	X      int
	Y      int
	height int
	width  int
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func parseInput(inp string) Panel {
	r, err := regexp.Compile(`#(\d+) @ (\d+),(\d+): (\d+)x(\d+)`)
	if err != nil {
		log.Fatal(err)
	}
	out := r.FindAllStringSubmatch(inp, -1)
	id, _ := strconv.Atoi(out[0][1])
	X, _ := strconv.Atoi(out[0][2])
	Y, _ := strconv.Atoi(out[0][3])
	width, _ := strconv.Atoi(out[0][4])
	height, _ := strconv.Atoi(out[0][5])

	panel := Panel{id, X, Y, height, width}
	return panel
}

func solutionPart1(scanner *bufio.Scanner) int {

	var grid = map[int]map[int]int{}
	var overlapBlock = 0

	for scanner.Scan() {
		inp := scanner.Text()
		panel := parseInput(inp)
		for x := panel.X; x < (panel.X + panel.width); x++ {
			for y := panel.Y; y < (panel.Y + panel.height); y++ {
				if grid[x] == nil {
					grid[x] = make(map[int]int)
				}
				grid[x][y]++
				// only count it as an overlap once e.g. when the second grid is identified
				// it's counted
				if grid[x][y] == 2 {
					overlapBlock++
				}
			}
		}
	}
	return overlapBlock
}

func solutionPart2(scanner *bufio.Scanner) int {
	var hasOverlap = make(map[int]bool)
	var grid = map[int]map[int]int{}

	for scanner.Scan() {
		inp := scanner.Text()
		panel := parseInput(inp)
		hasOverlap[panel.id] = false
		for x := panel.X; x < (panel.X + panel.width); x++ {
			for y := panel.Y; y < (panel.Y + panel.height); y++ {

				if grid[x] == nil {
					grid[x] = make(map[int]int)
				}

				//if the grid block has not been claimed then claim it to the first person
				if grid[x][y] == 0 {
					grid[x][y] = panel.id
				}

				//if a user wants to claim the grid that has already been claimed then flag that user and the user that claimed the block
				//as both having an overlapping claim
				if grid[x][y] != panel.id {
					// fmt.Println(x, y, grid[x][y], panel.id)
					hasOverlap[panel.id] = true
					hasOverlap[grid[x][y]] = true
				}
			}
		}
	}

	PrintMemUsage()

	for k, v := range hasOverlap {
		if !v {
			fmt.Println("No Overlap", k, v)
			return k
		}
	}
	// fmt.Println(hasOverlap)
	return -1
}

func getScanner(fname string) (*os.File, *bufio.Scanner) {
	fmt.Println("File Name", fname)
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
		file.Close()
	}

	scanner := bufio.NewScanner(file)
	return file, scanner
}

func main() {
	numPtr := flag.Int("part", 1, "which solution to use")
	isTest := flag.Bool("test", false, "whether to use test file")

	flag.Parse()

	var fname = "input.txt"
	if *isTest {
		fname = "test-input.txt"
	}

	switch ptr := *numPtr; ptr {
	case 1:
		file, scanner := getScanner(fname)
		defer file.Close()
		fmt.Println("Solution Part 1:", solutionPart1(scanner))
	case 2:
		file, scanner := getScanner(fname)
		defer file.Close()
		fmt.Println("Solution Part 2:", solutionPart2(scanner))

	}

}
