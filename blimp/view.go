package blimp

import "github.com/rivo/tview"

type View interface {
	Start()
	Stop()
	Update() error
	Primitive() tview.Primitive
}
