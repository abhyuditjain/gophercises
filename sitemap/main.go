package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/abhyuditjain/gophercices/link"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "URL to build sitemap for")
	maxDepth := flag.Int("depth", 10, "the maximum number of links deep to traverse")
	flag.Parse()

	pages := bfs(*urlFlag, *maxDepth)

	toXml := urlset{
		Urls:  make([]loc, len(pages)),
		Xmlns: xmlns,
	}
	for i, page := range pages {
		toXml.Urls[i] = loc{Value: page}
	}

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	fmt.Print(xml.Header)
	if err := enc.Encode(toXml); err != nil {
		panic(err)
	}
	fmt.Println()
}

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()

	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	var ret []string
	links, _ := link.Parse(r)
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, strings.Join([]string{base, l.Href}, ""))
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}

	return ret
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string

	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}

	return ret
}

func withPrefix(prefix string) func(string) bool {
	return func(s string) bool {
		return strings.HasPrefix(s, prefix)
	}
}

func bfs(urlStr string, maxDepth int) []string {
	seen := make(map[string]struct{}) // set

	// Queue of urls to fetch
	var q map[string]struct{}
	// Queue to keep the urls parsed from urls in q
	nq := map[string]struct{}{
		urlStr: {},
	}

	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})
		if len(q) == 0 {
			break
		}
		for url := range q {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			for _, link := range get(url) {
				if _, ok := seen[link]; !ok {
					nq[link] = struct{}{}
				}
			}
		}
	}
	var ret []string
	for url := range seen {
		ret = append(ret, url)
	}
	return ret
}
