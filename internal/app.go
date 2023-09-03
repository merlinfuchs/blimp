package internal

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/merlinfuchs/blimp/internal/config"
	"github.com/merlinfuchs/blimp/internal/widgets/feeds"
	"github.com/merlinfuchs/blimp/internal/widgets/latency"
	"github.com/merlinfuchs/blimp/internal/widgets/status"
	"github.com/merlinfuchs/blimp/internal/widgets/transit"
	"github.com/merlinfuchs/blimp/internal/widgets/weather"
	"github.com/rivo/tview"
)

func AppEntry() error {
	app := tview.NewApplication()

	widgets := map[string]Widget{
		"latency": latency.New(),
		"status":  status.New(),
		"weather": weather.New(),
		"feeds":   feeds.New(),
		"transit": transit.New(),
	}

	pages, err := parsePagesFromConfig()
	if err != nil {
		return fmt.Errorf("failed to parse pages: %w", err)
	}

	if len(pages) == 0 {
		return fmt.Errorf("no pages defined")
	}

	pagesView, err := constructPages(pages, widgets)
	if err != nil {
		return fmt.Errorf("failed to construct pages: %w", err)
	}
	currentPage := 0

	frame := tview.NewFrame(pagesView).
		SetBorders(2, 2, 1, 0, 4, 4)

	go func() {
		for {
			frame.Clear()

			now := time.Now()
			frame.
				AddText("Blimp v0.1.0", true, tview.AlignLeft, tcell.ColorDimGray).
				AddText(now.Format("15:04:05"), true, tview.AlignCenter, tcell.ColorLightGray).
				AddText(now.Format("Monday, January 2, 2006"), true, tview.AlignCenter, tcell.ColorDimGray).
				AddText(pages[currentPage].Title, true, tview.AlignRight, tcell.ColorDimGray)

			app.QueueUpdateDraw(func() {
				for _, widget := range widgets {
					widget.Update()
				}
			})

			<-time.After(time.Duration(config.K.Int("update_interval")) * time.Millisecond)
		}
	}()

	go func() {
		for {
			<-time.After(time.Duration(config.K.Int("page_interval")) * time.Millisecond)
			currentPage = (currentPage + 1) % len(pages)
			pagesView.SwitchToPage(fmt.Sprintf("%d", currentPage))
		}
	}()

	defer func() {
		for _, widget := range widgets {
			widget.Stop()
		}
	}()

	if err := app.SetRoot(frame, true).EnableMouse(false).Run(); err != nil {
		return fmt.Errorf("failed to run app: %w", err)
	}

	return nil
}

func constructPages(pages []Page, widgets map[string]Widget) (*tview.Pages, error) {
	view := tview.NewPages()

	usedWidgets := make(map[string]bool, len(widgets))

	for i, page := range pages {
		if len(page.Layout) == 0 {
			return nil, fmt.Errorf("page #%d (%s) has no layout", i, page.Title)
		}

		rowValues := make([]int, len(page.Layout))
		for i := 0; i < len(rowValues); i++ {
			rowValues[i] = -1
		}

		colValues := make([]int, len(page.Layout[0]))
		for i := 0; i < len(colValues); i++ {
			colValues[i] = -1
		}

		grid := tview.NewGrid().
			SetGap(1, 2).
			SetRows(rowValues...).
			SetColumns(colValues...)

		for widgetName, widget := range widgets {
			found := false

			minRow := -1
			maxRow := -1
			minCol := -1
			maxCol := -1

			for r, cols := range page.Layout {
				for c, name := range cols {
					if name == widgetName {
						found = true
						maxRow = r
						maxCol = c
						if minRow == -1 {
							minRow = r
						}
						if minCol == -1 {
							minCol = c
						}
					}
				}
			}

			if found {
				rowSpan := maxRow - minRow + 1
				colSpan := maxCol - minCol + 1

				usedWidgets[widgetName] = true
				grid.AddItem(widget.Primitive(), minRow, minCol, rowSpan, colSpan, 0, 0, false)
			}
		}

		view.AddPage(fmt.Sprintf("%d", i), grid, true, i == 0)
	}

	for widgetName, widget := range widgets {
		if usedWidgets[widgetName] {
			widget.Start()
		}
	}

	return view, nil
}
