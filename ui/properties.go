package ui

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"walk-gui-designer-pro/models"
)

type PropertyEditorWidget struct {
	*walk.Composite
	mw *DesignerWindow
	db *walk.DataBinder

	// Property Groups
	groupCommon    *walk.GroupBox
	groupNumeric   *walk.GroupBox
	groupContainer *walk.GroupBox
	groupChoices   *walk.GroupBox
}

func CreatePropertyEditor(mw *DesignerWindow) Widget {
	pe := &PropertyEditorWidget{mw: mw}
	mw.PropertyEditor = pe

	return Composite{
		AssignTo: &pe.Composite,
		Layout:   VBox{MarginsZero: true},
		DataBinder: DataBinder{
			AssignTo:   &pe.db,
			Name:       "c",
			DataSource: &models.Component{},
		},
		Children: []Widget{
			Label{Text: "Properties", Font: Font{Bold: true}},

			GroupBox{
				AssignTo: &pe.groupCommon,
				Title:    "Common",
				Layout:   Grid{Columns: 2},
				Children: []Widget{
					Label{Text: "Name:"}, LineEdit{Text: Bind("c.Name")},
					Label{Text: "Text:"}, LineEdit{Text: Bind("c.Text")},
					Label{Text: "X:"}, NumberEdit{Value: Bind("c.X")},
					Label{Text: "Y:"}, NumberEdit{Value: Bind("c.Y")},
					Label{Text: "Width:"}, NumberEdit{Value: Bind("c.Width")},
					Label{Text: "Height:"}, NumberEdit{Value: Bind("c.Height")},
					Label{Text: "Enabled:"}, CheckBox{Checked: Bind("c.Enabled")},
					Label{Text: "Visible:"}, CheckBox{Checked: Bind("c.Visible")},
				},
			},

			GroupBox{
				AssignTo: &pe.groupNumeric,
				Title:    "Numeric/Range",
				Layout:   Grid{Columns: 2},
				Children: []Widget{
					Label{Text: "Min:"}, NumberEdit{Value: Bind("c.Min")},
					Label{Text: "Max:"}, NumberEdit{Value: Bind("c.Max")},
					Label{Text: "Value:"}, NumberEdit{Value: Bind("c.Value")},
				},
			},

			GroupBox{
				AssignTo: &pe.groupChoices,
				Title:    "Choices/State",
				Layout:   Grid{Columns: 2},
				Children: []Widget{
					Label{Text: "Checked:"}, CheckBox{Checked: Bind("c.Checked")},
					Label{Text: "ReadOnly:"}, CheckBox{Checked: Bind("c.ReadOnly")},
				},
			},

			GroupBox{
				AssignTo: &pe.groupContainer,
				Title:    "Container/Layout",
				Layout:   Grid{Columns: 2},
				Children: []Widget{
					Label{Text: "Layout:"},
					ComboBox{
						Model: []string{"", "HBox", "VBox", "Grid"},
						Value: Bind("c.Layout"),
					},
					Label{Text: "Columns:"}, NumberEdit{Value: Bind("c.Columns")},
				},
			},

			PushButton{
				Text: "Apply Changes",
				OnClicked: func() {
					pe.db.Submit()
					mw.RefreshAll()
				},
			},
		},
	}
}

func (pe *PropertyEditorWidget) Update() {
	c := pe.mw.SelectedComponent
	if c == nil {
		pe.SetVisible(false)
		return
	}
	pe.SetVisible(true)
	pe.db.SetDataSource(c)
	pe.db.Reset()

	// Filter visible groups
	isNumeric := false
	isContainer := false
	hasChoices := false

	switch c.Type {
	case "NumberEdit", "Slider", "ProgressBar":
		isNumeric = true
	case "Composite", "GroupBox", "Splitter":
		isContainer = true
	case "CheckBox", "RadioButton", "LineEdit", "TextEdit":
		hasChoices = true
	}

	pe.groupNumeric.SetVisible(isNumeric)
	pe.groupContainer.SetVisible(isContainer)
	pe.groupChoices.SetVisible(hasChoices)
}
