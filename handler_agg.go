package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Overhustler/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err := ScrapFeeds(s, cmd)
		if err != nil {
			return err
		}
	}
}

func ScrapFeeds(s *state, cmd command) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())

	if err != nil {
		return err
	}
	err = s.db.MarFeedFetched(context.Background(), database.MarFeedFetchedParams{LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true}, UpdatedAt: time.Now(), ID: nextFeed.ID})

	if err != nil {
		return err
	}

	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}
	layout := "01/02/2006 3:04:05 PM"
	for i := range feed.Channel.Item {
		pubDate, err := time.Parse(layout, feed.Channel.Item[i].PubDate)
		if err != nil {
			pubDate = time.Now()
		}
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Title: feed.Channel.Item[i].Title, Url: feed.Channel.Item[i].Link,
			Description: feed.Channel.Item[i].Description, PublishedAt: pubDate, FeedID: nextFeed.ID})
		if err != nil && err.Error() != "post with that URL already exists" {
			return err
		}
	}
	return nil
}
