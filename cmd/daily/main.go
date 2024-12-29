package main

import (
	"fmt"
	"time"
)

type DailyNote struct {
	Date    string
	Content string
}

func main() {
	var dailyPath = "/Users/hazadus/Library/Mobile Documents/iCloud~md~obsidian/Documents/Hazadus Vault/Daily/"
	var start = "2024-05-13"
	var end = "2024-05-16"

	// Convert dates from strings
	startDate, err := time.Parse("2006-01-02", start)
	if err != nil {
		panic("Incorrect start date")
	}
	endDate, err := time.Parse("2006-01-02", end)
	if err != nil {
		panic("Incorrect end date")
	}

	// Check if endDate >= startDate
	if !(endDate.After(startDate) || endDate.Equal(startDate)) {
		panic("End date must be greater or equal to start date")
	}

	fmt.Printf("go-daily-dayone\nDaily path: %s\nDates: %s - %s\n", dailyPath, start, end)
	readNotes(dailyPath, startDate, endDate)
}

func readNotes(path string, startDate, endDate time.Time) ([]*DailyNote, error) {
	for d := startDate; d.After(endDate) == false; d = d.AddDate(0, 0, 1) {
		fmt.Println(d.Format("2006-01-02"))
	}

	return nil, nil
}
