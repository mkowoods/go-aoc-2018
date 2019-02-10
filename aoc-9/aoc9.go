package main

import (
	"bytes"
	"flag"
	"fmt"
	"strconv"
)

var NUM_PLAYERS int
var LAST_MARBLE int
var PLAYER_SCORE = make(map[int]int)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

//Going to use a ring list for storing the marbles in Memory, this means each Marble will have a pointer to it's clockwise and counter clocwise neighbor

type Marble struct {
	value                  int
	clocwiseMarble         *Marble
	counterClockwiseMarble *Marble
}

func getMarbleNSteps(currentMarble *Marble, steps int) *Marble {
	//use postive one to move clockwise and negative  to move counter clockwise
	moveCounterClockwise := (steps < 0)
	steps = abs(steps)
	for i := 0; i < steps; i++ {
		if moveCounterClockwise {
			currentMarble = currentMarble.counterClockwiseMarble
		} else {
			currentMarble = currentMarble.clocwiseMarble
		}
	}
	return currentMarble
}

func printMarblesMovingClockwise(currentMarble *Marble) string {
	var buffer bytes.Buffer
	initMarbleValue := currentMarble.value
	buffer.WriteString("(" + strconv.Itoa(initMarbleValue) + ")")
	for {

		currentMarble = currentMarble.clocwiseMarble
		if currentMarble.value == initMarbleValue {
			break
		}
		valString := strconv.Itoa(currentMarble.value)
		buffer.WriteString(" " + valString)
	}
	// fmt.Println("Marbles From Current Moving Clockwise")
	return buffer.String()
}

func addMarble(currentMarble *Marble, newMarbleValue int, playerIdx int) *Marble {
	if (newMarbleValue % 23) == 0 {
		marbleToBeRemoved := getMarbleNSteps(currentMarble, -7)

		marbleCounterClockwise := getMarbleNSteps(marbleToBeRemoved, -1)
		marbleClockwise := getMarbleNSteps(marbleToBeRemoved, 1)

		marbleCounterClockwise.clocwiseMarble = marbleClockwise
		marbleClockwise.counterClockwiseMarble = marbleCounterClockwise

		currentMarble = marbleClockwise

		//update score & delete removed marble
		PLAYER_SCORE[playerIdx] += (newMarbleValue + marbleToBeRemoved.value)
		marbleToBeRemoved = nil

	} else {

		marbleCounterClockwise := getMarbleNSteps(currentMarble, 1)
		marbleClockwise := getMarbleNSteps(currentMarble, 2)

		//Create the new Marble and define it's neighbors
		newMarble := &Marble{
			value: newMarbleValue,
			counterClockwiseMarble: marbleCounterClockwise,
			clocwiseMarble:         marbleClockwise,
		}
		//update the neighbors for the marbles it's inserted between
		marbleClockwise.counterClockwiseMarble = newMarble
		marbleCounterClockwise.clocwiseMarble = newMarble
		currentMarble = newMarble

	}

	return currentMarble
}

func getHighestScore() int {
	var maxValue int
	for _, val := range PLAYER_SCORE {
		if val > maxValue {
			maxValue = val
		}
	}
	return maxValue
}

func solution1() int {
	//go run aoc9.go -players=441 -last-marble=71032
	//initGame
	currentMarble := &Marble{
		value: 0,
	}
	currentMarble.clocwiseMarble = currentMarble
	currentMarble.counterClockwiseMarble = currentMarble

	newMarbleValue := 1
	playerIdx := 1
	for i := 0; i < LAST_MARBLE; i++ {
		currentMarble = addMarble(currentMarble, newMarbleValue, playerIdx)
		// fmt.Println("Marble Counter Clockwise", getMarbleNSteps(currentMarble, -1).value)
		// fmt.Println("Marble Clockwise", getMarbleNSteps(currentMarble, 1).value)
		// fmt.Println("["+strconv.Itoa(playerIdx)+"]", newMarbleValue, printMarblesMovingClockwise(currentMarble))
		newMarbleValue++
		playerIdx++
		if (playerIdx % NUM_PLAYERS) == 0 {
			playerIdx = NUM_PLAYERS
		} else {
			playerIdx = (playerIdx % NUM_PLAYERS)
		}
	}

	return getHighestScore()
}

func solution2() int {
	return -1
}

func main() {
	numPlayers := flag.Int("players", 0, "Number of Players")
	lastMarble := flag.Int("last-marble", 0, "Index of Last Marble")

	flag.Parse()

	NUM_PLAYERS = *numPlayers
	LAST_MARBLE = *lastMarble

	fmt.Println("Solution 1", solution1())
	fmt.Println("Solution 2", solution2())

}
