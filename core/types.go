package core

import (
	"time"
)

type Log struct {
	Content string `json:"content"`
	Time time.Time `json:"time"`
}
