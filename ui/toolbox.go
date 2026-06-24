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
	treeView *walk.TreeView
	listBox *walk.ListBox // Kept for compatibility if needed, but we'll use treeView
}

type ToolboxItem struct {
	parent   *ToolboxItem
	children []*ToolboxItem
	text     string
	typ      string
}

func (i *ToolboxItem) Text() string { return i.text }
func (i *ToolboxItem) Parent() walk.TreeItem {
	if i.parent == nil { return nil }
	return i.parent
}
func (i *ToolboxItem) ChildCount() int { return len(i.children) }
func (i *ToolboxItem) ChildAt(index int) walk.TreeItem { return i.children[index] }
func (i *ToolboxItem) Image() interface{} { return nil }

type ToolboxModel struct {
	walk.TreeModelBase
	roots []*ToolboxItem
}

func (m *ToolboxModel) RootCount() int { return len(m.roots) }
func (m *ToolboxModel) RootAt(index int) walk.TreeItem { return m.roots[index] }

func CreateToolbox(mw *DesignerWindow) Widget {
	tb := &ToolboxWidget{mw: mw}
	mw.Toolbox = tb

	model := &ToolboxModel{}

	basic := &ToolboxItem{text: "Basic Controls"}
	basic.children = []*ToolboxItem{
		{parent: basic, text: "PushButton", typ: "PushButton"},
		{parent: basic, text: "SplitButton", typ: "SplitButton"},
		{parent: basic, text: "Label", typ: "Label"},
		{parent: basic, text: "LinkLabel", typ: "LinkLabel"},
		{parent: basic, text: "ImageView", typ: "ImageView"},
		{parent: basic, text: "LineEdit", typ: "LineEdit"},
		{parent: basic, text: "TextEdit", typ: "TextEdit"},
		{parent: basic, text: "CheckBox", typ: "CheckBox"},
		{parent: basic, text: "RadioButton", typ: "RadioButton"},
		{parent: basic, text: "ComboBox", typ: "ComboBox"},
		{parent: basic, text: "ListBox", typ: "ListBox"},
	}

	menus := &ToolboxItem{text: "Menus & Toolbars"}
	menus.children = []*ToolboxItem{
		{parent: menus, text: "ToolBar", typ: "ToolBar"},
		{parent: menus, text: "StatusBar", typ: "StatusBar"},
	}

	numeric := &ToolboxItem{text: "Numeric & Date"}
	numeric.children = []*ToolboxItem{
		{parent: numeric, text: "NumberEdit", typ: "NumberEdit"},
		{parent: numeric, text: "DateEdit", typ: "DateEdit"},
		{parent: numeric, text: "Slider", typ: "Slider"},
		{parent: numeric, text: "ProgressBar", typ: "ProgressBar"},
	}

	advanced := &ToolboxItem{text: "Advanced"}
	advanced.children = []*ToolboxItem{
		{parent: advanced, text: "TableView", typ: "TableView"},
		{parent: advanced, text: "TreeView", typ: "TreeView"},
		{parent: advanced, text: "TabWidget", typ: "TabWidget"},
	}

	containers := &ToolboxItem{text: "Containers"}
	containers.children = []*ToolboxItem{
		{parent: containers, text: "GroupBox", typ: "GroupBox"},
		{parent: containers, text: "Composite", typ: "Composite"},
		{parent: containers, text: "ScrollView", typ: "ScrollView"},
		{parent: containers, text: "Splitter", typ: "Splitter"},
	}

	spacers := &ToolboxItem{text: "Spacers"}
	spacers.children = []*ToolboxItem{
		{parent: spacers, text: "HSpacer", typ: "HSpacer"},
		{parent: spacers, text: "VSpacer", typ: "VSpacer"},
	}

	model.roots = []*ToolboxItem{basic, menus, numeric, advanced, containers, spacers}

	return Composite{
		AssignTo: &tb.Composite,
		Layout:   VBox{MarginsZero: true},
		Children: []Widget{
			Label{Text: "Toolbox", Font: Font{Bold: true}},
			TreeView{
				AssignTo: &tb.treeView,
				Model:    model,
				OnCurrentItemChanged: func() {
					item := tb.treeView.CurrentItem()
					if item != nil {
						ti := item.(*ToolboxItem)
						if ti.typ != "" {
							tb.mw.PendingTool = ti.typ
						} else {
							tb.mw.PendingTool = ""
						}
					}
				},
				OnItemActivated: func() {
					item := tb.treeView.CurrentItem()
					if item != nil {
						ti := item.(*ToolboxItem)
						if ti.typ != "" {
							tb.mw.AddComponent(ti.typ, 50, 50)
						}
					}
				},
			},
		},
	}
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
	case "PushButton":
		c.Text = "Button"
	case "SplitButton":
		c.Text = "SplitButton"
		c.Width = 120
	case "Label":
		c.Text = "Label"
		c.Width = 100
		c.Height = 20
	case "LinkLabel":
		c.Text = `<a>Link</a>`
		c.Width = 100
		c.Height = 20
	case "ImageView":
		c.Text = "Image"
		c.Width = 100
		c.Height = 100
	case "ToolBar":
		c.Width = 300
		c.Height = 30
	case "StatusBar":
		c.Width = 400
		c.Height = 25
	case "CheckBox":
		c.Text = "CheckBox"
	case "RadioButton":
		c.Text = "RadioButton"
	case "LineEdit":
		c.Height = 25
	case "TextEdit":
		c.Width = 200
		c.Height = 100
	case "ComboBox":
		c.Height = 25
		c.Items = []string{"Option 1", "Option 2"}
	case "ListBox":
		c.Width = 150
		c.Height = 120
		c.Items = []string{"Item 1", "Item 2"}
	case "NumberEdit":
		c.Height = 25
		c.MaxValue = 100
		c.Value = 50
	case "ProgressBar":
		c.Width = 200
		c.Height = 25
		c.MaxValue = 100
	case "Slider":
		c.Width = 200
		c.Height = 30
		c.MaxValue = 100
		c.Value = 50
		c.Orientation = "Horizontal"
	case "Composite", "GroupBox", "ScrollView":
		c.Width = 200
		c.Height = 150
		c.Layout = &models.LayoutConfig{Type: models.LayoutVBox}
	case "TabWidget":
		c.Width = 400
		c.Height = 300
		c.TabPages = []models.TabPageConfig{
			{Title: "Tab 1", Layout: &models.LayoutConfig{Type: models.LayoutVBox}},
			{Title: "Tab 2", Layout: &models.LayoutConfig{Type: models.LayoutVBox}},
		}
	case "Splitter":
		c.Width = 400
		c.Height = 300
		c.Orientation = "Horizontal"
	case "HSpacer", "VSpacer":
		c.Text = ""
		c.Width = 20
		c.Height = 20
	}

	mw.Project.Components = append(mw.Project.Components, c)
	mw.SelectComponent(c)
	return c
}
