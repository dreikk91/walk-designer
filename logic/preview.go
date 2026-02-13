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
		Layout:   getLayout(p.Layout),
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

			switch c.Type {
			case "PushButton":
				w = PushButton{Text: c.Text}
			case "Label":
				w = Label{Text: c.Text}
			case "LineEdit":
				w = LineEdit{Text: c.Text, ReadOnly: c.ReadOnly, PasswordMode: c.IsPassword}
			case "TextEdit":
				w = TextEdit{Text: c.Text, ReadOnly: c.ReadOnly}
			case "CheckBox":
				w = CheckBox{Text: c.Text, Checked: c.Checked}
			case "RadioButton":
				// RadioButton in declarative doesn't have Checked, it's usually handled by DataBinder or by being in a group
				w = RadioButton{Text: c.Text}
			case "ComboBox":
				w = ComboBox{Model: c.Items}
			case "ListBox":
				w = ListBox{Model: c.Items}
			case "NumberEdit":
				w = NumberEdit{Value: float64(c.Value), MinValue: float64(c.MinValue), MaxValue: float64(c.MaxValue)}
			case "DateEdit":
				w = DateEdit{}
			case "Slider":
				var orientation Orientation
				if c.Orientation == "Vertical" { orientation = Vertical } else { orientation = Horizontal }
				w = Slider{Value: c.Value, MinValue: c.MinValue, MaxValue: c.MaxValue, Orientation: orientation}
			case "ProgressBar":
				w = ProgressBar{Value: c.Value, MinValue: c.MinValue, MaxValue: c.MaxValue}
			case "Composite":
				w = Composite{
					Layout:   getLayout(c.Layout),
					Children: createWidgets(p, c.ID),
				}
			case "Splitter":
				if c.Orientation == "Vertical" {
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
					Pages: createTabPages(p, c.ID, c.TabPages),
				}
			case "TableView":
				w = TableView{
					Columns: createTableColumns(c.Columns),
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
				w = Label{Text: "[" + c.Type + "] " + c.Text}
			}

			if w != nil {
				widgets = append(widgets, w)
			}
		}
	}

	return widgets
}

func createTabPages(p *models.Project, parentID string, configs []models.TabPageConfig) []TabPage {
	pages := []TabPage{}
	if len(configs) > 0 {
		for _, tp := range configs {
			pages = append(pages, TabPage{
				Title:    tp.Title,
				Layout:   getLayout(tp.Layout),
				Children: createWidgets(p, ""),
			})
		}
	} else {
		for _, c := range p.Components {
			if c.ParentID == parentID {
				pages = append(pages, TabPage{
					Title:    c.Text,
					Layout:   getLayout(c.Layout),
					Children: createWidgets(p, c.ID),
				})
			}
		}
	}
	return pages
}

func createTableColumns(cols []string) []TableViewColumn {
	var tvcs []TableViewColumn
	for _, col := range cols {
		tvcs = append(tvcs, TableViewColumn{Title: col})
	}
	if len(tvcs) == 0 {
		tvcs = append(tvcs, TableViewColumn{Title: "Column 1"})
	}
	return tvcs
}

func getLayout(cfg *models.LayoutConfig) Layout {
	if cfg == nil { return nil }
	switch cfg.Type {
	case models.LayoutHBox:
		return HBox{MarginsZero: cfg.MarginsZero, SpacingZero: cfg.SpacingZero}
	case models.LayoutVBox:
		return VBox{MarginsZero: cfg.MarginsZero, SpacingZero: cfg.SpacingZero}
	case models.LayoutGrid:
		return Grid{Columns: cfg.Columns, MarginsZero: cfg.MarginsZero, SpacingZero: cfg.SpacingZero}
	}
	return nil
}
