package cmd

import (
	"bufio"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/vdlbk/playbtg/structs"
	"github.com/vdlbk/playbtg/utils"
	"github.com/vdlbk/playbtg/utils/consts"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

var rootCmd = &cobra.Command{
	Use:              consts.GlobalAppName,
	Short:            consts.GlobalAppName,
	Version:          "v0.1.0",
	TraverseChildren: true,
	Run:              root,
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute rootCmd: %s", err.Error())
	}
}

func root(cmd *cobra.Command, args []string) {
	fmt.Println("CTRL+C to exit the game")
	nbError := 0
	nbSuccess := 0
	var events = make([]structs.Event, 0)
	var doOnce sync.Once

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			printFinal(nbSuccess, nbError, events)
			os.Exit(0)
		}
	}()

	min := len(utils.WordSet)
	max := 0
	for k, _ := range utils.WordSet {
		if k < min {
			min = k
		}
		if k > max {
			max = k
		}
	}

	for {
		x := rand.Intn(max - min) + min
		words := utils.WordSet[x]
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
			fmt.Printf("Enter [%s]: ", expectedText)
			start := time.Now()
			text, _ := reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)
			if text == expectedText {
				events = append(events, structs.Event{expectedText, time.Since(start)})
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

func printFinal(nbSuccess int, nbError int, events []structs.Event) {
	fmt.Println()
	var result = float64(nbSuccess) / float64(nbSuccess+nbError) * 100
	total := int64(0)
	wordsResult := make([][]string, 0)
	for _, e := range events {
		wordsResult = append(wordsResult, []string{e.Word, e.Duration.String()})
		total += e.Duration.Nanoseconds()
	}
	avg := total / int64(len(events))
	tref := time.Unix(0, avg)
	dref := tref.Sub(time.Unix(0, 0))
	wordsResult = append(wordsResult, []string{"Average", dref.String()})

	data := [][]string{
		{"Success", fmt.Sprintf("%d", nbSuccess), fmt.Sprintf("%.2f", result)},
		{"Error", fmt.Sprintf("%d", nbError), fmt.Sprintf("%.2f", 100-result)},
		{"Total", fmt.Sprintf("%d", nbSuccess+nbError), "100.00%"},
	}

	printTable(data, []string{"", "Result", "%"})
	fmt.Println()
	printTable(wordsResult, []string{"Word", "Duration"})
}

func printTable(data [][]string, headers []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(data)
	table.Render()
}