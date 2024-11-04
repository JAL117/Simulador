package views

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type StatsPanel struct {
	Container     *fyne.Container
	TotalLabel    *canvas.Text
	OccupiedLabel *canvas.Text
	FreeLabel     *canvas.Text
	WaitingLabel  *canvas.Text
}

var (
	totalColor     = &color.NRGBA{R: 255, G: 250, B: 0, A: 255}
	occupiedColor  = &color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	availableColor = &color.NRGBA{R: 0, G: 255, B: 53, A: 255}
)

func NewStatsPanel(capacity int) *StatsPanel {
	stats := &StatsPanel{}

	stats.TotalLabel = canvas.NewText(fmt.Sprintf("%d", capacity), textColor)
	stats.OccupiedLabel = canvas.NewText("0", textColor)
	stats.FreeLabel = canvas.NewText(fmt.Sprintf("%d", capacity), textColor)
	stats.WaitingLabel = canvas.NewText("Esperando: Ninguno", textColor)

	for _, label := range []*canvas.Text{stats.TotalLabel, stats.OccupiedLabel, stats.FreeLabel, stats.WaitingLabel} {
		label.TextSize = 24
		label.TextStyle = fyne.TextStyle{Bold: true}
		label.Alignment = fyne.TextAlignCenter
	}

	totalBox := createStatsBox("TOTAL", stats.TotalLabel, totalColor)
	occupiedBox := createStatsBox("OCUPADOS", stats.OccupiedLabel, occupiedColor)
	freeBox := createStatsBox("DISPONIBLES", stats.FreeLabel, availableColor)

	waitingBox := createStatsBox("EN ESPERA", stats.WaitingLabel, availableColor)

	stats.Container = container.NewVBox(
		widget.NewSeparator(),
		container.NewHBox(
			layout.NewSpacer(),
			totalBox,
			layout.NewSpacer(),
			occupiedBox,
			layout.NewSpacer(),
			freeBox,
			layout.NewSpacer(),
			waitingBox,
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)

	return stats
}

func (s *StatsPanel) UpdateWaitingCars(waitingCars []int) {
	if len(waitingCars) == 0 {
		s.WaitingLabel.Text = "Carro #:"
	} else {
		waitingText := "Carro #:" + fmt.Sprint(waitingCars)
		s.WaitingLabel.Text = waitingText
	}
	s.WaitingLabel.Refresh()
}

func (s *StatsPanel) UpdateStats(occupied, capacity int) {
	s.OccupiedLabel.Text = fmt.Sprintf("%d", occupied)
	s.FreeLabel.Text = fmt.Sprintf("%d", capacity-occupied)
	s.OccupiedLabel.Refresh()
	s.FreeLabel.Refresh()
}
