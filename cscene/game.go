package cscene

import (
	"fmt"
	"github.com/tubelz/macaw"
	"github.com/tubelz/macaw/entity"
	"github.com/tubelz/macaw/math"
	"github.com/tubelz/macaw/system"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"time"
)

// GameScene is responsible to manage the content of the game scene
type GameScene struct {
	Scene           *macaw.Scene
	entityIDs       []uint16
	EntityManager   *entity.Manager
	SceneManager    *macaw.SceneManager
	CollisionSystem *system.CollisionSystem
	ticker          *time.Ticker
}

// Init initialize this scene
func (g *GameScene) Init(renderSystem *system.RenderSystem, font *ttf.Font) {
	initFunc := g.initializeEntities(renderSystem, font)
	scene := &macaw.Scene{
		Name:         "game",
		InitFunc:     initFunc,
		ExitFunc:     g.Exit,
		SceneOptions: macaw.SceneOptions{HideCursor: true, Music: "assets/sound/bgmusic.mp3"},
	}
	scene.AddRenderSystem(renderSystem)
	g.Scene = scene
}

// Exit clear the entities created for this scene
func (g *GameScene) Exit() {
	g.ticker.Stop()

	for id, obj := range g.EntityManager.GetAll() {
		if id > 0 && obj != nil {
			g.EntityManager.Delete(obj.GetID())
		}
	}

	g.CollisionSystem.ClearEvents()
}

// GetScene returns the scene from GameScene
func (g *GameScene) GetScene() *macaw.Scene {
	return g.Scene
}

func (g *GameScene) addEntity(ent *entity.Entity) {
	g.entityIDs = append(g.entityIDs, ent.GetID())
}

// func initializeEntities
func (g *GameScene) initializeEntities(renderSystem *system.RenderSystem, font *ttf.Font) func() {
	em := g.EntityManager
	return func() {
		background := em.Create("background")
		player := em.Create("player")
		playerScore := em.Create("score")
		timer := em.Create("timer")
		// grid := em.Create("grid")

		// g.addEntity(background)
		// g.addEntity(player)
		// g.addEntity(playerScore)
		// g.addEntity(timer)
		// g.addEntity(grid)

		//load sprite
		g.createLife(renderSystem)

		spritesheet1 := &entity.Spritesheet{}
		spritesheet1.Init(renderSystem.Renderer, "assets/img/background.png")
		crop1 := &sdl.Rect{0, 0, 800, 600}
		sprite1 := spritesheet1.LoadSprite(crop1)
		background.AddComponent("render", &sprite1)
		background.AddComponent("position", &entity.PositionComponent{&sdl.Point{0, 0}})

		// grid.AddComponent("position", &entity.PositionComponent{&sdl.Point{0, 0}})
		// grid.AddComponent("grid", &entity.GridComponent{Size: &sdl.Point{20, 20}})
		// grid.AddComponent("geometry", &entity.RectangleComponent{
		// 	Size:   &sdl.Point{10, 80},
		// 	Color:  &sdl.Color{0x66, 0x66, 0x66, 0xFF},
		// 	Filled: true,
		// })

		playerSpritesheet := &entity.Spritesheet{}
		playerSpritesheet.Init(renderSystem.Renderer, "assets/img/macaw_pixel.png")
		crop := &sdl.Rect{0, 200, 112, 100}
		sprite := playerSpritesheet.LoadSprite(crop)
		sprite.Flip = 1
		player.AddComponent("render", &sprite)
		player.AddComponent("physics", &entity.PhysicsComponent{
			Vel:       &math.FPoint{0, 0},
			Acc:       &math.FPoint{0, 0},
			FuturePos: &math.FPoint{550, 50},
		})
		player.AddComponent("position", &entity.PositionComponent{&sdl.Point{550, 50}})
		player.AddComponent("collision", &entity.CollisionComponent{
			CollisionAreas: []sdl.Rect{sdl.Rect{80, 22, 25, 19}},
		})
		player.AddComponent("animation", &entity.AnimationComponent{
			InitialPos:     sdl.Point{0, 0},
			AnimationSpeed: 7,
			Current:        0,
			Frames:         5,
			RowLength:      2,
		})

		playerScore.AddComponent("position", &entity.PositionComponent{&sdl.Point{20, 20}})
		playerScore.AddComponent("font", &entity.FontComponent{Text: "score: 0", Modified: true, Font: font})
		playerScore.AddComponent("render", &entity.RenderComponent{Renderer: renderSystem.Renderer})

		timer.AddComponent("position", &entity.PositionComponent{&sdl.Point{300, 20}})
		timer.AddComponent("font", &entity.FontComponent{Text: "00:00:00", Modified: true, Font: font})
		timer.AddComponent("render", &entity.RenderComponent{Renderer: renderSystem.Renderer})

		g.ticker = startTimer(em)
	}
}

func (g *GameScene) createLife(render *system.RenderSystem) {
	em := g.EntityManager
	objSpritesheet := &entity.Spritesheet{}
	objSpritesheet.Init(render.Renderer, "assets/img/objects.png")
	// Load sprites from spritesheet
	heartCrop := &sdl.Rect{60, 0, 32, 38}
	sprite := objSpritesheet.LoadSprite(heartCrop)
	for i := 2; i >= 0; i-- {
		obj := em.Create("life")
		g.addEntity(obj)
		x := 620 + i*40
		obj.AddComponent("render", &sprite)
		obj.AddComponent("position", &entity.PositionComponent{&sdl.Point{int32(x), 10}})
	}
}

func startTimer(em *entity.Manager) *time.Ticker {
	// ticker is a channel that receives periodic info / second
	ticker := time.NewTicker(time.Second)
	var timer *entity.Entity
	it := em.IterAvailable()
	for obj, itok := it(); itok; obj, itok = it() {
		if obj.GetType() == "timer" {
			timer = obj
			break
		}
	}
	start := time.Now()
	// update start time using ticker
	go func() {
		for t := range ticker.C {
			if component, ok := timer.GetComponent("font"); ok {
				font := component.(*entity.FontComponent)
				// fmt.Println(t.Format("15:04:05"))
				now := t.Sub(start)
				fmtNow := fmt.Sprintf("%02d:%02d:%02d", math.Round64(now.Hours()), math.Round64(now.Minutes()), math.Round64(now.Seconds()))
				font.Text = fmtNow
				font.Modified = true
			}
		}
	}()
	return ticker
}
