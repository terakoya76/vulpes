package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/terakoya76/vulpes/parser"
)

const (
	innodbStatus56 string = `
*************************** 1. row ***************************
  Type: InnoDB
  Name:
Status:
=====================================
2015-03-09 20:11:22 7f6c0c845700 INNODB MONITOR OUTPUT
=====================================
Per second averages calculated from the last 6 seconds
-----------------
BACKGROUND THREAD
-----------------
srv_master_thread loops: 178 srv_active, 0 srv_shutdown, 1244368 srv_idle
srv_master_thread log flush and writes: 1244546
----------
SEMAPHORES
----------
OS WAIT ARRAY INFO: reservation count 227
OS WAIT ARRAY INFO: signal count 220
Mutex spin waits 923, rounds 9442, OS waits 193
RW-shared spins 19, rounds 538, OS waits 16
RW-excl spins 5, rounds 476, OS waits 13
Spin rounds per wait: 10.23 mutex, 28.32 RW-shared, 95.20 RW-excl
------------
TRANSACTIONS
------------
Trx id counter 1093821584
Purge done for trx's n:o < 1093815563 undo n:o < 0 state: running but idle
History list length 649
LIST OF TRANSACTIONS FOR EACH SESSION:
---TRANSACTION 0, not started
MySQL thread id 27954, OS thread handle 0x7f6c0c845700, query id 90345 localhost root init
SHOW /*!50000 ENGINE*/ INNODB STATUS
---TRANSACTION 1093821554, not started
MySQL thread id 27893, OS thread handle 0x7f6c0c886700, query id 90144 127.0.0.1 cactiuser cleaning up
---TRANSACTION 1093821583, not started
MySQL thread id 27888, OS thread handle 0x7f6c0c8c7700, query id 90175 127.0.0.1 cactiuser cleaning up
---TRANSACTION 1093811214, not started
MySQL thread id 27887, OS thread handle 0x7f6c0c98a700, query id 80071 127.0.0.1 cactiuser cleaning up
---TRANSACTION 1093820819, not started
MySQL thread id 27886, OS thread handle 0x7f6c0c949700, query id 89403 127.0.0.1 cactiuser cleaning up
---TRANSACTION 1093811160, not started
MySQL thread id 27885, OS thread handle 0x7f6c0c908700, query id 80015 127.0.0.1 cactiuser cleaning up
--------
FILE I/O
--------
I/O thread 0 state: waiting for completed aio requests (insert buffer thread)
I/O thread 1 state: waiting for completed aio requests (log thread)
I/O thread 2 state: waiting for completed aio requests (read thread)
I/O thread 3 state: waiting for completed aio requests (read thread)
I/O thread 4 state: waiting for completed aio requests (read thread)
I/O thread 5 state: waiting for completed aio requests (read thread)
I/O thread 6 state: waiting for completed aio requests (write thread)
I/O thread 7 state: waiting for completed aio requests (write thread)
I/O thread 8 state: waiting for completed aio requests (write thread)
I/O thread 9 state: waiting for completed aio requests (write thread)
Pending normal aio reads: 0 [0, 0, 0, 0] , aio writes: 0 [0, 0, 0, 0] ,
 ibuf aio reads: 0, log i/o's: 0, sync i/o's: 0
Pending flushes (fsync) log: 0; buffer pool: 0
124669 OS file reads, 4457 OS file writes, 3498 OS fsyncs
0.00 reads/s, 0 avg bytes/read, 0.00 writes/s, 0.00 fsyncs/s
-------------------------------------
INSERT BUFFER AND ADAPTIVE HASH INDEX
-------------------------------------
Ibuf: size 1, free list len 63, seg size 65, 2 merges
merged operations:
 insert 48, delete mark 0, delete 0
discarded operations:
 insert 0, delete mark 0, delete 0
Hash table size 34679, node heap has 1 buffer(s)
0.00 hash searches/s, 0.00 non-hash searches/s
---
LOG
---
Log sequence number 53339891261
Log flushed up to   53339891261
Pages flushed up to 53339891261
Last checkpoint at  53339891261
0 pending log writes, 0 pending chkp writes
3395 log i/o's done, 0.00 log i/o's/second
----------------------
BUFFER POOL AND MEMORY
----------------------
Total memory allocated 17170432; in additional pool allocated 0
Dictionary memory allocated 318159
Buffer pool size   1024
Free buffers       755
Database pages     256
Old database pages 0
Modified db pages  0
Pending reads 0
Pending writes: LRU 0, flush list 0, single page 0
Pages made young 6, not young 751793
0.00 youngs/s, 0.00 non-youngs/s
Pages read 124617, created 40, written 1020
0.00 reads/s, 0.00 creates/s, 0.00 writes/s
No buffer pool page gets since the last printout
Pages read ahead 0.00/s, evicted without access 0.00/s, Random read ahead 0.00/s
LRU len: 256, unzip_LRU len: 0
I/O sum[0]:cur[0], unzip sum[0]:cur[0]
--------------
ROW OPERATIONS
--------------
0 queries inside InnoDB, 0 queries in queue
0 read views open inside InnoDB
Main thread process no. 1968, id 140101998331648, state: sleeping
Number of rows inserted 3089, updated 220, deleted 212, read 2099881
0.00 inserts/s, 0.00 updates/s, 0.00 deletes/s, 0.00 reads/s
----------------------------
END OF INNODB MONITOR OUTPUT
`

	// nolint:lll
	innodbStatus57 string = `
*************************** 1. row ***************************
  Type: InnoDB
  Name:
Status:
=====================================
2016-02-22 19:08:31 0x700000eda000 INNODB MONITOR OUTPUT
=====================================
Per second averages calculated from the last 4 seconds
-----------------
BACKGROUND THREAD
-----------------
srv_master_thread loops: 1 srv_active, 0 srv_shutdown, 2 srv_idle
srv_master_thread log flush and writes: 3
----------
SEMAPHORES
----------
OS WAIT ARRAY INFO: reservation count 63
OS WAIT ARRAY INFO: signal count 111
RW-shared spins 0, rounds 85, OS waits 22
RW-excl spins 0, rounds 4705, OS waits 17
RW-sx spins 0, rounds 0, OS waits 0
Spin rounds per wait: 85.00 RW-shared, 4705.00 RW-excl, 0.00 RW-sx
------------
TRANSACTIONS
------------
Trx id counter 507
Purge done for trx's n:o < 505 undo n:o < 0
History list length 1
LIST OF TRANSACTIONS FOR EACH SESSION:
---TRANSACTION 0, not started
MySQL thread id 8, OS thread handle 0x7efe7cb12700, query id 52 localhost root
SHOW ENGINE INNODB STATUS
---TRANSACTION 506, ACTIVE 804 sec starting index read
mysql tables in use 1, locked 1
LOCK WAIT 2 lock struct(s), heap size 376, 1 row lock(s)
MySQL thread id 3, OS thread handle 0x7efe7cb5b700, query id 47 localhost root statistics
SELECT * FROM test WHERE id = 1 LOCK IN SHARE MODE
------- TRX HAS BEEN WAITING 22 SEC FOR THIS LOCK TO BE GRANTED:
RECORD LOCKS space id 0 page no 307 n bits 72 index ` + "`PRIMARY` of table `test`.`test`" + ` trx id 506 lock mode S locks rec but not gap waiting
------------------
---TRANSACTION 505, ACTIVE 815 sec
2 lock struct(s), heap size 376, 1 row lock(s), undo log entries 1
MySQL thread id 2, OS thread handle 0x7efe7cba4700, query id 35 localhost root
---TRANSACTION 421551662098272, not started
0 lock struct(s), heap size 1136, 0 row lock(s)
--------
FILE I/O
--------
I/O thread 0 state: waiting for i/o request (insert buffer thread)
I/O thread 1 state: waiting for i/o request (log thread)
I/O thread 2 state: waiting for i/o request (read thread)
I/O thread 3 state: waiting for i/o request (read thread)
I/O thread 4 state: waiting for i/o request (read thread)
I/O thread 5 state: waiting for i/o request (read thread)
I/O thread 6 state: waiting for i/o request (write thread)
I/O thread 7 state: waiting for i/o request (write thread)
I/O thread 8 state: waiting for i/o request (write thread)
I/O thread 9 state: waiting for i/o request (write thread)
Pending normal aio reads: [0, 0, 0, 0] , aio writes: [0, 0, 0, 0] ,
 ibuf aio reads:, log i/o's:, sync i/o's:
Pending flushes (fsync) log: 0; buffer pool: 0
516 OS file reads, 55 OS file writes, 9 OS fsyncs
128.97 reads/s, 20393 avg bytes/read, 13.75 writes/s, 2.25 fsyncs/s
-------------------------------------
INSERT BUFFER AND ADAPTIVE HASH INDEX
-------------------------------------
Ibuf: size 1, free list len 0, seg size 2, 0 merges
merged operations:
 insert 0, delete mark 0, delete 0
discarded operations:
 insert 0, delete mark 0, delete 0
Hash table size 276671, node heap has 2 buffer(s)
Hash table size 276671, node heap has 0 buffer(s)
Hash table size 276671, node heap has 0 buffer(s)
Hash table size 276671, node heap has 0 buffer(s)
Hash table size 276671, node heap has 1 buffer(s)
Hash table size 276671, node heap has 1 buffer(s)
Hash table size 276671, node heap has 0 buffer(s)
Hash table size 276671, node heap has 4 buffer(s)
276.93 hash searches/s, 835.29 non-hash searches/s
---
LOG
---
Log sequence number 379575319
Log flushed up to   379575319
Pages flushed up to 379575319
Last checkpoint at  379575310
0 pending log flushes, 0 pending chkp writes
12 log i/o's done, 3.00 log i/o's/second
----------------------
BUFFER POOL AND MEMORY
----------------------
Total large memory allocated 1099431936
Dictionary memory allocated 312184
Buffer pool size   65528
Free buffers       64999
Database pages     521
Old database pages 0
Modified db pages  0
Pending reads 0
Pending writes: LRU 0, flush list 0, single page 0
Pages made young 0, not young 0
0.00 youngs/s, 0.00 non-youngs/s
Pages read 487, created 34, written 36
121.72 reads/s, 8.50 creates/s, 9.00 writes/s
Buffer pool hit rate 974 / 1000, young-making rate 0 / 1000 not 0 / 1000
Pages read ahead 0.00/s, evicted without access 0.00/s, Random read ahead 0.00/s
LRU len: 521, unzip_LRU len: 0
I/O sum[0]:cur[0], unzip sum[0]:cur[0]
----------------------
INDIVIDUAL BUFFER POOL INFO
----------------------
---BUFFER POOL 0
Buffer pool size   16382
Free buffers       16228
Database pages     152
Old database pages 0
Modified db pages  0
Pending reads 0
Pending writes: LRU 0, flush list 0, single page 0
Pages made young 0, not young 0
0.00 youngs/s, 0.00 non-youngs/s
Pages read 152, created 0, written 2
37.99 reads/s, 0.00 creates/s, 0.50 writes/s
Buffer pool hit rate 976 / 1000, young-making rate 0 / 1000 not 0 / 1000
Pages read ahead 0.00/s, evicted without access 0.00/s, Random read ahead 0.00/s
LRU len: 152, unzip_LRU len: 0
I/O sum[0]:cur[0], unzip sum[0]:cur[0]
---BUFFER POOL 1
Buffer pool size   16382
Free buffers       16244
Database pages     136
Old database pages 0
Modified db pages  0
Pending reads 0
Pending writes: LRU 0, flush list 0, single page 0
Pages made young 0, not young 0
0.00 youngs/s, 0.00 non-youngs/s
Pages read 136, created 0, written 0
33.99 reads/s, 0.00 creates/s, 0.00 writes/s
Buffer pool hit rate 978 / 1000, young-making rate 0 / 1000 not 0 / 1000
Pages read ahead 0.00/s, evicted without access 0.00/s, Random read ahead 0.00/s
LRU len: 136, unzip_LRU len: 0
I/O sum[0]:cur[0], unzip sum[0]:cur[0]
---BUFFER POOL 2
Buffer pool size   16382
Free buffers       16313
Database pages     67
Old database pages 0
Modified db pages  0
Pending reads 0
Pending writes: LRU 0, flush list 0, single page 0
Pages made young 0, not young 0
0.00 youngs/s, 0.00 non-youngs/s
Pages read 67, created 0, written 0
16.75 reads/s, 0.00 creates/s, 0.00 writes/s
Buffer pool hit rate 975 / 1000, young-making rate 0 / 1000 not 0 / 1000
Pages read ahead 0.00/s, evicted without access 0.00/s, Random read ahead 0.00/s
LRU len: 67, unzip_LRU len: 0
I/O sum[0]:cur[0], unzip sum[0]:cur[0]
---BUFFER POOL 3
Buffer pool size   16382
Free buffers       16214
Database pages     166
Old database pages 0
Modified db pages  0
Pending reads 0
Pending writes: LRU 0, flush list 0, single page 0
Pages made young 0, not young 0
0.00 youngs/s, 0.00 non-youngs/s
Pages read 132, created 34, written 34
32.99 reads/s, 8.50 creates/s, 8.50 writes/s
Buffer pool hit rate 963 / 1000, young-making rate 0 / 1000 not 0 / 1000
Pages read ahead 0.00/s, evicted without access 0.00/s, Random read ahead 0.00/s
LRU len: 166, unzip_LRU len: 0
I/O sum[0]:cur[0], unzip sum[0]:cur[0]
--------------
ROW OPERATIONS
--------------
0 queries inside InnoDB, 0 queries in queue
0 read views open inside InnoDB
Process ID=28837, Main thread ID=123145312497664, state: sleeping
Number of rows inserted 0, updated 0, deleted 0, read 8
0.00 inserts/s, 0.00 updates/s, 0.00 deletes/s, 2.00 reads/s
----------------------------
END OF INNODB MONITOR OUTPUT
============================
`
)

var (
	innodbStatusResult56 map[string]interface{} = map[string]interface{}{
		// background_thread
		"background_active_thread_loops":         178.0,
		"background_idle_thread_loops":           1.244368e+06,
		"background_shutdown_thread_loops":       0.0,
		"background_thread_log_flush_and_writes": 1.244546e+06,
		// buffer_pool_and_memory
		"buffer_pool_database_pages":                  256.0,
		"buffer_pool_dictionary_memory_allocated":     318159.0,
		"buffer_pool_free_buffers":                    755.0,
		"buffer_pool_io_cur":                          0.0,
		"buffer_pool_io_sum":                          0.0,
		"buffer_pool_io_unzip_cur":                    0.0,
		"buffer_pool_io_unzip_sum":                    0.0,
		"buffer_pool_lru_len":                         256.0,
		"buffer_pool_modified_db_pages":               0.0,
		"buffer_pool_old_database_pages":              0.0,
		"buffer_pool_page_evicted_wo_access_per_sec":  0.0,
		"buffer_pool_page_random_reads_ahead_per_sec": 0.0,
		"buffer_pool_page_reads_ahead_per_sec":        0.0,
		"buffer_pool_pages_created":                   40.0,
		"buffer_pool_pages_creates_per_sec":           0.0,
		"buffer_pool_pages_made_young":                6.0,
		"buffer_pool_pages_made_young_per_sec":        0.0,
		"buffer_pool_pages_not_made_young":            751793.0,
		"buffer_pool_pages_not_made_youngs_per_sec":   0.0,
		"buffer_pool_pages_read":                      124617.0,
		"buffer_pool_pages_reads_per_sec":             0.0,
		"buffer_pool_pages_writes_per_sec":            0.0,
		"buffer_pool_pages_written":                   1020.0,
		"buffer_pool_pending_reads":                   0.0,
		"buffer_pool_size":                            1024.0,
		"buffer_pool_unzip_lru_len":                   0.0,
		// file_io
		"insert_buffer_threads_num": 1,
		"log_threads_num":           1,
		"os_file_bytes_per_read":    0.0,
		"os_file_fsyncs":            3498.0,
		"os_file_fsyncs_per_sec":    0.0,
		"os_file_reads":             124669.0,
		"os_file_reads_per_sec":     0.0,
		"os_file_writes":            4457.0,
		"os_file_writes_per_sec":    0.0,
		"pending_buf_pool_flushes":  0.0,
		"pending_log_flushes":       0.0,
		"pending_normal_reads_aio":  0.0,
		"pending_normal_writes_aio": 0.0,
		"read_threads_num":          4,
		"write_threads_num":         4,
		// insert_buffer_and_adaptive_hash_index
		"hash_searches_per_sec":                0.0,
		"insert_buffer_discarded_delete_marks": 0.0,
		"insert_buffer_discarded_deletes":      0.0,
		"insert_buffer_discarded_inserts":      0.0,
		"insert_buffer_free_list_len":          63.0,
		"insert_buffer_merged_delete_marks":    0.0,
		"insert_buffer_merged_deletes":         0.0,
		"insert_buffer_merged_inserts":         48.0,
		"insert_buffer_merges":                 2.0,
		"insert_buffer_seg_size":               65.0,
		"insert_buffer_size":                   1.0,
		"non_hash_searches_per_sec":            0.0,
		// log
		"log_flushed_up_to":      5.3339891261e+10,
		"log_iops":               0.0,
		"log_ios_done":           3395.0,
		"log_last_checkpoint_at": 5.3339891261e+10,
		"log_sequence_number":    5.3339891261e+10,
		// row_operations
		"queries_in_queue":         0.0,
		"queries_inside_innodb":    0.0,
		"read_views_inside_innodb": 0.0,
		"rows_deleted":             212.0,
		"rows_deleted_per_sec":     0.0,
		"rows_inserted":            3089.0,
		"rows_inserted_per_sec":    0.0,
		"rows_read":                2.099881e+06,
		"rows_read_per_sec":        0.0,
		"rows_updated":             220.0,
		"rows_updated_per_sec":     0.0,
		// semaphores
		"os_waits_sync_array_reservations": 227.0,
		"os_waits_sync_array_signals":      220.0,
		"rw_excl_os_waits":                 13.0,
		"rw_excl_spin_rounds":              476.0,
		"rw_excl_spin_waits":               5.0,
		"rw_shared_os_waits":               16.0,
		"rw_shared_spin_rounds":            538.0,
		"rw_shared_spin_waits":             19.0,
		// start_of_innodb_monitor_output
		"calculated_from_last_secs": 6.0,
		// transactions
		"current_transactions":    6,
		"history_list_length":     649.0,
		"purge_done_for_trx_num":  1.093815563e+09,
		"purge_done_for_undo_num": 0.0,
		"trx_id_counter":          1.093821584e+09,
	}

	innodbStatusResult57 map[string]interface{} = map[string]interface{}{
		// background_thread
		"background_active_thread_loops":         1.0,
		"background_idle_thread_loops":           2.0,
		"background_shutdown_thread_loops":       0.0,
		"background_thread_log_flush_and_writes": 3.0,
		// buffer_pool_and_memory
		"buffer_pool_cache_hit_make_young_rate":       0.0,
		"buffer_pool_cache_hit_not_make_young_rate":   0.0,
		"buffer_pool_cache_hit_rate":                  0.974,
		"buffer_pool_dictionary_memory_allocated":     312184.0,
		"buffer_pool_free_buffers":                    64999.0,
		"buffer_pool_database_pages":                  521.0,
		"buffer_pool_io_cur":                          0.0,
		"buffer_pool_io_sum":                          0.0,
		"buffer_pool_io_unzip_cur":                    0.0,
		"buffer_pool_io_unzip_sum":                    0.0,
		"buffer_pool_lru_len":                         521.0,
		"buffer_pool_modified_db_pages":               0.0,
		"buffer_pool_old_database_pages":              0.0,
		"buffer_pool_page_evicted_wo_access_per_sec":  0.0,
		"buffer_pool_page_random_reads_ahead_per_sec": 0.0,
		"buffer_pool_page_reads_ahead_per_sec":        0.0,
		"buffer_pool_pages_created":                   34.0,
		"buffer_pool_pages_creates_per_sec":           8.5,
		"buffer_pool_pages_made_young":                0.0,
		"buffer_pool_pages_made_young_per_sec":        0.0,
		"buffer_pool_pages_not_made_young":            0.0,
		"buffer_pool_pages_not_made_youngs_per_sec":   0.0,
		"buffer_pool_pages_read":                      487.0,
		"buffer_pool_pages_reads_per_sec":             121.72,
		"buffer_pool_pages_written":                   36.0,
		"buffer_pool_pages_writes_per_sec":            9.0,
		"buffer_pool_pending_reads":                   0.0,
		"buffer_pool_size":                            65528.0,
		"buffer_pool_total_large_memory_allocated":    1.099431936e+09,
		"buffer_pool_unzip_lru_len":                   0.0,
		// file_io
		"insert_buffer_threads_num": 1,
		"log_threads_num":           1,
		"os_file_bytes_per_read":    20393.0,
		"os_file_fsyncs":            9.0,
		"os_file_fsyncs_per_sec":    2.25,
		"os_file_reads":             516.0,
		"os_file_reads_per_sec":     128.97,
		"os_file_writes":            55.0,
		"os_file_writes_per_sec":    13.75,
		"pending_buf_pool_flushes":  0.0,
		"pending_log_flushes":       0.0,
		"read_threads_num":          4,
		"write_threads_num":         4,
		// insert_buffer_and_adaptive_hash_index
		"hash_searches_per_sec":                276.93,
		"insert_buffer_discarded_delete_marks": 0.0,
		"insert_buffer_discarded_deletes":      0.0,
		"insert_buffer_discarded_inserts":      0.0,
		"insert_buffer_free_list_len":          0.0,
		"insert_buffer_merged_delete_marks":    0.0,
		"insert_buffer_merged_deletes":         0.0,
		"insert_buffer_merged_inserts":         0.0,
		"insert_buffer_merges":                 0.0,
		"insert_buffer_seg_size":               2.0,
		"insert_buffer_size":                   1.0,
		"non_hash_searches_per_sec":            835.29,
		// log
		"log_flushed_up_to":             3.79575319e+08,
		"log_iops":                      3.0,
		"log_ios_done":                  12.0,
		"log_last_checkpoint_at":        3.7957531e+08,
		"log_pending_checkpoint_writes": 0.0,
		"log_pending_log_flushes":       0.0,
		"log_sequence_number":           3.79575319e+08,
		// row_operations
		"queries_in_queue":         0.0,
		"queries_inside_innodb":    0.0,
		"read_views_inside_innodb": 0.0,
		"rows_deleted":             0.0,
		"rows_deleted_per_sec":     0.0,
		"rows_inserted":            0.0,
		"rows_inserted_per_sec":    0.0,
		"rows_read":                8.0,
		"rows_read_per_sec":        2.0,
		"rows_updated":             0.0,
		"rows_updated_per_sec":     0.0,
		// semaphores
		"os_waits_sync_array_reservations": 63.0,
		"os_waits_sync_array_signals":      111.0,
		"rw_excl_os_waits":                 17.0,
		"rw_excl_spin_rounds":              4705.0,
		"rw_excl_spin_waits":               0.0,
		"rw_shared_os_waits":               22.0,
		"rw_shared_spin_rounds":            85.0,
		"rw_shared_spin_waits":             0.0,
		"rw_sx_os_waits":                   0.0,
		"rw_sx_spin_rounds":                0.0,
		"rw_sx_spin_waits":                 0.0,
		// start_of_innodb_monitor_output
		"calculated_from_last_secs": 4.0,
		// transactions
		"active_transactions":     2,
		"current_transactions":    4,
		"history_list_length":     1.0,
		"locked_transactions":     1,
		"purge_done_for_trx_num":  505.0,
		"purge_done_for_undo_num": 0.0,
		"trx_id_counter":          507.0,
	}
)

func TestParseInnodbStatus(t *testing.T) {
	cases := []struct {
		name     string
		str      string
		expected map[string]interface{}
		err      error
	}{
		{
			name:     "mysql5.6",
			str:      innodbStatus56,
			expected: innodbStatusResult56,
			err:      nil,
		},
		{
			name:     "mysql5.7",
			str:      innodbStatus57,
			expected: innodbStatusResult57,
			err:      nil,
		},
	}

	for _, c := range cases {
		actual, err := parser.ParseInnodbStatus(c.str)
		if err != nil {
			t.Errorf("err: %s\n", err.Error())
		}

		if !assert.Equal(t, c.expected, actual) {
			t.Errorf("case: %s is failed, expected: %+v, actual: %+v\n", c.name, c.expected, actual)
		}
	}
}
