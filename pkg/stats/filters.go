package stats

// FilterShort removes word with less than 5 characters
func FilterShort(word string) bool {
	return len(word) > 4
}

// FilterMedium removes words with less than 57characters
func FilterMedium(word string) bool {
	return len(word) > 7
}

// FilterListInclude includes the words in words, but excludes all other words
func FilterListInclude(words ...string) FilterFunc {
	return func(word string) bool {
		for _, w := range words {
			if w == word {
				return true
			}
		}
		return false
	}
}

// FilterListExclude excludes all words in words
func FilterListExclude(words ...string) FilterFunc {
	return func(word string) bool {
		for _, w := range words {
			if w == word {
				return false
			}
		}
		return true
	}
}
