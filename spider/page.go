// Author Daniel Schlaug
// Written at Hong Kong University of Science and Technology in 2013

package spider

import (
	"gopher/stemmer"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	DefaultPositionsLength = 10
	DefaultLinksLength
)

type RawWord []byte
type positionList []int
type Link struct {
	URL        string
	AnchorText string
}

// Size is the size of the page in bytes.
// Modified is the Last-Modified field or Date field of the page's http header.
type Page struct {
	words     map[string]positionList
	wordCount int
	links     []Link
	Size      int64
	Title     string
	URL       string
	Modified  time.Time
	wordValid func(string) bool
}

func NewPage(url string, modified time.Time, wordValid func(string) bool) *Page {
	url = strings.Trim(url, "/")
	return &Page{
		make(map[string]positionList),
		0,
		make([]Link, 0, DefaultLinksLength),
		0,
		"",
		url,
		modified,
		wordValid,
	}
}

// Sequentially adds word to the list of words for the page.
// The word will be assigned the next position.
// If the word already exists only the position is added.
// Duplicates and stopwords are ignored.
func (p *Page) AddWord(word string) {
	newWord := strings.TrimSpace(word)
	if newWord == "" {
		return
	}
	newWord = strings.ToLower(newWord)
	if !p.wordValid(newWord) {
		return
	}
	newWord = string(stemmer.Stem([]byte(newWord)))

	positions, exists := p.words[newWord]
	position := p.wordCount + 1
	if !exists {
		positions = make(positionList, 1, DefaultPositionsLength)
		positions[0] = position
	} else {
		positions = append(positions, position)
	}
	p.words[newWord] = positions
	p.wordCount++
}

func (p *Page) AddText(text string) {
	text = strings.TrimSpace(text)
	whiteSpace, _ := regexp.Compile("[[:space:][:punct:][:cntrl:]]+")
	words := whiteSpace.Split(text, -1)
	for _, word := range words {
		p.AddWord(word)
	}
}

// Adds links with their associated anchor text.
func (p *Page) AddLink(relativeURL string, text string) {
	baseURI, err1 := url.Parse(p.URL)
	parsedRelativeURL, err2 := url.Parse(relativeURL)
	if err1 != nil || err2 != nil {
		return
	}
	newURI := baseURI.ResolveReference(parsedRelativeURL)
	if strings.HasPrefix(newURI.String(), "http") {
		p.links = append(p.links, Link{newURI.String(), text})
	}
}

func (p *Page) Words() []Word {
	n := len(p.words)
	wordSlice := make([]Word, n)
	i := 0
	for word, positions := range p.words {
		wordSlice[i] = Word{word, positions}
		i++
	}
	return wordSlice
}

func (p *Page) Links() []Link {
	return p.links
}

type Word struct {
	Word      string
	positions []int
}

func (w *Word) Positions() []int {
	return w.positions
}

func (w *Word) TF() int {
	return len(w.positions)
}
