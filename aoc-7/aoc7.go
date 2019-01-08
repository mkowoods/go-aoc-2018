package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func solution1() {
	// create a directed graph and then the reverse graph
	// determine the parent vertex (with no children)
	// starting with that vertex add it's children to the frontier
	// explore the frontier in alphabetical order (using a heap??)
	// and only explore a node  if all of it's children have also been explored.
}

func solution(scanner *bufio.Scanner) int {
	for scanner.Scan() {
		inp := scanner.Text()
		fmt.Println("input", inp)
	}

	fmt.Println("Solution 1", 1)
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
	fmt.Println("Solution Part 1:", solution(scanner))

}
