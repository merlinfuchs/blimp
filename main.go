package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog/log"

	"github.com/merlinfuchs/blimp/internal"
	"github.com/merlinfuchs/blimp/internal/config"
	"github.com/merlinfuchs/blimp/internal/views/latency"
	"github.com/merlinfuchs/blimp/internal/views/status"
	"github.com/merlinfuchs/blimp/internal/views/weather"
	"github.com/rivo/tview"
)

func main() {
	config.InitConfig()

	app := tview.NewApplication()

	views := map[string]internal.View{
		"latency": latency.New(),
		"status":  status.New(),
		"weather": weather.New(),
	}

	layout := make([][]string, 0)
	if err := config.K.Unmarshal("layout", &layout); err != nil {
		log.Panic().Err(err).Msgf("Failed to unmarshal layout")
	}

	rowValues := make([]int, len(layout))
	for i := 0; i < len(rowValues); i++ {
		rowValues[i] = -1
	}

	colValues := make([]int, len(layout[0]))
	for i := 0; i < len(colValues); i++ {
		colValues[i] = -1
	}

	grid := tview.NewGrid().
		SetGap(1, 2).
		SetRows(rowValues...).
		SetColumns(colValues...)

	for viewName, view := range views {
		found := false

		minRow := -1
		maxRow := -1
		minCol := -1
		maxCol := -1

		for r, cols := range layout {
			for c, name := range cols {
				if name == viewName {
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

			grid.AddItem(view.Primitive(), minRow, minCol, rowSpan, colSpan, 0, 0, false)

			view.Start()
			defer view.Stop()
		}
	}

	frame := tview.NewFrame(grid).
		SetBorders(2, 2, 1, 0, 4, 4)

	go func() {
		for {
			frame.Clear()

			now := time.Now()
			frame.
				AddText("Blimp v0.1.0", true, tview.AlignLeft, tcell.ColorDimGray).
				AddText(now.Format("15:04:05"), true, tview.AlignCenter, tcell.ColorLightGray).
				AddText(now.Format("Monday, January 2, 2006"), true, tview.AlignCenter, tcell.ColorDimGray)

			app.QueueUpdateDraw(func() {
				for _, view := range views {
					view.Update()
				}
			})

			<-time.After(time.Duration(config.K.Int("update_interval")) * time.Millisecond)
		}
	}()

	if err := app.SetRoot(frame, true).EnableMouse(false).Run(); err != nil {
		log.Fatal().Err(err).Msgf("Failed to run app")
	}
}
