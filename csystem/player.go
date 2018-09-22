package csystem

import (
	"fmt"
	"github.com/tubelz/macaw"
	"github.com/tubelz/macaw/entity"
	"github.com/tubelz/macaw/input"
	"github.com/tubelz/macaw/math"
	"github.com/tubelz/macaw/system"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"math/rand"
	"time"
)

const (
	// MoveNone represents the flag that the bird has no movement
	MoveNone = 0x00
	// MoveLeft represents the flag to move the bird to the left
	MoveLeft = 0x01
	// MoveRight represents the flag to move the bird to the right
	MoveRight = 0x02
	// MoveUp represents the flag to move the bird to the up
	MoveUp = 0x04
	// MoveDown represents the flag to move the bird to the down
	MoveDown = 0x08
)

// PlayerSystem is the struct that contains the controllable stick
type PlayerSystem struct {
	EntityManager   *entity.Manager
	Name            string
	InputManager    *input.Manager
	CollisionSystem *system.CollisionSystem
	RenderSystem    *system.RenderSystem
	SceneManager    *macaw.SceneManager
	score           int
	lifes           int
	crazy           bool
	crazydirection  uint8
	crazyspeed      *math.FPoint
	sprites         map[string]entity.RenderComponent
}

// Init adds the collision handler
func (p *PlayerSystem) Init() {
	log.Print("Starting player system")
	collectFunc := func(event system.Event) {
		collectObject(p, event)
	}
	p.CollisionSystem.AddHandler("collision event", collectFunc)
	p.CollisionSystem.AddHandler("border event", stopPlayer)
	p.lifes = 3
	p.score = 0
	p.crazy = false
	objSpritesheet := &entity.Spritesheet{Renderer: p.RenderSystem.Renderer, Filepath: "assets/img/objects.png"}
	objSpritesheet.Init()
	// Load sprites from spritesheet
	p.sprites = make(map[string]entity.RenderComponent)
	lifelessCrop := &sdl.Rect{94, 0, 32, 38}
	lifelessSprite := objSpritesheet.LoadSprite(lifelessCrop)
	p.sprites["lifeless"] = lifelessSprite
}

// SetSceneManager sets the scene manager in the system
func (p *PlayerSystem) SetSceneManager(sm *macaw.SceneManager) {
	p.SceneManager = sm
}

// Update handle the input event
func (p *PlayerSystem) Update() {
	if p.crazy {
		var player *entity.Entity
		it := p.EntityManager.IterAvailable(-1)
		for tmpObj, entIndex := it(); entIndex != -1; tmpObj, entIndex = it() {
			if tmpObj.GetType() == "player" {
				player = tmpObj
				break
			}
		}
		p.moveCrazy(player)
	} else if button := p.InputManager.Button(); button != (sdl.KeyboardEvent{}) {
		for _, obj := range p.EntityManager.GetAll() {
			if obj == nil {
				continue
			}
			if obj.GetType() != "player" {
				continue
			}
			physicsComponent := obj.GetComponent(&entity.PhysicsComponent{})
			physics := physicsComponent.(*entity.PhysicsComponent)
			positionComponent := obj.GetComponent(&entity.PositionComponent{})
			position := positionComponent.(*entity.PositionComponent)

			var value float32
			if button.State == sdl.RELEASED {
				value = 0
			} else {
				value = 3
			}

			switch button.Keysym.Sym {
			case sdl.K_LEFT:
				if position.Pos.X > 3 {
					physics.Vel.X = value * -1
				}
			case sdl.K_RIGHT:
				if position.Pos.X < 800 {
					physics.Vel.X = value
				}
			case sdl.K_UP:
				if position.Pos.Y > 3 {
					physics.Vel.Y = value * -1
				}
			case sdl.K_DOWN:
				{
					if position.Pos.Y >= 450 {
						physics.Vel.Y = 0
					} else {
						physics.Vel.Y = value
					}
				}
			}
		}
	}
}

// stops the player once it hits the border
func stopPlayer(event system.Event) {
	border := event.(*system.BorderEvent)
	if border.Ent.GetID() != 1 {
		return
	}
	component := border.Ent.GetComponent(&entity.PositionComponent{})
	position := component.(*entity.PositionComponent)

	component = border.Ent.GetComponent(&entity.PhysicsComponent{})
	physics := component.(*entity.PhysicsComponent)

	switch border.Side {
	case "top":
		position.Pos.Y = 1
		physics.FuturePos.Y = 1
		physics.Vel.Y *= 0
	case "bottom":
		// size := collision.Size.Y
		// position.Pos.Y = 599 - size
		// physics.FuturePos.Y = float32(599 - size)
		physics.Vel.Y *= 0
	case "left":
		position.Pos.X = 1
		physics.FuturePos.X = 1
		physics.Vel.X *= 0
	case "right":
		// size := collision.Size.X
		// position.Pos.X = 799 - size
		// physics.FuturePos.X = float32(799 - size)
		physics.Vel.X *= 0
	}
}

func collectObject(ps *PlayerSystem, event system.Event) {
	collision := event.(*system.CollisionEvent)
	if collision.Ent.GetType() != "player" {
		return
	}
	em := ps.EntityManager
	objCollided := collision.With
	log.Printf("Obj %d collected", objCollided.GetID())

	if objCollided.GetType() == "fruit" {
		// score is wrong here... we should get the right type
		var playerScore *entity.Entity
		it := ps.EntityManager.IterAvailable(-1)
		for tmpObj, entIndex := it(); entIndex != -1; tmpObj, entIndex = it() {
			if tmpObj.GetType() == "score" {
				playerScore = tmpObj
				break
			}
		}
		f := playerScore.GetComponent(&entity.FontComponent{}).(*entity.FontComponent)
		ps.score = ps.score + 1
		f.Text = fmt.Sprintf("score: %d", ps.score)
		f.Modified = true
	} else if objCollided.GetType() == "badfruit" {
		ps.lifes--
		lostLife(ps)
		if ps.lifes == 0 {
			ps.SceneManager.ChangeScene("menu")
		}
	} else if objCollided.GetType() == "mushroom" {
		c := collision.Ent.GetComponent(&entity.RenderComponent{})
		render := c.(*entity.RenderComponent)
		render.Texture.SetColorMod(0xFF, 0, 0x6F)
		ps.crazy = true
		defineDirection(ps)
		// Speed of the movement
		ps.crazyspeed = &math.FPoint{float32(rand.Int31n(5) + 3), float32(rand.Int31n(5) + 3)}
		_ = time.AfterFunc(time.Second*3, func() {
			if collision.Ent != nil {
				ps.crazy = false
				render.Texture.SetColorMod(0xFF, 0xFF, 0xFF)
				physicsComponent := collision.Ent.GetComponent(&entity.PhysicsComponent{})
				if physicsComponent != nil {
					physics := physicsComponent.(*entity.PhysicsComponent)
					physics.Vel.X = 0
					physics.Vel.Y = 0
				}
			}
		})

	}
	em.Delete(objCollided.GetID())
}

// defines the direction the bird will go after it eats the mushroom
func defineDirection(ps *PlayerSystem) {
	// directions
	verticalFlag := rand.Int31n(2)
	horizontalFlag := rand.Int31n(2)
	var upOrDown int
	var leftOrRight int

	if verticalFlag == 1 {
		upOrDown = MoveUp
	} else {
		upOrDown = MoveDown
	}

	if horizontalFlag == 1 {
		leftOrRight = MoveLeft
	} else {
		leftOrRight = MoveRight
	}

	ps.crazydirection = uint8(upOrDown | leftOrRight)
}

// move the bird according to the crazydirection flag
func (p *PlayerSystem) moveCrazy(player *entity.Entity) {
	if player == nil {
		return
	}
	physicsComponent := player.GetComponent(&entity.PhysicsComponent{})
	physics := physicsComponent.(*entity.PhysicsComponent)
	positionComponent := player.GetComponent(&entity.PositionComponent{})
	position := positionComponent.(*entity.PositionComponent)

	speed := p.crazyspeed.X
	if (p.crazydirection & MoveLeft) != 0 {
		if position.Pos.X > 6 {
			physics.Vel.X = -speed
		} else {
			physics.Vel.X = speed
			p.crazydirection ^= MoveLeft
			p.crazydirection |= MoveRight
		}
	} else if (p.crazydirection & MoveRight) != 0 {
		if position.Pos.X < 700 {
			physics.Vel.X = speed
		} else {
			physics.Vel.X = -speed
			p.crazydirection ^= MoveRight
			p.crazydirection |= MoveLeft
		}
	}

	speed = p.crazyspeed.Y
	if (p.crazydirection & MoveUp) != 0 {
		if position.Pos.Y > 6 {
			physics.Vel.Y = -speed
		} else {
			physics.Vel.Y = speed
			p.crazydirection ^= MoveUp
			p.crazydirection |= MoveDown
		}
	} else if (p.crazydirection & MoveDown) != 0 {
		if position.Pos.Y < 450 {
			physics.Vel.Y = speed
		} else {
			physics.Vel.Y = -speed
			p.crazydirection ^= MoveDown
			p.crazydirection |= MoveUp
		}
	}
}

// takes care when the player loses a life
func lostLife(ps *PlayerSystem) {
	// remove life
	it := ps.EntityManager.IterAvailable(-1)
	var lifeIds [3]uint16
	var count int
	for obj, entIndex := it(); entIndex != -1; obj, entIndex = it() {
		if obj.GetType() == "life" {
			lifeIds[count] = obj.GetID()
			count++
		}
	}
	ps.EntityManager.Delete(lifeIds[0])
	// add lost life image
	obj := ps.EntityManager.Create("lifeless")
	sprite, _ := ps.sprites["lifeless"]
	x := 620 + ps.lifes*40
	obj.AddComponent(&sprite)
	obj.AddComponent(&entity.PositionComponent{&sdl.Point{int32(x), 10}})
}
