package cscene

import (
	// "github.com/tubelz/crazybird/csystem"
	"github.com/tubelz/macaw"
	"github.com/tubelz/macaw/entity"
	"github.com/tubelz/macaw/system"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"log"
)

// MenuScene is responsible to manage the content of the menu scene
type MenuScene struct {
	Scene         *macaw.Scene
	entityIDs     []uint16
	EntityManager *entity.Manager
}

// Init initialize this scene
func (m *MenuScene) Init(renderSystem *system.RenderSystem, font *ttf.Font) {
	initFunc := initializeEntities(m.EntityManager, renderSystem, font)
	bgColor := sdl.Color{32, 180, 230, 255}
	scene := &macaw.Scene{
		Name:         "menu",
		InitFunc:     initFunc,
		ExitFunc:     m.Exit,
		SceneOptions: macaw.SceneOptions{BgColor: bgColor, HideCursor: true},
	}
	scene.AddRenderSystem(renderSystem)
	m.Scene = scene
}

// Exit clear the entities created for this scene
func (m *MenuScene) Exit() {
	for id, obj := range m.EntityManager.GetAll() {
		if id > 0 && obj != nil {
			log.Printf("delete: %v", obj.GetID())
			m.EntityManager.Delete(obj.GetID())
		}
	}
	// var id uint16
	// for _, id = range m.entityIDs {
	// 	m.EntityManager.Delete(id)
	// }
}

// GetScene returns the scene from MenuScene
func (m *MenuScene) GetScene() *macaw.Scene {
	return m.Scene
}

func (m *MenuScene) addEntity(ent entity.Entity) {
	m.entityIDs = append(m.entityIDs, ent.GetID())
}

// func initializeEntities
func initializeEntities(em *entity.Manager, renderSystem *system.RenderSystem, font *ttf.Font) func() {
	return func() {
		title := em.Create("title")
		selectbox := em.Create("menuselectbox")
		start := em.Create("menuoption")
		credits := em.Create("menuoption")

		title.AddComponent(&entity.PositionComponent{&sdl.Point{280, 20}})
		title.AddComponent(&entity.FontComponent{Text: "crazy bird", Modified: true, Font: font})
		title.AddComponent(&entity.RenderComponent{RenderType: entity.RTFont})

		selectbox.AddComponent(&entity.PositionComponent{&sdl.Point{280, 399}})
		selectbox.AddComponent(&entity.RenderComponent{RenderType: entity.RTGeometry})
		selectbox.AddComponent(&entity.RectangleComponent{
			Size:   &sdl.Point{140, 22},
			Color:  &sdl.Color{0xC0, 0xC0, 0xC0, 0x99},
			Filled: true,
		})

		start.AddComponent(&entity.PositionComponent{&sdl.Point{300, 400}})
		start.AddComponent(&entity.FontComponent{Text: "start", Modified: true, Font: font})
		start.AddComponent(&entity.RenderComponent{RenderType: entity.RTFont})

		credits.AddComponent(&entity.PositionComponent{&sdl.Point{300, 430}})
		credits.AddComponent(&entity.FontComponent{Text: "credits", Modified: true, Font: font})
		credits.AddComponent(&entity.RenderComponent{RenderType: entity.RTFont})
	}
}
