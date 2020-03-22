package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"

	"github.com/terakoya76/vulpes/config"
	"github.com/terakoya76/vulpes/source"
)

func main() {
	var dbConf config.Database
	if err := envconfig.Process("db", &dbConf); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}

	db := mysql.New("tcp", "", fmt.Sprintf("%s:%d", dbConf.Hostname, dbConf.Port), dbConf.Username, dbConf.Password, "")
	if err := db.Connect(); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	defer db.Close()

	gStatus, err := source.FetchGlobalStatus(db)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}

	res, err := json.Marshal(gStatus)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	fmt.Printf("%s\n", res)

	gVar, err := source.FetchGlobalVariables(db)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}

	res, err = json.Marshal(gVar)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	fmt.Printf("%s\n", res)

	iStatus, err := source.FetchInnodbStatus(db)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}

	res, err = json.Marshal(iStatus)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	fmt.Printf("%s\n", res)
}
