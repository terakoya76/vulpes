package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	GlobalVariablesColCount = 2
)

// JSONizeGlobalVariables returns result
func JSONizeGlobalVariables(str string) {
	gVar, err := ParseGlobalVariables(str)
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

// ParseGlobalVariables returns result
func ParseGlobalVariables(str string) (map[string]interface{}, error) {
	gVar := make(map[string]interface{})

	normStr := strings.TrimPrefix(str, "\n")
	lines := strings.Split(normStr, "\n")
	for _, line := range lines {
		row := strings.Fields(line)

		if len(row) != GlobalVariablesColCount {
			continue
		}

		varName := strings.ToLower(row[0])
		if strings.HasPrefix(varName, "variable_name") {
			continue
		}

		val := row[1]
		num, err := strconv.ParseFloat(val, 64)
		if err != nil {
			gVar[varName] = val
		} else {
			gVar[varName] = num
		}
	}

	return gVar, nil
}
