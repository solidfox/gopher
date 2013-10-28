// Author Daniel Schlaug
// Written at Hong Kong University of Science and Technology in 2013

package spider

import (
	"code.google.com/p/go.net/html"
	// "fmt"
	// "io/ioutil"
	"bytes"
	"container/list"
	//"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	//"sort"
	"strconv"
	"strings"
	"time"
)

const (
	ChannelBuffer    = 100
	MaxBytes         = 1000000
	FetchRoutines    = 200
	ConnectTimeout   = 30
	ReadWriteTimeout = 30
	RootPage         = "http://www.cse.ust.hk/"
)

type dnsCache struct {
	isParsed map[string]bool
	dnsMap   map[string]string
	urlList  *list.List
}

func newDNSCache() *dnsCache {
	return &dnsCache{
		make(map[string]bool),
		make(map[string]string),
		list.New(),
	}
}

func (d *dnsCache) addURL(pageUrl string) {
	if !d.isParsed[pageUrl] {
		d.isParsed[pageUrl] = true
		urlObj, _ := url.Parse(pageUrl)
		hostparts := strings.Split(urlObj.Host, ":")
		ip, exists := d.dnsMap[hostparts[0]]
		if !exists {
			rawIp, err := net.LookupIP(hostparts[0])
			if err != nil {
				return
			}
			ip = rawIp[0].String()
			d.dnsMap[hostparts[0]] = ip
		}
		hostparts[0] = ip
		urlObj.Host = strings.Join(hostparts, ":")
		pageUrl = urlObj.String()
		d.urlList.PushBack(pageUrl)
	}
}

func (d *dnsCache) hasURL() bool {
	return d.urlList.Len() > 0
}

func (d *dnsCache) getURL() string {
	first := d.urlList.Remove(d.urlList.Front())
	if theUrl, ok := first.(string); ok {
		return theUrl
	} else {
		log.Fatalln("dnsCache: Tried to get a URL that was not a string.")
	}
	return ""
}

func Get30Pages() []*Page {
	pages := make([]*Page, 30)
	pageChannel := GetPages()
	for i, _ := range pages {
		pages[i] = <-pageChannel
	}
	return pages
}

func GetPages() <-chan *Page {
	urlChannel := make(chan string, ChannelBuffer)
	pageChannel := make(chan *Page, ChannelBuffer)
	fetchChannel := make(chan string, ChannelBuffer)
	urls := newDNSCache()

	go func() {
		for {
			if urls.hasURL() {
				select {
				case pageURL := <-urlChannel:
					urls.addURL(pageURL)
				case fetchChannel <- urls.getURL():
				}
			} else {
				select {
				case newURL := <-urlChannel:
					urls.addURL(newURL)
				case <-time.After((ConnectTimeout + ReadWriteTimeout) * time.Second):
					log.Fatalln("Reached the end of the internet. (No new links to follow.)")
				}
			}

		}
	}()
	for i := FetchRoutines; i > 0; i-- {
		go func() {
			for {
				pageUrl := <-fetchChannel
				page := ParsePage(pageUrl)
				if page == nil {
					continue
				}
				for _, link := range page.Links() {
					urlChannel <- link.URL
				}
				// log.Println("sent")
				pageChannel <- page
			}
		}()
	}

	urlChannel <- RootPage
	return pageChannel
}

func timeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}

func newTimeoutClient(connectTimeout time.Duration, readWriteTimeout time.Duration) *http.Client {

	return &http.Client{
		Transport: &http.Transport{
			Dial: timeoutDialer(connectTimeout, readWriteTimeout),
		},
	}
}

func ParsePage(pageUrl string) *Page {
	client := newTimeoutClient(ConnectTimeout*time.Second, ReadWriteTimeout*time.Second)
	// log.Println(" " + pageUrl + " pre")
	res, err := client.Get(pageUrl)
	// log.Println(" " + pageUrl + " post")
	if err != nil {
		//log.Println(err)
		return nil
	} else if res.StatusCode != 200 {
		//log.Println(res)
		return nil
	}
	// res.Body != nil when err == nil
	defer res.Body.Close()

	pageDate := time.Now()
	dateHeaderList := res.Header["Last-Modified"]
	if dateHeaderList == nil {
		dateHeaderList = res.Header["Date"]
	}
	if dateHeaderList != nil {
		pageDate, _ = time.Parse("Mon, 2 Jan 2006 15:04:05 MST", dateHeaderList[0])
	}
	p := NewPage()
	p.URL = pageUrl
	p.Modified = pageDate

	// log.Println("Will trigger tokenizer.")
	body := NewCountingReader(res.Body)
	z := html.NewTokenizer(body)
	// log.Println("Tokenizer triggered.")
	for {
		if body.ByteCount() > MaxBytes {
			log.Println("Skipping page larger than " + MaxBytes + ": " + strconv.FormatInt(body.ByteCount(), 10) + " (" + pageUrl + ")")
			return nil
		}
		token := z.Next()
		switch token {
		case html.ErrorToken:
			p.Size = body.ByteCount()
			// log.Println(" " + pageUrl + " success")
			return p
		case html.StartTagToken:
			tagName, hasAttr := z.TagName()
			switch string(tagName) {
			case "script":
				burnTokensUntilEndTag(z, "script")
			case "a":
				href := []byte("href")
				var rawAnchorText []byte
				var rawUrl string
				for hasAttr {
					key, val, moreAttr := z.TagAttr()
					hasAttr = moreAttr
					if bytes.Equal(key, href) {
						rawUrl = string(val)
						break
					}
				}
				rawAnchorText = textUpToEndTag(z, "a")
				anchorText := string(rawAnchorText)
				p.AddLink(rawUrl, anchorText)
				p.AddText(anchorText)
			case "title":
				titleText := string(textUpToEndTag(z, "title"))
				p.Title = titleText
				p.AddText(titleText)
			}
		case html.TextToken:
			p.AddText(string(z.Text()))
		}
	}
}

func textUpToEndTag(tokenizer *html.Tokenizer, tagName string) []byte {
	var textBuffer bytes.Buffer
	rawTagName := []byte(tagName)
	for done := false; !done; {
		token := tokenizer.Next()
		switch token {
		case html.TextToken:
			textBuffer.Write(tokenizer.Text())
		case html.EndTagToken:
			name, _ := tokenizer.TagName()
			if bytes.Equal(rawTagName, name) {
				done = true
			}
		case html.ErrorToken:
			done = true
		}
	}
	return textBuffer.Bytes()
}

func burnTokensUntilEndTag(firewood *html.Tokenizer, tagName string) {
	rawTagName := []byte(tagName)
	for {
		token := firewood.Next()
		switch token {
		case html.ErrorToken:
			return
		case html.EndTagToken:
			name, _ := firewood.TagName()
			// log.Println("Struck token " + string(name))
			if bytes.Equal(name, rawTagName) {
				// log.Println("Extinguishing token fire.")
				return
			}
		}
	}
}
