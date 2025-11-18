package core

import (
	"bufio"
	"encoding/json"
	"os"
)

func (state *AtomLogState) ReadFile() error {
	var logs []Log
	file, err := os.OpenFile(state.FileName, os.O_CREATE|os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var log Log
		json.Unmarshal(scanner.Bytes() ,&log)
		logs = append(logs, log)
	}
	state.Logs = logs
	return nil
}

func (state *AtomLogState) WriteFile() error {
	file, err := os.OpenFile(state.FileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, log := range state.Logs {
		out, _ := json.Marshal(&log)
		file.WriteString(string(out))
		file.WriteString("\n")
	}
	return nil
}

func (state *AtomLogState) CreateLogsDir() error {
	if _, err := os.Stat(state.LogsDir); os.IsNotExist(err) {
		err := os.MkdirAll(state.LogsDir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
