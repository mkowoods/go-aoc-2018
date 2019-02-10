package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func max(a int, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a int, b int) int {
	return -max(-a, -b)
}

type Point struct {
	posX int
	posY int
	velX int
	velY int
}

type Boundary struct {
	minX int
	maxX int
	minY int
	maxY int
}

var initPointsArray = make([]*Point, 0)

func getBounds(pointsArray []*Point) *Boundary {
	minX := pointsArray[0].posX
	maxX := pointsArray[0].posX
	minY := pointsArray[0].posY
	maxY := pointsArray[0].posY
	for _, pt := range pointsArray {
		minX = min(minX, pt.posX)
		maxX = max(maxX, pt.posX)
		minY = min(minY, pt.posY)
		maxY = max(maxY, pt.posY)
	}
	return &Boundary{minX, maxX, minY, maxY}
}

func printPoint(point *Point) {
	fmt.Println("position", point.posX, point.posY, "velocity", point.velX, point.velY)
}

func printMessage(arryOfPoints []*Point) {
	//used to keep track of the current position of all of the points
	pointMap := make(map[int]map[int]bool)
	for _, pt := range arryOfPoints {
		if pointMap[pt.posY] == nil {
			pointMap[pt.posY] = make(map[int]bool)
		}
		pointMap[pt.posY][pt.posX] = true
	}

	bounds := getBounds(arryOfPoints)
	fmt.Println(bounds)
	for i := bounds.minY; i <= bounds.maxY; i++ {
		var buffer bytes.Buffer
		for j := bounds.minX; j <= bounds.maxX; j++ {
			if pointMap[i][j] {
				buffer.WriteString("#")
			} else {
				buffer.WriteString(".")
			}
		}
		fmt.Println(buffer.String())
	}
}

func movePointsOneTimeStep(pointsArray []*Point) []*Point {
	newPointsArray := make([]*Point, len(pointsArray))
	for idx, pt := range pointsArray {
		newPointsArray[idx] = &Point{
			posX: (pt.posX + pt.velX),
			posY: (pt.posY + pt.velY),
			velX: pt.velX,
			velY: pt.velY,
		}
	}
	return newPointsArray
}

func entropy(pointsArray []*Point) int {

	bounds := getBounds(pointsArray)
	width := abs(bounds.minX - bounds.maxX)
	height := abs(bounds.minY - bounds.maxY)
	return width * height
	// return float64(len(pointsArray)) / (float64(width) * float64(height))
}

func solution1() int {
	bestEntropy := math.MaxInt64
	var bestPointsArray []*Point

	pointsArray := initPointsArray
	seconds := 0
	for {
		newEntropy := entropy(pointsArray)
		if newEntropy > bestEntropy {
			//entropy measured as the relative number of dpi
			//should steadily decrease until it forms the message and then will start to increase
			//since the number of pixels doesnt change we can just measure the area of the points
			break
		}
		bestEntropy = newEntropy
		bestPointsArray = pointsArray
		pointsArray = movePointsOneTimeStep(pointsArray)
		seconds++
	}
	fmt.Println("Num Seconds", (seconds - 1))
	printMessage(bestPointsArray)
	return -1
}

func parseRow(row string) *Point {
	r, err := regexp.Compile(`position=<\s?(\-?\d+), \s?(\-?\d+)> velocity=<\s?(\-?\d+), \s?(\-?\d+)>`)
	if err != nil {
		log.Fatal(err)
	}

	rowMatch := r.FindAllStringSubmatch(row, -1)
	posX, _ := strconv.Atoi(rowMatch[0][1])
	posY, _ := strconv.Atoi(rowMatch[0][2])
	velX, _ := strconv.Atoi(rowMatch[0][3])
	velY, _ := strconv.Atoi(rowMatch[0][4])

	return &Point{
		posX: posX,
		posY: posY,
		velX: velX,
		velY: velY,
	}
}

func initDataStructures(scanner *bufio.Scanner) {
	for scanner.Scan() {
		inp := scanner.Text()
		point := parseRow(inp)
		initPointsArray = append(initPointsArray, point)
	}
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
	isTest := flag.Bool("test", true, "whether to use test file")

	flag.Parse()

	var fname = "test-input.txt"
	if !(*isTest) {
		fname = "input.txt"
	}

	file, scanner := getScanner(fname)
	defer file.Close()
	initDataStructures(scanner)

	fmt.Println("Solution 1", solution1())
	// fmt.Println("Solution 2", solution2())

}
