package tokenizer

import (
	"errors"
	"strings"
)

type StringTokenizer struct {
	currentPosition   int
	newPosition       int
	maxPosition       int
	str               []rune
	delimiters        []rune
	retDelims         bool
	delimsChanged     bool
	maxDelimCodePoint int32
}

func (s *StringTokenizer) setMaxDelimCodePoint() {
	var m int32
	for _, d := range s.delimiters {
		if d > m {
			m = d
		}
	}
	s.maxDelimCodePoint = m
}

func NewStringTokenizerWithRetDelims(str string, delim string, retDelims bool) *StringTokenizer {
	s := &StringTokenizer{
		currentPosition: 0,
		newPosition:     -1,
		maxPosition:     len([]rune(str)),
		str:             []rune(str),
		delimiters:      []rune(delim),
		retDelims:       retDelims,
		delimsChanged:   false,
	}
	s.setMaxDelimCodePoint()
	return s
}

func NewStringTokenizerWithDelim(str string, delim string) *StringTokenizer {
	return NewStringTokenizerWithRetDelims(str, delim, false)
}

func NewStringTokenizer(str string) *StringTokenizer {
	return NewStringTokenizerWithDelim(str, " \t\n\r\f")
}

func (s *StringTokenizer) skipDelimiters(startPos int) int {
	if len(s.delimiters) == 0 {
		panic(errors.New("delimiters is a null string"))
	}
	position := startPos
	for !s.retDelims && position < s.maxPosition {
		r := s.str[position]
		if r > s.maxDelimCodePoint || strings.IndexRune(string(s.delimiters), r) < 0 {
			break
		}
		position++
	}
	return position
}

func (s *StringTokenizer) scanToken(startPos int) int {
	position := startPos
	for position < s.maxPosition {
		r := s.str[position]
		if r <= s.maxDelimCodePoint && strings.IndexRune(string(s.delimiters), r) >= 0 {
			break
		}
		position++
	}
	if s.retDelims && startPos == position {
		r := s.str[position]
		if r <= s.maxDelimCodePoint && strings.IndexRune(string(s.delimiters), r) >= 0 {
			position++
		}
	}
	return position
}

func (s *StringTokenizer) HasMoreTokens() bool {
	s.newPosition = s.skipDelimiters(s.currentPosition)
	return s.newPosition < s.maxPosition
}

func (s *StringTokenizer) NextToken() string {
	if s.newPosition >= 0 && !s.delimsChanged {
		s.currentPosition = s.newPosition
	} else {
		s.currentPosition = s.skipDelimiters(s.currentPosition)
	}
	s.delimsChanged = false
	s.newPosition = -1
	if s.currentPosition >= s.maxPosition {
		panic(errors.New("no such element"))
	}
	start := s.currentPosition
	s.currentPosition = s.scanToken(s.currentPosition)
	return string(s.str[start:s.currentPosition])
}

func (s *StringTokenizer) NextTokenByNewDelim(delim string) string {
	s.delimiters = []rune(delim)
	s.delimsChanged = true
	s.setMaxDelimCodePoint()
	return s.NextToken()
}

func (s *StringTokenizer) CountTokens() int {
	count := 0
	currpos := s.currentPosition
	for currpos < s.maxPosition {
		currpos = s.skipDelimiters(currpos)
		if currpos >= s.maxPosition {
			break
		}
		currpos = s.scanToken(currpos)
		count++
	}
	return count
}
