package csystem

import (
	// "fmt"
	"github.com/tubelz/macaw"
	"github.com/tubelz/macaw/entity"
	"github.com/tubelz/macaw/input"
	"github.com/tubelz/macaw/system"
	"github.com/veandco/go-sdl2/sdl"
	"log"
)

// CreditsSystem is the struct that contains the controllable stick
type CreditsSystem struct {
	Name          string
	EntityManager *entity.Manager
	InputManager  *input.Manager
	RenderSystem  *system.RenderSystem
	SceneManager  *macaw.SceneManager
}

// Init adds the collision handler
func (m *CreditsSystem) Init() {
	log.Print("Starting credit system")
}

// SetSceneManager sets the scene manager in the system
func (m *CreditsSystem) SetSceneManager(sm *macaw.SceneManager) {
	m.SceneManager = sm
}

// Update handle the input event
func (m *CreditsSystem) Update() {
	if button := m.InputManager.Button(); button != (sdl.KeyboardEvent{}) {
		if button.Keysym.Sym == sdl.K_RETURN && button.State == sdl.PRESSED {
			m.SceneManager.ChangeScene("menu")
		}
	}
}
