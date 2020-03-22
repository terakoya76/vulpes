package source

import (
	"strconv"
	"strings"

	"github.com/ziutek/mymysql/mysql"
)

// FetchGlobalStatus returns result
func FetchGlobalStatus(db mysql.Conn) (map[string]interface{}, error) {
	gStatus := make(map[string]interface{})
	rows, _, err := db.Query("SHOW GLOBAL STATUS")
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		varName := strings.ToLower(row.Str(0))
		if strings.HasPrefix(varName, "rsa_public_key") {
			continue
		}

		val := row.Str(1)
		num, err := strconv.ParseFloat(val, 64)
		if err != nil {
			gStatus[varName] = val
		} else {
			gStatus[varName] = num
		}
	}

	return gStatus, nil
}
