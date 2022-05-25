package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/terakoya76/vulpes/parser"
)

const (
	slaveStatus string = `
*************************** 1. row ***************************
               Slave_IO_State: Waiting for master to send event
                  Master_Host: db
                  Master_User: mysql
                  Master_Port: 3306
                Connect_Retry: 60
              Master_Log_File: bin-log.000002
          Read_Master_Log_Pos: 154
               Relay_Log_File: 29dc768b9532-relay-bin.000004
                Relay_Log_Pos: 318
        Relay_Master_Log_File: bin-log.000002
             Slave_IO_Running: Yes
            Slave_SQL_Running: Yes
              Replicate_Do_DB:
          Replicate_Ignore_DB:
           Replicate_Do_Table:
       Replicate_Ignore_Table:
      Replicate_Wild_Do_Table:
  Replicate_Wild_Ignore_Table:
                   Last_Errno: 0
                   Last_Error:
                 Skip_Counter: 0
          Exec_Master_Log_Pos: 154
              Relay_Log_Space: 532
              Until_Condition: None
               Until_Log_File:
                Until_Log_Pos: 0
           Master_SSL_Allowed: No
           Master_SSL_CA_File:
           Master_SSL_CA_Path:
              Master_SSL_Cert:
            Master_SSL_Cipher:
               Master_SSL_Key:
        Seconds_Behind_Master: 0
Master_SSL_Verify_Server_Cert: No
                Last_IO_Errno: 0
                Last_IO_Error:
               Last_SQL_Errno: 0
               Last_SQL_Error:
  Replicate_Ignore_Server_Ids:
             Master_Server_Id: 1
                  Master_UUID: f8b8d2e3-7c1c-11ea-a3c4-0242ac120002
             Master_Info_File: /var/lib/mysql/master.info
                    SQL_Delay: 0
          SQL_Remaining_Delay: NULL
      Slave_SQL_Running_State: Slave has read all relay log; waiting for more updates
           Master_Retry_Count: 86400
                  Master_Bind:
      Last_IO_Error_Timestamp:
     Last_SQL_Error_Timestamp:
               Master_SSL_Crl:
           Master_SSL_Crlpath:
           Retrieved_Gtid_Set:
            Executed_Gtid_Set:
                Auto_Position: 0
         Replicate_Rewrite_DB:
                 Channel_Name:
           Master_TLS_Version:
`
)

var (
	slaveStatusResult = map[string]interface{}{
		"auto_position":                 0.0,
		"connect_retry":                 60.0,
		"exec_master_log_pos":           154.0,
		"last_errno":                    0.0,
		"last_io_errno":                 0.0,
		"last_sql_errno":                0.0,
		"master_host":                   "db",
		"master_info_file":              "/var/lib/mysql/master.info",
		"master_log_file":               "bin-log.000002",
		"master_port":                   3306.0,
		"master_retry_count":            86400.0,
		"master_server_id":              1.0,
		"master_ssl_allowed":            "No",
		"master_ssl_verify_server_cert": "No",
		"master_user":                   "mysql",
		"master_uuid":                   "f8b8d2e3-7c1c-11ea-a3c4-0242ac120002",
		"read_master_log_pos":           154.0,
		"relay_log_file":                "29dc768b9532-relay-bin.000004",
		"relay_log_pos":                 318.0,
		"relay_log_space":               532.0,
		"relay_master_log_file":         "bin-log.000002",
		"seconds_behind_master":         0.0,
		"skip_counter":                  0.0,
		"slave_io_running":              "Yes",
		"slave_io_state":                "Waiting for master to send event",
		"slave_sql_running":             "Yes",
		"slave_sql_running_state":       "Slave has read all relay log; waiting for more updates",
		"sql_delay":                     0.0,
		"sql_remaining_delay":           "NULL",
		"until_condition":               "None",
		"until_log_pos":                 0.0,
	}
)

// nolint:dupl
func TestParseSlaveStatus(t *testing.T) {
	cases := []struct {
		name     string
		str      string
		expected map[string]interface{}
		err      error
	}{
		{
			name:     "slave status",
			str:      slaveStatus,
			expected: slaveStatusResult,
			err:      nil,
		},
	}

	for _, c := range cases {
		actual, err := parser.ParseSlaveStatus(c.str)
		if err != nil {
			t.Errorf("err: %s\n", err.Error())
		}

		if !assert.Equal(t, c.expected, actual) {
			t.Errorf("case: %s is failed, expected: %+v, actual: %+v\n", c.name, c.expected, actual)
		}
	}
}
