package ui

import (
	"fmt"
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

	// Groups for dynamic visibility
	grpSpecific   *walk.GroupBox
	grpNumeric    *walk.GroupBox
	grpItems      *walk.GroupBox
	grpLayout     *walk.GroupBox
	grpAppearance *walk.GroupBox
	grpTableCols  *walk.GroupBox
	grpTabPages   *walk.GroupBox

	// Appearance
	fontFamily    *walk.LineEdit
	fontSize      *walk.NumberEdit
	fontBold      *walk.CheckBox
	fontItalic    *walk.CheckBox

	// Table Cols / Tabs
	tableColsEdit *walk.TextEdit
	tabPagesEdit  *walk.TextEdit
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
						if c != nil {
							if c.Type == "ComboBox" || c.Type == "ListBox" || c.Type == "TreeView" {
								lines := strings.Split(pe.itemsEdit.Text(), "\n")
								c.Items = nil
								for _, l := range lines {
									if strings.TrimSpace(l) != "" {
										c.Items = append(c.Items, strings.TrimSpace(l))
									}
								}
							}

							if c.Type == "TableView" {
								lines := strings.Split(pe.tableColsEdit.Text(), "\n")
								c.TableCols = nil
								for _, l := range lines {
									l = strings.TrimSpace(l)
									if l == "" { continue }
									parts := strings.Split(l, ",")
									if len(parts) > 0 {
										col := models.TableColumnConfig{Name: parts[0], Title: parts[0], Width: 100, Alignment: "Near"}
										if len(parts) > 1 { col.Title = parts[1] }
										if len(parts) > 2 {
											var w int
											fmt.Sscanf(parts[2], "%d", &w)
											if w > 0 { col.Width = w }
										}
										if len(parts) > 3 { col.Alignment = parts[3] }
										c.TableCols = append(c.TableCols, col)
									}
								}
							}

							if c.Type == "TabWidget" {
								lines := strings.Split(pe.tabPagesEdit.Text(), "\n")
								c.TabPages = nil
								for _, l := range lines {
									l = strings.TrimSpace(l)
									if l != "" {
										c.TabPages = append(c.TabPages, models.TabPageConfig{Title: l, Layout: &models.LayoutConfig{Type: models.LayoutVBox}})
									}
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
							Label{Text: "Min Width:"}, NumberEdit{Value: Bind("c.MinWidth")},
							Label{Text: "Min Height:"}, NumberEdit{Value: Bind("c.MinHeight")},
							Label{Text: "Max Width:"}, NumberEdit{Value: Bind("c.MaxWidth")},
							Label{Text: "Max Height:"}, NumberEdit{Value: Bind("c.MaxHeight")},
							Label{Text: "Alignment:"}, ComboBox{
								Model: []string{"", "Left", "Center", "Right"},
								Value: Bind("c.Alignment"),
							},
							Label{Text: "Enabled:"}, CheckBox{Checked: Bind("c.Enabled")},
							Label{Text: "Visible:"}, CheckBox{Checked: Bind("c.Visible")},
						},
					},

					GroupBox{
						AssignTo: &pe.grpAppearance,
						Title:    "Appearance",
						Layout:   Grid{Columns: 2},
						Children: []Widget{
							Label{Text: "Font Family:"}, LineEdit{AssignTo: &pe.fontFamily},
							Label{Text: "Font Size:"}, NumberEdit{AssignTo: &pe.fontSize},
							Label{Text: "Bold:"}, CheckBox{AssignTo: &pe.fontBold},
							Label{Text: "Italic:"}, CheckBox{AssignTo: &pe.fontItalic},
							Label{Text: "Bg Color:"}, LineEdit{Text: Bind("c.BgColor"), ToolTipText: "e.g. 255,255,255"},
							Label{Text: "Text Color:"}, LineEdit{Text: Bind("c.TextColor"), ToolTipText: "e.g. 0,0,0"},
						},
					},

					GroupBox{
						AssignTo: &pe.grpSpecific,
						Title:    "Specific Settings",
						Layout:   Grid{Columns: 2},
						Children: []Widget{
							Label{Text: "Checked:"}, CheckBox{Checked: Bind("c.Checked")},
							Label{Text: "ReadOnly:"}, CheckBox{Checked: Bind("c.ReadOnly")},
							Label{Text: "Password:"}, CheckBox{Checked: Bind("c.IsPassword")},
							Label{Text: "Max Length:"}, NumberEdit{Value: Bind("c.MaxLength")},
							Label{Text: "Orientation:"}, ComboBox{
								Model: []string{"Horizontal", "Vertical"},
								Value: Bind("c.Orientation"),
							},
							Label{Text: "Image Mode:"}, ComboBox{
								Model: []string{"Ideal", "Corner", "Center", "Shrink", "Zoom", "Stretch"},
								Value: Bind("c.ImageMode"),
							},
							Label{Text: "ImagePath:"}, LineEdit{Text: Bind("c.ImagePath")},
							Label{Text: "Marquee:"}, CheckBox{Checked: Bind("c.Marquee")},
						},
					},

					GroupBox{
						AssignTo: &pe.grpNumeric,
						Title:    "Numeric/Range",
						Layout:   Grid{Columns: 2},
						Children: []Widget{
							Label{Text: "Min:"}, NumberEdit{Value: Bind("c.MinValue")},
							Label{Text: "Max:"}, NumberEdit{Value: Bind("c.MaxValue")},
							Label{Text: "Value:"}, NumberEdit{Value: Bind("c.Value")},
						},
					},

					GroupBox{
						AssignTo: &pe.grpItems,
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
						AssignTo: &pe.grpTableCols,
						Title:    "Table Columns (Name,Title,Width,Align)",
						Layout:   VBox{},
						Children: []Widget{
							TextEdit{
								AssignTo: &pe.tableColsEdit,
								VScroll: true,
								MinSize: Size{Height: 80},
								ToolTipText: "e.g.: col1,First Name,100,Near",
							},
						},
					},

					GroupBox{
						AssignTo: &pe.grpTabPages,
						Title:    "Tab Pages (one per line)",
						Layout:   VBox{},
						Children: []Widget{
							TextEdit{
								AssignTo: &pe.tabPagesEdit,
								VScroll: true,
								MinSize: Size{Height: 80},
								ToolTipText: "e.g.: Tab 1\nTab 2",
							},
						},
					},

					GroupBox{
						AssignTo: &pe.grpLayout,
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
							c := mw.SelectedComponent
							if c != nil {
								// Sync Layout
								if pe.grpLayout.Visible() {
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
								// Sync Font
								fFam := pe.fontFamily.Text()
								fSz := int(pe.fontSize.Value())
								if fFam != "" || fSz > 0 || pe.fontBold.Checked() || pe.fontItalic.Checked() {
									if c.Font == nil { c.Font = &models.FontConfig{} }
									c.Font.Family = fFam
									c.Font.PointSize = fSz
									c.Font.Bold = pe.fontBold.Checked()
									c.Font.Italic = pe.fontItalic.Checked()
								} else {
									c.Font = nil
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

	// Update Group Visibility based on type
	isContainer := c.Type == "Composite" || c.Type == "GroupBox" || c.Type == "Splitter" || c.Type == "ScrollView" || c.Type == "TabWidget"
	isNumeric := c.Type == "Slider" || c.Type == "ProgressBar" || c.Type == "NumberEdit"
	isList := c.Type == "ComboBox" || c.Type == "ListBox" || c.Type == "TreeView"
	isTable := c.Type == "TableView"
	isTabWidget := c.Type == "TabWidget"

	pe.grpLayout.SetVisible(isContainer)
	pe.grpNumeric.SetVisible(isNumeric)
	pe.grpItems.SetVisible(isList)
	pe.grpTableCols.SetVisible(isTable)
	pe.grpTabPages.SetVisible(isTabWidget)

	// Ensure Layout exists for binding if we were to bind to it
	if isContainer && c.Layout == nil {
		c.Layout = &models.LayoutConfig{Type: models.LayoutVBox}
	}

	pe.db.SetDataSource(c)
	pe.db.Reset()

	// Sync Lists/Tables/Tabs
	if isList && c.Items != nil {
		pe.itemsEdit.SetText(strings.Join(c.Items, "\n"))
	} else {
		pe.itemsEdit.SetText("")
	}

	if isTable && c.TableCols != nil {
		lines := []string{}
		for _, col := range c.TableCols {
			lines = append(lines, col.Name+","+col.Title+",100,"+col.Alignment)
		}
		pe.tableColsEdit.SetText(strings.Join(lines, "\n"))
	} else {
		pe.tableColsEdit.SetText("")
	}

	if isTabWidget && c.TabPages != nil {
		lines := []string{}
		for _, page := range c.TabPages {
			lines = append(lines, page.Title)
		}
		pe.tabPagesEdit.SetText(strings.Join(lines, "\n"))
	} else {
		pe.tabPagesEdit.SetText("")
	}

	// Sync Font
	if c.Font != nil {
		pe.fontFamily.SetText(c.Font.Family)
		pe.fontSize.SetValue(float64(c.Font.PointSize))
		pe.fontBold.SetChecked(c.Font.Bold)
		pe.fontItalic.SetChecked(c.Font.Italic)
	} else {
		pe.fontFamily.SetText("")
		pe.fontSize.SetValue(0)
		pe.fontBold.SetChecked(false)
		pe.fontItalic.SetChecked(false)
	}

	// Sync Layout
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
