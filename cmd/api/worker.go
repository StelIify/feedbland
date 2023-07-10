package main

import (
	"context"
	"errors"
	"time"

	"github.com/StelIify/feedbland/internal/database"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

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

	rssFeed, err := urlToFeed(feed.Url)
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
