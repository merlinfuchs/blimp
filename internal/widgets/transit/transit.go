package transit

import (
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/merlinfuchs/blimp/internal/config"
	"github.com/rivo/tview"
)

type TransitView struct {
	view   *tview.Flex
	ticker *time.Ticker
	items  []DepartureItem
}

func New() *TransitView {
	view := tview.NewFlex().SetDirection(tview.FlexRow)
	view.SetBorder(true).
		SetBorderColor(tcell.ColorGray).
		SetTitle("Next Departures").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 1, 1)

	return &TransitView{
		view: view,
	}
}

func (l *TransitView) Start() {
	l.ticker = time.NewTicker(1 * time.Minute)
	err := l.updateItems()
	if err != nil {
		slog.With("error", err).Error("Failed to update items")
		panic(err)
	}
}

func (l *TransitView) Stop() {
	if l.ticker != nil {
		l.ticker.Stop()
	}
}

func (l *TransitView) Update() error {
	if l.ticker == nil {
		return nil
	}

	select {
	case <-l.ticker.C:
		err := l.updateItems()
		if err != nil {
			return fmt.Errorf("failed to update items: %w", err)
		}
	default:
	}

	l.updateView()
	return nil
}

type DepartureItem struct {
	PlaceName string
	Name      string
	HeadSign  string
	ShortName string
	Platform  string
	Category  string
	Departure time.Time
}

func (l *TransitView) updateItems() error {
	nextDepartures, err := getNextDeparturesData()
	if err != nil {
		return fmt.Errorf("failed to get next departures data: %w", err)
	}

	items := make([]DepartureItem, 0)
	for _, board := range nextDepartures.Boards {
		for _, departure := range board.Departures {
			items = append(items, DepartureItem{
				PlaceName: board.Place.Name,
				Name:      departure.Transport.Name,
				HeadSign:  departure.Transport.HeadSign,
				ShortName: departure.Transport.ShortName,
				Platform:  departure.Platform,
				Category:  departure.Transport.Category,
				Departure: departure.Time,
			})
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Departure.Before(items[j].Departure)
	})
	l.items = items

	return nil
}

func (l *TransitView) updateView() error {
	l.view.Clear()

	maxItems := config.K.Int("widgets.transit.max_items")
	if maxItems == 0 {
		_, _, _, height := l.view.GetRect()
		maxItems = height - 4
	}

	items := l.items
	if len(items) > maxItems {
		items = items[:maxItems]
	}

	for _, item := range items {
		departureMinutes := item.Departure.Sub(time.Now()).Minutes()
		if departureMinutes < 0 {
			departureMinutes = 0
		}

		l.view.AddItem(
			tview.NewTextView().
				SetText(fmt.Sprintf(
					"[gray]- [yellowgreen]%d min (%s) [white]%s [lightgray]%s [gray]%s",
					int(departureMinutes),
					item.Departure.Format("15:04"),
					item.Name,
					item.HeadSign,
					item.PlaceName,
				)).
				SetDynamicColors(true).
				SetWrap(false),
			1, 1, false)
	}

	return nil
}

func (l *TransitView) Primitive() tview.Primitive {
	return l.view
}
