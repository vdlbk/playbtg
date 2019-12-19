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
	"strconv"
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

var gameConfig = structs.GameConfig{}

func init() {
	rand.Seed(time.Now().UnixNano())

	rootCmd.Flags().BoolVarP(&gameConfig.UpperMode, consts.ParamUpperMode, "u", false, "The words will be displayed in uppercase")
	rootCmd.Flags().BoolVarP(&gameConfig.MixUpperLowerMode, consts.ParamMixUpperLowerMode, "m", false, "The words will be displayed with a mix of character in uppercase and lowercase")
	rootCmd.Flags().BoolVarP(&gameConfig.NumberMode, consts.ParamNumberMode, "n", false, "The words will be replaced by numbers")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Failed to execute rootCmd: %s", err.Error())
		os.Exit(1)
	}
}

func displayStart() {
	fmt.Println(consts.AppTag)
	gameConfig.Render()
	fmt.Println("Press 'CTRL+C' to exit the game")
}

func root(cmd *cobra.Command, args []string) {
	displayStart()
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

	gameConfig.WordSetMinLength, gameConfig.WordSetMaxLength = utils.ComputeBounds(utils.WordSet)

	for {
		expectedText := generateWord(gameConfig)
		reader := bufio.NewReader(os.Stdin)

		doOnce.Do(func() {
			fmt.Println("Ready? Press any key to start the game")
			if _, _, err := reader.ReadRune(); err != nil {
				fmt.Printf("Failed to start the game: %v\n", err.Error())
				os.Exit(1)
			}
		})

		event := structs.Event{Word: expectedText}
		for {
			fmt.Printf("Enter [%s]: ", expectedText)
			start := time.Now()
			text, _ := reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)
			if text == expectedText {
				event.Duration = time.Since(start)
				events = append(events, event)
				nbSuccess++
				break
			} else {
				event.Attempts = append(event.Attempts, text)
				nbError++
			}

			// Safety circuit breaker
			if nbError > consts.SafetyCircuitBreakerLimit {
				log.Println("An unknown issue occurred. The game has stopped to prevent any other problem")
				os.Exit(1)
			}
		}
	}
}
func generateWord(config structs.GameConfig) string {
	if config.NumberMode {
		number := strconv.FormatInt(time.Now().UnixNano(), 10)
		x := rand.Intn(consts.NumberModeMaxLength) + 1
		return number[len(number)-x:]
	}

	x := rand.Intn(config.WordSetMaxLength-config.WordSetMinLength) + config.WordSetMinLength
	words := utils.WordSet[x]
	y := rand.Intn(len(words))

	return transformWord(words[y], gameConfig)
}

func transformWord(word string, config structs.GameConfig) string {
	if config.UpperMode {
		word = strings.ToUpper(word)
	} else if config.MixUpperLowerMode {
		transformedWord := ""
		for _, c := range word {
			if utils.FiftyFifty() {
				transformedWord += strings.ToUpper(string(c))
			} else {
				transformedWord += string(c)
			}
		}
		word = transformedWord
	}
	return word
}

func printFinal(nbSuccess int, nbError int, events []structs.Event) {
	fmt.Println()
	if len(events) == 0 {
		return
	}

	var result = float64(nbSuccess) / float64(nbSuccess+nbError) * 100
	total := int64(0)
	wordsResult := make([][]string, 0)
	for _, e := range events {

		// Compute avg time by letter
		duration := int64(e.Duration / time.Millisecond)
		avg := float64(duration) / float64(len(e.Word))
		avgDuration := time.Duration(avg) * time.Millisecond

		wordsResult = append(wordsResult, []string{e.Word, e.Duration.String(), avgDuration.String(), e.Attempts.String()})
		total += e.Duration.Nanoseconds()
	}
	avg := total / int64(len(events))
	tref := time.Unix(0, avg)
	dref := tref.Sub(time.Unix(0, 0))
	wordsResult = append(wordsResult, []string{"Average", dref.String(), ""})

	data := [][]string{
		{"Success", fmt.Sprintf("%d", nbSuccess), fmt.Sprintf("%.2f", result)},
		{"Error", fmt.Sprintf("%d", nbError), fmt.Sprintf("%.2f", 100-result)},
		{"Total", fmt.Sprintf("%d", nbSuccess+nbError), "100.00%"},
	}

	printTable(data, []string{"", "Result", "%"})
	fmt.Println()
	printTable(wordsResult, []string{"Word", "Duration", "Duration/letter", "Attempts"})
}

func printTable(data [][]string, headers []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	//table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(data)
	table.Render()
}
