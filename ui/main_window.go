package ui

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"walk-gui-designer-pro/logic"
	"walk-gui-designer-pro/models"
)

type DesignerWindow struct {
	*walk.MainWindow
	Project           *models.Project
	SelectedComponent *models.Component
	PendingTool       string

	// Actions
	ActionShowGrid   *walk.Action
	ActionSnapToGrid *walk.Action

	// UI Components (to be implemented in other files)
	Canvas         *CanvasWidget
	PropertyEditor *PropertyEditorWidget
	ObjectList     *ObjectListWidget
	Toolbox        *ToolboxWidget
}

func (mw *DesignerWindow) NewProject() {
	mw.Project = &models.Project{
		Title:      "Untitled",
		Width:      800,
		Height:     600,
		Components: []*models.Component{},
	}
	mw.SelectedComponent = nil
	mw.RefreshAll()
}

func (mw *DesignerWindow) OpenProject() {
	dlg := new(walk.FileDialog)
	dlg.Title = "Open Project"
	dlg.Filter = "JSON files (*.json)|*.json|All files (*.*)|*.*"

	if ok, err := dlg.ShowOpen(mw); err != nil {
		return
	} else if !ok {
		return
	}

	p, err := logic.LoadProject(dlg.FilePath)
	if err != nil {
		walk.MsgBox(mw, "Error", "Failed to load project: "+err.Error(), walk.MsgBoxIconError)
		return
	}

	mw.Project = p
	mw.SelectedComponent = nil
	mw.RefreshAll()
}

func (mw *DesignerWindow) SaveProject() {
	if mw.Project == nil {
		return
	}

	dlg := new(walk.FileDialog)
	dlg.Title = "Save Project"
	dlg.Filter = "JSON files (*.json)|*.json|All files (*.*)|*.*"

	if ok, err := dlg.ShowSave(mw); err != nil {
		return
	} else if !ok {
		return
	}

	if err := logic.SaveProject(mw.Project, dlg.FilePath); err != nil {
		walk.MsgBox(mw, "Error", "Failed to save project: "+err.Error(), walk.MsgBoxIconError)
	}
}

func (mw *DesignerWindow) ImportUI() {
	dlg := new(walk.FileDialog)
	dlg.Title = "Import Qt .ui File"
	dlg.Filter = "Qt UI files (*.ui)|*.ui|All files (*.*)|*.*"

	if ok, err := dlg.ShowOpen(mw); err != nil {
		return
	} else if !ok {
		return
	}

	p, err := logic.ImportFromUI(dlg.FilePath)
	if err != nil {
		walk.MsgBox(mw, "Error", "Failed to import UI file: "+err.Error(), walk.MsgBoxIconError)
		return
	}

	mw.Project = p
	mw.SelectedComponent = nil
	mw.RefreshAll()
}

func (mw *DesignerWindow) ExportUI() {
	if mw.Project == nil {
		return
	}

	dlg := new(walk.FileDialog)
	dlg.Title = "Export to Qt .ui File"
	dlg.Filter = "Qt UI files (*.ui)|*.ui|All files (*.*)|*.*"

	if ok, err := dlg.ShowSave(mw); err != nil {
		return
	} else if !ok {
		return
	}

	if err := logic.ExportToUI(mw.Project, dlg.FilePath); err != nil {
		walk.MsgBox(mw, "Error", "Failed to export UI file: "+err.Error(), walk.MsgBoxIconError)
	} else {
		walk.MsgBox(mw, "Success", "Project successfully exported to .ui format.", walk.MsgBoxIconInformation)
	}
}

func (mw *DesignerWindow) RefreshCanvas() {
	if mw.Canvas != nil {
		mw.Canvas.Invalidate()
	}
}

func (mw *DesignerWindow) RefreshAll() {
	if mw.Project != nil {
		logic.ApplyLayout(mw.Project)
	}
	mw.RefreshCanvas()
	if mw.ObjectList != nil {
		mw.ObjectList.Refresh()
	}
	if mw.PropertyEditor != nil {
		mw.PropertyEditor.Update()
	}
}

func (mw *DesignerWindow) ShowPreview() {
	logic.ShowPreview(mw, mw.Project)
}

func (mw *DesignerWindow) GenerateCode() {
	if mw.Project == nil {
		return
	}
	code := logic.GenerateGoCode(mw.Project)

	// Show in a dialog
	var dlg *walk.Dialog

	(Dialog{
		AssignTo: &dlg,
		Title:    "Generated Go Code",
		Size:     Size{Width: 600, Height: 400},
		Layout:   VBox{},
		Children: []Widget{
			TextEdit{
				Text:     code,
				ReadOnly: true,
				VScroll:  true,
			},
			PushButton{
				Text:      "Close",
				OnClicked: func() { dlg.Accept() },
			},
		},
	}).Run(mw)
}

func (mw *DesignerWindow) DeleteSelected() {
	if mw.SelectedComponent == nil {
		return
	}

	// Collect IDs to delete (self + recursive children)
	toDelete := make(map[string]bool)
	var collect func(string)
	collect = func(id string) {
		toDelete[id] = true
		for _, c := range mw.Project.Components {
			if c.ParentID == id {
				collect(c.ID)
			}
		}
	}
	collect(mw.SelectedComponent.ID)

	// Filter project components
	newComps := []*models.Component{}
	for _, c := range mw.Project.Components {
		if !toDelete[c.ID] {
			newComps = append(newComps, c)
		}
	}
	mw.Project.Components = newComps
	mw.SelectedComponent = nil
	mw.RefreshAll()
}

func (mw *DesignerWindow) ShowAbout() {
	walk.MsgBox(mw, "About", "Walk GUI Designer Pro v1.0\nA senior Go desktop engineer production.", walk.MsgBoxIconInformation)
}

func (mw *DesignerWindow) SelectComponent(c *models.Component) {
	mw.SelectedComponent = c
	mw.RefreshAll()
}
