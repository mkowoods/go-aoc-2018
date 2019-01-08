package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type ShiftFeed struct {
	year       int
	month      int
	day        int
	hour       int
	minute     int
	guardID    int
	asleep     bool
	awake      bool
	shiftStart bool
}

type SleepSegment struct {
	start int
	end   int
}

func argMaxMap(intMap map[int]int) int {
	var best int
	var maxValue int
	for k := range intMap {
		if intMap[k] >= maxValue {
			best = k
			maxValue = intMap[k]
		}
	}
	return best
}

func getTotalAsleepMinutesByGuard(feed []ShiftFeed) map[int]int {

	var sleepSum = make(map[int]int)
	var lastAsleepMinute int
	for _, sf := range feed {
		if sf.asleep {
			lastAsleepMinute = sf.minute
		}
		if sf.awake {
			sleepSum[sf.guardID] += (sf.minute - lastAsleepMinute)
		}
	}
	return sleepSum
}

func getMaxAsleepMinute(guardID int, feed []ShiftFeed) int {
	sleepMinuteHist := getSleepHist(guardID, feed)
	return argMaxMap(sleepMinuteHist)
}

func getSleepHist(guardID int, feed []ShiftFeed) map[int]int {

	var sleepSegmentArray []SleepSegment

	var lastSegment SleepSegment
	for _, sf := range feed {
		if sf.guardID != guardID {
			continue
		}
		if sf.asleep {
			lastSegment = SleepSegment{sf.minute, -1}
		}
		if sf.awake {
			lastSegment.end = sf.minute //they stop sleeping on the prior minute
			sleepSegmentArray = append(sleepSegmentArray, lastSegment)
		}
	}

	sleepMinuteHist := make(map[int]int)
	for _, seg := range sleepSegmentArray {
		for min := seg.start; min < seg.end; min++ {
			sleepMinuteHist[min]++
		}
	}
	return sleepMinuteHist
}

func parseInput(inp string) ShiftFeed {
	guard, err := regexp.Compile(`\[(\d+)-(\d+)-(\d+) (\d+):(\d+)\] Guard #(\d+) begins shift`)
	asleepOrAwake, err := regexp.Compile(`\[(\d+)-(\d+)-(\d+) (\d+):(\d+)\] (falls asleep|wakes up)`)

	if err != nil {
		log.Fatal(err)
	}
	guardMatch := guard.FindAllStringSubmatch(inp, -1)
	asleepOrAwakeMatch := asleepOrAwake.FindAllStringSubmatch(inp, -1)

	var year, month, day, hour, minute int
	var asleep, awake, shiftStart bool
	var guardID int
	if len(guardMatch) > 0 {
		year, _ = strconv.Atoi(guardMatch[0][1])
		month, _ = strconv.Atoi(guardMatch[0][2])
		day, _ = strconv.Atoi(guardMatch[0][3])
		hour, _ = strconv.Atoi(guardMatch[0][4])
		minute, _ = strconv.Atoi(guardMatch[0][5])
		guardID, _ = strconv.Atoi(guardMatch[0][6])
		shiftStart = true
	}
	if len(asleepOrAwakeMatch) > 0 {
		if asleepOrAwakeMatch[0][6] == "falls asleep" {
			asleep = true
		}
		if asleepOrAwakeMatch[0][6] == "wakes up" {
			awake = true
		}
		year, _ = strconv.Atoi(asleepOrAwakeMatch[0][1])
		month, _ = strconv.Atoi(asleepOrAwakeMatch[0][2])
		day, _ = strconv.Atoi(asleepOrAwakeMatch[0][3])
		hour, _ = strconv.Atoi(asleepOrAwakeMatch[0][4])
		minute, _ = strconv.Atoi(asleepOrAwakeMatch[0][5])
	}
	return ShiftFeed{year, month, day, hour, minute, guardID, asleep, awake, shiftStart}
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

func solutionPart1(scanner *bufio.Scanner) int {
	var feed []ShiftFeed
	var lastSeenGuardID int
	for scanner.Scan() {
		inp := scanner.Text()
		out := parseInput(inp)
		if out.guardID == 0 {
			out.guardID = lastSeenGuardID
		} else {
			lastSeenGuardID = out.guardID
		}
		feed = append(feed, out)
	}

	sleepSum := getTotalAsleepMinutesByGuard(feed)
	mostAsleepGuard := argMaxMap(sleepSum)
	fmt.Println("Most Asleep Guard", mostAsleepGuard)

	mostAsleepMinute := getMaxAsleepMinute(mostAsleepGuard, feed)
	fmt.Println("Max Asleep Minute", mostAsleepMinute)
	return mostAsleepGuard * mostAsleepMinute
}

func solutionPart2(scanner *bufio.Scanner) int {
	var feed []ShiftFeed
	var lastSeenGuardID int
	for scanner.Scan() {
		inp := scanner.Text()
		out := parseInput(inp)
		if out.guardID == 0 {
			out.guardID = lastSeenGuardID
		} else {
			lastSeenGuardID = out.guardID
		}
		feed = append(feed, out)
	}

	sleepSum := getTotalAsleepMinutesByGuard(feed)
	var bestGuard, highestFrequency, bestMinute int
	for guardID := range sleepSum {
		sleepHist := getSleepHist(guardID, feed)
		maxMinute := argMaxMap(sleepHist)
		if sleepHist[maxMinute] > highestFrequency {
			bestGuard = guardID
			highestFrequency = sleepHist[maxMinute]
			bestMinute = maxMinute
		}
	}
	return bestGuard * bestMinute
}

func main() {
	// file sorted with the below command:
	// sort -h input.txt > tmp && mv tmp input.txt
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
