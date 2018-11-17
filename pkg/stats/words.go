package stats

import (
	"bufio"
	"log"
	"sort"
	"strings"
	"unicode"
)

// FilterFunc should return if the word should remain in the list
type FilterFunc func(string) bool

// NormaliseText removes the non letter text and converts to lowercase
func NormaliseText(text string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsSpace(r) {
			return unicode.ToLower(r)
		}
		return -1
	}, text)
}

// Words returns a list of words in the text filtered by each filter function
func Words(text string, filters []FilterFunc) []string {
	var words []string

	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanWords)

SCANWORD:
	for scanner.Scan() {
		word := scanner.Text()
		for _, filter := range filters {
			// skip words that the filter doesn't want
			if !filter(word) {
				continue SCANWORD
			}
		}
		words = append(words, word)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln("failed to scan words: ", err)
	}

	return words
}

// Frequency calculates the frequency of each word in words and adds it to the freq parameter map
func Frequency(words []string, freq map[string]int) {
	for _, word := range words {
		freq[word]++
	}
}

// WordFreq is a key value pair for word and frequency
type WordFreq struct {
	Word string
	Freq int
}

// RankedList is used to hold a ordered array of WordFreq
type RankedList []WordFreq

func (r RankedList) Len() int {
	return len(r)
}

func (r RankedList) Less(i int, j int) bool {
	return r[i].Freq < r[j].Freq
}

func (r RankedList) Swap(i int, j int) {
	r[i], r[j] = r[j], r[i]
}

// Rank orders the list in decending order and returns a RankedList
func Rank(freq map[string]int) RankedList {
	rankValues := make(RankedList, 0, len(freq))
	for k, v := range freq {
		rankValues = append(rankValues, WordFreq{k, v})
	}
	sort.Sort(sort.Reverse(rankValues))
	return rankValues
}
