// Copyright (c) 2019 Maximilian Schoenenberg
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"testing"
)

func TestShuffleWord(t *testing.T) {
	word := "Programming"
	shuffledWord := shuffleWord(word)

	wordCharCount := 0
	for _, c := range word {
		wordCharCount += int(c)
	}

	shuffledWordCharCount := 0
	for _, c := range shuffledWord {
		shuffledWordCharCount += int(c)
	}

	if wordCharCount != shuffledWordCharCount {
		t.Error(wordCharCount, "!=", shuffledWordCharCount)
	} else if word == shuffledWord {
		t.Error(word, "==", shuffledWord)
	} else if len(word) != len(shuffledWord) {
		t.Error("You have lost some chars..")
	}
}
