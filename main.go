package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kshrs/atomLog/core"
	"github.com/kshrs/atomLog/ansi_colors"
)

var currentDate string
var fileName string
var logs []core.Log

func RefreshDate() {
	if time.Now().Format("02-01-2006") == currentDate {
		return
	}
	currentDate = time.Now().Format("02-01-2006")
}

func main() {
	currentDate = time.Now().Format("02-01-2006")
	RefreshDate()

	fileName = currentDate + ".json"

	fileLogs, err := ReadFile(fileName)
	logs = fileLogs
	if err != nil {
		fmt.Println("Error: ", err)
	}

	err = MainLoop()
	if err != nil {
		log.Fatal(err)
	}
	

	fmt.Println("Saved the logs to the log file")
	WriteFile(fileName, logs)

}

func MainLoop() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		RefreshDate()
		fmt.Print(ansi_colors.BgBrightBlack + ansi_colors.White, "logger: ", ansi_colors.Reset)
		input, err := reader.ReadString('\n')
		if err != nil {
			return errors.New("Failed to read StdIn")
		}

		code := ParseLog(strings.Trim(input, " \n"))
		if code == "exit" {
			fmt.Println("Exiting...")
			return nil
		}
	}
}

func ParseLog(input string) string {
	var log core.Log

	switch strings.ToLower(input) {

	case "exit", ":q":
		return "exit"

	case "":
		fmt.Println("Null")

	case "print":
		PrintLogs(logs)

	default:
		log.Content = input
		log.Time = time.Now()
		logs = append(logs, log)
	}
	return ""
}

func ReadFile(fileName string) ([]core.Log, error) {
	var logs []core.Log
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var log core.Log
		json.Unmarshal(scanner.Bytes() ,&log)
		logs = append(logs, log)
	}
	return logs, nil

}

func WriteFile(fileName string, logs []core.Log) (error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, log := range logs {
		out, _ := json.Marshal(&log)
		file.WriteString(string(out))
		file.WriteString("\n")
	}
	return nil
}

func PrintLogs(logs []core.Log) {
	for _, log := range logs {
		fmt.Println()
		fmt.Println("Content: ", log.Content)
		fmt.Println("Time: ", log.Time)
	}
}
