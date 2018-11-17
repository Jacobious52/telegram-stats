package stats

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/Jacobious52/telegram-stats/pkg/importer"
)

// Stats struct contains all the statistics calculated from a table and a person
// Includes numeric values, freq table, and rankedlist
// Create and calculate with CalculateStats()
type Stats struct {
	From    string
	Filters []FilterFunc

	TotalWords    int
	TotalMessages int
	TotalVocab    int
	AverageWords  float64

	WordFreq map[string]int
	WordRank RankedList
}

// CalculateStats calculates the stats for a person with the filter funcs and returns a new stats struct to hold them
func CalculateStats(from string, table *importer.Table, filters ...FilterFunc) *Stats {
	// keep a running average of words by sentence
	var totalWords, totalMessages float64
	// create a map to add the frequencies to
	freq := make(map[string]int)
	for _, row := range table.Rows {
		// ingore people where's not looking for
		if row.From != from {
			continue
		}

		// remove the non alpha text from the strings
		norm := NormaliseText(row.Text)
		// separate into words, filtering out any that match the filter func
		words := Words(norm, filters)
		// calculate the frequency of each word and add it to the running total
		Frequency(words, freq)

		// update the average words and total messages
		totalWords += float64(len(words))
		totalMessages++
	}
	rank := Rank(freq)

	var averageWords float64
	if totalMessages > 0 {
		averageWords = totalWords / totalMessages
	}

	return &Stats{
		From: from,

		TotalWords:    int(totalWords),
		TotalMessages: int(totalMessages),
		TotalVocab:    len(freq),
		AverageWords:  averageWords,
		WordFreq:      freq,
		WordRank:      rank,
	}
}

// WriteStats writes out the stats in a "pretty" format for printing
func (s *Stats) WriteStats(w io.Writer, topN int) {
	fmt.Fprintln(w, "person:", s.From)
	fmt.Fprintln(w, "total vocabulary:", len(s.WordRank), "words")
	fmt.Fprintln(w, "total words:", s.TotalWords, "words")
	fmt.Fprintln(w, "total messages:", s.TotalMessages, "messages")
	fmt.Fprintln(w, "average message length:", s.AverageWords, "words")
	fmt.Fprintln(w, ".")
	fmt.Fprintln(w, "------------------------------------------")
	fmt.Fprintln(w, "top", topN, "most frequent words")
	fmt.Fprintln(w, "==========================================")

	// pretty print in table format
	tabWriter := tabwriter.NewWriter(w, 0, 0, 8, ' ', tabwriter.TabIndent)
	fmt.Fprintln(tabWriter, "rank\tword\tfrequency")
	tabWriter.Flush()
	fmt.Fprintln(w, "------------------------------------------")
	for i := 0; i < topN && i < len(s.WordRank); i++ {
		rankVal := s.WordRank[i]
		fmt.Fprintf(tabWriter, "%d\t%s\t%d\n", i+1, rankVal.Word, rankVal.Freq)
	}
	tabWriter.Flush()
	fmt.Fprintln(w, "------------------------------------------")
}
