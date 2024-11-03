package views

import (
	"Parking-Simulator/src/models"
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type ParkingSpace struct {
	Container   *fyne.Container
	Background  *canvas.Rectangle
	CarImage    *canvas.Image
	NumberLabel *canvas.Text
	StatusText  *canvas.Text
}

var (
	backgroundColor = &color.NRGBA{R: 40, G: 44, B: 52, A: 255}
	availableImage  = "src/assets/si.png"
	occupiedImages  = []string{"src/assets/cars/car-black.png", "src/assets/cars/car-blue.png", "src/assets/cars/car-green.png", "src/assets/cars/car-red.png", "src/assets/cars/car-purple.png"}
	borderColor     = &color.NRGBA{R: 255, G: 250, B: 0, A: 255}
	textColor       = &color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	totalColor      = &color.NRGBA{R: 255, G: 250, B: 0, A: 255}
	occupiedColor   = &color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	availableColor  = &color.NRGBA{R: 0, G: 255, B: 53, A: 255}
)

func NewParkingSpace(number int) *ParkingSpace {
	space := &ParkingSpace{}

	space.Background = canvas.NewRectangle(backgroundColor)
	space.Background.SetMinSize(fyne.NewSize(50, 50))
	space.Background.Resize(fyne.NewSize(60, 100))
	space.Background.StrokeWidth = 1.5
	space.Background.StrokeColor = borderColor

	space.CarImage = canvas.NewImageFromFile(availableImage)
	space.CarImage.SetMinSize(fyne.NewSize(90, 150))
	space.CarImage.FillMode = canvas.ImageFillContain
	space.NumberLabel = canvas.NewText(fmt.Sprintf("%d", number), textColor)
	space.NumberLabel.TextSize = 16
	space.NumberLabel.TextStyle = fyne.TextStyle{Bold: true}
	space.NumberLabel.Alignment = fyne.TextAlignCenter

	space.StatusText = canvas.NewText("LIBRE", textColor)
	space.StatusText.TextSize = 12
	space.StatusText.Alignment = fyne.TextAlignCenter

	space.Container = container.NewStack(
		space.Background,
		container.NewPadded(
			container.NewVBox(
				container.NewCenter(space.NumberLabel),
				container.NewCenter(space.CarImage),
				container.NewCenter(space.StatusText),
			),
		),
	)

	return space
}

func (p *ParkingSpace) UpdateStatus(occupied bool, carID int) {
	if occupied {

		p.CarImage.File = occupiedImages[rand.Intn(len(occupiedImages))]
		p.StatusText.Text = fmt.Sprintf("Auto #%d", carID)
	} else {
		p.CarImage.File = availableImage
		p.StatusText.Text = "LIBRE"
	}
	p.CarImage.Refresh()
	p.StatusText.Refresh()
}

type StatsPanel struct {
	Container     *fyne.Container
	TotalLabel    *canvas.Text
	OccupiedLabel *canvas.Text
	FreeLabel     *canvas.Text
	WaitingLabel  *canvas.Text
}

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
		),
	)
}

func (s *StatsPanel) UpdateStats(occupied, capacity int) {
	s.OccupiedLabel.Text = fmt.Sprintf("%d", occupied)
	s.FreeLabel.Text = fmt.Sprintf("%d", capacity-occupied)
	s.OccupiedLabel.Refresh()
	s.FreeLabel.Refresh()
}

func CreateWindow(app fyne.App, parkingLot *models.Parking, duration float64, totalCars int) fyne.Window {
	window := app.NewWindow("Simulador con Concurrencia")
	window.SetPadded(false)

	background := canvas.NewRectangle(&color.NRGBA{R: 0, G: 85, B: 19, A: 255})
	background.SetMinSize(fyne.NewSize(1200, 800))

	spaces := make([]*ParkingSpace, parkingLot.Capacity())
	for i := 0; i < parkingLot.Capacity(); i++ {
		spaces[i] = NewParkingSpace(i + 1)
	}

	spacesPerRow := parkingLot.Capacity() / 2
	topRow := container.NewGridWithColumns(spacesPerRow)
	bottomRow := container.NewGridWithColumns(spacesPerRow)

	for i := 0; i < spacesPerRow; i++ {
		topRow.Add(spaces[i].Container)
		bottomRow.Add(spaces[i+spacesPerRow].Container)
	}

	roadColor := &color.NRGBA{R: 58, G: 58, B: 58, A: 255}
	road := canvas.NewRectangle(roadColor)
	road.SetMinSize(fyne.NewSize(1600, 80))

	roadHorizontal := canvas.NewRectangle(roadColor)
	roadHorizontal.SetMinSize(fyne.NewSize(700, 100))

	roadHorizontal2 := canvas.NewRectangle(roadColor)
	roadHorizontal2.SetMinSize(fyne.NewSize(700, 100))

	roadVertical := canvas.NewRectangle(roadColor)
	roadVertical.SetMinSize(fyne.NewSize(80, 600))

	stats := NewStatsPanel(parkingLot.Capacity())

	parkingLayout := container.NewVBox(
		container.NewHBox(
			roadVertical,
			container.NewVBox(
				topRow,
				container.NewCenter(road),
				bottomRow,
			),
		),
		roadHorizontal2,
	)

	scrollContainer := container.NewScroll(parkingLayout)
	scrollContainer.Resize(fyne.NewSize(800, 600))

	mainContainer := container.NewBorder(
		stats.Container,
		nil, nil, nil,
		scrollContainer,
	)

	content := container.NewMax(background, mainContainer)
	window.SetContent(content)
	window.Resize(fyne.NewSize(1200, 800))
	window.CenterOnScreen()

	go func() {
		ticker := time.NewTicker(700 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			occupiedSpaces, carIDs := parkingLot.OccupiedSpaces()
			occupied := 0
			waitingCars := make([]int, 0)
			for i, isOccupied := range occupiedSpaces {
				spaces[i].UpdateStatus(isOccupied, carIDs[i])
				if isOccupied {
					occupied++
				}
			}
			stats.UpdateStats(occupied, parkingLot.Capacity())

			select {
			case car := <-parkingLot.Queue:
				waitingCars = append(waitingCars, car.ID)
			default:

			}

			stats.UpdateWaitingCars(waitingCars)
		}
	}()

	go func() {
		for i := 1; i <= totalCars; i++ {
			time.Sleep(time.Duration(rand.ExpFloat64()*duration) * time.Millisecond)
			car := &models.Car{ID: i}
			go parkingLot.Enter(car)
		}
	}()

	return window
}
