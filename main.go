package main

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func getLineWithContaining(data []string, search string) string {
	for _, str := range data {
		if strings.Contains(str, search) {
			return str
		}
	}

	return ""
}

func getDataFromField(data []string, search string) uint64 {
	line := getLineWithContaining(data, search)
	if line == "" {
		log.Fatalf("Data field %s not found", search)
	}

	dataFields := strings.Fields(line)
	if len(dataFields) != 3 {
		log.Fatalf("Cannot extract data from %d fields, needed 3", len(dataFields))
	}

	dataField := dataFields[2]

	dataValue, err := strconv.ParseUint(dataField, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	return dataValue
}

func getColorForFilled(filled int64) *color.Color {
	if filled > 6 {
		return color.New(color.FgGreen)
	} else if filled > 4 {
		return color.New(color.FgYellow)
	} else {
		return color.New(color.FgRed)
	}
}

func main() {
	cmd := exec.Command("ioreg", "-rc", "AppleSmartBattery")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Split it by newline
	dataLines := strings.Split(out.String(), "\n")

	maxCapacity := getDataFromField(dataLines, "MaxCapacity")
	currentCapacity := getDataFromField(dataLines, "CurrentCapacity")
	minutesRemaining := getDataFromField(dataLines, "TimeRemaining")
	timeRemaining := time.Duration(minutesRemaining) * time.Minute

	charge := float64(currentCapacity) / float64(maxCapacity)
	chargeThreshold := int64(math.Ceil(10.0 * charge))

	filledArrow := string('▸')
	emptyArrow := string('▹')

	var totalSlots int64 = 10

	filledAmount := int64(math.Ceil(float64(chargeThreshold) * (float64(totalSlots) / 10.0)))
	emptyAmount := totalSlots - filledAmount

	var i int64

	var output string

	for i = 0; i < filledAmount; i++ {
		output += filledArrow
	}

	for i = 0; i < emptyAmount; i++ {
		output += emptyArrow
	}

	output += fmt.Sprintf(" (%s)", timeRemaining.String())

	// color := getColorForFilled(filledAmount)
	fmt.Println(output)
}
