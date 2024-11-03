package main

import (
	"Parking-Simulator/src/models"
	view "Parking-Simulator/src/views"
	"log"

	"fyne.io/fyne/v2/app"
)

const (
	capacidadEstacionamiento = 20
	duracion                 = 245.0
	totalCarros              = 100
)

func main() {
	application := app.New()

	parkingLot := models.NewParking(capacidadEstacionamiento)
	if parkingLot == nil {
		log.Fatalf("Error al crear el estacionamiento con capacidad %d", capacidadEstacionamiento)
	}

	window := view.CreateWindow(application, parkingLot, duracion, totalCarros)
	if window == nil {
		log.Fatalf("Error al crear la ventana de la aplicación")
	}

	window.ShowAndRun()
}
