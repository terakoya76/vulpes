package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/imdario/mergo"
)

// JSONizeInnodbStatus returns result.
func JSONizeInnodbStatus(str string) {
	iStatus, err := ParseInnodbStatus(str)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	} else {
		res, err := json.Marshal(iStatus)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		} else {
			fmt.Printf("%s\n", res)
		}
	}
}

// ParseInnodbStatus returns result.
func ParseInnodbStatus(str string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	lines := strings.Split(str, "\n")

	p := InnodbStatusParser{
		source:      lines,
		content:     []string{lines[0]},
		currPos:     0,
		currSection: StartOfInnodbMonitorOutput,
	}

	for p.currPos < (len(p.source)-1) && p.currSection != EndOfInnodbMonitorOutput {
		p.scanToNextSection()
		metrics := p.parseContent(p.currSection)

		if err := mergo.Merge(&result, metrics); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		}

		p.stepNextSection()
	}

	return result, nil
}

type Section int

const (
	StartOfInnodbMonitorOutput Section = iota + 1
	BackgroundThread
	Semaphores
	Deadlocks
	ForeignKeyError
	Transactions
	FileIO
	InsertBufferAndAdaptiveHashIndex
	Log
	BufferPoolAndMemory
	IndividualBufferPoolInfo
	RowOperations
	EndOfInnodbMonitorOutput
)

func (s Section) String() string {
	switch s {
	case StartOfInnodbMonitorOutput:
		return ""
	case BackgroundThread:
		return "BACKGROUND THREAD"
	case Semaphores:
		return "SEMAPHORES"
	case Deadlocks:
		return "LATEST DETECTED DEADLOCK"
	case ForeignKeyError:
		return "LATEST FOREIGN KEY ERROR"
	case Transactions:
		return "TRANSACTIONS"
	case FileIO:
		return "FILE I/O"
	case InsertBufferAndAdaptiveHashIndex:
		return "INSERT BUFFER AND ADAPTIVE HASH INDEX"
	case Log:
		return "LOG"
	case BufferPoolAndMemory:
		return "BUFFER POOL AND MEMORY"
	case IndividualBufferPoolInfo:
		return "INDIVIDUAL BUFFER POOL INFO"
	case RowOperations:
		return "ROW OPERATIONS"
	case EndOfInnodbMonitorOutput:
		return "END OF INNODB MONITOR OUTPUT"
	default:
		return ""
	}
}

func isSectionLabel(sectionLabel string) bool {
	for i := StartOfInnodbMonitorOutput; i <= EndOfInnodbMonitorOutput; i++ {
		if sectionLabel == i.String() {
			return true
		}
	}

	return false
}

type InnodbStatusParser struct {
	source      []string
	content     []string
	currPos     int
	currSection Section
}

func (p *InnodbStatusParser) scanToNextSection() {
	for !p.checkNextSection() {
		p.next()
	}
}

/*
 * Sections are divided by the blocks like below.
 * ------------
 * xxx yyyy zzz.
 * ------------
 *.
 */
func (p *InnodbStatusParser) checkNextSection() bool {
	return isSectionLabel(p.maybeLabel())
}

func (p *InnodbStatusParser) maybeLabel() string {
	return p.source[p.currPos+2]
}

func (p *InnodbStatusParser) next() {
	p.currPos++
	p.content = append(p.content, p.source[p.currPos])
}

func (p *InnodbStatusParser) parseContent(s Section) map[string]interface{} {
	switch s {
	case StartOfInnodbMonitorOutput:
		return parseStartOfInnodbMonitorOutput(p.content)
	case BackgroundThread:
		return parseBackgroundThreadContent(p.content)
	case Semaphores:
		return parseSemaphoresContent(p.content)
	case Deadlocks:
		return parseDeadlocksContent(p.content)
	case ForeignKeyError:
		return parseForeignKeyErrorContent(p.content)
	case Transactions:
		return parseTransactionsContent(p.content)
	case FileIO:
		return parseFileIoContent(p.content)
	case InsertBufferAndAdaptiveHashIndex:
		return parseInsertBufferAndAdaptiveHashIndexContent(p.content)
	case Log:
		return parseLogContent(p.content)
	case BufferPoolAndMemory:
		return parseBufferPoolAndMemoryContent(p.content)
	case IndividualBufferPoolInfo:
		return parseIndividualBufferPoolInfoContent(p.content)
	case RowOperations:
		return parseRowOperationsContent(p.content)
	case EndOfInnodbMonitorOutput:
		return make(map[string]interface{})
	default:
		return make(map[string]interface{})
	}
}

func fillMetric(key string, value interface{}, result map[string]interface{}) map[string]interface{} {
	switch v := value.(type) {
	case string:
		num, err := strconv.ParseFloat(v, 64)
		if err != nil {
			result[key] = v
		} else {
			result[key] = num
		}
	default:
		result[key] = v
	}

	return result
}

func countUpMetric(key string, increment int, result map[string]interface{}) map[string]interface{} {
	if val, ok := result[key]; ok {
		if i, ok := val.(int); ok {
			result[key] = i + increment
		}
	} else {
		result[key] = increment
	}

	return result
}

func calcRate(a, b string) (float64, error) {
	aFloat, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return 0, err
	}

	bFloat, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return 0, err
	}

	return (aFloat / bFloat), nil
}

func parseStartOfInnodbMonitorOutput(content []string) map[string]interface{} {
	result := make(map[string]interface{})

	for _, line := range content {
		item := "Per second averages calculated from the last "
		if strings.HasPrefix(line, item) {
			/*
			 * The decent period of time that the per second values are sampled over
			 */
			record := strings.Fields(line)
			result = fillMetric("calculated_from_last_secs", record[7], result)
		}
	}

	return result
}

func parseBackgroundThreadContent(content []string) map[string]interface{} {
	result := make(map[string]interface{})

	for _, line := range content {
		item := "srv_master_thread loops: "
		if strings.HasPrefix(line, item) {
			parts := strings.Split(line, ", ")
			active, shutdown, idle := parts[0], parts[1], parts[2]

			/*
			 * The total number of background active thread loops
			 */
			metricActive := strings.Fields(active)
			result = fillMetric("background_active_thread_loops", metricActive[2], result)

			/*
			 * The total number of background shutdown thread loops
			 */
			metricShutdown := strings.Fields(shutdown)
			result = fillMetric("background_shutdown_thread_loops", metricShutdown[0], result)

			/*
			 * The total number of background idle thread loops
			 */
			metricIdle := strings.Fields(idle)
			result = fillMetric("background_idle_thread_loops", metricIdle[0], result)
		}

		item = "srv_master_thread log flush and writes: "
		if strings.HasPrefix(line, item) {
			/*
			 * The total number of background thread's log flush/write operations
			 */
			record := strings.Fields(line)
			result = fillMetric("background_thread_log_flush_and_writes", record[5], result)
		}
	}

	return result
}

func parseSemaphoresContent(content []string) map[string]interface{} {
	result := make(map[string]interface{})

	for _, line := range content {
		line = strings.ReplaceAll(line, "OS waits", "os_waits")

		item := "OS WAIT ARRAY INFO: reservation count "
		if strings.HasPrefix(line, item) {
			/*
			 * The total number of securing cells from the sync_array to wait for mutex acquisition
			 */
			record := strings.Fields(line)
			result = fillMetric("os_waits_sync_array_reservations", record[6], result)
		}

		item = "OS WAIT ARRAY INFO: signal count "
		if strings.HasPrefix(line, item) {
			/*
			 * The total number of unlocking mutex which threads held
			 */
			record := strings.Fields(line)
			result = fillMetric("os_waits_sync_array_signals", record[6], result)
		}

		item = "RW-shared "
		if strings.HasPrefix(line, item) {
			parts := strings.Split(line, ", ")
			spin, rounds, osWaits := parts[0], parts[1], parts[2]

			/*
			 * The total number of #mutex_spin_wait() call for RW-Shared lock
			 */
			metricSpin := strings.Fields(spin)
			result = fillMetric("rw_shared_spin_waits", metricSpin[2], result)

			/*
			 * The total number of spin loops in #mutex_spin_wait() for RW-Shared lock
			 */
			metricRounds := strings.Fields(rounds)
			result = fillMetric("rw_shared_spin_rounds", metricRounds[1], result)

			/*
			 * The total number of the sync_array event_wait in #mutex_spin_wait() for RW-Shared lock
			 */
			metricOsWaits := strings.Fields(osWaits)
			result = fillMetric("rw_shared_os_waits", metricOsWaits[1], result)
		}

		item = "RW-excl "
		if strings.HasPrefix(line, item) {
			parts := strings.Split(line, ", ")
			spin, rounds, osWaits := parts[0], parts[1], parts[2]

			/*
			 * The total number of #mutex_spin_wait() call for RW-Exclusive lock
			 */
			metricSpin := strings.Fields(spin)
			result = fillMetric("rw_excl_spin_waits", metricSpin[2], result)

			/*
			 * The total number of spin loops in #mutex_spin_wait() for RW-Exclusive lock
			 */
			metricRounds := strings.Fields(rounds)
			result = fillMetric("rw_excl_spin_rounds", metricRounds[1], result)

			/*
			 * The total number of the sync_array event_wait in #mutex_spin_wait() for RW-Exclusive lock
			 */
			metricOsWaits := strings.Fields(osWaits)
			result = fillMetric("rw_excl_os_waits", metricOsWaits[1], result)
		}

		item = "RW-sx "
		if strings.HasPrefix(line, item) {
			parts := strings.Split(line, ", ")
			spin, rounds, osWaits := parts[0], parts[1], parts[2]

			/*
			 * The total number of #mutex_spin_wait() call for RW-SX lock
			 */
			metricSpin := strings.Fields(spin)
			result = fillMetric("rw_sx_spin_waits", metricSpin[2], result)

			/*
			 * The total number of spin loops in #mutex_spin_wait() for RW-SX lock
			 */
			metricRounds := strings.Fields(rounds)
			result = fillMetric("rw_sx_spin_rounds", metricRounds[1], result)

			/*
			 * The total number of the sync_array event_wait in #mutex_spin_wait() for RW-SX lock
			 */
			metricOsWaits := strings.Fields(osWaits)
			result = fillMetric("rw_sx_os_waits", metricOsWaits[1], result)
		}
	}

	return result
}

// Not support currently.
func parseDeadlocksContent(_ []string) map[string]interface{} {
	result := make(map[string]interface{})
	return result
}

// Not support currently.
func parseForeignKeyErrorContent(_ []string) map[string]interface{} {
	result := make(map[string]interface{})
	return result
}

//nolint:gocyclo
func parseTransactionsContent(content []string) map[string]interface{} {
	result := make(map[string]interface{})

	for _, line := range content {
		item := "Trx id counter "
		if strings.HasPrefix(line, item) {
			/*
			 * The current transaction identifier which is incremented for each transaction
			 */
			record := strings.Fields(line)
			result = fillMetric("trx_id_counter", record[3], result)
		}

		item = "Purge done for "
		if strings.HasPrefix(line, item) {
			record := strings.Fields(line)
			/*
			 * The number of transaction to which purge is done
			 */
			result = fillMetric("purge_done_for_trx_num", record[6], result)
			/*
			 * The undo log record number which purge is currently processing
			 */
			result = fillMetric("purge_done_for_undo_num", record[10], result)
		}

		item = "History list length "
		if strings.HasPrefix(line, item) {
			/*
			 * The number of unpurged transactions in undo space
			 * It is increased as transactions which have done updates are committed and decreased as purge runs
			 */
			record := strings.Fields(line)
			result = fillMetric("history_list_length", record[3], result)
		}

		item = "---TRANSACTION"
		if strings.HasPrefix(line, item) {
			result = countUpMetric("current_transactions", 1, result)
			if strings.Contains(line, "ACTIVE") {
				result = countUpMetric("active_transactions", 1, result)
			}
		}

		item = "------- TRX HAS BEEN"
		if strings.HasPrefix(line, item) {
			/*
			 * The total time of lock waits
			 */
			record := strings.Fields(line)
			if num, err := strconv.Atoi(record[5]); err != nil {
				result = countUpMetric("transaction_lock_wait_secs", num, result)
			}
		}

		item = "mysql tables in use "
		if strings.HasPrefix(line, item) {
			parts := strings.Split(line, ", ")
			use, lock := parts[0], parts[1]

			/*
			 * The number of tables used by the transaction in question
			 */
			metricsUse := strings.Fields(use)
			if num, err := strconv.Atoi(metricsUse[4]); err != nil {
				countUpMetric("transaction_tables_in_use", num, result)
			}

			/*
			 * The number of tables locked by transactions
			 */
			metricsLock := strings.Fields(lock)
			if num, err := strconv.Atoi(metricsLock[1]); err != nil {
				countUpMetric("transaction_locked_tables", num, result)
			}
		}

		item = " lock struct(s), "
		if strings.Contains(line, item) {
			/*
			 * The total number of lock structs in row lock hash table is the number of row lock structures allocated by all transactions
			 * Not same as the number of locked rows, there are normally many rows for each lock structure
			 */
			record := strings.Fields(line)

			if strings.HasPrefix(line, "LOCK WAIT") {
				if num, err := strconv.Atoi(record[2]); err != nil {
					countUpMetric("transaction_lock_structs", num, result)
				}

				countUpMetric("locked_transactions", 1, result)
			} else if num, err := strconv.Atoi(record[0]); err != nil {
				countUpMetric("transaction_lock_structs", num, result)
			}
		}
	}

	return result
}

//nolint:funlen
func parseFileIoContent(content []string) map[string]interface{} {
	result := make(map[string]interface{})

	for _, line := range content {
		item := "(insert buffer thread)"
		if strings.Contains(line, item) {
			/*
			 * The total number of the insert_buffer threads
			 */
			result = countUpMetric("insert_buffer_threads_num", 1, result)
		}

		item = "(log thread)"
		if strings.Contains(line, item) {
			/*
			 * The total number of the log threads
			 */
			result = countUpMetric("log_threads_num", 1, result)
		}

		item = "(read thread)"
		if strings.Contains(line, item) {
			/*
			 * The total number of the read threads
			 */
			result = countUpMetric("read_threads_num", 1, result)
		}

		item = "(write thread)"
		if strings.Contains(line, item) {
			/*
			 * The total number of the write threads
			 */
			result = countUpMetric("write_threads_num", 1, result)
		}

		item = "Pending normal aio reads:"
		if strings.HasPrefix(line, item) {
			record := strings.Fields(line)
			//nolint:gomnd
			if len(record) >= 17 {
				/*
				 * The total number of the pending async I/O in normal reads
				 */
				result = fillMetric("pending_normal_reads_aio", record[4], result)

				/*
				 * The total number of the pending async I/O in normal writes
				 */
				result = fillMetric("pending_normal_writes_aio", record[12], result)
			}
		}

		item = "ibuf aio reads"
		if strings.HasPrefix(line, item) {
			record := strings.Fields(line)
			//nolint:gomnd
			if len(record) >= 10 {
				/*
				 * The total number of the pending async I/O in insert_buffer reads
				 */
				result = fillMetric("pending_ibuf_reads_aio", record[3], result)

				/*
				 * The total number of the pending async I/O in insert_buffer logs
				 */
				result = fillMetric("pending_log_ios_aio", record[6], result)

				/*
				 * The total number of the pending async I/O in insert_buffer syncs
				 */
				result = fillMetric("pending_sync_ios_aio", record[9], result)
			}
		}

		item = "Pending flushes (fsync)"
		if strings.HasPrefix(line, item) {
			parts := strings.Split(line, "; ")
			log, bufPool := parts[0], parts[1]

			/*
			 * The total number of the pending I/O in log flushes
			 */
			metricsLog := strings.Fields(log)
			result = fillMetric("pending_log_flushes", metricsLog[4], result)

			/*
			 * The total number of the pending I/O in buffer pool flushes
			 */
			metricsBufPool := strings.Fields(bufPool)
			result = fillMetric("pending_buf_pool_flushes", metricsBufPool[2], result)
		}

		item = " OS file reads"
		if strings.Contains(line, item) {
			parts := strings.Split(line, ", ")
			reads, writes, fsyncs := parts[0], parts[1], parts[2]

			/*
			 * The total number of OS file reads
			 */
			metricsReads := strings.Fields(reads)
			result = fillMetric("os_file_reads", metricsReads[0], result)

			/*
			 * The total number of OS file writes
			 */
			metricsWrites := strings.Fields(writes)
			result = fillMetric("os_file_writes", metricsWrites[0], result)

			/*
			 * The total number of OS file fsyncs
			 */
			metricsFsyncs := strings.Fields(fsyncs)
			result = fillMetric("os_file_fsyncs", metricsFsyncs[0], result)
		}

		item = " reads/s"
		if strings.Contains(line, item) {
			parts := strings.Split(line, ", ")
			reads, avgReads, writes, fsyncs := parts[0], parts[1], parts[2], parts[3]

			/*
			 * The per second average of OS file reads
			 */
			metricsReads := strings.Fields(reads)
			result = fillMetric("os_file_reads_per_sec", metricsReads[0], result)

			/*
			 * The per operation average bytes of OS file reads
			 */
			metricAvgReads := strings.Fields(avgReads)
			result = fillMetric("os_file_bytes_per_read", metricAvgReads[0], result)

			/*
			 * The per second average of OS file writes
			 */
			metricsWrites := strings.Fields(writes)
			result = fillMetric("os_file_writes_per_sec", metricsWrites[0], result)

			/*
			 * The per second average of OS file fsyncs
			 */
			metricsFsyncs := strings.Fields(fsyncs)
			result = fillMetric("os_file_fsyncs_per_sec", metricsFsyncs[0], result)
		}
	}

	return result
}

func parseInsertBufferAndAdaptiveHashIndexContent(content []string) map[string]interface{} {
	result := make(map[string]interface{})

	for i, line := range content {
		item := "Ibuf: size "
		if strings.HasPrefix(line, item) {
			parts := strings.Split(line, ", ")
			size, free, seg, merge := parts[0], parts[1], parts[2], parts[3]

			/*
			 * The size of insert_buffer
			 */
			metricsSize := strings.Fields(size)
			result = fillMetric("insert_buffer_size", metricsSize[2], result)

			/*
			 * The length of the free list in insert_buffer
			 */
			metricsFree := strings.Fields(free)
			result = fillMetric("insert_buffer_free_list_len", metricsFree[3], result)

			/*
			 * The segment size of insert_buffer
			 */
			metricsSeg := strings.Fields(seg)
			result = fillMetric("insert_buffer_seg_size", metricsSeg[2], result)

			/*
			 * The segment size of insert_buffer
			 */
			metricsMerge := strings.Fields(merge)
			result = fillMetric("insert_buffer_merges", metricsMerge[0], result)
		}

		item = "merged operations:"
		if strings.HasPrefix(line, item) {
			stmt := content[i+1]
			parts := strings.Split(stmt, ", ")
			insert, deleteMark, deletes := parts[0], parts[1], parts[2]

			/*
			 * The total number of the merged insert operations in insert_buffer
			 */
			metricsInsert := strings.Fields(insert)
			result = fillMetric("insert_buffer_merged_inserts", metricsInsert[1], result)

			/*
			 * The total number of the merged delete mark operations in insert_buffer
			 */
			metricsDeleteMark := strings.Fields(deleteMark)
			result = fillMetric("insert_buffer_merged_delete_marks", metricsDeleteMark[2], result)

			/*
			 * The total number of the merged delete operations in insert_buffer
			 */
			metricsDelete := strings.Fields(deletes)
			result = fillMetric("insert_buffer_merged_deletes", metricsDelete[1], result)
		}

		item = "discarded operations:"
		if strings.HasPrefix(line, item) {
			stmt := content[i+1]
			parts := strings.Split(stmt, ", ")
			insert, deleteMark, deletes := parts[0], parts[1], parts[2]

			/*
			 * The total number of the discarded insert operations in insert_buffer
			 */
			metricsInsert := strings.Fields(insert)
			result = fillMetric("insert_buffer_discarded_inserts", metricsInsert[1], result)

			/*
			 * The total number of the discarded delete mark operations in insert_buffer
			 */
			metricsDeleteMark := strings.Fields(deleteMark)
			result = fillMetric("insert_buffer_discarded_delete_marks", metricsDeleteMark[2], result)

			/*
			 * The total number of the discarded delete operations in insert_buffer
			 */
			metricsDelete := strings.Fields(deletes)
			result = fillMetric("insert_buffer_discarded_deletes", metricsDelete[1], result)
		}

		item = "hash searches/s, "
		if strings.Contains(line, item) {
			parts := strings.Split(line, ", ")
			hash, nonHash := parts[0], parts[1]

			/*
			 * The per second average of the adaptive hash searches
			 */
			metricsHash := strings.Fields(hash)
			result = fillMetric("hash_searches_per_sec", metricsHash[0], result)

			/*
			 * The per second average of the non hash searches
			 */
			metricsNonHash := strings.Fields(nonHash)
			result = fillMetric("non_hash_searches_per_sec", metricsNonHash[0], result)
		}
	}

	return result
}

func parseLogContent(content []string) map[string]interface{} {
	result := make(map[string]interface{})

	for _, line := range content {
		item := "Log sequence number "
		if strings.HasPrefix(line, item) {
			/*
			 * Current Log Sequence Number which is amount of bytes InnoDB has written in log files since system tablespace creation
			 */
			record := strings.Fields(line)
			result = fillMetric("log_sequence_number", record[3], result)
		}

		item = "Log flushed up to "
		if strings.HasPrefix(line, item) {
			/*
			 * The range of logs which have been flushed to the storage
			 */
			record := strings.Fields(line)
			result = fillMetric("log_flushed_up_to", record[4], result)
		}

		item = "Last checkpoint at "
		if strings.HasPrefix(line, item) {
			/*
			 * LSN of which point logs have been written by the last checkpoint
			 */
			record := strings.Fields(line)
			result = fillMetric("log_last_checkpoint_at", record[3], result)
		}

		item = " pending log flushes, "
		if strings.Contains(line, item) {
			parts := strings.Split(line, ", ")
			pendLog, pendCheckpoint := parts[0], parts[1]

			/*
			 * The total number of pending normal log writes
			 */
			metricsPendLog := strings.Fields(pendLog)
			result = fillMetric("log_pending_log_flushes", metricsPendLog[0], result)

			/*
			 * The total number of checkpoint log writes
			 */
			metricPendCheckpoint := strings.Fields(pendCheckpoint)
			result = fillMetric("log_pending_checkpoint_writes", metricPendCheckpoint[0], result)
		}

		item = " log i/o's done, "
		if strings.Contains(line, item) {
			parts := strings.Split(line, ", ")
			doneLog, ioPS := parts[0], parts[1]

			/*
			 * The total number of log I/O operations
			 */
			metricsDoneLog := strings.Fields(doneLog)
			result = fillMetric("log_ios_done", metricsDoneLog[0], result)

			/*
			 * The per second average of log I/O operations
			 */
			metricIoPS := strings.Fields(ioPS)
			result = fillMetric("log_iops", metricIoPS[0], result)
		}
	}

	return result
}

//nolint:funlen,gocyclo
func parseBufferPoolAndMemoryContent(content []string) map[string]interface{} {
	result := make(map[string]interface{})

	for _, line := range content {
		item := "Total large memory allocated "
		if strings.HasPrefix(line, item) {
			/*
			 * The total memory allocated for the buffer pool in bytes
			 */
			record := strings.Fields(line)
			result = fillMetric("buffer_pool_total_large_memory_allocated", record[4], result)
		}

		item = "Dictionary memory allocated "
		if strings.HasPrefix(line, item) {
			/*
			 * The total memory allocated for the InnoDB data dictionary in bytes
			 */
			record := strings.Fields(line)
			result = fillMetric("buffer_pool_dictionary_memory_allocated", record[3], result)
		}

		item = "Buffer pool size "
		if strings.HasPrefix(line, item) {
			/*
			 * The total size in pages allocated to the buffer pool
			 */
			record := strings.Fields(line)
			result = fillMetric("buffer_pool_size", record[3], result)
		}

		item = "Free buffers "
		if strings.HasPrefix(line, item) {
			/*
			 * The total size in pages of the buffer pool free list
			 */
			record := strings.Fields(line)
			result = fillMetric("buffer_pool_free_buffers", record[2], result)
		}

		item = "Database pages "
		if strings.HasPrefix(line, item) {
			/*
			 * The total size in pages of the buffer pool LRU list
			 */
			record := strings.Fields(line)
			result = fillMetric("buffer_pool_database_pages", record[2], result)
		}

		item = "Old database pages "
		if strings.HasPrefix(line, item) {
			/*
			 * The total size in pages of the buffer pool old LRU sublist
			 */
			record := strings.Fields(line)
			result = fillMetric("buffer_pool_old_database_pages", record[3], result)
		}

		item = "Modified db pages "
		if strings.HasPrefix(line, item) {
			/*
			 * The current number of pages modified in the buffer pool
			 */
			record := strings.Fields(line)
			result = fillMetric("buffer_pool_modified_db_pages", record[3], result)
		}

		item = "Pending reads "
		if strings.HasPrefix(line, item) {
			/*
			 * The number of of buffer pool pages waiting to be read into the buffer pool
			 */
			record := strings.Fields(line)
			result = fillMetric("buffer_pool_pending_reads", record[2], result)
		}

		item = "Pending writes "
		if strings.HasPrefix(line, item) {
			parts := strings.Split(line, ", ")
			lru, flushList, singlePage := parts[0], parts[1], parts[2]

			/*
			 * The number of old dirty pages within the buffer pool to be written from the bottom of the LRU list
			 */
			metricLru := strings.Fields(lru)
			result = fillMetric("buffer_pool_pending_writes_lru", metricLru[3], result)

			/*
			 * The number of buffer pool pages to be flushed during checkpoint
			 */
			metricFlushList := strings.Fields(flushList)
			result = fillMetric("buffer_pool_pending_writes_flush_list", metricFlushList[2], result)

			/*
			 * The number of pending independent page writes within the buffer pool
			 */
			metricSinglePage := strings.Fields(singlePage)
			result = fillMetric("buffer_pool_pending_writes_single_page", metricSinglePage[2], result)
		}

		item = "Pages made young "
		if strings.HasPrefix(line, item) {
			parts := strings.Split(line, ", ")
			young, notYoung := parts[0], parts[1]

			/*
			 * The total number of pages made young in the buffer pool LRU list
			 * moved to the head of sublist of new pages
			 */
			metricYoung := strings.Fields(young)
			result = fillMetric("buffer_pool_pages_made_young", metricYoung[3], result)

			/*
			 * The total number of pages not made young in the buffer pool LRU list
			 * pages that have remained in the old sublist without being made young
			 */
			metricNotYoung := strings.Fields(notYoung)
			result = fillMetric("buffer_pool_pages_not_made_young", metricNotYoung[2], result)
		}

		item = " youngs/s, "
		if strings.Contains(line, item) {
			parts := strings.Split(line, ", ")
			young, notYoung := parts[0], parts[1]

			/*
			 * The per second average of accesses to old pages in the buffer pool LRU list that have resulted in making pages young
			 */
			metricYoung := strings.Fields(young)
			result = fillMetric("buffer_pool_pages_made_young_per_sec", metricYoung[0], result)

			/*
			 * The per second average of accesses to old pages in the buffer pool LRU list that have resulted in not making pages young
			 */
			metricNotYoung := strings.Fields(notYoung)
			result = fillMetric("buffer_pool_pages_not_made_youngs_per_sec", metricNotYoung[0], result)
		}

		item = ", created "
		if strings.Contains(line, item) {
			parts := strings.Split(line, ", ")
			read, created, written := parts[0], parts[1], parts[2]

			/*
			 * The total number of pages read from the buffer pool
			 */
			metricRead := strings.Fields(read)
			result = fillMetric("buffer_pool_pages_read", metricRead[2], result)

			/*
			 * The total number of pages created within the buffer pool
			 */
			metricCreated := strings.Fields(created)
			result = fillMetric("buffer_pool_pages_created", metricCreated[1], result)

			/*
			 * The total number of pages written from the buffer pool
			 */
			metricWritten := strings.Fields(written)
			result = fillMetric("buffer_pool_pages_written", metricWritten[1], result)
		}

		item = " reads/s, "
		if strings.Contains(line, item) {
			parts := strings.Split(line, ", ")
			read, created, written := parts[0], parts[1], parts[2]

			/*
			 * The average number of buffer pool page reads per second
			 */
			metricRead := strings.Fields(read)
			result = fillMetric("buffer_pool_pages_reads_per_sec", metricRead[0], result)

			/*
			 * The average number of buffer pool page creates per second
			 */
			metricCreated := strings.Fields(created)
			result = fillMetric("buffer_pool_pages_creates_per_sec", metricCreated[0], result)

			/*
			 * The average number of buffer pool page writes per second
			 */
			metricWritten := strings.Fields(written)
			result = fillMetric("buffer_pool_pages_writes_per_sec", metricWritten[0], result)
		}

		item = "Buffer pool hit rate "
		if strings.HasPrefix(line, item) {
			parts := strings.Split(line, ", ")
			cacheHit, youngCreate := parts[0], parts[1]

			/*
			 * The buffer pool page hit rate for pages read from the buffer pool memory vs from disk storage
			 */
			metricCacheHit := strings.Fields(cacheHit)
			rate, err := calcRate(metricCacheHit[4], metricCacheHit[6])

			if err == nil {
				result = fillMetric("buffer_pool_cache_hit_rate", rate, result)
			}

			metricYoungCreate := strings.Fields(youngCreate)
			/*
			 * The average hit rate at which page accesses have resulted in making pages young
			 */
			rate, err = calcRate(metricYoungCreate[2], metricYoungCreate[4])
			if err == nil {
				result = fillMetric("buffer_pool_cache_hit_make_young_rate", rate, result)
			}

			/*
			 * The average hit rate at which page accesses have not resulted in making pages young
			 */
			rate, err = calcRate(metricYoungCreate[6], metricYoungCreate[8])
			if err == nil {
				result = fillMetric("buffer_pool_cache_hit_not_make_young_rate", rate, result)
			}
		}

		item = "Pages read ahead "
		if strings.HasPrefix(line, item) {
			parts := strings.Split(line, ", ")
			read, evict, randRead := parts[0], parts[1], parts[2]

			/*
			 * The per second average of read ahead operations
			 */
			metricRead := strings.Fields(read)
			result = fillMetric("buffer_pool_page_reads_ahead_per_sec", strings.TrimSuffix(metricRead[3], "/s"), result)

			/*
			 * The per second average of the pages evicted without being accessed from the buffer pool
			 */
			metricEvict := strings.Fields(evict)
			result = fillMetric("buffer_pool_page_evicted_wo_access_per_sec", strings.TrimSuffix(metricEvict[3], "/s"), result)

			/*
			 * The per second average of random read ahead operations
			 */
			metricRandRead := strings.Fields(randRead)
			result = fillMetric("buffer_pool_page_random_reads_ahead_per_sec", strings.TrimSuffix(metricRandRead[3], "/s"), result)
		}

		item = "LRU len: "
		if strings.HasPrefix(line, item) {
			parts := strings.Split(line, ", ")
			lru, unzipLru := parts[0], parts[1]

			/*
			 * The total size in pages of the buffer pool LRU list
			 */
			metricLru := strings.Fields(lru)
			result = fillMetric("buffer_pool_lru_len", metricLru[2], result)

			/*
			 * The total size in pages of the buffer pool unzip_LRU list
			 */
			metricUnzipLru := strings.Fields(unzipLru)
			result = fillMetric("buffer_pool_unzip_lru_len", metricUnzipLru[2], result)
		}

		item = "I/O sum"
		if strings.HasPrefix(line, item) {
			parts := strings.Split(line, ", ")
			ioPart, unzipPart := parts[0], parts[1]

			ioRecord := strings.Fields(ioPart)
			io := strings.Split(ioRecord[1], ":")
			ioSum, ioCur := io[0], io[1]

			/*
			 * The total number of buffer pool LRU list pages accessed for the last 50 seconds
			 */
			metricIoSum := strings.TrimSuffix(strings.TrimPrefix(ioSum, "sum["), "]")
			result = fillMetric("buffer_pool_io_sum", metricIoSum, result)

			/*
			 * The total number of buffer pool LRU list pages accessed
			 */
			metricIoCur := strings.TrimSuffix(strings.TrimPrefix(ioCur, "cur["), "]")
			result = fillMetric("buffer_pool_io_cur", metricIoCur, result)

			unzipRecord := strings.Fields(unzipPart)
			unzip := strings.Split(unzipRecord[1], ":")
			unzipSum, unzipCur := unzip[0], unzip[1]

			/*
			 * The total number of buffer pool unzip_LRU list pages accessed
			 */
			metricUnzipSum := strings.TrimSuffix(strings.TrimPrefix(unzipSum, "sum["), "]")
			result = fillMetric("buffer_pool_io_unzip_sum", metricUnzipSum, result)

			/*
			 * The total number of buffer pool unzip_LRU list pages accessed
			 */
			metricUnzipCur := strings.TrimSuffix(strings.TrimPrefix(unzipCur, "cur["), "]")
			result = fillMetric("buffer_pool_io_unzip_cur", metricUnzipCur, result)
		}
	}

	return result
}

// Not support currently.
func parseIndividualBufferPoolInfoContent(_ []string) map[string]interface{} {
	result := make(map[string]interface{})
	return result
}

func parseRowOperationsContent(content []string) map[string]interface{} {
	result := make(map[string]interface{})

	for _, line := range content {
		item := " queries inside InnoDB, "
		if strings.Contains(line, item) {
			parts := strings.Split(line, ", ")
			innodb, queue := parts[0], parts[1]

			/*
			 * The total number of queries inside the InnoDB
			 */
			metricInnodb := strings.Fields(innodb)
			result = fillMetric("queries_inside_innodb", metricInnodb[0], result)

			/*
			 * The total number of queries in the queue
			 */
			metricQueue := strings.Fields(queue)
			result = fillMetric("queries_in_queue", metricQueue[0], result)
		}

		item = " read views open inside InnoDB"
		if strings.Contains(line, item) {
			record := strings.Fields(line)
			result = fillMetric("read_views_inside_innodb", record[0], result)
		}

		item = "Number of rows inserted "
		if strings.HasPrefix(line, item) {
			parts := strings.Split(line, ", ")
			inserted, updated, deleted, read := parts[0], parts[1], parts[2], parts[3]

			/*
			 * The total number of inserted rows
			 */
			metricInserted := strings.Fields(inserted)
			result = fillMetric("rows_inserted", metricInserted[4], result)

			/*
			 * The total number of updated rows
			 */
			metricUpdated := strings.Fields(updated)
			result = fillMetric("rows_updated", metricUpdated[1], result)

			/*
			 * The total number of deleted rows
			 */
			metricDeleted := strings.Fields(deleted)
			result = fillMetric("rows_deleted", metricDeleted[1], result)

			/*
			 * The total number of read rows
			 */
			metricRead := strings.Fields(read)
			result = fillMetric("rows_read", metricRead[1], result)
		}

		item = "inserts/s, "
		if strings.Contains(line, item) {
			parts := strings.Split(line, ", ")
			inserted, updated, deleted, read := parts[0], parts[1], parts[2], parts[3]

			/*
			 * The per second average of inserted rows
			 */
			metricInserted := strings.Fields(inserted)
			result = fillMetric("rows_inserted_per_sec", metricInserted[0], result)

			/*
			 * The per second average of updated rows
			 */
			metricUpdated := strings.Fields(updated)
			result = fillMetric("rows_updated_per_sec", metricUpdated[0], result)

			/*
			 * The per second average of deleted rows
			 */
			metricDeleted := strings.Fields(deleted)
			result = fillMetric("rows_deleted_per_sec", metricDeleted[0], result)

			/*
			 * The per second average of read rows
			 */
			metricRead := strings.Fields(read)
			result = fillMetric("rows_read_per_sec", metricRead[0], result)
		}
	}

	return result
}

func (p *InnodbStatusParser) stepNextSection() {
	p.content = make([]string, 0)

	for i := p.currSection; i <= EndOfInnodbMonitorOutput; i++ {
		if p.maybeLabel() == i.String() {
			p.currSection = i
			break
		}
	}

	p.currPos += 3
}
