package core

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	colors "github.com/kshrs/atomLog/ansi_colors"
)

var countOfLines int

func (state *AtomLogState) MainLoop() error {
	err := state.CreateLogsDir()
	if err != nil {
		return err
	}
	if state.Prompt == "" {
		state.Prompt = "AtomLog: "
	}

	err = state.ReadFile()
	if err != nil {
		fmt.Println("Error during read file")
		return err
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		state.RefreshDate()
		fmt.Print(colors.BgBrightBlack + colors.White, state.Prompt, colors.Reset)
		input, err := reader.ReadString('\n')
		if err != nil {
			return errors.New("Failed to read StdIn")
		}

		code := state.ParseLog(strings.Trim(input, " \n"))
		if code == "exit" {
			fmt.Println(colors.Red + "Exiting..." + colors.Reset)
			return nil
		}
		err = state.WriteFile()
		if err != nil {
			return err
		}
	}
}

func (state *AtomLogState) RefreshDate() {
	if time.Now().Format("02-01-2006") == state.CurrentDate {
		return
	}
	state.CurrentDate = time.Now().Format("02-01-2006")
}

func (state *AtomLogState) ParseLog(input string) string {
	var log Log

	switch strings.TrimSpace(strings.ToLower(input)) {

	case "exit", "q":
		state.ClearPreviousLine(1)
		return "exit"

	case "\\!", "\\imp", "\\important", "\\imps":
		state.FindAndPrintImportant()

	case "\\@", "\\mention", "\\mentions":
		state.FindAndPrintMentions()

	case "\\-", "\\li", "\\list", "\\lists":
		state.FindAndPrintLists()
	case "\\#", "\\head", "\\heads", "\\heading", "\\headings":
		state.FindAndPrintHeadings()

	case "":
		state.ClearPreviousLine(1)
		fmt.Println(colors.Magenta, time.Now().Format("15-04-05"), ">", colors.Reset)
		countOfLines += 1

	case "print":
		state.ClearPreviousLine(1)
		state.PrintLogs(0)

	default:
		log.Content = input
		log.Time = time.Now()
		state.Logs = append(state.Logs, log)
		state.ClearPreviousLine(1)
		state.PrettyPrintLog(state.Logs[len(state.Logs)-1])
	}
	return ""
}

func (state *AtomLogState) PrintLogs(count int) {
	defaultCount := 5
	if count == 0 {
		count = defaultCount
	}
	if count >= len(state.Logs) {
		count = len(state.Logs) - 1
	}
	for _, log := range state.Logs[len(state.Logs)-count-1:len(state.Logs)] {
		fmt.Println()
		fmt.Println("Content: ", log.Content)
		fmt.Println("Time: ", log.Time)
		countOfLines += 3
	}
}

func (state *AtomLogState) ClearPreviousLine(count int) {
	for _ = range count {
		fmt.Printf("\r\033[F\033[2K")
	}
}

func (state *AtomLogState) PrettyPrintLog(log Log) {
	fmt.Print(colors.Magenta, log.Time.Format("15:04:05"),"> ", colors.Reset)
	
	var coloredContent string = log.Content

	// @ Flag
	re := regexp.MustCompile(`@[A-Za-z0-9_]+`)
	coloredContent = re.ReplaceAllStringFunc(coloredContent, func(match string) string {
		return colors.Cyan + match + colors.Reset
	})
	// ! Flag
	re = regexp.MustCompile(`![A-Za-z0-9_]+`)
	coloredContent = re.ReplaceAllStringFunc(coloredContent, func(match string) string {
		return colors.Red + match + colors.Reset
	})
	// # Flag
	if strings.HasPrefix(coloredContent, "#") {
	coloredContent = colors.Bold + colors.Green + coloredContent + colors.Reset
	}
	// - Flag
	coloredContent = strings.ReplaceAll(coloredContent, "-", colors.Yellow + "-" + colors.Reset)

	// @date Flag
	coloredContent = strings.ReplaceAll(strings.ToLower(coloredContent), "@date", colors.Yellow + time.Now().Format("02-01-2006") + colors.Reset)

	// @time Flag
	coloredContent = strings.ReplaceAll(strings.ToLower(coloredContent), "@time", colors.Yellow + time.Now().Format("15-04-05") + colors.Reset)

	fmt.Println(coloredContent)
	countOfLines += 1
}

func (state *AtomLogState) ClearScreen() {
	state.ClearPreviousLine(countOfLines+1)
}

func (state *AtomLogState) FindAndPrintImportant() {
	// ! Flag
	re := regexp.MustCompile(`![A-Za-z0-9_]+`)
	for _, log := range state.Logs {
		matched := re.MatchString(log.Content)
		if matched {
			state.PrettyPrintLog(log)
		}
	}
}

func (state *AtomLogState) FindAndPrintMentions() {
	// @ Flag
	re := regexp.MustCompile(`@[A-Za-z0-9_]+`)
	for _, log := range state.Logs {
		matched := re.MatchString(log.Content)
		if matched {
			state.PrettyPrintLog(log)
		}
	}
}

func (state *AtomLogState) FindAndPrintLists() {
	// - Flag
	for _, log := range state.Logs {
		matched := strings.Contains(log.Content, "-")
		if matched {
			state.PrettyPrintLog(log)
		}
	}
}

func (state *AtomLogState) FindAndPrintHeadings() {
	// # Flag
	re := regexp.MustCompile(`#[A-Za-z0-9_]+`)
	for _, log := range state.Logs {
		matched := re.MatchString(log.Content)
		if matched {
			state.PrettyPrintLog(log)
		}
	}
}
