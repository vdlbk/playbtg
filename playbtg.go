package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

type Event struct {
	Word     string
	Duration time.Duration
}

func (e *Event) String() string {
	return fmt.Sprintf("%s: %s", e.Duration.String(), e.Word)
}

func main() {
	fmt.Println("CTRL+C to exit the game")
	nbError := 0
	nbSuccess := 0
	var events = make([]Event, 0)
	var doOnce sync.Once

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Println(sig.String())
			PrintFinal(nbSuccess, nbError, events)
		}
	}()

	for {
		rand.Seed(time.Now().UnixNano())
		x := rand.Intn(6) + 2
		words := WordSet[x]
		y := rand.Intn(len(words))

		expectedText := words[y]
		reader := bufio.NewReader(os.Stdin)

		doOnce.Do(func() {
			fmt.Println("Ready? Press any key")
			if _, _, err := reader.ReadRune(); err != nil {
				os.Exit(1)
			}
		})

		for {
			start := time.Now()
			fmt.Printf("Enter [%s]: ", expectedText)
			text, _ := reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)
			if text == expectedText {
				events = append(events, Event{expectedText, time.Since(start)})
				nbSuccess++
				break
			} else {
				nbError++
			}

			if nbError > 2000000 {
				os.Exit(1)
			}
		}
	}
}

func PrintFinal(nbSuccess int, nbError int, events []Event) {
	var result = float64(nbSuccess) / float64(nbSuccess+nbError) * 100
	fmt.Printf("Success: %d | Error: %d (%.2f)\n", nbSuccess, nbError, result)
	total := int64(0)
	for _, e := range events {
		total += e.Duration.Nanoseconds()
		fmt.Println(e.String())
	}
	avg := total / int64(len(events))
	tref := time.Unix(0, avg)
	dref := tref.Sub(time.Unix(0, 0))
	fmt.Printf("Moyenne: %s\n", dref.String())
}
