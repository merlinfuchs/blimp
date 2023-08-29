package weather

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/merlinfuchs/blimp/internal/config"
	"github.com/merlinfuchs/blimp/internal/widgets/weather/icons"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

type StatusView struct {
	stopped         chan struct{}
	view            *tview.Flex
	currentWeather  *CurrentWeatherData
	forecastWeather *ForecastWeatherData
}

func New() *StatusView {
	view := tview.NewFlex().SetDirection(tview.FlexColumn)
	view.
		SetBorderColor(tcell.ColorGray).
		SetTitle("Weather").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 2, 2)

	return &StatusView{
		stopped: make(chan struct{}),
		view:    view,
	}
}

func (l *StatusView) Stop() {
	close(l.stopped)
}

func (l *StatusView) Primitive() tview.Primitive {
	return l.view
}

func (l *StatusView) Start() {
	l.updateData()

	go func() {
		for {
			select {
			case <-l.stopped:
				break
			case <-time.After(time.Duration(config.K.Int("widgets.weather.update_interval")) * time.Millisecond):
				l.updateData()
			}
		}
	}()
}

func (l *StatusView) updateData() {
	var err error

	currentWeather, err := getCurrentWeatherData()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get current weather data")
		return
	}
	l.currentWeather = &currentWeather

	forecastWeather, err := getWeatherForecast()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get forecast weather data")
		return
	}
	l.forecastWeather = &forecastWeather

	l.updateView()
}

func (l *StatusView) updateView() error {
	l.view.Clear()

	unitText := ""
	switch config.K.String("widgets.weather.owm_unit") {
	case "metric":
		unitText = "°C"
	case "imperial":
		unitText = "°F"
	case "standard":
		unitText = "°K"
	}

	if l.currentWeather != nil {
		weather := l.currentWeather.Weather[0]

		icon, err := icons.IconImage(weather.Icon)
		if err != nil {
			return fmt.Errorf("Failed to get icon image, %w", err)
		}

		image := tview.NewImage()
		image.SetImage(icon).SetAlign(tview.AlignCenter, tview.AlignCenter)

		sunset := time.Unix(int64(l.currentWeather.Sys.Sunset), 0).Format("15:04")
		sunrise := time.Unix(int64(l.currentWeather.Sys.Sunrise), 0).Format("15:04")

		maxTemperature := l.currentWeather.Main.TempMax
		minTemperature := l.currentWeather.Main.TempMin
		if l.forecastWeather != nil {
			for _, entry := range l.forecastWeather.List {
				if entry.Main.TempMax > maxTemperature {
					maxTemperature = entry.Main.TempMax
				}
				if entry.Main.TempMin < minTemperature {
					minTemperature = entry.Main.TempMin
				}
			}
		}

		table := tview.NewTable().
			SetCell(0, 0, tview.NewTableCell("Current Temperature").SetExpansion(1)).
			SetCell(0, 1, tview.NewTableCell(fmt.Sprintf("%.1f %s", l.currentWeather.Main.Temp, unitText)).SetTextColor(tcell.ColorBlue)).
			SetCell(1, 0, tview.NewTableCell("Max Temperature").SetExpansion(1)).
			SetCell(1, 1, tview.NewTableCell(fmt.Sprintf("%.1f %s", maxTemperature, unitText)).SetTextColor(tcell.ColorBlue)).
			SetCell(2, 0, tview.NewTableCell("Min Temperature").SetExpansion(1)).
			SetCell(2, 1, tview.NewTableCell(fmt.Sprintf("%.1f %s", minTemperature, unitText)).SetTextColor(tcell.ColorBlue)).
			SetCell(3, 0, tview.NewTableCell("Sunrise").SetExpansion(1)).
			SetCell(3, 1, tview.NewTableCell(fmt.Sprintf("%s", sunrise)).SetTextColor(tcell.ColorBlue)).
			SetCell(4, 0, tview.NewTableCell("Sunset").SetExpansion(1)).
			SetCell(4, 1, tview.NewTableCell(fmt.Sprintf("%s", sunset)).SetTextColor(tcell.ColorBlue))
		table.SetBorderPadding(1, 0, 0, 0)

		flexView := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(image, 0, 1, false).
			AddItem(tview.NewTextView().
				SetText(l.currentWeather.Name).
				SetDynamicColors(true).
				SetTextAlign(tview.AlignCenter).
				SetTextStyle(tcell.StyleDefault.Attributes(tcell.AttrBold)).SetTextColor(tcell.ColorYellowGreen), 1, 1, false).
			AddItem(tview.NewTextView().
				SetText(fmt.Sprintf("%s [gray]- %s", weather.Main, weather.Description)).
				SetDynamicColors(true).
				SetTextAlign(tview.AlignCenter).
				SetTextStyle(tcell.StyleDefault.Attributes(tcell.AttrBold)), 1, 1, false).
			AddItem(table, 6, 1, false)

		l.view.AddItem(flexView, 0, 1, false)

	}

	if l.forecastWeather != nil {
		grid := tview.NewGrid().
			SetRows(0, 0).
			SetColumns(0, 0, 0).
			SetGap(0, 1)
		grid.SetBorderPadding(0, 0, 3, 0)

		for i, entry := range l.forecastWeather.List {
			if i > 5 {
				break
			}

			weather := entry.Weather[0]

			icon, err := icons.IconImage(weather.Icon)
			if err != nil {
				return fmt.Errorf("Failed to get icon image, %w", err)
			}

			flex := tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(tview.NewImage().
					SetImage(icon).
					SetAlign(tview.AlignCenter, tview.AlignCenter), 0, 1, false).
				AddItem(tview.
					NewTextView().
					SetText(fmt.Sprintf("%s [gray]- %s", weather.Main, weather.Description)).
					SetDynamicColors(true).
					SetTextAlign(tview.AlignCenter).
					SetTextStyle(tcell.StyleDefault.Attributes(tcell.AttrBold)), 1, 1, false)

			flex.SetBorder(true).SetBorderColor(tcell.ColorGray).SetTitle(time.Unix(int64(entry.Dt), 0).Format("15:04"))

			row := 0
			col := i
			if i > 2 {
				row = 1
				col = i - 3
			}

			grid.AddItem(flex, row, col, 1, 1, 0, 0, false)
		}

		l.view.AddItem(grid, 0, 3, false)
	}

	return nil
}

func (l *StatusView) Update() error {
	return nil
}
