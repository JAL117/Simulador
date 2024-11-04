package views

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func createStatsBox(title string, valueLabel *canvas.Text, bgColor color.Color) *fyne.Container {
	bg := canvas.NewRectangle(bgColor)
	bg.SetMinSize(fyne.NewSize(150, 80))

	titleLabel := canvas.NewText(title, textColor)
	titleLabel.TextSize = 16
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.Alignment = fyne.TextAlignCenter

	return container.NewStack(
		bg,
		container.NewVBox(
			container.NewCenter(titleLabel),
			container.NewCenter(valueLabel),
			container.NewCenter(),
		),
	)
}
