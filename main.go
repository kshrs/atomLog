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
	"path/filepath"

	"github.com/kshrs/atomLog/core"
	"github.com/kshrs/atomLog/ansi_colors"
)

var currentDate string
var fileName string
var logs []core.Log
var logsDir string

func RefreshDate() {
	if time.Now().Format("02-01-2006") == currentDate {
		return
	}
	currentDate = time.Now().Format("02-01-2006")
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err)
	}
	logsDir = filepath.Join(home, "logs")
	CreateLogsDir()

	currentDate = time.Now().Format("02-01-2006")
	RefreshDate()

	fileName = filepath.Join(logsDir,currentDate + ".json")

	fileLogs, err := ReadFile(fileName)
	logs = fileLogs
	if err != nil {
		fmt.Println("Error: ", err)
	}

	err = MainLoop()
	if err != nil {
		log.Panic(err)
	}
	

	// fmt.Println("Saved the logs to the log file")
	// WriteFile(fileName, logs)

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
	    WriteFile(fileName, logs)
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
		ClearPreviousLine(1)
		PrettyPrintLog(logs[len(logs)-1])
	}
	return ""
}

func ReadFile(fileName string) ([]core.Log, error) {
	var logs []core.Log
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDONLY, 0755)
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

func PrettyPrintLog(log core.Log) {
	fmt.Print(ansi_colors.Magenta, log.Time.Format("15:04:05"),"> ", ansi_colors.Reset)
	fmt.Println(log.Content)
}

func ClearPreviousLine(count int) {
	for _ = range count {
		fmt.Printf("\r\033[F\033[2K")
	}
}

func CreateLogsDir() {
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		err := os.MkdirAll(logsDir, 0755)
		if err != nil {
			log.Panic(err)
		}
	}
}
