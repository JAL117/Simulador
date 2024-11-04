package views

import (
	"Parking-Simulator/src/models"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func CreateWindow(app fyne.App, parkingLot *models.Parking, totalCars int) fyne.Window {
	window := app.NewWindow("Simulador de estacionamiento")
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
	roadHorizontal.SetMinSize(fyne.NewSize(700, 80))

	roadHorizontal2 := canvas.NewRectangle(roadColor)
	roadHorizontal2.SetMinSize(fyne.NewSize(700, 80))

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
			roadVertical.FillColor = parkingLot.EntryColor
			roadVertical.Refresh()
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
				roadHorizontal2.FillColor = parkingLot.WaitColor
				roadHorizontal2.Refresh()
			default:
				roadHorizontal2.FillColor = roadColor
			}

			stats.UpdateWaitingCars(waitingCars)
		}
	}()

	go func() {
		for i := 1; i <= totalCars; i++ {
			time.Sleep(200 * time.Millisecond)
			car := &models.Car{ID: i}
			go parkingLot.Enter(car)
		}
	}()

	return window
}
