package latency

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/guptarohit/asciigraph"
	"github.com/merlinfuchs/blimp/internal/config"
	ping "github.com/prometheus-community/pro-bing"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

type LatencyView struct {
	stopped chan struct{}
	data    []float64
	view    *tview.TextView
}

func New() *LatencyView {
	textView := tview.NewTextView().SetWrap(false)
	textView.SetBorder(true).
		SetBorderColor(tcell.ColorGray).
		SetTitle("Network Latency (ms)").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 0, 0, 0)

	return &LatencyView{
		stopped: make(chan struct{}),
		data:    make([]float64, config.K.Int("widgets.latency.history_length")),
		view:    textView,
	}
}

func (l *LatencyView) Start() {
	var pinger *ping.Pinger

	go func() {
		<-l.stopped
		if pinger != nil {
			pinger.Stop()
		}
	}()

	go func() {
		for {
			select {
			case <-l.stopped:
				return
			default:
				var err error
				pinger, err = ping.NewPinger(config.K.String("widgets.latency.target_host"))
				if err != nil {
					log.Error().Err(err).Msgf("Failed to create pinger")
					return
				}

				pinger.SetPrivileged(false)
				pinger.Interval = time.Duration(config.K.Int("widgets.latency.update_interval")) * time.Millisecond

				pinger.OnRecv = func(pkt *ping.Packet) {
					newValue := float64(pkt.Rtt.Microseconds()) / 1000
					l.data = append(l.data[1:], newValue)
				}

				err = pinger.Run()
				if err != nil {
					log.Error().Err(err).Msgf("Failed to run pinger, latency won't be displayed")
				}
			}
		}
	}()
}

func (l *LatencyView) Stop() {
	close(l.stopped)
}

func (l *LatencyView) Update() error {
	_, _, width, height := l.view.GetRect()
	graph := asciigraph.Plot(l.data, asciigraph.Precision(1), asciigraph.Width(width-12), asciigraph.Height(height-4))
	l.view.SetText(graph)
	return nil
}

func (l *LatencyView) Primitive() tview.Primitive {
	return l.view
}
