package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/hokiegeek/life"
	"os"
	"time"
)

func displaypond(strategy *life.Life, rate time.Duration, iterations int, static, paused bool) {
	// Clear the screen and put the cursor on the top left
	if static {
		fmt.Print("\033[2J")
		fmt.Print("\033[H")
	}

	// Print the seed
	fmt.Print(strategy)

	if paused {
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
	}

	updates := make(chan bool)
	stop := strategy.Start(updates, rate)

	countGenerations := 1
	for {
		select {
		case <-updates:
			if static {
				fmt.Print("\033[H")
			}
			fmt.Print(strategy)

			if iterations >= 0 {
				countGenerations++
				if countGenerations >= iterations {
					stop()
					break
				}
			}
		}
	}
}

func displayFromFile(filename string) {
	_, err := life.Read(filename)
	if err != nil {
		fmt.Printf("Well, crap: %s\n", err)
	}
}

func displayTestpond(width int, height int, rate time.Duration, initializer func(life.Dimensions, life.Location) []life.Location) {
	strategy, err := life.New("",
		life.Dimensions{Height: height, Width: width},
		life.NEIGHBORS_ALL,
		initializer,
		life.ConwayTester(),
		life.SimultaneousProcessor)
	if err == nil {
		displaypond(strategy, rate, -1, true, true)
	} else {
		fmt.Printf("Could not create: %s\n", err)
	}
}

func main() {

	patternPtr := flag.String("pattern", "random", "Specify the pattern to run")
	widthPtr := flag.Int("width", 1, "Width of the Life board")
	heightPtr := flag.Int("height", 1, "Height of the Life board")
	ratePtr := flag.Duration("rate", -1, "Rate at which the board should be updated in milliseconds")

	flag.Parse()

	switch *patternPtr {
	case "blinkers":
		width := 9
		if *widthPtr > width {
			width = *widthPtr
		}
		height := 9
		if *heightPtr > height {
			height = *heightPtr
		}

		displayTestpond(width, height, *ratePtr, life.Blinkers)
	case "toads":
		width := 10
		if *widthPtr > width {
			width = *widthPtr
		}
		height := 10
		if *heightPtr > height {
			height = *heightPtr
		}

		displayTestpond(width, height, *ratePtr, life.Toads)
	case "glider":
		width := 30
		if *widthPtr > width {
			width = *widthPtr
		}
		height := 30
		if *heightPtr > height {
			height = *heightPtr
		}

		displayTestpond(width, height, *ratePtr,
			func(dimensions life.Dimensions, offset life.Location) []life.Location {
				return life.Gliders(life.Dimensions{Height: 4, Width: 4}, offset)
			})
	case "pulsar":
		width := 15
		if *widthPtr > width {
			width = *widthPtr
		}
		height := 15
		if *heightPtr > height {
			height = *heightPtr
		}

		displayTestpond(width, height, *ratePtr, life.Pulsar)
	case "random":
		width := 120
		if *widthPtr > width {
			width = *widthPtr
		}
		height := 45
		if *heightPtr > height {
			height = *heightPtr
		}

		displayTestpond(width, height, *ratePtr,
			func(dimensions life.Dimensions, offset life.Location) []life.Location {
				return life.Random(dimensions, offset, 85)
			})
	default:
		if _, err := os.Stat(*patternPtr); err == nil {
			displayFromFile(*patternPtr)
		} else {
			fmt.Println("Did not recognize pattern")
		}
	}
}
