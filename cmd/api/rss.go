package main

import (
	"encoding/xml"
	"net/http"
	"time"
)

type RssFeed struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Atom    string   `xml:"atom,attr"`
	Channel struct {
		Text  string `xml:",chardata"`
		Title string `xml:"title"`
		Link  struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Description   string `xml:"description"`
		Generator     string `xml:"generator"`
		Language      string `xml:"language"`
		LastBuildDate string `xml:"lastBuildDate"`
		Item          []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
			Guid        string `xml:"guid"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

func urlToFeed(url string) (RssFeed, error) {
	client := http.Client{Timeout: time.Second * 20}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RssFeed{}, err
	}
	response, err := client.Do(req)
	if err != nil {
		return RssFeed{}, err
	}
	defer response.Body.Close()

	var rssFeed RssFeed

	err = xml.NewDecoder(response.Body).Decode(&rssFeed)
	if err != nil {
		return RssFeed{}, err
	}
	return rssFeed, nil
}
