package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/tubelz/crazybird/cscene"
	"github.com/tubelz/crazybird/csystem"
	"github.com/tubelz/macaw"
	"github.com/tubelz/macaw/entity"
	"github.com/tubelz/macaw/input"
	"github.com/tubelz/macaw/system"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	fmt.Println("Crazy Bird!")
	var err error
	err = macaw.Initialize()
	if err != nil {
		fmt.Println("Macaw could not initialize")
	}
	defer macaw.Quit()

	mfont := entity.MFont{File: "assets/font/manaspc.ttf", Size: uint8(22)}
	font := mfont.Open()
	defer mfont.Close()

	input := &input.Manager{}
	em := &entity.Manager{}

	rand.Seed(time.Now().UnixNano())
	systems := initializeSystems(input, em)
	initializeEntities(em)

	systems[0].(*system.RenderSystem).SetCamera(em.Get(0))

	gameLoop := initializeGameLoop(systems, em, input, font)
	addSceneManagerToSystems(systems, &gameLoop.SceneManager)
	// startTimer(em)
	gameLoop.Run()
}

func addSceneManagerToSystems(systems []system.Systemer, sm *macaw.SceneManager) {
	menu := systems[1].(*csystem.MenuSystem)
	menu.SetSceneManager(sm)
	credits := systems[2].(*csystem.CreditsSystem)
	credits.SetSceneManager(sm)
	player := systems[5].(*csystem.PlayerSystem)
	player.SetSceneManager(sm)
}

func initializeSystems(im *input.Manager, em *entity.Manager) []system.Systemer {
	render := &system.RenderSystem{Name: "render system", Window: macaw.Window, EntityManager: em}
	menu := &csystem.MenuSystem{Name: "menu system", EntityManager: em, InputManager: im}
	credits := &csystem.CreditsSystem{Name: "credits system", EntityManager: em, InputManager: im}
	physics := &system.PhysicsSystem{Name: "physics system", EntityManager: em}
	collision := &system.CollisionSystem{Name: "collision system", EntityManager: em}
	player := &csystem.PlayerSystem{Name: "player system", CollisionSystem: collision, InputManager: im, EntityManager: em, RenderSystem: render}
	scroll := &csystem.ScrollSystem{Name: "scroll system", EntityManager: em}
	spawn := &csystem.SpawnSystem{Name: "spawn system", EntityManager: em, RenderSystem: render, CollisionSystem: collision}

	systems := []system.Systemer{
		render,
		menu,
		credits,
		physics,
		collision,
		player,
		scroll,
		spawn,
	}

	// initialize render
	render.Init()

	return systems
}

func initializeGameLoop(systems []system.Systemer, em *entity.Manager, im *input.Manager, font *ttf.Font) *macaw.GameLoop {
	log.Println("")
	gameLoop := &macaw.GameLoop{InputManager: im}
	renderSystem := systems[0].(*system.RenderSystem)
	// Menu
	sceneMenu := &cscene.MenuScene{EntityManager: em}
	sceneMenu.Init(renderSystem, font)
	sceneMenu.Scene.AddGameUpdateSystem(systems[1])

	// Game itself
	sceneGame := &cscene.GameScene{EntityManager: em}
	sceneGame.Init(renderSystem, font)
	for _, system := range systems[3:] {
		sceneGame.Scene.AddGameUpdateSystem(system)
	}
	sceneGame.CollisionSystem = systems[4].(*system.CollisionSystem)
	// Game itself
	sceneCredits := &cscene.CreditsScene{EntityManager: em}
	sceneCredits.Init(renderSystem, font)
	sceneCredits.Scene.AddGameUpdateSystem(systems[2])

	gameLoop.AddScene(sceneMenu.GetScene())
	gameLoop.AddScene(sceneGame.GetScene())
	gameLoop.AddScene(sceneCredits.GetScene())

	return gameLoop
}

func initializeEntities(em *entity.Manager) {
	camera := em.Create("camera")
	camera.AddComponent(&entity.PositionComponent{&sdl.Point{0, 0}})
	camera.AddComponent(&entity.CameraComponent{
		ViewportSize: sdl.Point{800, 600},
		WorldSize:    sdl.Point{1145, 600},
	})
}
