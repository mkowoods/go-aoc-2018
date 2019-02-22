package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
)

type bot struct {
	X      int
	Y      int
	Z      int
	radius int
}

type botNet struct {
	node      int
	neighbors []int
}

type region struct {
	XMin int
	XMax int
	YMin int
	YMax int
	ZMin int
	ZMax int
}

type serchCube struct {
	X     int
	Y     int
	Z     int
	width int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	return -min(-x, -y)
}

func parseInput(inp string) bot {
	// fmt.Println(inp)
	regexStr := `pos=<(-?\d+),(-?\d+),(-?\d+)>, r=(\d+)`
	r, err := regexp.Compile(regexStr)
	if err != nil {
		log.Fatal(err)
	}
	out := r.FindAllStringSubmatch(inp, -1)
	X, _ := strconv.Atoi(out[0][1])
	Y, _ := strconv.Atoi(out[0][2])
	Z, _ := strconv.Atoi(out[0][3])
	radius, _ := strconv.Atoi(out[0][4])

	bot := bot{X, Y, Z, radius}
	return bot
}

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func botManhattanDist(bot1 bot, bot2 bot) int {
	return abs(bot1.X-bot2.X) + abs(bot1.Y-bot2.Y) + abs(bot1.Z-bot2.Z)
}

func isSubset(set []int, superset []int) bool {
	// returns true if set is entirely contained in superset
	for _, k := range set {
		if !intInSlice(k, superset) {
			return false
		}
	}
	return true
}

func solutionPart1(scanner *bufio.Scanner) int {

	var botArray []bot
	var strongestBot bot //bot with largest radius
	var numBotsInRange int

	for scanner.Scan() {
		inp := scanner.Text()
		fmt.Println(inp)
		bot := parseInput(inp)

		botArray = append(botArray, bot)
		if strongestBot.radius < bot.radius {
			strongestBot = bot
		}
	}

	for _, bot := range botArray {
		dist := botManhattanDist(strongestBot, bot)
		if dist <= strongestBot.radius {
			fmt.Println(bot)
			numBotsInRange++
		}
	}

	return numBotsInRange
}

func buildBotAdjacencyGraph(botArray []bot) map[int][]int {
	//Note this includes the source vertex in the list of edges
	// O(n^2) -- building adjacency list of bots that can reach other bots
	var botNet = make(map[int][]int)
	for i := 0; i < len(botArray); i++ {
		for j := (i + 1); j < len(botArray); j++ {
			bot1 := botArray[i]
			bot2 := botArray[j]
			botDistance := botManhattanDist(bot1, bot2)
			totalRange := (bot1.radius + bot2.radius)
			if botDistance <= totalRange {
				botNet[i] = append(botNet[i], j)
				botNet[j] = append(botNet[j], i)
			}
		}
	}
	return botNet
}

func getSortedSubGraphs(botGraph map[int][]int) [][]int {
	// takes the botGraph and returns the set of vertices that are connected through at least one vertex in order
	// of the number of edges these will then be evaluated to determine the maximum clique
	var clusters [][]int
	keys := make([]int, 0, len(botGraph))
	for k := range botGraph {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i int, j int) bool {
		return len(botGraph[i]) > len(botGraph[j])
	})

	for k := range keys {
		vertices := botGraph[k]
		candidateCluster := make([]int, 0, len(vertices)+1)
		candidateCluster = append(candidateCluster, k)
		candidateCluster = append(candidateCluster, vertices...)
		clusters = append(clusters, candidateCluster)
	}

	return clusters
}

func getLargestBotClique(botGraph map[int][]int) []int {

	for _, subGraph := range getSortedSubGraphs(botGraph) {
		// fmt.Println(subGraph)
		allInRange := true
		for _, v := range subGraph {
			vertexSubGraph := botGraph[v]
			vertexSubGraph = append(vertexSubGraph, v)
			fullyConnected := isSubset(subGraph, vertexSubGraph)
			if !fullyConnected {
				allInRange = false
				break
			}
		}
		if allInRange {
			return subGraph
		}
	}
	return []int{}
}

func getBotRegion(bot bot) region {
	XMin := bot.X - bot.radius
	XMax := bot.X + bot.radius
	YMin := bot.Y - bot.radius
	YMax := bot.Y + bot.radius
	ZMin := bot.Z - bot.radius
	ZMax := bot.Z + bot.radius
	return region{XMin, XMax, YMin, YMax, ZMin, ZMax}
}

func mergeRegion(region1, region2 region) region {
	XMin := max(region1.XMin, region2.XMin)
	XMax := min(region1.XMax, region2.XMax)
	YMin := max(region1.YMin, region2.YMin)
	YMax := min(region1.YMax, region2.YMax)
	ZMin := max(region1.ZMin, region2.ZMin)
	ZMax := min(region1.ZMax, region2.ZMax)

	return region{XMin, XMax, YMin, YMax, ZMin, ZMax}
}

func getFeasibilityRegion(botArray []bot, botSubGroup []int) region {

	if len(botSubGroup) == 0 {
		log.Fatal("SubGroup is Zero")
	}
	fmt.Println("init bot", botArray[botSubGroup[0]])
	mergedRegion := getBotRegion(botArray[botSubGroup[0]])

	for _, botID := range botSubGroup {
		fmt.Println("bot", botID, getBotRegion(botArray[botID]))
		botRegion := getBotRegion(botArray[botID])

		mergedRegion = mergeRegion(mergedRegion, botRegion)
		if mergedRegion.XMax < mergedRegion.XMin {
			fmt.Println("botID", botID, "botRegion", botRegion)
			fmt.Println(mergedRegion)
		}

	}
	return mergedRegion
}

func getNearestDistanceToOrigin(bot1 bot) int {
	nullBot := bot{0, 0, 0, 0}
	return abs(botManhattanDist(nullBot, bot1) - bot1.radius)
}

func botIntersectsRegion(r region, bot bot) bool {
	//there exists a poin in the region that tbe bot also crosses

	// region block {12 12 13 14 10 12} {16 12 12 4} true

	x := (r.XMin >= bot.X-bot.radius) && (r.XMax <= bot.X+bot.radius)
	y := (r.YMin >= bot.Y-bot.radius) && (r.YMax <= bot.Y+bot.radius)
	z := (r.ZMin >= bot.Z-bot.radius) && (r.ZMax <= bot.Z+bot.radius)
	return x && y && z
}

func getNearestPointInFeasibilityRegion(fr region, botArray []bot, botSubGroup []int) (bot, error) {
	// fr = feasibilityRegion
	fmt.Println("Feasibility Region", fr)
	fmt.Println("Size of Bot Group", len(botSubGroup))

	xSplit := (fr.XMax + fr.XMin) / 2
	ySplit := (fr.YMax + fr.YMin) / 2
	zSplit := (fr.ZMax + fr.ZMin) / 2

	for _, x := range [2]int{fr.XMin, fr.XMax} {
		for _, y := range [2]int{fr.YMin, fr.YMax} {
			for _, z := range [2]int{fr.ZMin, fr.ZMax} {
				xmin := min(x, xSplit)
				xmax := max(x, xSplit)
				ymin := min(y, ySplit)
				ymax := max(y, ySplit)
				zmin := min(z, zSplit+1)
				zmax := max(z, zSplit)
				r := region{xmin, xmax, ymin, ymax, zmin, zmax}
				for _, botID := range botSubGroup {
					fmt.Println("region block", r, botArray[botID], botIntersectsRegion(r, botArray[botID]))
				}
			}
		}
	}

	pointBot := bot{xSplit, ySplit, zSplit, 0}
	// pointFound := true

	fmt.Println("pointBot", pointBot)
	for _, botID := range botSubGroup {
		bot := botArray[botID]
		inRadius := botManhattanDist(pointBot, bot) > bot.radius
		fmt.Println(botID, bot, inRadius)
		// if botManhattanDist(pointBot, bot) > bot.radius {
		// 	pointFound = false
		// 	break
		// }
	}
	return bot{-1, -1, -1, -1}, errors.New("Bot Not Found")
}

func solutionPart2(scanner *bufio.Scanner) int {
	var botArray []bot

	for scanner.Scan() {
		inp := scanner.Text()
		// fmt.Println(inp)
		bot := parseInput(inp)

		botArray = append(botArray, bot)
	}

	botGraph := buildBotAdjacencyGraph(botArray)
	// fmt.Println(botGraph)
	bestsubGraph := getLargestBotClique(botGraph)
	for _, id := range bestsubGraph {

		fmt.Println(id, botArray[id], getNearestDistanceToOrigin(botArray[id]))
	}
	fmt.Println("Size of Best SubGraph", len(bestsubGraph))

	mergedRegion := getFeasibilityRegion(botArray, bestsubGraph)
	fmt.Println(mergedRegion)
	bestPoint, err := getNearestPointInFeasibilityRegion(mergedRegion, botArray, bestsubGraph)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(bestsubGraph)
	fmt.Println(mergedRegion)
	fmt.Println(bestPoint)
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
		fname = "test-input-2.txt"
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
		// region1 := region{10, 10, -20, 20, -30, 30}
		// region2 := region{-9, -5, -10, 10, -30, 30}
		// fmt.Println("r1", region1)
		// fmt.Println("r2", region2)
		// fmt.Println("mr", mergeRegion(region1, region2))

	}

}
