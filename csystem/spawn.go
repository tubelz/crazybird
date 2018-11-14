package csystem

import (
	"github.com/tubelz/macaw/entity"
	"github.com/tubelz/macaw/math"

	// "github.com/tubelz/macaw/cmd"
	"log"
	"math/rand"

	"github.com/tubelz/macaw/system"
	"github.com/veandco/go-sdl2/sdl"
)

// SpawnSystem is the struct that contains the controllable stick
type SpawnSystem struct {
	EntityManager   *entity.Manager
	Name            string
	RenderSystem    *system.RenderSystem
	CollisionSystem *system.CollisionSystem
	Spritesheet     *entity.Spritesheet
	Sprites         map[string]entity.RenderComponent
	count           uint
}

// Init loads the spritesheet
func (s *SpawnSystem) Init() {
	log.Print("Starting spawn system")
	// Event to delete objects after they cross the left border
	removeFunc := func(event system.Event) {
		removeObj(s, event)
	}
	s.CollisionSystem.AddHandler("border event", removeFunc)

	objSpritesheet := &entity.Spritesheet{Renderer: s.RenderSystem.Renderer, Filepath: "assets/img/objects.png"}
	objSpritesheet.Init()
	s.Spritesheet = objSpritesheet
	// Load sprites from spritesheet
	s.Sprites = make(map[string]entity.RenderComponent)
	mushroomCrop := &sdl.Rect{0, 0, 30, 38}
	mushroomSprite := objSpritesheet.LoadSprite(mushroomCrop)
	s.Sprites["mushroom"] = mushroomSprite

	fruitCrop := &sdl.Rect{30, 0, 30, 38}
	fruitSprite := objSpritesheet.LoadSprite(fruitCrop)
	s.Sprites["fruit"] = fruitSprite

	modifySpritesheet := &entity.Spritesheet{Renderer: s.RenderSystem.Renderer, Filepath: "assets/img/objects.png"}
	modifySpritesheet.Init()
	modifySpritesheet.Texture.SetColorMod(128, 128, 0)
	badSprite := modifySpritesheet.LoadSprite(fruitCrop)
	s.Sprites["badfruit"] = badSprite

}

// Update handle the input event
func (s *SpawnSystem) Update() {
	s.count++
	if s.count%100 == 0 {
		s.count = 0
		obj := s.EntityManager.Create("badfruit")
		sprite, _ := s.Sprites["badfruit"]
		obj.AddComponent(&sprite)
		// obj.AddComponent("geometry", &entity.RectangleComponent{
		// 	Size:   &sdl.Point{40, 40},
		// 	Color:  &sdl.Color{0xFF, 0x00, 0x00, 0xFF},
		// 	Filled: true,
		// })
		y := rand.Int31n(470) + 30
		obj.AddComponent(&entity.PositionComponent{Pos: &sdl.Point{740, y}})
		obj.AddComponent(&entity.CollisionComponent{[]sdl.Rect{sdl.Rect{0, 0, 40, 32}}})
		obj.AddComponent(&entity.PhysicsComponent{
			Vel:       &math.FPoint{-1, 0},
			Acc:       &math.FPoint{0, 0},
			FuturePos: &math.FPoint{740, float32(y)},
		})
		return
	}
	if s.count%20 == 0 {
		if rand.Int31n(3) == 0 {
			obj := s.EntityManager.Create("fruit")
			sprite, _ := s.Sprites["fruit"]
			obj.AddComponent(&sprite)
			// obj.AddComponent("geometry", &entity.RectangleComponent{
			// 	Size:   &sdl.Point{40, 40},
			// 	Color:  &sdl.Color{0x00, 0xFF, 0x00, 0xFF},
			// 	Filled: true,
			// })
			y := rand.Int31n(470) + 30
			obj.AddComponent(&entity.PositionComponent{Pos: &sdl.Point{740, y}, Z: 1})
			obj.AddComponent(&entity.CollisionComponent{[]sdl.Rect{sdl.Rect{0, 0, 40, 32}}})
			obj.AddComponent(&entity.PhysicsComponent{
				Vel:       &math.FPoint{-1, 0},
				Acc:       &math.FPoint{0, 0},
				FuturePos: &math.FPoint{740, float32(y)},
			})
		} else {
			obj := s.EntityManager.Create("mushroom")
			sprite, _ := s.Sprites["mushroom"]
			obj.AddComponent(&sprite)
			y := rand.Int31n(470) + 30
			obj.AddComponent(&entity.PositionComponent{Pos: &sdl.Point{740, y}})
			obj.AddComponent(&entity.CollisionComponent{[]sdl.Rect{sdl.Rect{0, 0, 40, 32}}})
			obj.AddComponent(&entity.PhysicsComponent{
				Vel:       &math.FPoint{-1, 0},
				Acc:       &math.FPoint{0, 0},
				FuturePos: &math.FPoint{740, float32(y)},
			})
		}
	}
}

func removeObj(s *SpawnSystem, event system.Event) {
	em := s.EntityManager
	border := event.(*system.BorderEvent)
	if border.Ent.GetType() == "player" {
		return
	}
	component := border.Ent.GetComponent(&entity.PositionComponent{})
	position := component.(*entity.PositionComponent)

	if border.Side == "left" && position.Pos.X < -30 {
		em.Delete(border.Ent.GetID())
	}
}
