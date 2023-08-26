package status

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/gdamore/tcell/v2"
	"github.com/merlinfuchs/blimp/blimp/config"
	"github.com/rivo/tview"
)

type StatusView struct {
	stopped chan struct{}
	view    *tview.Flex
	targets []StatusTarget
	data    []StatusEntry
}

func New() *StatusView {
	targets := make([]StatusTarget, 0)
	if err := config.K.Unmarshal("views.status.targets", &targets); err != nil {
		log.Panic().Err(err).Msgf("Failed to unmarshal status targets")
	}

	view := tview.NewFlex().SetDirection(tview.FlexRow)
	view.SetBorder(true).
		SetBorderColor(tcell.ColorGray).
		SetTitle("Application Status").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 1, 1)

	return &StatusView{
		stopped: make(chan struct{}),
		targets: targets,
		view:    view,
		data:    make([]StatusEntry, len(targets)),
	}
}

func (l *StatusView) Start() {
	l.updateData()

	go func() {
		for {
			select {
			case <-l.stopped:
				break
			case <-time.After(time.Duration(config.K.Int("views.status.update_interval")) * time.Millisecond):
				l.updateData()
			}
		}
	}()
}

func (l *StatusView) updateData() {
	for i, target := range l.targets {
		switch target.Type {
		case "http", "https":
			uri := fmt.Sprintf("%s://%s", target.Type, target.Host)
			resp, err := http.Get(uri)
			if err != nil {
				l.data[i] = StatusEntry{
					Target:     target,
					Online:     false,
					HTTPStatus: 0,
				}
			} else {
				l.data[i] = StatusEntry{
					Target:     target,
					Online:     resp.StatusCode >= 200 && resp.StatusCode < 300,
					HTTPStatus: resp.StatusCode,
				}
			}
		default:
			log.Error().Msgf("Unknown status target type %s", target.Type)
		}
	}
}

func (l *StatusView) Stop() {
	close(l.stopped)
}

func (l *StatusView) Update() error {
	l.view.Clear()

	for _, entry := range l.data {
		status := "ONLINE"
		statusColor := "green"
		if !entry.Online {
			status = "OFFLINE"
			statusColor = "red"
		}

		extraInfo := "unreachable"
		if entry.HTTPStatus != 0 {
			extraInfo = fmt.Sprintf("HTTP %d", entry.HTTPStatus)
		}

		l.view.AddItem(
			tview.NewTextView().
				SetText(fmt.Sprintf("[%s]⏺ %s [white]%s (%s) [gray] %s", statusColor, status, entry.Target.Name, entry.Target.Host, extraInfo)).
				SetDynamicColors(true).
				SetWrap(true),
			2, 1, false)
	}

	return nil
}

func (l *StatusView) Primitive() tview.Primitive {
	return l.view
}
