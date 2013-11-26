// Author Daniel Schlaug
// Written at Hong Kong University of Science and Technology in 2013

package spider

import (
	"gopher/stemmer"
	"io/ioutil"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	DefaultPositionsLength = 10
	DefaultLinksLength     = 10
)

type positionList []int
type Link struct {
	Title      string
	URL        string
	AnchorText string
}

// PageID is a database's unique integer identifier for the page. If none is set the value is -1.
// Size is the size of the page in bytes.
// Modified is the Last-Modified field or Date field of the page's http header.
type Page struct {
	PageID    int64
	words     map[string]*Word
	wordCount int
	links     []*Link
	Size      int64
	Title     string
	URL       string
	Modified  time.Time
	Parents   []*Link
	Children  []*Link
	wordValid func(string) bool
}

var stopwords []string

func stopwordFilter() func(string) bool {
	if len(stopwords) == 0 {
		stopwordFile, _ := ioutil.ReadFile("stopwords.txt")
		stopwords = sort.StringSlice(strings.Fields(string(stopwordFile)))
	}

	return func(word string) bool {
		index := sort.SearchStrings(stopwords, word)
		return index >= len(stopwords) || stopwords[index] != word
	}
}

func NewPage() *Page {
	return &Page{
		PageID:    -1,
		words:     make(map[string]*Word),
		wordCount: 0,
		links:     make([]*Link, 0, DefaultLinksLength),
		wordValid: stopwordFilter(),
	}
}

// Sequentially adds word to the list of words for the page.
// The word will be assigned the next position.
// If the word already exists only the position is added.
// Duplicates and stopwords are ignored.
func (p *Page) addWordFromText(word string) {
	word = strings.TrimSpace(word)
	if word == "" {
		return
	}
	word = strings.ToLower(word)
	if !p.wordValid(word) {
		return
	}
	word = string(stemmer.Stem([]byte(word)))

	wordObj, exists := p.words[word]
	position := p.wordCount + 1
	if !exists {
		wordObj = NewWord(word)
	}
	wordObj.positions = append(wordObj.positions, position)
	p.words[word] = wordObj
	p.wordCount++
}

// A
func (p *Page) AddWords(words []*Word) {
	for _, word := range words {
		p.words[word.Word] = word
	}
}

// Adds all words in the text, separated by any match of "[[:space:][:punct:][:cntrl:]]+"
func (p *Page) AddText(text string) {
	text = strings.TrimSpace(text)
	whiteSpace, _ := regexp.Compile("[[:space:][:punct:][:cntrl:]]+")
	words := whiteSpace.Split(text, -1)
	for _, word := range words {
		p.addWordFromText(word)
	}
}

func (p *Page) AddQueryWord(word string) {
	word = strings.TrimSpace(word)
	if word == "" {
		return
	}
	word = strings.ToLower(word)
	words := strings.Fields(word)

	for _, aWord := range words {
		if !p.wordValid(aWord) {
			aWord = ""
		}

		aWord = string(stemmer.Stem([]byte(aWord)))

	}

	word = strings.Join(words, " ")

	wordObj, exists := p.words[word]
	position := p.wordCount + 1
	if !exists {
		wordObj = NewWord(word)
	}
	wordObj.positions = append(wordObj.positions, position)
	p.words[word] = wordObj
	p.wordCount++
}

// Adds links with their associated anchor text.
func (p *Page) AddLink(relativeURL string, text string) {
	baseURI, err1 := url.Parse(p.URL)
	parsedRelativeURL, err2 := url.Parse(relativeURL)
	if err1 != nil || err2 != nil {
		return
	}
	newURI := baseURI.ResolveReference(parsedRelativeURL)
	if isValidHtmlLink(newURI.String()) {
		p.links = append(p.links, &Link{URL: newURI.String(), AnchorText: text})
	}
}

func isValidHtmlLink(link string) bool {
	return strings.HasPrefix(link, "http") &&
		!strings.HasSuffix(link, ".pdf") &&
		!strings.HasSuffix(link, ".doc") &&
		!strings.HasSuffix(link, ".docx") &&
		!strings.HasSuffix(link, ".ppt") &&
		!strings.HasSuffix(link, ".jpg") &&
		!strings.HasSuffix(link, ".bmp") &&
		!strings.HasSuffix(link, ".png")
}

func (p *Page) Words() []*Word {
	n := len(p.words)
	wordSlice := make([]*Word, n)
	i := 0
	for _, wordObj := range p.words {
		wordSlice[i] = wordObj
		i++
	}
	return wordSlice
}

func (p *Page) Links() []*Link {
	return p.links
}
