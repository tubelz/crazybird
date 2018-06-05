package csystem

import (
	"github.com/tubelz/macaw/entity"
	"log"
)

// ScrollSystem is the struct that contains the controllable stick
type ScrollSystem struct {
	EntityManager *entity.Manager
	Name          string
}

// Init initializes the render system using the current window
func (s *ScrollSystem) Init() {
	log.Print("Starting scroll system")
}

// Update handle the input event
func (s *ScrollSystem) Update() {
	for _, e := range s.EntityManager.GetAll() {
		if e == nil {
			continue
		}
		if e.GetType() != "background" {
			continue
		}
		r, _ := e.GetComponent("render")
		render := r.(*entity.RenderComponent)
		render.Crop.X++
		if render.Crop.X > 1942-800 {
			render.Crop.X = 0
		}
	}
}
