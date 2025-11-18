package core

import (
	"time"
)

type Log struct {
	Content string `json:"content"`
	Time time.Time `json:"time"`
}

type AtomLogState struct {
	CurrentDate string
	LogsDir string
	FileName string
	Logs []Log
	Prompt string
}

