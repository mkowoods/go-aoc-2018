package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
)

var GRAPH = make(map[string][]string)
var INVERTED_GRAPH = make(map[string][]string)

type Item struct {
	value    string
	priority int
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

type TimedNode struct {
	label          string
	startTime      int
	endTime        int
	workerIsActive bool
}

func getSleepTimeFromRune(c byte) int {
	return int(c) - 64
}

func parseInput(inp string) [2]string {
	text, err := regexp.Compile(`Step ([A-Z]) must be finished before step ([A-Z]) can begin.`)
	if err != nil {
		log.Fatal(err)
	}

	textMatch := text.FindAllStringSubmatch(inp, -1)
	// fmt.Println(textMatch)

	return [2]string{textMatch[0][1], textMatch[0][2]}
}

func invertGraph(graph map[string][]string) map[string][]string {
	invertedGraph := make(map[string][]string)
	for from, toArray := range graph {
		_, exists := invertedGraph[from]
		if !exists {
			invertedGraph[from] = make([]string, 0)
		}
		for _, v := range toArray {
			invertedGraph[v] = append(invertedGraph[v], from)
		}
	}
	return invertedGraph
}

func getFirstInstructions(invertedGraph map[string][]string) []string {
	var orphanNodes []string //no parents
	for k, v := range invertedGraph {
		if len(v) == 0 {
			orphanNodes = append(orphanNodes, k)
		}
	}

	return orphanNodes

}

func allParentsInSeenNodes(parents []string, seenNodes map[string]bool) bool {
	var allSeen = true
	for _, parent := range parents {
		allSeen = allSeen && seenNodes[parent]
	}
	return allSeen
}

func getItemsFromNode(node string) *Item {
	item := &Item{
		value:    node,
		priority: -int(node[0]),
	}
	return item
}

func solution1() string {
	// create a directed graph and then the reverse graph
	// determine the parent vertex (with no children)
	// starting with that vertex add it's children to the frontier
	// explore the frontier in alphabetical order (using a heap??)
	// and only explore a node  if all of it's children have also been explored.

	orphanNodes := getFirstInstructions(INVERTED_GRAPH)

	seenNodes := make(map[string]bool)

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	for _, orphanNode := range orphanNodes {
		headItem := getItemsFromNode(orphanNode)
		heap.Push(&pq, headItem)
	}

	var buffer bytes.Buffer
	for pq.Len() > 0 {

		item := heap.Pop(&pq).(*Item)
		node := item.value
		buffer.WriteString(node)
		seenNodes[node] = true
		children := GRAPH[node]
		sort.Strings(children)
		for _, childNode := range children {
			parents := INVERTED_GRAPH[childNode]
			if allParentsInSeenNodes(parents, seenNodes) && !seenNodes[childNode] {
				childItem := getItemsFromNode(childNode)
				heap.Push(&pq, childItem)
			}
		}
	}

	return buffer.String()
}

func updatePriorityQueue(pq *PriorityQueue, node string, seenNodes map[string]bool) {
	children := GRAPH[node]
	for _, childNode := range children {
		parents := INVERTED_GRAPH[childNode]
		if allParentsInSeenNodes(parents, seenNodes) && !seenNodes[childNode] {
			childItem := getItemsFromNode(childNode)
			heap.Push(pq, childItem)
		}
	}
}

func workerIsRunning(workers []*TimedNode) bool {
	var allDone bool
	for _, worker := range workers {
		allDone = allDone || worker.workerIsActive
	}
	return allDone
}

func solution2() string {
	//same as solution1, but with implementation details and ticker

	var numWorkers int = 5
	var delay int = 60

	var ticker int
	var workers = make([]*TimedNode, numWorkers)
	var buffer bytes.Buffer

	seenNodes := make(map[string]bool)
	orphanNodes := getFirstInstructions(INVERTED_GRAPH)
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	for _, orphanNode := range orphanNodes {
		headItem := getItemsFromNode(orphanNode)
		fmt.Println("head item", headItem)
		heap.Push(&pq, headItem)
	}
	for idx := range workers {
		workers[idx] = &TimedNode{}
	}

	for ticker == 0 || workerIsRunning(workers) && (ticker < 999) {

		completedNodes := make([]string, 0)
		//check if a job is finished if so pop the node from the graph and add to the priority queue
		//then reset the worker
		for idx, worker := range workers {
			if worker.endTime == ticker && worker.workerIsActive {
				node := worker.label
				seenNodes[node] = true
				completedNodes = append(completedNodes, node)
				updatePriorityQueue(&pq, node, seenNodes)
				workers[idx] = &TimedNode{}
			}
		}
		// update completed nodes after sorting
		sort.Strings(completedNodes)
		for _, node := range completedNodes {
			buffer.WriteString(node)
		}

		for idx, worker := range workers {
			if !worker.workerIsActive && pq.Len() > 0 {
				item := heap.Pop(&pq).(*Item)
				node := item.value
				// seenNodes[node] = true
				workers[idx] = &TimedNode{
					node,
					ticker,
					(ticker + getSleepTimeFromRune(node[0]) + delay),
					true,
				}
			}
		}
		fmt.Println("Ticker", ticker, workers[0].label, workers[1].label, buffer.String())
		// fmt.Println("Ticker", ticker, workers[0].label, workers[1].label, workers[2].label, workers[3].label, workers[4].label, buffer.String())

		ticker++

	}
	// fmt.Println((ticker))
	return buffer.String()

}

func initDataStructures(scanner *bufio.Scanner) {
	//
	for scanner.Scan() {
		inp := scanner.Text()
		parsedInput := parseInput(inp)
		from := parsedInput[0]
		to := parsedInput[1]
		GRAPH[from] = append(GRAPH[from], to)
	}
	INVERTED_GRAPH = invertGraph(GRAPH)
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
	initDataStructures(scanner)

	fmt.Println("Solution 1", solution1())
	fmt.Println("Solution 2", solution2())

}

// "PFKQWJSVUXEMNIHGTYDOZACRLB"

// "PQWFKJSVUXYEMZDNIHTAGOCRLB"
