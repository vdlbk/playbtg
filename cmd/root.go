package cmd

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/gosuri/uilive"
	"github.com/spf13/cobra"
	"github.com/vdlbk/playbtg/renderers"
	"github.com/vdlbk/playbtg/structs"
	"github.com/vdlbk/playbtg/utils"
	"github.com/vdlbk/playbtg/utils/consts"
)

const (
	DefaultNumberOfWordsGenerated = 50
	NumberOfWordsDisplayed        = 10
	DefaultOutput                 = "Console"
)

var rootCmd = &cobra.Command{
	Use:              consts.GlobalAppName,
	Short:            consts.GlobalAppName,
	Version:          "v0.2.1",
	TraverseChildren: true,
	Run:              root,
}

var gameConfig = structs.GameConfig{}

func init() {
	rand.Seed(time.Now().UnixNano())

	rootCmd.Flags().BoolVarP(&gameConfig.UpperMode, consts.ParamUpperMode, "u", false, "The words will be displayed in uppercase")
	rootCmd.Flags().BoolVarP(&gameConfig.MixUpperLowerMode, consts.ParamMixUpperLowerMode, "m", false, "The words will be displayed with a mix of character in uppercase and lowercase")
	rootCmd.Flags().BoolVarP(&gameConfig.NumberMode, consts.ParamNumberMode, "n", false, "The words will be replaced by numbers")
	rootCmd.Flags().BoolVarP(&gameConfig.InfiniteAttempts, consts.ParamInfiniteAttempts, "i", false, "You have an infinite numbers of attempts for each words (By default, you only have 1 attempt)")
	rootCmd.Flags().StringVarP(&gameConfig.Output, consts.ParamOutput, "o", DefaultOutput, "Specify the folder in which it will create a save the result into a file")
}

func checkConfig(gameConfig structs.GameConfig) error {
	if gameConfig.Output != DefaultOutput {
		fileInfo, err := os.Stat(gameConfig.Output)
		if err != nil {
			return err
		}
		if !fileInfo.IsDir() {
			return fmt.Errorf("%s is not a valid folder path", gameConfig.Output)
		}
	}
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Failed to execute %s: %s", consts.GlobalAppName, err.Error())
		os.Exit(1)
	}
}

func displayStart() {
	fmt.Println(consts.AppTag)
	gameConfig.Render()
	fmt.Println("Press 'CTRL+C' or 'ESC' to exit the game")
}

func root(_ *cobra.Command, _ []string) {
	if err := checkConfig(gameConfig); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	displayStart()
	nbError := 0
	nbSuccess := 0
	var events = make([]structs.Event, 0)

	// init Writer for instructions
	writer := uilive.New()
	writer.Start()
	defer writer.Stop()

	// init channel to handle signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, os.Kill)
	go func() {
		for range signals {
			finalStep(nbSuccess, nbError, events)
		}
	}()

	gameConfig.WordSetMinLength, gameConfig.WordSetMaxLength = utils.ComputeBounds(utils.WordSet)

	// open keyboard event listener
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	runGame(writer, &events, &nbSuccess, &nbError)

	finalStep(nbSuccess, nbError, events)
}

func runGame(writer *uilive.Writer, events *[]structs.Event, nbSuccess, nbError *int) {
	var doOnce sync.Once
	words := generateWords(gameConfig)

	for idx, expectedText := range words {
		reader := bufio.NewReader(os.Stdin)

		doOnce.Do(func() {
			fmt.Println("Ready? Press any key to start the game")
			if _, _, err := reader.ReadRune(); err != nil {
				fmt.Printf("Failed to start the game: %v\n", err.Error())
				os.Exit(1)
			}
		})

		end := idx + NumberOfWordsDisplayed
		if end >= len(words) {
			end = len(words)
		}

		event := structs.Event{Word: expectedText}
		start := time.Now()
		fmt.Fprintf(writer, "%s %s\n", utils.PrintYellow(expectedText), strings.Join(words[idx+1:end], " "))
		for {
			text, stop := readWord(&event)
			if stop {
				return
			}

			// compare user input with expected text
			if text == expectedText {
				event.Duration = time.Since(start)
				*events = append(*events, event)
				*nbSuccess++
				break
			} else {
				if len(text) > 0 && text[0] != ' ' {
					fmt.Fprintf(writer, "%s\n", utils.PrintRed(text))
				}
				event.Attempts = append(event.Attempts, text)
				*nbError++
				if !gameConfig.InfiniteAttempts {
					event.Duration = time.Since(start)
					*events = append(*events, event)
					break
				}
			}

			// Safety circuit breaker
			if *nbError > consts.SafetyCircuitBreakerLimit {
				fmt.Println("An unknown issue occurred. The game has stopped to prevent any other problem")
				os.Exit(1)
			}
		}
	}
}

func finalStep(nbSuccess int, nbError int, events []structs.Event) {
	if err := keyboard.Close(); err != nil {
		fmt.Println(err)
	}
	printResults(nbSuccess, nbError, events)
	os.Exit(0)
}

func readWord(event *structs.Event) (string, bool) {
	writer := uilive.New()
	writer.Start()
	defer writer.Flush()
	word := ""
	// Remove +1 if you want to create a mode without having to use space or enter between words
	for len(word) < len(event.Word)+1 {
		if event.Word == word {
			fmt.Fprintf(writer, "%s\n", utils.PrintBlue(word))
		} else if strings.HasPrefix(event.Word, word) {
			fmt.Fprintf(writer, "%s\n", word)
		} else {
			fmt.Fprintf(writer, "%s\n", utils.PrintRed(word))
		}

		//fmt.Println(word)
		char, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Println(err)
			return "", true
		}
		//fmt.Printf("You pressed: rune %q, key %X. cursor %d\r\n", char, key, cursor)

		switch key {
		case keyboard.KeyEsc, keyboard.KeyCtrlC:
			fmt.Println("exiting...")
			return "", true
		case keyboard.KeyBackspace, keyboard.KeyBackspace2:
			if len(word) > 0 {
				r := []rune(word)
				word = string(r[:len(r)-1])
				event.Deletion++
			}
			continue
		case keyboard.KeySpace, keyboard.KeyEnter:
			return word, false
		}

		word += string(char)
	}
	return word, false
}

func generateWord(config structs.GameConfig) string {
	if config.NumberMode {
		number := strconv.FormatInt(rand.Int63(), 10)
		x := rand.Intn(consts.NumberModeMaxLength) + 1
		return number[len(number)-x:]
	}

	x := rand.Intn(config.WordSetMaxLength-config.WordSetMinLength) + config.WordSetMinLength
	words := utils.WordSet[x]
	y := rand.Intn(len(words))

	return transformWord(words[y], gameConfig)
}

func generateWords(config structs.GameConfig) []string {
	var words []string
	for i := 0; i < DefaultNumberOfWordsGenerated; i++ {
		words = append(words, generateWord(config))
	}
	return words
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

func printResults(nbSuccess int, nbError int, events []structs.Event) {
	var writer io.Writer = os.Stdout
	var renderer renderers.Renderer = &renderers.TextRenderer{
		Writer: writer,
	}

	if gameConfig.Output != DefaultOutput {
		writer, fileName := renderers.InitFile(gameConfig.Output)
		renderer = &renderers.MarkdownRenderer{
			Writer:   writer,
			FileName: fileName,
		}
	}
	renderer.RenderEmptyLines(1)
	if len(events) == 0 {
		fmt.Println("no entry")
		return
	}

	var percentageOfSuccess = float64(nbSuccess) / float64(nbSuccess+nbError) * 100
	totalDuration := time.Duration(0)
	totalDurationPerLetter := time.Duration(0)
	wordsResult := make([][]string, 0)
	deletions := 0
	charsToReview := make(structs.AnalysisMap)
	for _, event := range events {
		event.AnalyzeAttempts(&charsToReview)

		// Compute avg time by letter
		avgPerLetter := float64(event.Duration) / float64(len(event.Word))
		avgDuration := time.Duration(avgPerLetter)
		wordsResult = append(wordsResult, []string{event.Word, event.Duration.String(), avgDuration.String(), event.Attempts.String(), strconv.Itoa(event.Deletion)})

		totalDuration += event.Duration
		totalDurationPerLetter += avgDuration
		deletions += event.Deletion
	}
	avgDuration := totalDuration.Nanoseconds() / int64(len(events))
	totalAVGPerLetter := totalDurationPerLetter.Nanoseconds() / int64(len(events))
	dAVGDuration := time.Unix(0, avgDuration).Sub(time.Unix(0, 0))
	dTotalAVGPerLetter := time.Unix(0, totalAVGPerLetter).Sub(time.Unix(0, 0))
	totalTimeDuration := time.Unix(0, totalDuration.Nanoseconds()).Sub(time.Unix(0, 0))
	wpm := float64(len(events)) / totalDuration.Minutes()
	wordsResult = append(wordsResult, []string{" ", " ", " ", " ", " "})
	wordsResult = append(wordsResult, []string{"Total", totalTimeDuration.String(), "", "", strconv.Itoa(deletions)})
	wordsResult = append(wordsResult, []string{"Average", dAVGDuration.String(), dTotalAVGPerLetter.String(), "", fmt.Sprintf("%.2f", float64(deletions)/float64(len(events)))})
	wordsResult = append(wordsResult, []string{"WPM", fmt.Sprintf("%.2f", wpm), "", "", ""})

	resultData := [][]string{
		{"Success", strconv.Itoa(nbSuccess), fmt.Sprintf("%.2f%s", percentageOfSuccess, "%")},
		{"Error", strconv.Itoa(nbError), fmt.Sprintf("%.2f%s", 100-percentageOfSuccess, "%")},
		{"Total", strconv.Itoa(nbSuccess + nbError), "100.00%"},
	}

	renderer.RenderTitle(1, "Results")
	renderer.RenderTitle(2, "Main stats")
	renderer.RenderTable(resultData, []string{"", "Result", "%"})
	renderer.RenderEmptyLines(1)
	renderer.RenderTitle(2, "Words stats")
	renderer.RenderTable(wordsResult, []string{"Word", "Duration", "Duration/letter", "Errors", "Backspace"})
	renderer.RenderEmptyLines(1)

	if len(charsToReview) > 0 {
		charSwaps := charsToReview.ToCharSwaps()
		sort.Sort(charSwaps)
		var charSwapsResult [][]string
		for _, c := range charSwaps {
			charSwapsResult = append(charSwapsResult, c.ToStrings())
		}

		renderer.RenderTitle(2, "Chars errors")
		renderer.RenderTable(charSwapsResult, []string{"Expected char", "Input", "Number"})
		renderer.RenderEmptyLines(1)
	}
}
