package models

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
)

// Crear un logger personalizado sin fecha y hora
var logger = log.New(log.Writer(), "", 0) // El tercer parámetro 0 significa que no habrá prefijo de tiempo

// Parking es la estructura que representa el estacionamiento
type Parking struct {
	capacity       int           // Capacidad total del estacionamiento
	mutex          sync.RWMutex  // Mutex para acceso concurrente
	Queue          chan *Car     // Canal para esperar vehículos
	availableSpots chan struct{} // Canal que indica espacios disponibles
	entryExitMutex chan struct{} // Mutex para controlar la entrada/salida
	occupiedSpaces []bool        // Array para saber qué espacios están ocupados
	carIDs         []int         // Array para almacenar IDs de vehículos
	nextSpotIndex  int           // Índice para el próximo espacio disponible
}

// NewParking crea un nuevo estacionamiento con la capacidad que se pase como parámetro
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
	return parking
}

// Capacity devuelve la capacidad total del estacionamiento
func (p *Parking) Capacity() int {
	return p.capacity
}

// GetOccupiedSpaces devuelve los espacios ocupados y los IDs de los vehículos
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

// findNextAvailableSpot busca el próximo espacio libre
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

// Maneja la llegada de un vehículo al estacionamiento
func (p *Parking) Enter(car *Car) {
	if car == nil {
		logger.Println("==> No se puede procesar un carro.")
		return
	}

	select {
	case <-p.availableSpots: // Si hay un espacio disponible
		p.entryExitMutex <- struct{}{} // Bloqueamos la entrada/salida

		p.mutex.Lock()                // Bloqueamos para hacer cambios en el estado
		spotIndex := p.findNextSpot() // Buscamos un espacio libre
		if spotIndex != -1 {
			p.occupiedSpaces[spotIndex] = true // Marcamos el espacio como ocupado
			p.carIDs[spotIndex] = car.ID       // Guardamos el ID del vehículo
		} else {
			logger.Printf("==> No hay espacio disponible para el carro %d.\n", car.ID)
		}
		p.mutex.Unlock() // Desbloqueamos el acceso

		<-p.entryExitMutex // Liberamos el mutex de entrada/salida

		const minParkingDuration = 3
		const maxParkingDuration = 3
		time.Sleep(time.Duration(minParkingDuration+rand.Intn(maxParkingDuration)) * time.Second)
		// Simulamos el tiempo que el vehículo estará estacionado
		p.Exit(car) // El vehículo se va
	default:
		logger.Printf("== Carro %d esperando un espacio.\n", car.ID)
		p.Queue <- car // Si no hay espacio, el vehículo se agrega a la cola
	}
}

// Maneja la salida de un vehículo del estacionamiento
func (p *Parking) Exit(car *Car) {
	p.mutex.Lock() // Bloqueamos para hacer cambios
	spotFound := false

	// Buscamos el vehículo en los espacios ocupados
	for i := 0; i < p.capacity; i++ {
		if p.carIDs[i] == car.ID {
			spotFound = true            // Lo encontramos
			p.occupiedSpaces[i] = false // Marcamos el espacio como libre
			p.carIDs[i] = 0             // Limpiamos el ID
			break
		}
	}
	p.mutex.Unlock() // Desbloqueamos

	if spotFound {
		p.availableSpots <- struct{}{} // Liberamos el espacio
		select {
		case nextcar := <-p.Queue: // Si hay vehículos en la cola, intentamos que entre uno
			go p.Enter(nextcar)
		default:
		}
	} else {
		logger.Printf("==> El carro %d no estaba en el estacionamiento.\n", car.ID)
	}
}

// SimulateArrivals simula la llegada de vehículos al estacionamiento
func Simulate(parking *Parking, arrivalRate float64, ctx context.Context) {
	carID := 1
	for {
		select {
		case <-ctx.Done():
			return // Salir si el contexto se cancela
		default:
			time.Sleep(time.Duration(rand.ExpFloat64()/arrivalRate) * time.Second) // Generamos llegadas aleatorias
			car := &Car{ID: carID}                                                 // Creamos un nuevo vehículo
			go parking.Enter(car)                                                  // Intentamos que el vehículo llegue al estacionamiento
			carID++                                                                // Incrementamos el ID del vehículo
		}
	}
}