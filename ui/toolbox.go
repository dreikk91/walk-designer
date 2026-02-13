package ui

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"walk-gui-designer-pro/models"
)

type ToolboxWidget struct {
	*walk.Composite
	mw      *DesignerWindow
	listBox *walk.ListBox
}

func CreateToolbox(mw *DesignerWindow) Widget {
	tb := &ToolboxWidget{mw: mw}
	mw.Toolbox = tb
	return Composite{
		AssignTo: &tb.Composite,
		Layout:   VBox{MarginsZero: true},
		Children: []Widget{
			Label{Text: "Toolbox", Font: Font{Bold: true}},
			ListBox{
				AssignTo: &tb.listBox,
				Model: []string{
					"PushButton", "Label", "LineEdit", "TextEdit",
					"CheckBox", "RadioButton", "ComboBox", "ListBox",
					"NumberEdit", "DateEdit", "Slider", "ProgressBar",
					"TableView", "TreeView", "TabWidget",
					"Composite", "GroupBox", "Splitter",
					"HSpacer", "VSpacer",
				},
				OnItemActivated: tb.onItemActivated,
				OnCurrentIndexChanged: func() {
					idx := tb.listBox.CurrentIndex()
					if idx >= 0 {
						tb.mw.PendingTool = tb.listBox.Model().([]string)[idx]
					}
				},
			},
		},
	}
}

func (tb *ToolboxWidget) onItemActivated() {
	item := tb.listBox.Model().([]string)[tb.listBox.CurrentIndex()]
	tb.mw.AddComponent(item, 50, 50)
}

func (mw *DesignerWindow) AddComponent(typ string, x, y int) *models.Component {
	if mw.Project == nil {
		mw.NewProject()
	}

	id := fmt.Sprintf("%s_%d", typ, len(mw.Project.Components)+1)
	c := &models.Component{
		ID:      id,
		Type:    typ,
		Name:    id,
		Text:    typ,
		X:       x,
		Y:       y,
		Width:   100,
		Height:  30,
		Enabled: true,
		Visible: true,
	}

	// Adjust defaults based on type
	switch typ {
	case "Label":
		c.Width = 100
		c.Height = 20
	case "TextEdit":
		c.Width = 200
		c.Height = 100
	case "Composite", "GroupBox":
		c.Width = 200
		c.Height = 150
		c.Layout = models.LayoutVBox
	case "HSpacer", "VSpacer":
		c.Text = ""
	}

	mw.Project.Components = append(mw.Project.Components, c)
	mw.SelectComponent(c)
	return c
}
