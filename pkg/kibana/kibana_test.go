package kibana

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

const kibanaFetchExampleFilePath = "/ws/kibana/kibana.fetch"

func TestKibana(t *testing.T) {
	kibanaRequest, err := LoadFromFetchRequest(kibanaFetchExampleFilePath)
	if err != nil {
		t.Skip("Skipping test if file is not available")
	}
	filter := map[string]string{
		"application": "applikasjon",
		"team":        "team",
		"class":       "tm.class",
		"method":      "tm.method",
		"stacktrace":  "tm.stacktrace",
	}

	//	timeInterval, err := ExtractTimeIntervalFrom(kibanaRequest)
	expectedBuckets := createTestBuckets([]string{
		"2022-11-27T21:37:48.801Z,2022-12-02T13:19:46.59Z,1000",
		"2022-11-27T21:37:48.801Z,2022-12-02T11:43:33.929Z,1000",
	})

	timeInterval := TimeInterval{
		Gte: expectedBuckets[1].From,
		Lte: expectedBuckets[1].To,
	}

	if err != nil {
		t.Errorf("%v\n", err)
	}

	resultExists := make(map[string]bool)
	err, _, result, _ := ExecuteKibanaQuery(kibanaRequest, timeInterval, filter, resultExists, "")

	if err != nil {
		t.Errorf("%v\n", err)
	}

	fmt.Println("Unique result hits:", len(result))
}

func TestReadFetchKibana(t *testing.T) {
	kibanaRequest, err := LoadFromFetchRequest(kibanaFetchExampleFilePath)

	if err != nil {
		t.Skip("Skipping test if file is not available")
	}

	fmt.Println("request:")
	fmt.Println(kibanaRequest.Body)
}

func createTestBuckets(buckets []string) []TimeBucket {
	var timeBuckets []TimeBucket
	for _, bucket := range buckets {
		split := strings.Split(bucket, ",")
		gte, _ := time.Parse(time.RFC3339Nano, split[0])
		lte, _ := time.Parse(time.RFC3339Nano, split[1])
		count, _ := strconv.Atoi(split[2])
		timeBuckets = append(timeBuckets, TimeBucket{
			From:  gte,
			To:    lte,
			Count: count,
		})

	}
	return timeBuckets
}
