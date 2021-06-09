package main

import (
	"fmt"
	"os"

	"github.com/JulieWinters/population-health-designer/cmd/events"
	"github.com/JulieWinters/population-health-designer/cmd/maps"
	"github.com/JulieWinters/population-health-designer/cmd/stats"
)

const (
	STATS  = "stats"
	EVENTS = "events"
	MAPS   = "maps"
)

func main() {
	args := os.Args[1:]
	actions, err := GetArgMap(args)

	if err != nil {
		fmt.Println(err)
		return
	}

	var message string
	if val, ok := actions[STATS]; ok {
		message, err = stats.Execute(val)
	}
	if err != nil {
		fmt.Printf("Failed in Population Statistics step: %v\n", err)
		return
	}

	if val, ok := actions[EVENTS]; ok {
		message, err = events.Execute(val)
	}
	if err != nil {
		fmt.Printf("Failed in Event Generation step: %v\n", err)
		return
	}

	if val, ok := actions[MAPS]; ok {
		message, err = maps.Execute(val)
	}
	if err != nil {
		fmt.Printf("Failed in Event Data Mapping step: %v\n", err)
		return
	}

	fmt.Println(message)
}

func GetArgMap(args []string) (map[string]string, error) {
	var actions map[string]string = make(map[string]string)
	for i := 0; i < len(args); {
		if args[i] == "-s" || args[i] == "-stats" {
			actions[STATS] = args[i+1]
			i += 2
		} else if args[i] == "-e" || args[i] == "-events" {
			actions[EVENTS] = args[i+1]
			i += 2
		} else if args[i] == "-m" || args[i] == "-maps" {
			actions[MAPS] = args[i+1]
			i += 2
		} else if args[i] == "--" {
			// ignore this argument
			i += 1
		} else {
			return nil, fmt.Errorf("unrecognized argument '%v', please see help for more details", args[i])
		}
	}
	return actions, nil
}
