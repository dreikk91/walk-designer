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

			var minSize, maxSize Size
			if c.MinWidth > 0 || c.MinHeight > 0 { minSize = Size{Width: c.MinWidth, Height: c.MinHeight} }
			if c.MaxWidth > 0 || c.MaxHeight > 0 { maxSize = Size{Width: c.MaxWidth, Height: c.MaxHeight} }

			switch c.Type {
			case "PushButton":
				w = PushButton{Text: c.Text, MinSize: minSize, MaxSize: maxSize}
			case "Label":
				w = Label{Text: c.Text, MinSize: minSize, MaxSize: maxSize}
			case "LineEdit":
				w = LineEdit{Text: c.Text, ReadOnly: c.ReadOnly, PasswordMode: c.IsPassword, MinSize: minSize, MaxSize: maxSize}
			case "TextEdit":
				w = TextEdit{Text: c.Text, ReadOnly: c.ReadOnly, MinSize: minSize, MaxSize: maxSize}
			case "CheckBox":
				w = CheckBox{Text: c.Text, Checked: c.Checked, MinSize: minSize, MaxSize: maxSize}
			case "RadioButton":
				w = RadioButton{Text: c.Text, MinSize: minSize, MaxSize: maxSize}
			case "ComboBox":
				w = ComboBox{Model: c.Items, MinSize: minSize, MaxSize: maxSize}
			case "ListBox":
				w = ListBox{Model: c.Items, MinSize: minSize, MaxSize: maxSize}
			case "NumberEdit":
				w = NumberEdit{Value: float64(c.Value), MinValue: float64(c.MinValue), MaxValue: float64(c.MaxValue), MinSize: minSize, MaxSize: maxSize}
			case "DateEdit":
				w = DateEdit{MinSize: minSize, MaxSize: maxSize}
			case "Slider":
				var orientation Orientation
				if c.Orientation == "Vertical" { orientation = Vertical } else { orientation = Horizontal }
				w = Slider{Value: c.Value, MinValue: c.MinValue, MaxValue: c.MaxValue, Orientation: orientation, MinSize: minSize, MaxSize: maxSize}
			case "ProgressBar":
				w = ProgressBar{Value: c.Value, MinValue: c.MinValue, MaxValue: c.MaxValue, MinSize: minSize, MaxSize: maxSize}
			case "Composite":
				w = Composite{
					Layout:   getLayout(c.Layout),
					Children: createWidgets(p, c.ID),
					MinSize: minSize, MaxSize: maxSize,
				}
			case "Splitter":
				if c.Orientation == "Vertical" {
					w = VSplitter{
						Children: createWidgets(p, c.ID),
						MinSize: minSize, MaxSize: maxSize,
					}
				} else {
					w = HSplitter{
						Children: createWidgets(p, c.ID),
						MinSize: minSize, MaxSize: maxSize,
					}
				}
			case "TabWidget":
				w = TabWidget{
					Pages: createTabPages(p, c.ID, c.TabPages),
					MinSize: minSize, MaxSize: maxSize,
				}
			case "TableView":
				w = TableView{
					Columns: createTableColumns(c.TableCols),
					MinSize: minSize, MaxSize: maxSize,
				}
			case "TreeView":
				w = TreeView{MinSize: minSize, MaxSize: maxSize}
			case "GroupBox":
				w = GroupBox{
					Title:    c.Text,
					Layout:   getLayout(c.Layout),
					Children: createWidgets(p, c.ID),
					MinSize: minSize, MaxSize: maxSize,
				}
			case "ScrollView":
				w = ScrollView{
					Layout:   getLayout(c.Layout),
					Children: createWidgets(p, c.ID),
					MinSize: minSize, MaxSize: maxSize,
				}
			case "SplitButton":
				w = SplitButton{Text: c.Text, MinSize: minSize, MaxSize: maxSize}
			case "LinkLabel":
				w = LinkLabel{Text: c.Text, MinSize: minSize, MaxSize: maxSize}
			case "ImageView":
				w = ImageView{MinSize: minSize, MaxSize: maxSize}
			case "ToolBar":
				// Empty placeholder for Toolbar to not crash
				w = Composite{Layout: HBox{MarginsZero: true, SpacingZero: true}, Children: []Widget{Label{Text: "ToolBar Placeholder"}}, MinSize: minSize, MaxSize: maxSize}
			case "StatusBar":
				// Empty placeholder
				w = Composite{Layout: HBox{MarginsZero: true, SpacingZero: true}, Children: []Widget{Label{Text: "StatusBar Placeholder"}}, MinSize: minSize, MaxSize: maxSize}
			case "HSpacer":
				w = HSpacer{}
			case "VSpacer":
				w = VSpacer{}
			default:
				w = Label{Text: "[" + c.Type + "] " + c.Text, MinSize: minSize, MaxSize: maxSize}
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

func createTableColumns(cols []models.TableColumnConfig) []TableViewColumn {
	var tvcs []TableViewColumn
	for _, col := range cols {
		// Try to parse alignment, default to Near
		align := AlignNear
		switch col.Alignment {
		case "Center": align = AlignCenter
		case "Far": align = AlignFar
		}

		tvcs = append(tvcs, TableViewColumn{
			Name: col.Name,
			Title: col.Title,
			Width: col.Width,
			Alignment: align,
		})
	}
	if len(tvcs) == 0 {
		tvcs = append(tvcs, TableViewColumn{Title: "Column 1", Width: 100})
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
