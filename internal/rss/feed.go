package rss

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type Feed struct {
	Channel struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Item        []Item `xml:"item"`
	} `xml:"channel"`
}

func FetchFeed(ctx context.Context, feedUrl string) (*Feed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var feed Feed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for idx, _ := range feed.Channel.Item {
		feed.Channel.Item[idx].Description = html.UnescapeString(feed.Channel.Item[idx].Description)
		feed.Channel.Item[idx].Title = html.UnescapeString(feed.Channel.Item[idx].Title)
	}
	return &feed, nil
}
