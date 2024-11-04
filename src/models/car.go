package models

import (
	"image/color"
	"math/rand"
	"time"
)

type Car struct {
	ID int
}

func (p *Parking) Enter(car *Car) {
	if car == nil {
		logger.Println("==> No se puede procesar un carro.")
		return
	}

	select {
	case <-p.availableSpots: // Si hay un espacio disponible
		p.mutex.Lock()                // Bloqueamos para hacer cambios en el estado
		spotIndex := p.findNextSpot() // Buscamos un espacio libre
		p.EntryColor = color.NRGBA{R: 0, G: 255, B: 0, A: 255}

		if spotIndex != -1 {
			p.occupiedSpaces[spotIndex] = true // Marcamos el espacio como ocupado
			p.carIDs[spotIndex] = car.ID       // Guardamos el ID del vehículo
			logger.Printf("### Carro %d ocupó el espacio %d.\n", car.ID, spotIndex)

			// Añadimos al WaitGroup
			p.wg.Add(1)

			// Iniciamos una goroutine separada para simular el tiempo de estacionamiento
			go func(car *Car, spot int) {
				defer p.wg.Done() // Marcamos como hecho cuando termina
				const minParkingDuration = 3
				const maxParkingDuration = 5 // Cambié el máximo a 5 segundos
				time.Sleep(time.Duration(minParkingDuration+rand.Intn(maxParkingDuration-minParkingDuration+1)) * time.Second)
				p.Exit(car) // Llamamos a Exit después del tiempo de estacionamiento
			}(car, spotIndex) // Pasamos spotIndex a la goroutine

		} else {
			logger.Printf("##### No hay espacio disponible para el carro %d.\n", car.ID)
		}
		p.mutex.Unlock() // Desbloqueamos el acceso

	default:
		// Si no hay espacio, el vehículo se agrega a la cola
		logger.Printf("== Carro %d esperando un espacio.\n", car.ID)
		p.WaitColor = color.NRGBA{R: 255, G: 250, B: 0, A: 255}
		p.Queue <- car
	}
}

func (p *Parking) Exit(car *Car) {
	select {
	case p.entryExitMutex <- struct{}{}: // Intentamos bloquear la entrada/salida
		p.mutex.Lock()
		spotFound := false
		for i := 0; i < p.capacity; i++ {
			if p.carIDs[i] == car.ID {
				spotFound = true
				p.occupiedSpaces[i] = false
				p.carIDs[i] = 0
				logger.Printf("=======> Carro %d salió del espacio %d.\n", car.ID, i)
				p.EntryColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
				break
			}
		}
		p.mutex.Unlock()
		<-p.entryExitMutex // Liberamos la entrada/salida después de salir

		if spotFound {
			p.availableSpots <- struct{}{} // Liberamos un espacio
			select {
			case nextCar := <-p.Queue: // Si hay vehículos esperando en la cola
				go p.Enter(nextCar)
			default:
			}
		} else {
			logger.Printf("==> El carro %d no estaba en el estacionamiento.\n", car.ID)
		}
	default:
		logger.Printf("==> Carro %d no pudo salir, entrada/salida ocupada.\n", car.ID)
	}
}
