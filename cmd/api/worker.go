package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/StelIify/feedbland/internal/database"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// todo some feeds could have image field in the channel, if image url != nil we save it to the db
func (app *App) fetchFeedsWorker(concurrentWorkers int, fetchInterval time.Duration) {
	app.infoLog.Printf("Fetching feeds on %d workers every %v duration", concurrentWorkers, fetchInterval)

	ticker := time.NewTicker(fetchInterval)
	for ; ; <-ticker.C {
		feeds, err := app.db.GenerateNextFeedsToFetch(context.Background(), int32(concurrentWorkers))
		if err != nil {
			app.errorLog.Println(err)
			continue
		}
		for _, feed := range feeds {
			app.wg.Add(1)
			go app.fetchFeed(feed)
		}
		app.wg.Wait()
	}
}

func (app *App) fetchFeed(feed database.GenerateNextFeedsToFetchRow) {
	defer app.wg.Done()

	rssFeed, err := UrlToFeed(feed.Url)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.infoLog.Printf("on feed %s found %d posts", feed.Name, len(rssFeed.Channel.Item))
	err = app.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	for _, item := range rssFeed.Channel.Item {
		pubDate, err := parseDate(item.PubDate)
		if err != nil {
			app.errorLog.Println(err)
			continue
		}
		err = app.db.CreatePost(context.Background(), database.CreatePostParams{
			FeedID:      feed.ID,
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: pubDate,
		})
		if err != nil {
			var pg_err *pgconn.PgError
			if !errors.As(err, &pg_err) && pg_err.Code == pgerrcode.UniqueViolation {
				app.errorLog.Println(err)
				continue
			}
		}
	}
}

func parseDate(date string) (time.Time, error) {
	layout1 := time.RFC1123Z
	pubDate, err := time.Parse(layout1, date)
	if err == nil {
		return pubDate, nil
	}

	layout2 := "Mon, 02 Jan 2006 15:04:05 MST"

	pubDate, err = time.Parse(layout2, date)
	if err == nil {
		return pubDate, nil
	}

	return time.Time{}, err
}

func (app *App) fetchImage() {
	response, err := http.Get("https://blog.boot.dev/")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	baseURL, err := url.Parse(response.Request.URL.String())
	if err != nil {
		log.Fatal(err)
	}

	// Find all image elements in the document
	doc.Find("img").Each(func(index int, element *goquery.Selection) {
		imageURL, exists := element.Attr("src")
		if exists {
			absoluteURL := resolveURL(baseURL, imageURL)
			if strings.Contains(strings.ToLower(absoluteURL), "logo") {
				fmt.Println("Image URL:", absoluteURL)
				response, err := http.Get(absoluteURL)
				if err != nil {
					fmt.Println("get request error", err)
				}
				defer response.Body.Close()
				imgBytes, err := io.ReadAll(response.Body)
				if err != nil {
					fmt.Println("read response body error", err)
				}
				parsedUrl, err := url.Parse(absoluteURL)
				if err != nil {
					fmt.Println("url parsing error", err)
				}
				filePath := path.Base(parsedUrl.Path)
				result, err := app.uploader.Upload(context.TODO(), &s3.PutObjectInput{
					Bucket: aws.String("feebland"),
					Key:    aws.String(filePath),
					Body:   bytes.NewReader(imgBytes),
					ACL:    "public-read",
				})
				if err != nil {
					fmt.Println("file upload error", err)
				}
				fmt.Println(result)
			}
		}
	})
}
