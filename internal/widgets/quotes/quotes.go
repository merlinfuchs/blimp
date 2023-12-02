package quotes

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/merlinfuchs/blimp/internal/config"
	"github.com/rivo/tview"
)

type QuotesView struct {
	view         *tview.Flex
	ticker       *time.Ticker
	currentQuote QuoteData
}

func New() *QuotesView {
	view := tview.NewFlex().SetDirection(tview.FlexRow)
	view.SetBorder(true).
		SetBorderColor(tcell.ColorGray).
		SetTitle("Random Quote").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 1, 1)

	return &QuotesView{
		view: view,
	}
}

func (l *QuotesView) Start() {
	l.ticker = time.NewTicker(time.Duration(config.K.Int("widgets.quotes.update_interval")) * time.Millisecond)
	err := l.updateQuote()
	if err != nil {
		slog.With("error", err).Error("Failed to update current quote")
		panic(err)
	}
}

func (l *QuotesView) updateQuote() error {
	quote, err := getRandomQuote(
		config.K.Strings("widgets.quotes.tags"),
		config.K.Bool("widgets.quotes.tags_require_all"),
	)
	if err != nil {
		return err
	}

	l.currentQuote = quote
	return nil
}

func (l *QuotesView) Stop() {
	if l.ticker != nil {
		l.ticker.Stop()
	}
}

func (l *QuotesView) Update() error {
	if l.ticker == nil {
		return nil
	}

	select {
	case <-l.ticker.C:
		err := l.updateQuote()
		if err != nil {
			return fmt.Errorf("failed to update items: %w", err)
		}
	default:
	}

	l.updateView()
	return nil
}

func (l *QuotesView) updateView() {
	l.view.Clear()

	l.view.
		AddItem(tview.NewTextView().
			SetText("[yellowgreen]"+l.currentQuote.Content).
			SetDynamicColors(true).
			SetWrap(true),
			0, 10, false).
		AddItem(tview.NewTextView().
			SetText("[gray]"+l.currentQuote.Author).
			SetDynamicColors(true).
			SetWrap(false).
			SetTextAlign(tview.AlignRight),
			1, 1, false)
}

func (l *QuotesView) Primitive() tview.Primitive {
	return l.view
}
