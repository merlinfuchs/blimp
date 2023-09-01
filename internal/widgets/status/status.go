package status

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/merlinfuchs/blimp/internal/config"
	ping "github.com/prometheus-community/pro-bing"
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
	if err := config.K.Unmarshal("widgets.status.targets", &targets); err != nil {
		slog.With("error", err).Error("Failed to unmarshal status targets")
		panic(err)
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
	go func() {
		l.updateData()

		for {
			select {
			case <-l.stopped:
				return
			case <-time.After(time.Duration(config.K.Int("widgets.status.update_interval")) * time.Millisecond):
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
		case "ping":
			pinger, err := ping.NewPinger(target.Host)
			if err != nil {
				slog.With("error", err).Error("Failed to create pinger, latency won't be displayed")
				panic(err)
			}

			pinger.SetPrivileged(false)
			pinger.Count = 1

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			err = pinger.RunWithContext(ctx)
			if err != nil {
				l.data[i] = StatusEntry{
					Target: target,
					Online: false,
				}
			} else {
				stats := pinger.Statistics()
				l.data[i] = StatusEntry{
					Target:      target,
					Online:      stats.PacketLoss == 0,
					PingLatency: stats.AvgRtt,
				}
			}
		default:
			slog.Error(fmt.Sprintf("Unknown status target type %s", target.Type))
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
		} else if entry.PingLatency != 0 {
			extraInfo = fmt.Sprintf("%.1fms", float64(entry.PingLatency.Microseconds())/1000)
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
