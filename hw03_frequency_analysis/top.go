package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const (
	numTopWords int    = 10
	dash        string = "-"
)

var punctuationRegExp = regexp.MustCompile(`[,.\?\!'"]`)

type wordFreq struct {
	word string
	freq int
}

func calculateWordsFreq(s []string) map[string]int {
	wordsFreqs := make(map[string]int)

	for _, word := range s {
		wordWithoutPunctuation := punctuationRegExp.ReplaceAllString(word, "")
		if word != dash {
			wordsFreqs[strings.ToLower(wordWithoutPunctuation)]++
		}
	}

	return wordsFreqs
}

func sortWordsFreqs(wordsFreqs map[string]int) []wordFreq {
	freqs := make([]wordFreq, 0, len(wordsFreqs))
	for word, freq := range wordsFreqs {
		freqs = append(freqs, wordFreq{
			word: word,
			freq: freq,
		})
	}
	sort.Slice(freqs, func(i, j int) bool {
		if freqs[i].freq > freqs[j].freq {
			return true
		}
		if freqs[i].freq < freqs[j].freq {
			return false
		}

		var result bool
		switch strings.Compare(freqs[i].word, freqs[j].word) {
		case -1:
			result = true
		case 1:
			result = false
		}

		return result
	})

	return freqs
}

func Top10(s string) []string {
	if len(s) == 0 {
		return []string{}
	}
	words := strings.Fields(s)
	wordsFreqs := calculateWordsFreq(words)
	sortedWordFreqs := sortWordsFreqs(wordsFreqs)

	result := make([]string, 0, numTopWords)
	sliceSize := numTopWords
	if len(sortedWordFreqs) < sliceSize {
		sliceSize = len(sortedWordFreqs)
	}
	for _, wordsFreq := range sortedWordFreqs[:sliceSize] {
		result = append(result, wordsFreq.word)
	}
	return result
}
