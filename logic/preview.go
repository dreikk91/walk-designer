package logic

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"walk-gui-designer-pro/models"
)

func ShowPreview(owner walk.Form, p *models.Project) {
	if p == nil {
		return
	}

	var dlg *walk.Dialog

	children := createWidgets(p, "")

	err := (Dialog{
		AssignTo: &dlg,
		Title:    "Preview - " + p.Title,
		MinSize:  Size{Width: p.Width, Height: p.Height},
		Layout:   VBox{}, // Default or p.Layout
		Children: children,
	}.Create(owner))

	if err != nil {
		walk.MsgBox(owner, "Error", "Failed to create preview: "+err.Error(), walk.MsgBoxIconError)
		return
	}

	dlg.Run()
}

func createWidgets(p *models.Project, parentID string) []Widget {
	widgets := []Widget{}

	for _, c := range p.Components {
		if c.ParentID == parentID {
			var w Widget

			// Basic mapping
			switch c.Type {
			case "PushButton":
				w = PushButton{Text: c.Text}
			case "Label":
				w = Label{Text: c.Text}
			case "LineEdit":
				w = LineEdit{Text: c.Text, ReadOnly: c.ReadOnly}
			case "TextEdit":
				w = TextEdit{Text: c.Text, ReadOnly: c.ReadOnly}
			case "CheckBox":
				w = CheckBox{Text: c.Text, Checked: c.Checked}
			case "RadioButton":
				w = RadioButton{Text: c.Text}
			case "ComboBox":
				w = ComboBox{Model: c.Items}
			case "ListBox":
				w = ListBox{Model: c.Items}
			case "NumberEdit":
				w = NumberEdit{Value: float64(c.Value)}
			case "Slider":
				w = Slider{Value: c.Value, MinValue: c.Min, MaxValue: c.Max}
			case "ProgressBar":
				w = ProgressBar{Value: c.Value, MinValue: c.Min, MaxValue: c.Max}
			case "Composite":
				w = Composite{
					Layout:   getLayout(c.Layout),
					Children: createWidgets(p, c.ID),
				}
			case "Splitter":
				if c.Orientation == 2 { // Vertical
					w = VSplitter{
						Children: createWidgets(p, c.ID),
					}
				} else {
					w = HSplitter{
						Children: createWidgets(p, c.ID),
					}
				}
			case "TabWidget":
				w = TabWidget{
					Pages: createTabPages(p, c.ID),
				}
			case "TableView":
				w = TableView{
					Columns: []TableViewColumn{{Title: "Column 1"}},
				}
			case "TreeView":
				w = TreeView{}
			case "GroupBox":
				w = GroupBox{
					Title:    c.Text,
					Layout:   getLayout(c.Layout),
					Children: createWidgets(p, c.ID),
				}
			case "HSpacer":
				w = HSpacer{}
			case "VSpacer":
				w = VSpacer{}
			default:
				// Fallback
				w = Label{Text: "[" + c.Type + "] " + c.Text}
			}

			if w != nil {
				widgets = append(widgets, w)
			}
		}
	}

	return widgets
}

func createTabPages(p *models.Project, parentID string) []TabPage {
	pages := []TabPage{}
	for _, c := range p.Components {
		if c.ParentID == parentID {
			pages = append(pages, TabPage{
				Title:    c.Text,
				Layout:   getLayout(c.Layout),
				Children: createWidgets(p, c.ID),
			})
		}
	}
	return pages
}

func getLayout(lt models.LayoutType) Layout {
	switch lt {
	case models.LayoutHBox:
		return HBox{}
	case models.LayoutVBox:
		return VBox{}
	case models.LayoutGrid:
		return Grid{Columns: 2}
	}
	return nil
}
