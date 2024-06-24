package feeds

import (
	"context"
	"database/sql"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mloughton/blog-agg/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(feedURL string) (*RSSFeed, error) {
	resp, err := http.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rssFeed RSSFeed
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return nil, err
	}

	return &rssFeed, nil
}

func StartScraping(db *database.Queries, amount int32, interval time.Duration) {
	t := time.NewTicker(interval)
	for {
		<-t.C
		log.Println("recieved tick")
		fetchFeeds, err := db.GetNextFeedsToFetch(context.Background(), amount)
		if err != nil {
			log.Println(err)
		}
		wg := &sync.WaitGroup{}
		for i, feed := range fetchFeeds {
			wg.Add(1)
			log.Printf("spawning go routine %d", i)
			go scrapeFeed(db, wg, feed)
		}
		log.Printf("waiting for go routines to finish")
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	feedData, err := fetchFeed(feed.Url)
	if err != nil {
		log.Println(err)
	}
	params := database.MarkFeedFetchedParams{
		UpdatedAt:     time.Now(),
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
		ID:            feed.ID,
	}
	err = db.MarkFeedFetched(context.Background(), params)
	if err != nil {
		log.Println(err)
	}
	for _, post := range feedData.Channel.Item {
		newUUID, err := uuid.NewUUID()
		if err != nil {
			log.Println(err)
			continue
		}
		pubDate, err := time.Parse(time.RFC1123Z, post.PubDate)
		if err != nil {
			log.Println(err)
			continue
		}
		params := database.CreatePostParams{
			ID:          newUUID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       post.Title,
			Url:         post.Link,
			Description: post.Description,
			PublishedAt: pubDate,
			FeedID:      feed.ID,
		}
		db.CreatePost(context.Background(), params)
	}
}
