// Author Daniel Schlaug
// Written at Hong Kong University of Science and Technology in 2013

package spider

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"container/list"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	ChannelBuffer    = 100
	MaxBytes         = 1000000
	FetchRoutines    = 1
	ConnectTimeout   = 30
	ReadWriteTimeout = 30
	RootPage         = "http://www.cse.ust.hk/~ericzhao/COMP4321/TestPages/testpage.htm"
)

// TODO dnsCache does two things. Should be split into a URLKeeper and a DNSCache.
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
		d.urlList.PushBack(pageUrl)
	}
}

func (d *dnsCache) hasURL() bool {
	return d.urlList.Len() > 0
}

// Returns the passed url with the hostname switched for its IP
func (d *dnsCache) lookupURL(hostURL string) string {
	urlObj, _ := url.Parse(hostURL)
	hostparts := strings.Split(urlObj.Host, ":")
	ip, exists := d.dnsMap[hostparts[0]]
	if !exists {
		ipList, err := net.LookupIP(hostparts[0])
		if err != nil {
			return hostURL
		}
		ip = ipList[0].String()
		d.dnsMap[hostparts[0]] = ip
	}
	hostparts[0] = ip
	urlObj.Host = strings.Join(hostparts, ":")
	return urlObj.String()
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

func Get300Pages() []*Page {
	pages := make([]*Page, 100)
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
					log.Println("Reached the end of the internet. (No new links to follow.)")
					return
				}
			}

		}
	}()
	for i := FetchRoutines; i > 0; i-- {
		go func() {
			for pageURL := range fetchChannel {
				page := parsePage(pageURL, urls)
				if page == nil {
					continue
				}
				page.URL = pageURL
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

func parsePage(pageUrl string, dns *dnsCache) *Page {
	client := newTimeoutClient(ConnectTimeout*time.Second, ReadWriteTimeout*time.Second)
	// log.Println(" " + pageUrl + " pre")
	ipURL := dns.lookupURL(pageUrl)
	res, err := client.Get(ipURL)
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
			log.Println("Skipping page larger than " + strconv.FormatInt(body.ByteCount(), 10) + ": " + pageUrl)
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
