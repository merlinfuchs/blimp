package feeds

import "github.com/mmcdole/gofeed"

type FeedTarget struct {
	URL string `koanf:"url"`
}

type FeedItem struct {
	Feed *gofeed.Feed
	Item *gofeed.Item
}
