package main

import (
	"fmt"
	"math/rand"
	"time"
)

const rows = 15
const cols = 30

var gameboard [rows][cols]int

type Universe struct {
	rows      int
	cols      int
	gameboard [][]int
	ruleset   func()
}

// func (t* Universe) create(rows int, cols int, rules func()) {
// }

func displayGameboard() { //[][]int) {
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if gameboard[i][j] >= 0 {
				fmt.Printf("%d ", gameboard[i][j])
			} else {
				fmt.Printf("- ")
			}
		}
		fmt.Printf("\n")
	}
}

func randomInit() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if rand.Intn(1000) > 970 {
				gameboard[i][j] = 0
			} else {
				gameboard[i][j] = -1
			}
		}
	}
}

func ruleset1() {
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			gameboard[i][j] = 2
		}
	}
}

// func applyRuleset(ruleset func()) {
// }

func main() {
	randomInit()
	displayGameboard()
}
