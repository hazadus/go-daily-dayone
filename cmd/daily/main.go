package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type DailyNote struct {
	Date    time.Time
	Content string
}

func main() {
	var dailyPath = "/Users/hazadus/Library/Mobile Documents/iCloud~md~obsidian/Documents/Hazadus Vault/Daily/"
	var start = "2024-12-29"
	var end = "2024-12-29"

	startFlag := flag.String("start", start, "Start date as YYYY-MM-DD")
	endFlag := flag.String("end", end, "End date as YYYY-MM-DD")
	flag.Parse()

	startDate, err := time.Parse("2006-01-02", *startFlag)
	if err != nil {
		panic("Incorrect start date")
	}
	endDate, err := time.Parse("2006-01-02", *endFlag)
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
		err := createDayOneNote(note)
		if err != nil {
			fmt.Printf("failed to create dayone note: %s\n", err)
			os.Exit(1)
		}
	}
}

// readNotes читает ежедневные заметки из файлов за указанные интервал дат,
// и сохраняет их в массиве.
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
			Content: clear(dailyInfo),
		})
	}

	return dailyNotes, nil
}

// getDailyInfo достаёт нужную информацию из ежедневной заметки.
func getDailyInfo(note []byte) (string, error) {
	re, err := regexp.Compile(`(?s)## События дня(.*?)----\n# План на сегодня(.*?)----(?s)`)
	if err != nil {
		return "", fmt.Errorf("error compiling regexp: %w", err)
	}
	res := re.Find(note)
	return string(res), nil
}

// createDayOneNote вызывает dayone2 CLI для создания заметки в DayOne.
func createDayOneNote(note *DailyNote) error {
	weekdays := map[string]string{
		"Monday":    "Понедельник",
		"Tuesday":   "Вторник",
		"Wednesday": "Среда",
		"Thursday":  "Четверг",
		"Friday":    "Пятница",
		"Saturday":  "Суббота",
		"Sunday":    "Воскресенье",
	}

	//nolint:all
	dayOneCmd := exec.Command(
		"dayone2",
		"-j",
		"Журнал",
		fmt.Sprintf("--date=%s 21:00:00", note.Date.Format("2006-01-02")),
		"new",
		weekdays[note.Date.Weekday().String()]+"\n",
		note.Content,
	)
	if output, err := dayOneCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error running dayone2: %w, %s", err, string(output))
	}

	return nil
}

// clear удаляет из заметки нежелательные строки.
func clear(content string) string {
	res := strings.Replace(content, "📚 **Книга**: –", "", 1)
	res = strings.Replace(res, "🎮 **Игра**: –", "", 1)
	res = strings.Replace(res, "🎥 **Сериал**: –", "", 1)
	res = strings.Replace(res, "🏒 **Спорт**: –", "", 1)
	res = strings.Replace(res, "🍿 **Фильм**: –", "", 1)
	return res
}
