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

// CreditsScene is responsible to manage the content of the menu scene
type CreditsScene struct {
	Scene         *macaw.Scene
	entityIDs     []uint16
	EntityManager *entity.Manager
}

// Init initialize this scene
func (m *CreditsScene) Init(renderSystem *system.RenderSystem, font *ttf.Font) {
	initFunc := initializeCreditEntities(m.EntityManager, renderSystem, font)
	bgColor := sdl.Color{32, 180, 230, 255}
	scene := &macaw.Scene{
		Name:         "credits",
		InitFunc:     initFunc,
		ExitFunc:     m.Exit,
		SceneOptions: macaw.SceneOptions{BgColor: bgColor, HideCursor: true},
	}
	scene.AddRenderSystem(renderSystem)
	m.Scene = scene
}

// Exit clear the entities created for this scene
func (m *CreditsScene) Exit() {
	for id, obj := range m.EntityManager.GetAll() {
		if id > 0 && obj != nil {
			log.Printf("delete: %v", obj.GetID())
			m.EntityManager.Delete(obj.GetID())
		}
	}
}

// GetScene returns the scene from CreditsScene
func (m *CreditsScene) GetScene() *macaw.Scene {
	return m.Scene
}

func (m *CreditsScene) addEntity(ent entity.Entity) {
	m.entityIDs = append(m.entityIDs, ent.GetID())
}

// func initializeEntities
func initializeCreditEntities(em *entity.Manager, renderSystem *system.RenderSystem, font *ttf.Font) func() {
	return func() {
		title := em.Create("title")
		developerSession := em.Create("session")
		developer := em.Create("content")
		musicSession := em.Create("session")
		music := em.Create("content")
		extraSession := em.Create("session")
		extra := em.Create("content")

		selectbox := em.Create("creditsselectbox")
		backButton := em.Create("select")

		selectbox.AddComponent("position", &entity.PositionComponent{&sdl.Point{280, 399}})
		selectbox.AddComponent("geometry", &entity.RectangleComponent{
			Size:   &sdl.Point{140, 22},
			Color:  &sdl.Color{0xC0, 0xC0, 0xC0, 0x99},
			Filled: true,
		})

		backButton.AddComponent("position", &entity.PositionComponent{&sdl.Point{300, 400}})
		backButton.AddComponent("font", &entity.FontComponent{Text: "back", Modified: true, Font: font})
		backButton.AddComponent("render", &entity.RenderComponent{Renderer: renderSystem.Renderer})

		title.AddComponent("font", &entity.FontComponent{Text: "Credits", Modified: true, Font: font})
		developerSession.AddComponent("font", &entity.FontComponent{Text: "Developer, Artist and Producer", Modified: true, Font: font})
		developer.AddComponent("font", &entity.FontComponent{Text: "Marcus Renno - @marcusrenno - marcusrenno.tech", Modified: true, Font: font})
		musicSession.AddComponent("font", &entity.FontComponent{Text: "Music", Modified: true, Font: font})
		music.AddComponent("font", &entity.FontComponent{Text: "Avaren - @avarenmusic", Modified: true, Font: font})
		extraSession.AddComponent("font", &entity.FontComponent{Text: "Extra", Modified: true, Font: font})
		extra.AddComponent("font", &entity.FontComponent{Text: "font by codeman38 - cody@zone38.net", Modified: true, Font: font})

		objects := []*entity.Entity{
			title, developerSession, developer, musicSession, music, extraSession, extra}

		for i, obj := range objects {
			pos := int32(i)
			if i == 0 || i%2 == 1 {
				addPos(obj, pos, 0)
			} else {
				addPos(obj, pos, 1)
			}
			addRender(obj, renderSystem.Renderer)
		}
	}
}

func addPos(obj *entity.Entity, ypos int32, xpos int32) {
	obj.AddComponent("position", &entity.PositionComponent{&sdl.Point{20 + xpos*20, 20 + ypos*40}})
}

func addRender(obj *entity.Entity, render *sdl.Renderer) {
	obj.AddComponent("render", &entity.RenderComponent{Renderer: render})
}
