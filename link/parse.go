package link

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

// Link represents a link (<a href="...">) in an HTML
type Link struct {
	Href string
	Text string
}

// Parse will take in an HTML document and will return a
// slice of links parsed from it
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nodes := linkNodes(doc)
	var links []Link
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}
	return links, nil
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}

	var ret []*html.Node

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}

	return ret
}

func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}

	var ret strings.Builder

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret.WriteString(text(c))
	}

	return strings.Join(strings.Fields(ret.String()), " ")
}

func buildLink(n *html.Node) Link {
	var ret Link
	for _, att := range n.Attr {
		if att.Key == "href" {
			ret.Href = att.Val
			break
		}
	}
	ret.Text = text(n)
	return ret
}