package main

import (
	"encoding/xml"
	"net/http"
	"time"
)

type RssFeed struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Atom    string   `xml:"atom,attr"`
	Channel struct {
		Title       string `xml:"title"`
		Description string `xml:"description"`
		Image       struct {
			URL   string `xml:"url"`
			Title string `xml:"title"`
		} `xml:"image"`
		Item []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

type RssFeedName struct {
	Channel struct {
		Title       string `xml:"title"`
		Description string `xml:"description"`
	} `xml:"channel"`
}

func UrlToFeed(url string) (RssFeed, error) {
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

func UrlToRssFeedName(url string) (RssFeedName, error) {
	client := http.Client{Timeout: time.Second * 20}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RssFeedName{}, err
	}
	response, err := client.Do(req)
	if err != nil {
		return RssFeedName{}, err
	}
	defer response.Body.Close()

	var rssFeedName RssFeedName

	err = xml.NewDecoder(response.Body).Decode(&rssFeedName)
	if err != nil {
		return RssFeedName{}, err
	}
	return rssFeedName, nil
}
