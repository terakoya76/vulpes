package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/terakoya76/vulpes/parser"
)

func main() {
	cobra.OnInitialize()
	rootCmd.DisableSuggestions = false
	rootCmd.AddCommand(globalStatusCmd, globalVariablesCmd, innodbStatusCmd, slaveStatusCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "vulpes",
	Short: "vulpes parse MySQL status output and output it as JSON",
	Long:  "",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var globalStatusCmd = &cobra.Command{
	Use:   "global_status",
	Short: "JSONize SHOW GLOBAL STATUS OUTPUT from stdin",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		stdin, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		} else {
			parser.JSONizeGlobalStatus(string(stdin))
		}
	},
}

var globalVariablesCmd = &cobra.Command{
	Use:   "global_variables",
	Short: "JSONize SHOW GLOBAL VARIABLES OUTPUT from stdin",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		stdin, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		} else {
			parser.JSONizeGlobalVariables(string(stdin))
		}
	},
}

var innodbStatusCmd = &cobra.Command{
	Use:   "innodb_status",
	Short: "JSONize SHOW INNODB STATUS OUTPUT from stdin",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		stdin, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		} else {
			parser.JSONizeInnodbStatus(string(stdin))
		}
	},
}

var slaveStatusCmd = &cobra.Command{
	Use:   "slave_status",
	Short: "JSONize SHOW SLAVE STATUS OUTPUT from stdin",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		stdin, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		} else {
			parser.JSONizeSlaveStatus(string(stdin))
		}
	},
}
