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
		Item  []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
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
