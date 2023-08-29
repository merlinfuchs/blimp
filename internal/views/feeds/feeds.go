package feeds

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"

	"github.com/gdamore/tcell/v2"
	"github.com/merlinfuchs/blimp/internal/config"
	"github.com/rivo/tview"
)

type FeedsView struct {
	stopped chan struct{}
	view    *tview.Flex
	targets []FeedTarget
	items   []FeedItem
}

func New() *FeedsView {
	targets := make([]FeedTarget, 0)
	if err := config.K.Unmarshal("views.feeds.targets", &targets); err != nil {
		log.Panic().Err(err).Msgf("Failed to unmarshal status targets")
	}

	view := tview.NewFlex().SetDirection(tview.FlexRow)
	view.SetBorder(true).
		SetBorderColor(tcell.ColorGray).
		SetTitle("Feeds").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 1, 1)

	return &FeedsView{
		stopped: make(chan struct{}),
		view:    view,
		targets: targets,
	}
}

func (l *FeedsView) Start() {
	go func() {
		l.updateItems()

		for {
			select {
			case <-l.stopped:
				return
			case <-time.After(time.Duration(config.K.Int("views.feeds.update_interval")) * time.Millisecond):
				l.updateItems()
			}
		}
	}()
}

func (l *FeedsView) updateItems() {
	fp := gofeed.NewParser()

	items := make([]FeedItem, 0)
	for _, target := range l.targets {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		feed, err := fp.ParseURLWithContext(target.URL, ctx)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch feed %s", target.URL)
			continue
		}

		for _, item := range feed.Items {
			items = append(items, FeedItem{
				Feed: feed,
				Item: item,
			})
		}
	}

	sort.Slice(items, func(i, j int) bool {
		publishedI := items[i].Item.PublishedParsed
		publishedJ := items[j].Item.PublishedParsed
		if publishedI == nil || publishedJ == nil {
			log.Warn().Msgf("Failed to sort feeds, missing published date")
			return true
		}
		return publishedI.After(*publishedJ)
	})

	l.items = items
}

func (l *FeedsView) Stop() {
	close(l.stopped)
}

func (l *FeedsView) Update() error {
	l.view.Clear()

	items := l.items
	maxItems := config.K.Int("views.feeds.max_items")
	if maxItems == 0 {
		_, _, _, height := l.view.GetRect()
		maxItems = height - 4
	}

	if len(items) > maxItems {
		items = items[:maxItems]
	}

	for _, item := range items {
		text := "[gray]-"
		if config.K.Bool("views.feeds.show_published_time") {
			published := item.Item.PublishedParsed
			if published != nil {
				text += fmt.Sprintf(" [yellowgreen]%s", published.Format("2006-01-02 15:04"))
			}
		}

		if config.K.Bool("views.feeds.show_feed_title") {
			text += fmt.Sprintf(" [gray]%s", item.Feed.Title)
		}

		text += fmt.Sprintf(" [white]%s", item.Item.Title)

		l.view.AddItem(
			tview.NewTextView().
				SetText(text).
				SetDynamicColors(true).
				SetWrap(false),
			1, 1, false)
	}

	return nil
}

func (l *FeedsView) Primitive() tview.Primitive {
	return l.view
}
