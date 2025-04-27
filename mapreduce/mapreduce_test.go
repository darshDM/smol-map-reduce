package mapreduce_test

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	mr "github.com/darshDM/smol-map-reduce/mapreduce"
)

type Logtype int

const (
	FAILED Logtype = iota
	SUCCESS
	DISCONNECTED
)

type LogRecord struct {
	Month    string
	Date     int
	Time     time.Time
	LogLevel Logtype
	IP       string
}
type IPStats struct {
	IPCounts map[string]int
}

func mapper(data string) *LogRecord {
	FailedPassRegex := regexp.MustCompile(`^([A-Z][a-z]{2})\s+(\d+)\s+([\d:]+).*Failed password.*from ([\d\.]+)`)
	matches := FailedPassRegex.FindStringSubmatch(data)

	if len(matches) == 5 {
		date, err := strconv.Atoi(matches[2])
		if err != nil {
			fmt.Print("Error converting date to int: ", err)
		}

		time, err := time.Parse("15:04:05", matches[3])
		if err != nil {
			fmt.Print("Error converting time to time.Time: ", err)
		}
		LogRecord := &LogRecord{
			Month:    matches[1],
			Date:     date,
			Time:     time,
			LogLevel: FAILED,
			IP:       matches[4],
		}
		return LogRecord
	}
	return nil
}

func reduce(log []*LogRecord) IPStats {
	ipCounter := make(map[string]int)
	for _, record := range log {
		ipCounter[record.IP]++
	}
	return IPStats{IPCounts: ipCounter}
}
func BenchmarkMapReduce(b *testing.B) {
	mapReduce := mr.NewMapReduce[string, LogRecord, IPStats](8, mapper, reduce)
	mapReduce.Start()

	fl, err := os.Open("../data/SSH.log")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer fl.Close()

	reader := bufio.NewReader(fl)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading line: %v\n", err)
			break
		}
		mapReduce.Pool <- string(line)
	}

	mapReduce.Stop()
	result := mapReduce.GetResult()
	for ip, count := range result.IPCounts {
		if count > 50 {
			fmt.Printf("IP: %s, [%d] failed attempt (suspicious)\n", ip, count)
		} else {
			fmt.Printf("IP: %s, [%d] failed attempt\n", ip, count)
		}
	}

}
