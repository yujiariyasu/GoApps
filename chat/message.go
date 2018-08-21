package main

import (
	"time"
)

type message struct {
	Name      string
	AvatarURL string
	Message   string
	When      time.Time
}
