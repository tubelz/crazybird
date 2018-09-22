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

// MenuSystem is the struct that contains the controllable stick
type MenuSystem struct {
	Name          string
	EntityManager *entity.Manager
	InputManager  *input.Manager
	RenderSystem  *system.RenderSystem
	SceneManager  *macaw.SceneManager
	pos           uint16
}

// Init adds the collision handler
func (m *MenuSystem) Init() {
	log.Print("Starting menu system")
	m.pos = 0
}

// SetSceneManager sets the scene manager in the system
func (m *MenuSystem) SetSceneManager(sm *macaw.SceneManager) {
	m.SceneManager = sm
}

// Update handle the input event
func (m *MenuSystem) Update() {
	if button := m.InputManager.Button(); button != (sdl.KeyboardEvent{}) {
		if button.Keysym.Sym == sdl.K_RETURN && button.State == sdl.PRESSED {
			switch m.pos {
			case 0:
				m.SceneManager.ChangeScene("game")
			case 1:
				m.SceneManager.ChangeScene("credits")
			}
		}
		// move the selectbox
		if (button.Keysym.Sym == sdl.K_DOWN || button.Keysym.Sym == sdl.K_UP) && button.State == sdl.PRESSED {
			if button.Keysym.Sym == sdl.K_DOWN {
				m.pos = 1
			} else {
				m.pos = 0
			}
			log.Println(m.pos)
			for _, obj := range m.EntityManager.GetAll() {
				if obj == nil {
					continue
				}
				if obj.GetType() != "menuselectbox" {
					continue
				}
				position := obj.GetComponent(&entity.PositionComponent{}).(*entity.PositionComponent)
				position.Pos.Y = 399 + int32(m.pos)*30
			}
		}
	}
}
