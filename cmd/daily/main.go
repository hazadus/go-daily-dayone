package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type DailyNote struct {
	Date    time.Time
	Content string
}

func main() {
	var dailyPath = "/Users/hazadus/Library/Mobile Documents/iCloud~md~obsidian/Documents/Hazadus Vault/Daily/"
	var start = "2024-05-13"
	var end = "2024-05-20"

	startDate, err := time.Parse("2006-01-02", start)
	if err != nil {
		panic("Incorrect start date")
	}
	endDate, err := time.Parse("2006-01-02", end)
	if err != nil {
		panic("Incorrect end date")
	}

	if !(endDate.After(startDate) || endDate.Equal(startDate)) {
		panic("End date must be greater or equal to start date")
	}

	dailyNotes, err := readNotes(dailyPath, startDate, endDate)
	if err != nil {
		fmt.Printf("error reading notes: %s", err)
		os.Exit(1)
	}

	for _, note := range dailyNotes {
		fmt.Printf("Date: %s\n%s\n\n", note.Date, note.Content)
	}
}

func readNotes(path string, startDate, endDate time.Time) ([]*DailyNote, error) {
	fmt.Printf("Daily notes path: %s\nDates: %s - %s\n",
		path, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	dailyNotes := []*DailyNote{}

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		fileName := fmt.Sprintf("%s.md", d.Format("2006-01-02"))
		filePath := filepath.Join(path, fileName)
		_, err := os.Stat(filePath)
		if err != nil {
			fmt.Printf("%s does not exist\n", fileName)
			continue
		}

		fmt.Printf("%s found, reading...\n", fileName)
		//nolint:all
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
		}
		fmt.Printf("    %d characters\n", len(content))

		dailyInfo, err := getDailyInfo(content)
		if err != nil {
			return nil, fmt.Errorf("error reading daily info from note: %w", err)
		}

		dailyNotes = append(dailyNotes, &DailyNote{
			Date:    d,
			Content: dailyInfo,
		})
	}

	return dailyNotes, nil
}

func getDailyInfo(note []byte) (string, error) {
	re, err := regexp.Compile(`(?s)## События дня(.*?)----(?s)`)
	if err != nil {
		return "", fmt.Errorf("error compiling regexp: %w", err)
	}
	res := re.Find(note)
	return string(res), nil
}
