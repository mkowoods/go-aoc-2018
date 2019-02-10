package main

import (
	"bufio"
	"container/list"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Coordinate struct {
	id int
	x  int
	y  int
}

type Range struct {
	minX int
	maxX int
	minY int
	maxY int
}

// var MIN_X, MAX_X, MIN_Y, MAX_Y int
var POINT_RANGE Range
var MAX_DIST_FOR_PROB_2 int

func parseInput(id int, inp string) *Coordinate {
	splitInput := strings.Split(inp, ", ")
	x, _ := strconv.Atoi(splitInput[0])
	y, _ := strconv.Atoi(splitInput[1])
	return &Coordinate{id, x, y}
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	return -min(-a, -b)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func getBoundedPoints(coordintates []*Coordinate) []*Coordinate {
	//get the points that have non-inifinte areas. e.g. they are bounded by some other point above or below
	var minX, maxX, minY, maxY int
	for idx, point := range coordintates {
		if idx == 0 {
			minX = point.x
			maxX = point.x
			minY = point.y
			maxY = point.y
		}
		minX = min(minX, point.x)
		maxX = max(maxX, point.x)
		minY = min(minY, point.y)
		maxY = max(maxY, point.y)
	}

	POINT_RANGE = Range{minX, maxX, minY, maxY}

	boundedPoints := make([]*Coordinate, 0)
	for _, point := range coordintates {
		if point.x > minX && point.x < maxX && point.y > minY && point.y < maxY {
			boundedPoints = append(boundedPoints, point)
		}
	}
	fmt.Println("range", minX, maxX, minY, maxY)
	return boundedPoints
}

func getNeighboringPoints(coord *Coordinate) []*Coordinate {
	return []*Coordinate{
		&Coordinate{-1, coord.x + 1, coord.y},
		&Coordinate{-1, coord.x - 1, coord.y},
		&Coordinate{-1, coord.x, coord.y + 1},
		&Coordinate{-1, coord.x, coord.y - 1},
	}
}

func getNearestCoords(coord *Coordinate, coordinatesArray []*Coordinate) []*Coordinate {
	var minDistance int
	var distMap = make(map[int][]*Coordinate)
	for idx, center := range coordinatesArray {
		dist := abs(center.x-coord.x) + abs(center.y-coord.y)
		distMap[dist] = append(distMap[dist], center)
		if idx == 0 {
			minDistance = dist
		}
		minDistance = min(minDistance, dist)
	}
	return distMap[minDistance]
}

func getAreaAroundBoundedPoint(id int, coorindates []*Coordinate) int {

	queue := list.New()
	initCoord := coorindates[id]
	var seenCoordinates = make(map[Coordinate]bool)
	var inArea = make(map[Coordinate]bool)
	queue.PushBack(&Coordinate{-1, initCoord.x, initCoord.y})
	numberOfPointsChecked := 1
	for queue.Len() > 0 {
		// fmt.Println("inArea", inArea)
		// fmt.Println("seenCoordinates", seenCoordinates)
		e := queue.Front()
		queue.Remove(e)
		coord := e.Value.(*Coordinate)
		if coord.x < POINT_RANGE.minX || coord.x > POINT_RANGE.maxX || coord.y < POINT_RANGE.minY || coord.y > POINT_RANGE.maxY {
			return -1
		}
		inArea[*coord] = true
		for _, neighbor := range getNeighboringPoints(coord) {
			if seenCoordinates[*neighbor] {
				continue
			}
			nearestPoints := getNearestCoords(neighbor, coorindates)
			numberOfPointsChecked++

			if len(nearestPoints) == 1 && nearestPoints[0].id == id {
				queue.PushBack(neighbor)
			}
			seenCoordinates[*neighbor] = true
		}
	}
	return len(inArea)
}

func getSumOfDistances(coord *Coordinate, coordinatesArray []*Coordinate) int {
	totalDistance := 0
	for _, center := range coordinatesArray {
		dist := abs(center.x-coord.x) + abs(center.y-coord.y)
		totalDistance += dist

	}
	return totalDistance
}

func getSizeOfRegion(coords []*Coordinate, maxDistance int) int {
	centerOfPoints := &Coordinate{-1, (POINT_RANGE.minX + POINT_RANGE.maxX) / 2, (POINT_RANGE.minY + POINT_RANGE.maxY) / 2}

	fmt.Println("Sum of Distances from Center", getSumOfDistances(centerOfPoints, coords))

	queue := list.New()
	var seenCoordinates = make(map[Coordinate]bool)
	var inArea = make(map[Coordinate]bool)
	queue.PushBack(centerOfPoints)

	for queue.Len() > 0 {
		e := queue.Front()
		queue.Remove(e)
		coord := e.Value.(*Coordinate)
		if coord.x < POINT_RANGE.minX || coord.x > POINT_RANGE.maxX || coord.y < POINT_RANGE.minY || coord.y > POINT_RANGE.maxY {
			return -1
		}
		inArea[*coord] = true
		for _, neighbor := range getNeighboringPoints(coord) {
			if seenCoordinates[*neighbor] {
				continue
			}

			if getSumOfDistances(neighbor, coords) < maxDistance {
				queue.PushBack(neighbor)
			}

			seenCoordinates[*neighbor] = true
		}
	}
	return len(inArea)

}

func solution(scanner *bufio.Scanner) int {
	coordinates := make([]*Coordinate, 0)
	i := 0
	for scanner.Scan() {
		inp := scanner.Text()
		coords := parseInput(i, inp)
		coordinates = append(coordinates, coords)
		i++
	}

	maxArea := 0
	for _, pt := range getBoundedPoints(coordinates) {
		fmt.Println("Point", pt, *pt)
		maxArea = max(maxArea, getAreaAroundBoundedPoint(pt.id, coordinates))
	}

	fmt.Println("Solution to Part 1, Maximum Finite Area", maxArea)
	fmt.Println("Solution to Part 2", getSizeOfRegion(coordinates, MAX_DIST_FOR_PROB_2))
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

	var fname string
	if *isTest {
		fname = "test-input.txt"
		MAX_DIST_FOR_PROB_2 = 32
	} else {
		fname = "input.txt"
		MAX_DIST_FOR_PROB_2 = 10000
	}

	file, scanner := getScanner(fname)
	defer file.Close()
	fmt.Println("Solution Part 1:", solution(scanner))

}
