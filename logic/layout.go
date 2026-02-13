package logic

import (
	"walk-gui-designer-pro/models"
)

func ApplyLayout(p *models.Project) {
	if p == nil {
		return
	}

	// Apply layout to root components if project has a layout
	if p.Layout != nil && p.Layout.Type != models.LayoutNone {
		roots := getChildren(p, "")
		applyLayoutToGroup(p.Layout, 0, 0, p.Width, p.Height, roots)
	}

	// Recursively apply to children of all containers
	for _, c := range p.Components {
		if c.Layout != nil && c.Layout.Type != models.LayoutNone {
			children := getChildren(p, c.ID)
			applyLayoutToGroup(c.Layout, c.X, c.Y, c.Width, c.Height, children)
		}
	}
}

func getChildren(p *models.Project, parentID string) []*models.Component {
	children := []*models.Component{}
	for _, c := range p.Components {
		if c.ParentID == parentID {
			children = append(children, c)
		}
	}
	return children
}

func applyLayoutToGroup(cfg *models.LayoutConfig, px, py, pw, ph int, children []*models.Component) {
	if len(children) == 0 {
		return
	}

	margin := 10
	spacing := 5
	if cfg.MarginsZero {
		margin = 0
	}
	if cfg.SpacingZero {
		spacing = 0
	}

	switch cfg.Type {
	case models.LayoutVBox:
		y := py + margin
		availW := pw - 2*margin
		if availW < 10 { availW = 10 }

		childH := (ph - 2*margin - (len(children)-1)*spacing) / len(children)
		if childH < 10 { childH = 10 }

		for _, c := range children {
			c.X = px + margin
			c.Y = y
			c.Width = availW
			c.Height = childH
			y += childH + spacing
		}

	case models.LayoutHBox:
		x := px + margin
		availH := ph - 2*margin
		if availH < 10 { availH = 10 }

		childW := (pw - 2*margin - (len(children)-1)*spacing) / len(children)
		if childW < 10 { childW = 10 }

		for _, c := range children {
			c.X = x
			c.Y = py + margin
			c.Width = childW
			c.Height = availH
			x += childW + spacing
		}

	case models.LayoutGrid:
		cols := cfg.Columns
		if cols <= 0 { cols = 2 }
		rows := (len(children) + cols - 1) / cols

		childW := (pw - 2*margin - (cols-1)*spacing) / cols
		childH := (ph - 2*margin - (rows-1)*spacing) / rows
		if childW < 10 { childW = 10 }
		if childH < 10 { childH = 10 }

		for i, c := range children {
			r, col := i/cols, i%cols
			c.X = px + margin + col*(childW+spacing)
			c.Y = py + margin + r*(childH+spacing)
			c.Width = childW
			c.Height = childH
		}
	}
}
