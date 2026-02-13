package logic

import (
	"walk-gui-designer-pro/models"
)

func ApplyLayout(p *models.Project) {
	if p == nil {
		return
	}

	// First pass: identify root components and containers
	// Actually we should do it recursively starting from those with ParentID == ""
	roots := []*models.Component{}
	for _, c := range p.Components {
		if c.ParentID == "" {
			roots = append(roots, c)
		}
	}

	// If the project itself has a layout, it applies to roots
	if p.Layout != models.LayoutNone {
		applyLayoutToGroup(p.Layout, 0, 0, p.Width, p.Height, roots)
	}

	// Recursively apply to children of containers
	for _, c := range p.Components {
		if c.Layout != models.LayoutNone {
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

func applyLayoutToGroup(layout models.LayoutType, px, py, pw, ph int, children []*models.Component) {
	if len(children) == 0 {
		return
	}

	// For simplicity, we check the first child's parent or the project
	// but better to pass the container component itself
	margin := 10
	spacing := 5
	// In a full implementation we would check c.MarginsZero and c.SpacingZero

	switch layout {
	case models.LayoutVBox:
		y := py + margin
		availW := pw - 2*margin
		childH := (ph - 2*margin - (len(children)-1)*spacing) / len(children)
		if childH < 10 {
			childH = 10
		}

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
		childW := (pw - 2*margin - (len(children)-1)*spacing) / len(children)
		if childW < 10 {
			childW = 10
		}

		for _, c := range children {
			c.X = x
			c.Y = py + margin
			c.Width = childW
			c.Height = availH
			x += childW + spacing
		}

	case models.LayoutGrid:
		// Simple grid implementation
		cols := 2 // Default
		// ... logic for grid ...
		rows := (len(children) + cols - 1) / cols
		childW := (pw - 2*margin - (cols-1)*spacing) / cols
		childH := (ph - 2*margin - (rows-1)*spacing) / rows

		for i, c := range children {
			r, col := i/cols, i%cols
			c.X = px + margin + col*(childW+spacing)
			c.Y = py + margin + r*(childH+spacing)
			c.Width = childW
			c.Height = childH
		}
	}
}
