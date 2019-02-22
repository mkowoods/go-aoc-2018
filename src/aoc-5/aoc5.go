package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

func shortestPolymer(polymer string) int {
	const alpha = "abcdefghijklmnopqrstuvwxyz"
	bestLength := math.MaxInt64
	// polymerLength := len(polymer)
	for _, char := range alpha {
		var charCt int
		lower := string(char)
		upper := strings.ToUpper(lower)
		var buffer bytes.Buffer
		for _, char := range polymer {
			c1 := string(char)
			if c1 == lower || c1 == upper {
				charCt++
				continue
			} else {
				buffer.WriteString(c1)
			}

		}
		polymerCopy := buffer.String()
		processedPoly := processPolymer(polymerCopy)
		processedPolyLength := len(processedPoly)
		if processedPolyLength <= bestLength {
			bestLength = processedPolyLength
		}
		fmt.Println(string(char), processedPolyLength, charCt, bestLength)

	}
	return bestLength
}

func processPolymer(polymer string) string {

	hasChange := true
	polymerLength := len(polymer)
	for hasChange {
		hasChange = false
		var buffer bytes.Buffer
		var bufferLength int
		for i := 0; i < (polymerLength); i++ {
			//if only 1 character remains then write it out and break the loop
			if i == (polymerLength - 1) {
				// fmt.Println("writing to string", string(priorPolymer[i]))
				buffer.WriteString(string(polymer[i]))
				bufferLength++
				break
			}
			c1 := string(polymer[i])
			c2 := string(polymer[i+1])
			//if you encounter the first in a reacting pair move the cursor past the pair
			//otherwise append the character to the new polymer string
			// fmt.Println(i, c1, c2)
			if c1 != c2 && strings.ToLower(c1) == strings.ToLower(c2) {
				i++
				hasChange = true
			} else {
				bufferLength++
				buffer.WriteString(c1)
			}
		}
		polymer = buffer.String()
		polymerLength = bufferLength
	}
	return polymer
}

func solutionPart1(scanner *bufio.Scanner) int {
	var processedPolymer string
	for scanner.Scan() {
		inp := scanner.Text()
		processedPolymer = processPolymer(inp)
		// fmt.Println("input", inp)
		fmt.Println("Solution 1: Length of Processed Polymber ", len(processedPolymer))
		fmt.Println("Solution 2: Shortest Polymner", shortestPolymer(inp))
	}
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
	// file sorted with the below command:
	// sort -h input.txt > tmp && mv tmp input.txt
	// numPtr := flag.Int("part", 1, "which solution to use")
	isTest := flag.Bool("test", false, "whether to use test file")

	flag.Parse()

	var fname = "input.txt"
	if *isTest {
		fname = "test-input.txt"
	}

	file, scanner := getScanner(fname)
	defer file.Close()
	fmt.Println("Solution Part 1:", solutionPart1(scanner))

}
