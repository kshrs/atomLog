package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/kshrs/atomLog/core"
)

func main() {
	var state core.AtomLogState
	state.RefreshDate()
	home, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err)
	}
	state.LogsDir = filepath.Join(home, "logs")
	state.FileName = filepath.Join(state.LogsDir,state.CurrentDate + ".json")
	state.Prompt = "logger: "

	err = state.MainLoop()
	if err != nil {
		log.Panic(err)
	}
}






