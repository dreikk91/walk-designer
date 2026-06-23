package ui

import (
	"strings"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"walk-gui-designer-pro/models"
)

type PropertyEditorWidget struct {
	*walk.Composite
	mw *DesignerWindow
	db *walk.DataBinder

	// Custom fields that might need manual handling
	itemsEdit      *walk.TextEdit
	layoutType     *walk.ComboBox
	layoutMargins  *walk.NumberEdit
	layoutSpacing  *walk.NumberEdit
	marginsZero    *walk.CheckBox
	spacingZero    *walk.CheckBox
}

func CreatePropertyEditor(mw *DesignerWindow) Widget {
	pe := &PropertyEditorWidget{mw: mw}
	mw.PropertyEditor = pe

	return ScrollView{
		HorizontalFixed: true,
		Layout: VBox{MarginsZero: true},
		Children: []Widget{
			Composite{
				AssignTo: &pe.Composite,
				Layout:   VBox{MarginsZero: true},
				DataBinder: DataBinder{
					AssignTo:   &pe.db,
					Name:       "c",
					DataSource: &models.Component{},
					OnSubmitted: func() {
						c := mw.SelectedComponent
						if c != nil && c.Type == "ComboBox" || c.Type == "ListBox" {
							lines := strings.Split(pe.itemsEdit.Text(), "\n")
							c.Items = nil
							for _, l := range lines {
								if strings.TrimSpace(l) != "" {
									c.Items = append(c.Items, strings.TrimSpace(l))
								}
							}
						}
						mw.RefreshAll()
					},
				},
				Children: []Widget{
					Label{Text: "Properties", Font: Font{Bold: true}},

					GroupBox{
						Title:    "Common",
						Layout:   Grid{Columns: 2},
						Children: []Widget{
							Label{Text: "Name:"}, LineEdit{Text: Bind("c.Name")},
							Label{Text: "Text:"}, LineEdit{Text: Bind("c.Text")},
							Label{Text: "ToolTip:"}, LineEdit{Text: Bind("c.ToolTip")},
							Label{Text: "X:"}, NumberEdit{Value: Bind("c.X")},
							Label{Text: "Y:"}, NumberEdit{Value: Bind("c.Y")},
							Label{Text: "Width:"}, NumberEdit{Value: Bind("c.Width")},
							Label{Text: "Height:"}, NumberEdit{Value: Bind("c.Height")},
							Label{Text: "Alignment:"}, ComboBox{
								Model: []string{"", "Left", "Center", "Right"},
								Value: Bind("c.Alignment"),
							},
							Label{Text: "Enabled:"}, CheckBox{Checked: Bind("c.Enabled")},
							Label{Text: "Visible:"}, CheckBox{Checked: Bind("c.Visible")},
						},
					},

					GroupBox{
						Title:    "Specific Settings",
						Layout:   Grid{Columns: 2},
						Children: []Widget{
							Label{Text: "Checked:"}, CheckBox{Checked: Bind("c.Checked")},
							Label{Text: "ReadOnly:"}, CheckBox{Checked: Bind("c.ReadOnly")},
							Label{Text: "Password:"}, CheckBox{Checked: Bind("c.IsPassword")},
							Label{Text: "Orientation:"}, ComboBox{
								Model: []string{"Horizontal", "Vertical"},
								Value: Bind("c.Orientation"),
							},
							Label{Text: "ImagePath:"}, LineEdit{Text: Bind("c.ImagePath")},
						},
					},

					GroupBox{
						Title:    "Numeric/Range",
						Layout:   Grid{Columns: 2},
						Children: []Widget{
							Label{Text: "Min:"}, NumberEdit{Value: Bind("c.MinValue")},
							Label{Text: "Max:"}, NumberEdit{Value: Bind("c.MaxValue")},
							Label{Text: "Value:"}, NumberEdit{Value: Bind("c.Value")},
						},
					},

					GroupBox{
						Title:    "Items (one per line)",
						Layout:   VBox{},
						Children: []Widget{
							TextEdit{
								AssignTo: &pe.itemsEdit,
								VScroll: true,
								MinSize: Size{Height: 100},
							},
						},
					},

					GroupBox{
						Title:    "Container/Layout",
						Layout:   Grid{Columns: 2},
						Children: []Widget{
							Label{Text: "Layout Type:"},
							ComboBox{
								AssignTo: &pe.layoutType,
								Model: []string{"None", "HBox", "VBox", "Grid"},
							},
							Label{Text: "Columns:"}, NumberEdit{Value: Bind("c.ColumnSpan")},
							Label{Text: "Margins:"}, NumberEdit{AssignTo: &pe.layoutMargins},
							Label{Text: "Spacing:"}, NumberEdit{AssignTo: &pe.layoutSpacing},
							Label{Text: "Margins Zero:"}, CheckBox{AssignTo: &pe.marginsZero},
							Label{Text: "Spacing Zero:"}, CheckBox{AssignTo: &pe.spacingZero},
						},
					},

					PushButton{
						Text: "Apply Changes",
						OnClicked: func() {
							pe.db.Submit()
							// Manual sync for Layout since it's a pointer
							c := mw.SelectedComponent
							if c != nil {
								lt := models.LayoutType(pe.layoutType.Text())
								if lt != models.LayoutNone && lt != "" {
									if c.Layout == nil {
										c.Layout = &models.LayoutConfig{}
									}
									c.Layout.Type = lt
									c.Layout.Margins = int(pe.layoutMargins.Value())
									c.Layout.Spacing = int(pe.layoutSpacing.Value())
									c.Layout.MarginsZero = pe.marginsZero.Checked()
									c.Layout.SpacingZero = pe.spacingZero.Checked()
								} else {
									c.Layout = nil
								}
							}
							mw.RefreshAll()
						},
					},
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

	// Ensure Layout exists for binding if we were to bind to it
	if (c.Type == "Composite" || c.Type == "GroupBox" || c.Type == "Splitter" || c.Type == "ScrollView") && c.Layout == nil {
		c.Layout = &models.LayoutConfig{Type: models.LayoutVBox}
	}

	pe.db.SetDataSource(c)
	pe.db.Reset()

	if c.Items != nil {
		pe.itemsEdit.SetText(strings.Join(c.Items, "\n"))
	} else {
		pe.itemsEdit.SetText("")
	}

	if c.Layout != nil {
		pe.layoutType.SetText(string(c.Layout.Type))
		pe.layoutMargins.SetValue(float64(c.Layout.Margins))
		pe.layoutSpacing.SetValue(float64(c.Layout.Spacing))
		pe.marginsZero.SetChecked(c.Layout.MarginsZero)
		pe.spacingZero.SetChecked(c.Layout.SpacingZero)
	} else {
		pe.layoutType.SetText("None")
		pe.layoutMargins.SetValue(0)
		pe.layoutSpacing.SetValue(0)
		pe.marginsZero.SetChecked(false)
		pe.spacingZero.SetChecked(false)
	}
}
