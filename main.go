package main

import (
	"log"
	"os"

	"github.com/Jacobious52/telegram-stats/pkg/stats"

	"github.com/Jacobious52/telegram-stats/pkg/importer"

	"gopkg.in/alecthomas/kingpin.v2"
)

var inputPath = kingpin.Arg("input", "location of csv file").Required().ExistingFile()
var fromTarget = kingpin.Flag("from", "who to calculate stats on").Required().String()
var topN = kingpin.Flag("top", "number of top words to print out").Default("10").Int()

func main() {
	kingpin.Parse()

	// open the file from disk
	file, err := os.Open(*inputPath)
	if err != nil {
		log.Fatalln("failed to open file:", err)
	}
	defer file.Close()

	// import the CSV into a table format
	table, err := importer.Import(file)
	if err != nil {
		log.Fatalln("failed to import file")
	}

	// calculate and return an object with the data we can print with
	calcStats := stats.CalculateStats(
		*fromTarget,
		table,
		stats.FilterShort,
		//Examples:
		//stats.FilterListInclude("games", "friends"),
		//stats.FilterListExclude("hunter2", "password1"),
	)

	// don't print if they didn't talk
	if calcStats.TotalMessages == 0 {
		return
	}
	calcStats.WriteStats(os.Stdout, *topN)
}
