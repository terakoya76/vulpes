package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	SlaveStatusColCount = 2
)

// JSONizeSlaveStatus returns result.
func JSONizeSlaveStatus(str string) {
	gVar, err := ParseSlaveStatus(str)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	} else {
		res, err := json.Marshal(gVar)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		} else {
			fmt.Printf("%s\n", res)
		}
	}
}

// ParseSlaveStatus returns result.
func ParseSlaveStatus(str string) (map[string]interface{}, error) {
	sStatus := make(map[string]interface{})
	normStr := strings.TrimPrefix(str, "\n")

	// skip header
	lines := strings.Split(normStr, "\n")[1:]
	for _, line := range lines {
		row := strings.Split(line, ": ")

		if len(row) != SlaveStatusColCount {
			continue
		}

		varName := strings.ToLower(strings.TrimSpace(row[0]))
		val := row[1]
		num, err := strconv.ParseFloat(val, 64)

		if err != nil {
			sStatus[varName] = val
		} else {
			sStatus[varName] = num
		}
	}

	return sStatus, nil
}
