package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

func solutionPart1(scanner *bufio.Scanner) int {
	/*
		Essentially just find the sum of a list of integers
	*/

	var total int

	for scanner.Scan() {
		x := scanner.Text()
		xParsed, err := strconv.Atoi(x)
		if err != nil {
			log.Fatal(err)
		}
		total += xParsed
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return total
}

func getFreqs(scanner *bufio.Scanner) []int {
	var freqs []int
	for scanner.Scan() {
		x := scanner.Text()
		xParsed, err := strconv.Atoi(x)
		if err != nil {
			log.Fatal(err)
		}
		freqs = append(freqs, xParsed)
	}
	fmt.Println("len of freqs", len(freqs))
	return freqs
}

func findDuplicateFrequency(freqs []int) int {
	var currentFrequency = 0
	var duplicateFrequencyFound = false
	var numPasses = 0
	set := make(map[int]bool)
	set[0] = true
	for !duplicateFrequencyFound {
		numPasses++
		for _, el := range freqs {
			currentFrequency += el
			if set[currentFrequency] {
				duplicateFrequencyFound = true
				break
			}
			set[currentFrequency] = true
		}
	}
	fmt.Println("Number of Passes", numPasses)
	return currentFrequency
}

func solutionPart2(scanner *bufio.Scanner) int {
	/*
		read a list of integers
		loop through the integers keeping track of the cummulative sum
		once the same cummulative sum is observer twice report the value and exit program
		keep track of the cummulative sum of integers
	*/

	freqs := getFreqs(scanner)
	return findDuplicateFrequency(freqs)
	// return len(freqs)
}

func main() {

	// fnamePtr := flag.String("fname", "", "fname")
	numPtr := flag.Int("part", 1, "fname")
	var fname = "input.txt"
	if *numPtr == 2 {
		fname = "input.txt"
	}
	flag.Parse()

	fmt.Println("File Name", fname)
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if *numPtr == 1 {
		total := solutionPart1(scanner)
		fmt.Println("Solution: ", total)
	}

	if *numPtr == 2 {
		// result := findDuplicateFrequency([]int{-6, +3, +8, +5, -6})
		// result := findDuplicateFrequency([]int{+1, -1})

		fmt.Println("Solutions: ", solutionPart2(scanner))

	}

}
