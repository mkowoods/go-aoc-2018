package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// 2 3 0 3 10 11 12 1 1 0 1 99 2 1 1 2
// A----------------------------------
//     B----------- C-----------
//                      D-----

// 2 3 1 1 0 1 99 2 0 3 10 11 12 1 1 2
// A----------------------------------
//     B----------- D-----------
//         C-----

var INPUT_ARRAY []int
var NODE_INDEX = 0
var HEAD *Node

func sum(arr []int) int {
	total := 0
	for _, val := range arr {
		total += val
	}
	return total
}

type Node struct {
	idx      int
	metaData []int
	children []*Node
}

func dfsSum(node *Node) int {
	metaDataTotal := sum(node.metaData)
	for _, childNode := range node.children {
		metaDataTotal += dfsSum(childNode)
	}
	return metaDataTotal
}

func getNodeIndex() int {
	nodeIndex := NODE_INDEX
	NODE_INDEX++
	return nodeIndex
}

func getNodeValue(node *Node) int {
	//base case / termintation case
	if len(node.children) == 0 {
		return sum(node.metaData)
	}

	childSum := 0
	for _, md := range node.metaData {
		childIdx := (md - 1)
		if (childIdx >= 0) && childIdx < len(node.children) {
			childSum += getNodeValue(node.children[childIdx])
		}
	}
	return childSum
}

func parseTree(idx int, parent *Node) int {
	nodeIndex := getNodeIndex()
	slice := INPUT_ARRAY[idx:len(INPUT_ARRAY)]
	numNodes := slice[0]
	lenMetaData := slice[1]
	node := &Node{idx: nodeIndex}
	idx += 2
	for i := 0; i < numNodes; i++ {
		idx = parseTree(idx, node)
	}

	node.metaData = INPUT_ARRAY[idx:(idx + lenMetaData)]

	if parent != nil {
		parent.children = append(parent.children, node)
	} else {
		HEAD = node
	}
	// fmt.Println(nodeIndex, INPUT_ARRAY[idx:(idx+lenMetaData)])
	idx += lenMetaData
	return idx
}

func solution1() int {
	return dfsSum(HEAD)
}

func solution2() int {
	return getNodeValue(HEAD)
}
func initDataStructures(scanner *bufio.Scanner) {
	for scanner.Scan() {
		inp := scanner.Text()
		arr := strings.Split(inp, " ")
		INPUT_ARRAY = make([]int, len(arr))
		for idx, val := range arr {
			INPUT_ARRAY[idx], _ = strconv.Atoi(val)
		}
	}
	parseTree(0, nil)
	fmt.Println(HEAD)

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
	isTest := flag.Bool("test", false, "whether to use test file")

	flag.Parse()

	var fname = "input.txt"
	if *isTest {
		fname = "test-input.txt"
	}

	file, scanner := getScanner(fname)
	defer file.Close()
	initDataStructures(scanner)

	fmt.Println("Solution 1", solution1())
	fmt.Println("Solution 2", solution2())

}
