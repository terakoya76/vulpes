package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/terakoya76/vulpes/parser"
)

const (
	subordinateStatus string = `
*************************** 1. row ***************************
               Subordinate_IO_State: Waiting for main to send event
                  Main_Host: db
                  Main_User: mysql
                  Main_Port: 3306
                Connect_Retry: 60
              Main_Log_File: bin-log.000002
          Read_Main_Log_Pos: 154
               Relay_Log_File: 29dc768b9532-relay-bin.000004
                Relay_Log_Pos: 318
        Relay_Main_Log_File: bin-log.000002
             Subordinate_IO_Running: Yes
            Subordinate_SQL_Running: Yes
              Replicate_Do_DB:
          Replicate_Ignore_DB:
           Replicate_Do_Table:
       Replicate_Ignore_Table:
      Replicate_Wild_Do_Table:
  Replicate_Wild_Ignore_Table:
                   Last_Errno: 0
                   Last_Error:
                 Skip_Counter: 0
          Exec_Main_Log_Pos: 154
              Relay_Log_Space: 532
              Until_Condition: None
               Until_Log_File:
                Until_Log_Pos: 0
           Main_SSL_Allowed: No
           Main_SSL_CA_File:
           Main_SSL_CA_Path:
              Main_SSL_Cert:
            Main_SSL_Cipher:
               Main_SSL_Key:
        Seconds_Behind_Main: 0
Main_SSL_Verify_Server_Cert: No
                Last_IO_Errno: 0
                Last_IO_Error:
               Last_SQL_Errno: 0
               Last_SQL_Error:
  Replicate_Ignore_Server_Ids:
             Main_Server_Id: 1
                  Main_UUID: f8b8d2e3-7c1c-11ea-a3c4-0242ac120002
             Main_Info_File: /var/lib/mysql/main.info
                    SQL_Delay: 0
          SQL_Remaining_Delay: NULL
      Subordinate_SQL_Running_State: Subordinate has read all relay log; waiting for more updates
           Main_Retry_Count: 86400
                  Main_Bind:
      Last_IO_Error_Timestamp:
     Last_SQL_Error_Timestamp:
               Main_SSL_Crl:
           Main_SSL_Crlpath:
           Retrieved_Gtid_Set:
            Executed_Gtid_Set:
                Auto_Position: 0
         Replicate_Rewrite_DB:
                 Channel_Name:
           Main_TLS_Version:
`
)

var (
	subordinateStatusResult map[string]interface{} = map[string]interface{}{
		"auto_position":                 0.0,
		"connect_retry":                 60.0,
		"exec_main_log_pos":           154.0,
		"last_errno":                    0.0,
		"last_io_errno":                 0.0,
		"last_sql_errno":                0.0,
		"main_host":                   "db",
		"main_info_file":              "/var/lib/mysql/main.info",
		"main_log_file":               "bin-log.000002",
		"main_port":                   3306.0,
		"main_retry_count":            86400.0,
		"main_server_id":              1.0,
		"main_ssl_allowed":            "No",
		"main_ssl_verify_server_cert": "No",
		"main_user":                   "mysql",
		"main_uuid":                   "f8b8d2e3-7c1c-11ea-a3c4-0242ac120002",
		"read_main_log_pos":           154.0,
		"relay_log_file":                "29dc768b9532-relay-bin.000004",
		"relay_log_pos":                 318.0,
		"relay_log_space":               532.0,
		"relay_main_log_file":         "bin-log.000002",
		"seconds_behind_main":         0.0,
		"skip_counter":                  0.0,
		"subordinate_io_running":              "Yes",
		"subordinate_io_state":                "Waiting for main to send event",
		"subordinate_sql_running":             "Yes",
		"subordinate_sql_running_state":       "Subordinate has read all relay log; waiting for more updates",
		"sql_delay":                     0.0,
		"sql_remaining_delay":           "NULL",
		"until_condition":               "None",
		"until_log_pos":                 0.0,
	}
)

func TestParseSubordinateStatus(t *testing.T) {
	cases := []struct {
		name     string
		str      string
		expected map[string]interface{}
		err      error
	}{
		{
			name:     "subordinate status",
			str:      subordinateStatus,
			expected: subordinateStatusResult,
			err:      nil,
		},
	}

	for _, c := range cases {
		actual, err := parser.ParseSubordinateStatus(c.str)
		if err != nil {
			t.Errorf("err: %s\n", err.Error())
		}
		if !assert.Equal(t, c.expected, actual) {
			t.Errorf("case: %s is failed, expected: %+v, actual: %+v\n", c.name, c.expected, actual)
		}
	}
}
