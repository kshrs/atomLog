package core

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"strings"
	"errors"
	colors "github.com/kshrs/atomLog/ansi_colors"
)

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
			fmt.Println("Exiting...")
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

	case "":
		state.ClearPreviousLine(1)
		fmt.Println(colors.Magenta, time.Now().Format("15-04-05"), ">", colors.Reset)

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
	}
}

func (state *AtomLogState) ClearPreviousLine(count int) {
	for _ = range count {
		fmt.Printf("\r\033[F\033[2K")
	}
}

func (state *AtomLogState) PrettyPrintLog(log Log) {
	fmt.Print(colors.Magenta, log.Time.Format("15:04:05"),"> ", colors.Reset)
	fmt.Println(log.Content)
}
