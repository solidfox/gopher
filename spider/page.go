// Author Daniel Schlaug
// Written at Hong Kong University of Science and Technology in 2013

package spider

import (
	"gopher/stemmer"
	"io/ioutil"
	"math"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	DefaultPositionsLength = 10
	DefaultLinksLength     = 10
	DefaultWordsLength     = 10
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
	PageID       int64
	wordMap      map[string]*Word
	words        []*Word
	vectorLength float64
	wordCount    int
	links        []*Link
	TFMax        int
	Size         int64
	Title        string
	URL          string
	Modified     time.Time
	Parents      []*Link
	Children     []*Link
	wordValid    func(string) bool
}

var stopwords []string

func (p *Page) VectorLength() float64 {
	if p.vectorLength != 0 {
		return p.vectorLength
	}

	sum := int64(0)
	for _, word := range p.words {
		tf := int64(word.TF())
		sum += tf * tf
	}
	return math.Sqrt(float64(sum))
}

func stopwordFilter() func(string) bool {
	if len(stopwords) == 0 {
		stopwordFile, _ := ioutil.ReadFile("stopwords.txt")
		stopwords := sort.StringSlice(strings.Fields(string(stopwordFile)))
		stopwords.Sort()
	}

	return func(word string) bool {
		index := sort.SearchStrings(stopwords, word)
		return index >= len(stopwords) || stopwords[index] != word
	}
}

func NewPage() *Page {
	return &Page{
		PageID:    -1,
		wordMap:   make(map[string]*Word),
		words:     make([]*Word, 0, DefaultWordsLength),
		wordCount: 0,
		links:     make([]*Link, 0, DefaultLinksLength),
		wordValid: stopwordFilter(),
	}
}

func (p *Page) getWord(wordString string) *Word {
	if len(p.wordMap) == len(p.words) {
		word, exists := p.wordMap[wordString]
		if exists {
			return word
		}
	} else {
		for _, word := range p.words {
			if word.Word == wordString {
				return word
			}
		}
	}
	return NewWord("")
}

func (p *Page) TFxIDF(wordString string, df int, n int64) float64 {
	tf := float64(p.getWord(wordString).TF())
	return tf * math.Log2(float64(n)/float64(df))
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
	p.addWordControlled(word)
}

func (p *Page) initWordMap() {
	for _, word := range p.words {
		p.wordMap[word.Word] = word
	}
}

func (p *Page) addWordControlled(word string) {
	if len(p.words) != 0 && len(p.wordMap) == 0 {
		p.initWordMap()
	} else if len(p.words) != len(p.wordMap) {
		panic("Page is in an illegal state; wordMap out of sync")
	}
	wordObj, exists := p.wordMap[word]
	position := p.wordCount + 1
	if !exists {
		wordObj = NewWord(word)
		p.words = append(p.words, wordObj)
	}
	wordObj.positions = append(wordObj.positions, position)
	if wordObj.TF() > p.TFMax {
		p.TFMax = wordObj.TF()
	}
	p.wordMap[word] = wordObj
	p.wordCount++
}

// Set words without controlling for duplicates
func (p *Page) SetWords(words []*Word) {
	p.words = words
}

// Add word without controlling for duplicates
func (p *Page) AddWord(word *Word) {
	p.words = append(p.words, word)
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
	return p.words
}

func (p *Page) Links() []*Link {
	return p.links
}
