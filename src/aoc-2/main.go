package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
)

func getCharFrequency(id string) map[string]int {
	ctr := make(map[string]int)
	for _, char := range id {
		ctr[string(char)]++
	}
	return ctr
}

func counterHasValue(ctr map[string]int, value int) bool {
	for _, v := range ctr {
		if v == value {
			return true
		}
	}
	return false
}

func solutionPart1(scanner *bufio.Scanner) int {
	var ct2 = 0
	var ct3 = 0

	for scanner.Scan() {
		id := scanner.Text()
		ctr := getCharFrequency(id)
		var has2 = counterHasValue(ctr, 2)
		var has3 = counterHasValue(ctr, 3)
		if has2 {
			ct2++
		}

		if has3 {
			ct3++
		}
	}

	return ct2 * ct3
}

func hammingDistance(str1 string, str2 string) int {
	var dist = 0
	for i := 0; i < len(str1); i++ {
		if str1[i] != str2[i] {
			dist++
		}
	}
	return dist
}

func removeDifferentCharacters(str1 string, str2 string) string {
	var buffer bytes.Buffer
	for i := 0; i < len(str1); i++ {
		if str1[i] == str2[i] {
			buffer.WriteByte(str1[i])
		}
	}
	return buffer.String()
}

func solutionPart2(scanner *bufio.Scanner) int {
	/*
		straight forward implementation of querying 2 most similar strings from a list
		just brute forced it, but a better approach would be to use a trie or substring matching to reduce candidates
	*/
	var ids []string
	for scanner.Scan() {
		id := scanner.Text()
		// fmt.Println(id)
		ids = append(ids, id)
	}
	// fmt.Println(ids)
	for i := 0; i < len(ids); i++ {
		for j := (i + 1); j < len(ids); j++ {
			id1 := ids[i]
			id2 := ids[j]
			if hammingDistance(id1, id2) == 1 {
				fmt.Println(id1, id2, hammingDistance(id1, id2))
				fmt.Println(removeDifferentCharacters(id1, id2))
			}
		}
	}
	return 0
}

// func getScanner(fname string) *bufio.Scanner {
// 	return scanner
// }

func main() {
	numPtr := flag.Int("part", 1, "fname")
	flag.Parse()

	var fname = "input.txt"
	if *numPtr == 0 {
		fname = "test-input.txt"
	}
	if *numPtr == 2 {
		fname = "input.txt"
	}

	fmt.Println("File Name", fname)
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// scanner := getScanner(fname)
	// fmt.Println("Solution Part 1:", solutionPart1(scanner))

	// scanner = getScanner(fname)
	fmt.Println("Solution Part 2:", solutionPart2(scanner))

}
