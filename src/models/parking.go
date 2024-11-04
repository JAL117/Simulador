package models

import (
	"context"
	"image/color"
	"log"
	"math/rand"
	"sync"
	"time"
)

var logger = log.New(log.Writer(), "", 0)

type Parking struct {
	capacity       int           // Capacidad total del estacionamiento
	mutex          sync.RWMutex  // Mutex para acceso concurrente
	Queue          chan *Car     // Canal para esperar vehículos
	availableSpots chan struct{} // Canal que indica espacios disponibles
	entryExitMutex chan struct{} // Mutex para controlar la entrada/salida
	occupiedSpaces []bool        // Array para saber qué espacios están ocupados
	carIDs         []int         // Array para almacenar IDs de vehículos
	nextSpotIndex  int           // Índice para el próximo espacio disponible
	EntryColor     color.Color
	WaitColor      color.Color
	wg             sync.WaitGroup // WaitGroup para sincronizar goroutines
}

func NewParking(capacity int) *Parking {
	if capacity <= 0 {
		logger.Fatalf("La capacidad del estacionamiento debe ser mayor que cero.")
	}

	parking := &Parking{
		capacity:       capacity,
		Queue:          make(chan *Car, capacity),     // Inicializa el canal para la cola
		availableSpots: make(chan struct{}, capacity), // Canal para espacios disponibles
		entryExitMutex: make(chan struct{}, 1),        // Mutex para entrada/salida
		occupiedSpaces: make([]bool, capacity),        // Array de espacios ocupados
		carIDs:         make([]int, capacity),         // Array para IDs de vehículos
		nextSpotIndex:  0,                             // Comenzamos con el primer espacio
	}
	// Llenamos el canal de espacios disponibles al iniciar
	for i := 0; i < capacity; i++ {
		parking.availableSpots <- struct{}{}
	}

	parking.EntryColor = color.NRGBA{R: 40, G: 43, B: 51, A: 255}
	return parking
}

// Capacity devuelve la capacidad total del estacionamiento
func (p *Parking) Capacity() int {
	return p.capacity
}

// Dvuelve los espacios ocupados y los IDs de los vehículos
func (p *Parking) OccupiedSpaces() ([]bool, []int) {
	p.mutex.RLock() // Usamos un bloqueo de lectura para no interferir con otros procesos
	defer p.mutex.RUnlock()

	occupiedSpaces := make([]bool, p.capacity)
	carIDs := make([]int, p.capacity)

	// Copiamos el estado actual de los espacios ocupados y los IDs
	for i := 0; i < p.capacity; i++ {
		occupiedSpaces[i] = p.occupiedSpaces[i]
		if p.occupiedSpaces[i] {
			carIDs[i] = p.carIDs[i]
		}
	}
	return occupiedSpaces, carIDs
}

// Busca el próximo espacio libre
func (p *Parking) findNextSpot() int {
	for i := 0; i < p.capacity; i++ {
		index := (p.nextSpotIndex + i) % p.capacity // Usamos el índice circular
		if !p.occupiedSpaces[index] {               // Si encontramos un espacio libre
			p.nextSpotIndex = (index + 1) % p.capacity // Actualizamos el índice para la próxima búsqueda
			return index                               // Retornamos el índice del espacio libre
		}
	}
	return -1 // Si no hay espacio disponible, devolvemos -1
}

/*
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
}*/

// Simulate simula la llegada de vehículos al estacionamiento
func Simulate(parking *Parking, arrivalRate float64, ctx context.Context) {
	carID := 1
	carChannel := make(chan *Car) // Canal para enviar nuevos vehículos

	// Goroutine para simular la llegada de vehículos
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(carChannel) // Cerramos el canal al cancelar el contexto
				return
			case <-time.After(time.Duration(rand.ExpFloat64()/arrivalRate) * time.Second):
				car := &Car{ID: carID}
				carChannel <- car // Enviamos el vehículo al canal
				carID++
			}
		}
	}()

	// Goroutine para manejar vehículos que llegan al estacionamiento
	go func() {
		for car := range carChannel {
			parking.Enter(car) // Intentamos que el vehículo entre al estacionamiento
		}
	}()

	// Esperamos a que todas las goroutines hayan terminado
	parking.wg.Wait()
}
