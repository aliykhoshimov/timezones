package main

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func Test_Timezone(t *testing.T) {

	// t1 is in the UTC timezone
	t1 := time.Date(2024, 01, 22, 13, 37, 0, 0, time.UTC)

	// load the Hong Kong location from the time zone database.
	hongKong, err := time.LoadLocation("Asia/Hong_Kong")
	if err != nil {
		log.Fatalf("failed to load location: %v", err)
	}

	t.Logf("Hong Kong timezone: %v", hongKong.String())

	// t2 represents the same time instant in the Hong Kong timezone (UTC+8)
	t2 := t1.In(hongKong)

	zoneName, offset := t2.Zone()

	t.Logf("Timezone: %v, Offset: %v", zoneName, offset)

	fmt.Println(t1.String())
	fmt.Println(t2.String())
	fmt.Printf("same time instant: %v\n", t1.Equal(t2))
}
