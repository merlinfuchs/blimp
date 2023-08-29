package internal

import "github.com/rivo/tview"

type Widget interface {
	Start()
	Stop()
	Update() error
	Primitive() tview.Primitive
}
