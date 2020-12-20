package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	GlobalStatusColCount = 2
)

// JSONizeGlobalStatus returns result
func JSONizeGlobalStatus(str string) {
	gStatus, err := ParseGlobalStatus(str)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	} else {
		res, err := json.Marshal(gStatus)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		} else {
			fmt.Printf("%s\n", res)
		}
	}
}

// ParseGlobalStatus returns result
func ParseGlobalStatus(str string) (map[string]interface{}, error) {
	gStatus := make(map[string]interface{})

	normStr := strings.TrimPrefix(str, "\n")
	lines := strings.Split(normStr, "\n")
	for _, line := range lines {
		row := strings.Fields(line)

		if len(row) != GlobalStatusColCount {
			continue
		}

		varName := strings.ToLower(row[0])
		if strings.HasPrefix(varName, "variable_name") {
			continue
		}
		if strings.HasPrefix(varName, "rsa_public_key") {
			continue
		}

		val := row[1]
		num, err := strconv.ParseFloat(val, 64)
		if err != nil {
			gStatus[varName] = val
		} else {
			gStatus[varName] = num
		}
	}

	return gStatus, nil
}
