package main

import (
	"log"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"walk-gui-designer-pro/ui"
)

func main() {
	mw := new(ui.DesignerWindow)

	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "Walk GUI Designer Pro",
		Size:     Size{Width: 1200, Height: 800},
		Layout:   VBox{MarginsZero: true, SpacingZero: true},
		MenuItems: []MenuItem{
			Menu{
				Text: "&File",
				Items: []MenuItem{
					Action{
						Text:        "&New",
						Shortcut:    Shortcut{Modifiers: walk.ModControl, Key: walk.KeyN},
						OnTriggered: mw.NewProject,
					},
					Action{
						Text:        "&Open...",
						Shortcut:    Shortcut{Modifiers: walk.ModControl, Key: walk.KeyO},
						OnTriggered: mw.OpenProject,
					},
					Action{
						Text:        "&Save",
						Shortcut:    Shortcut{Modifiers: walk.ModControl, Key: walk.KeyS},
						OnTriggered: mw.SaveProject,
					},
					Separator{},
					Action{
						Text:        "E&xit",
						OnTriggered: func() { mw.Close() },
					},
				},
			},
			Menu{
				Text: "&View",
				Items: []MenuItem{
					Action{
						AssignTo:    &mw.ActionShowGrid,
						Text:        "Show &Grid",
						Checkable:   true,
						Checked:     true,
						OnTriggered: mw.RefreshCanvas,
					},
					Action{
						AssignTo:  &mw.ActionSnapToGrid,
						Text:      "&Snap to Grid",
						Checkable: true,
						Checked:   true,
					},
				},
			},
			Menu{
				Text: "&Design",
				Items: []MenuItem{
					Action{
						Text:        "&Preview",
						Shortcut:    Shortcut{Modifiers: walk.ModControl, Key: walk.KeyP},
						OnTriggered: mw.ShowPreview,
					},
					Action{
						Text:        "&Generate Code",
						Shortcut:    Shortcut{Modifiers: walk.ModControl, Key: walk.KeyG},
						OnTriggered: mw.GenerateCode,
					},
					Separator{},
					Action{
						Text:        "&Delete Element",
						Shortcut:    Shortcut{Key: walk.KeyDelete},
						OnTriggered: mw.DeleteSelected,
					},
				},
			},
			Menu{
				Text: "&Help",
				Items: []MenuItem{
					Action{
						Text:        "&About",
						OnTriggered: mw.ShowAbout,
					},
				},
			},
		},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					ui.CreateToolbox(mw),
					HSplitter{
						StretchFactor: 3,
						Children: []Widget{
							ui.CreateCanvas(mw),
							VSplitter{
								Children: []Widget{
									ui.CreateObjectList(mw),
									ui.CreatePropertyEditor(mw),
								},
							},
						},
					},
				},
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}

	mw.Run()
}
