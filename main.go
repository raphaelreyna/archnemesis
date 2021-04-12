package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

var (
	filePath     string
	genCount     int
	timeIncrease int
)

func init() {
	flag.IntVar(&genCount, "g", 1, "Number of generations to simulate.")
	flag.IntVar(&timeIncrease, "t", 0, "Number of seconds to add to the simulation time each generation.")
	flag.StringVar(&filePath, "f", "", "The path to the save file to use.")
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().Unix())

	// Create a new game for the players to play
	game, err := NewGame()
	if err != nil {
		panic(err)
	}

	// Create initial generation of players
	players := make(Generation, 10)
	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}

		err = json.NewDecoder(file).Decode(&players)
		file.Close()
		if err != nil {
			panic(err)
		}

		fmt.Println("loaded initial generation from " + filePath)
	} else {
		for idx := range players {
			players[idx] = NewPlayer()
		}
	}

	// Graphics need to run on main goroutine so we run genetic algorithm logic in a seperate goroutine.
	go func() {
		// Loop over the number of generations we'll be simulating
		for n := 0; n < genCount; n++ {
			fmt.Printf("Starting gen %d\t %d sec/simulation\n", n+1, 5+timeIncrease*n)

			// Give each player in this generation a turn at playing the game for some amount of time.
			wg := sync.WaitGroup{}
			for idx := range players {
				wg.Add(1)

				// Start a timer to record the players score and reset the simulation for the next player
				time.AfterFunc(time.Duration(5+timeIncrease*n)*time.Second, func() {
					players[idx].Score = game.Score
					game.Reset()
					fmt.Printf(" scored: %v\n", players[idx].Score)
					wg.Done()
				})

				// Bring the player into the game world
				game.Player = &players[idx]

				fmt.Printf("%v", players[idx].Genes)

				// Wait for the players time in the game world to be over
				wg.Wait()
			}

			// Replace this generation of players with the next
			players = players.Next(0.3)
		}

		// Store the results
		if filePath == "" {
			return
		}

		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}

		err = json.NewEncoder(file).Encode(&players)
		file.Close()
		if err != nil {
			panic(err)
		}

		os.Exit(0)
	}()

	pixelgl.Run(game.Run)
}
