package source

import (
	"strconv"
	"strings"

	"github.com/ziutek/mymysql/mysql"
)

// FetchGlobalVariables returns result
func FetchGlobalVariables(db mysql.Conn) (map[string]interface{}, error) {
	gVar := make(map[string]interface{})
	rows, _, err := db.Query("SHOW GLOBAL VARIABLES")
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		varName := strings.ToLower(row.Str(0))
		val := row.Str(1)
		num, err := strconv.ParseFloat(val, 64)
		if err != nil {
			gVar[varName] = val
		} else {
			gVar[varName] = num
		}
	}

	return gVar, nil
}
