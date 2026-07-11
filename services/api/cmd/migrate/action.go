package main

import "errors"

type Action string

const (
	Up    Action = "up"
	Down  Action = "down"
	Fresh Action = "fresh"
)

func ParseAction(action string) (Action, error) {
	switch action {
	case "up":
		return Up, nil
	case "down":
		return Down, nil
	case "fresh":
		return Fresh, nil
	default:
		return "", errors.New("invalid action")
	}
}
