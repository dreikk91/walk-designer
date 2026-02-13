package ui

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"walk-gui-designer-pro/models"
)

type dragMode int

const (
	dragNone dragMode = iota
	dragMove
	dragResize
)

type CanvasWidget struct {
	*walk.CustomWidget
	mw *DesignerWindow

	mode       dragMode
	dragComp   *models.Component
	dragHandle int // 0-7
	startX     int // Logical
	startY     int // Logical
	origX      int // Logical
	origY      int // Logical
	origW      int // Logical
	origH      int // Logical
}

func (cw *CanvasWidget) toLogical(x, y int) (int, int) {
	dpi := cw.DPI()
	if dpi == 0 {
		dpi = 96
	}
	return x * 96 / dpi, y * 96 / dpi
}

func (cw *CanvasWidget) fromLogical(lx, ly int) (int, int) {
	dpi := cw.DPI()
	if dpi == 0 {
		dpi = 96
	}
	return lx * dpi / 96, ly * dpi / 96
}

func CreateCanvas(mw *DesignerWindow) Widget {
	cw := &CanvasWidget{mw: mw}
	mw.Canvas = cw

	return CustomWidget{
		AssignTo:            &cw.CustomWidget,
		ClearsBackground:    true,
		InvalidatesOnResize: true,
		Paint:               cw.onPaint,
		OnMouseDown:         cw.onMouseDown,
		OnMouseMove:         cw.onMouseMove,
		OnMouseUp:           cw.onMouseUp,
	}
}

func (cw *CanvasWidget) onPaint(canvas *walk.Canvas, updateBounds walk.Rectangle) error {
	if cw.mw.Project == nil {
		return nil
	}

	// Draw Grid
	if cw.mw.ActionShowGrid.Checked() {
		gridColor := walk.RGB(230, 230, 230)
		pen, _ := walk.NewCosmeticPen(walk.PenSolid, gridColor)
		defer pen.Dispose()
		
		// Use pixels for grid but align to logical 20px
		spacingPx, _ := cw.fromLogical(20, 20)
		for x := 0; x < cw.WidthPixels(); x += spacingPx {
			canvas.DrawLinePixels(pen, walk.Point{x, 0}, walk.Point{x, cw.HeightPixels()})
		}
		for y := 0; y < cw.HeightPixels(); y += 20 {
			canvas.DrawLinePixels(pen, walk.Point{0, y}, walk.Point{cw.WidthPixels(), y})
		}
	}

	// Draw Components
	for _, c := range cw.mw.Project.Components {
		cw.drawComponent(canvas, c)
	}

	// Draw Selection Handles
	if cw.mw.SelectedComponent != nil {
		cw.drawSelectionHandles(canvas, cw.mw.SelectedComponent)
	}

	return nil
}

func (cw *CanvasWidget) drawComponent(canvas *walk.Canvas, c *models.Component) {
	lx, ly := cw.fromLogical(c.X, c.Y)
	lw, lh := cw.fromLogical(c.Width, c.Height)
	rect := walk.Rectangle{lx, ly, lw, lh}

	brush, _ := walk.NewSolidColorBrush(walk.RGB(240, 240, 240))
	defer brush.Dispose()

	penColor := walk.RGB(100, 100, 100)
	if cw.mw.SelectedComponent == c {
		penColor = walk.RGB(0, 120, 215)
	}
	pen, _ := walk.NewCosmeticPen(walk.PenSolid, penColor)
	defer pen.Dispose()

	canvas.FillRectanglePixels(brush, rect)
	canvas.DrawRectanglePixels(pen, rect)

	// Draw Text
	font, _ := walk.NewFont("Segoe UI", 9, 0)
	defer font.Dispose()
	text := fmt.Sprintf("[%s] %s", c.Type, c.Text)
	canvas.DrawTextPixels(text, font, walk.RGB(0, 0, 0), rect, walk.TextVCenter|walk.TextCenter)
}

func (cw *CanvasWidget) drawSelectionHandles(canvas *walk.Canvas, c *models.Component) {
	handleSize := 6
	brush, _ := walk.NewSolidColorBrush(walk.RGB(0, 120, 215))
	defer brush.Dispose()

	handles := cw.getHandles(c, handleSize)
	for _, h := range handles {
		canvas.FillRectanglePixels(brush, h)
	}
}

func (cw *CanvasWidget) getHandles(c *models.Component, size int) []walk.Rectangle {
	lx, ly := cw.fromLogical(c.X, c.Y)
	lw, lh := cw.fromLogical(c.Width, c.Height)
	half := size / 2
	
	return []walk.Rectangle{
		{lx - half, ly - half, size, size},                   // TL
		{lx + lw/2 - half, ly - half, size, size},           // TC
		{lx + lw - half, ly - half, size, size},             // TR
		{lx + lw - half, ly + lh/2 - half, size, size},     // RC
		{lx + lw - half, ly + lh - half, size, size},       // BR
		{lx + lw/2 - half, ly + lh - half, size, size},     // BC
		{lx - half, ly + lh - half, size, size},             // BL
		{lx - half, ly + lh/2 - half, size, size},           // LC
	}
}

func (cw *CanvasWidget) onMouseDown(x, y int, button walk.MouseButton) {
	if button != walk.LeftButton || cw.mw.Project == nil {
		return
	}

	lx, ly := cw.toLogical(x, y)
	cw.startX, cw.startY = lx, ly

	// Handle pending tool (Stamp/Drop from toolbox)
	if cw.mw.PendingTool != "" {
		cw.mw.AddComponent(cw.mw.PendingTool, lx, ly)
		cw.mw.PendingTool = ""
		cw.mw.Toolbox.listBox.SetCurrentIndex(-1)
		return
	}

	// Check handles first
	if cw.mw.SelectedComponent != nil {
		if !cw.parentHasLayout(cw.mw.SelectedComponent) {
			handles := cw.getHandles(cw.mw.SelectedComponent, 8)
			for i, h := range handles {
				if x >= h.X && x <= h.X+h.Width && y >= h.Y && y <= h.Y+h.Height {
					cw.mode = dragResize
					cw.dragHandle = i
					cw.dragComp = cw.mw.SelectedComponent
					cw.origX, cw.origY = cw.dragComp.X, cw.dragComp.Y
					cw.origW, cw.origH = cw.dragComp.Width, cw.dragComp.Height
					return
				}
			}
		}
	}

	// Hit test (logical)
	var found *models.Component
	for i := len(cw.mw.Project.Components) - 1; i >= 0; i-- {
		c := cw.mw.Project.Components[i]
		if lx >= c.X && lx <= c.X+c.Width && ly >= c.Y && ly <= c.Y+c.Height {
			found = c
			break
		}
	}

	if found != nil {
		cw.mw.SelectComponent(found)
		if cw.parentHasLayout(found) {
			cw.mode = dragNone
		} else {
			cw.mode = dragMove
			cw.dragComp = found
			cw.origX, cw.origY = found.X, found.Y
		}
	} else {
		cw.mw.SelectComponent(nil)
		cw.mode = dragNone
	}
}

func (cw *CanvasWidget) onMouseMove(x, y int, button walk.MouseButton) {
	if cw.mode == dragNone || cw.dragComp == nil {
		return
	}

	lx, ly := cw.toLogical(x, y)
	dx, dy := lx-cw.startX, ly-cw.startY

	if cw.mode == dragMove && dx*dx+dy*dy < 4 {
		return
	}

	if cw.mode == dragMove {
		newX, newY := cw.origX+dx, cw.origY+dy
		if cw.mw.ActionSnapToGrid.Checked() {
			newX = (newX / 20) * 20
			newY = (newY / 20) * 20
		}
		cw.dragComp.X = newX
		cw.dragComp.Y = newY
	} else if cw.mode == dragResize {
		cw.applyResize(dx, dy)
	}

	cw.Invalidate()
	cw.mw.PropertyEditor.Update()
}

func (cw *CanvasWidget) applyResize(dx, dy int) {
	c := cw.dragComp
	nx, ny, nw, nh := cw.origX, cw.origY, cw.origW, cw.origH

	switch cw.dragHandle {
	case 0: nx, ny, nw, nh = cw.origX+dx, cw.origY+dy, cw.origW-dx, cw.origH-dy
	case 1: ny, nh = cw.origY+dy, cw.origH-dy
	case 2: ny, nw, nh = cw.origY+dy, cw.origW+dx, cw.origH-dy
	case 3: nw = cw.origW+dx
	case 4: nw, nh = cw.origW+dx, cw.origH+dy
	case 5: nh = cw.origH+dy
	case 6: nx, nw, nh = cw.origX+dx, cw.origW-dx, cw.origH+dy
	case 7: nx, nw = cw.origX+dx, cw.origW-dx
	}

	if nw < 10 { nw = 10 }
	if nh < 10 { nh = 10 }
	
	if cw.mw.ActionSnapToGrid.Checked() {
		nx = (nx / 20) * 20
		ny = (ny / 20) * 20
		nw = (nw / 20) * 20
		nh = (nh / 20) * 20
	}

	c.X, c.Y, c.Width, c.Height = nx, ny, nw, nh
}

func (cw *CanvasWidget) onMouseUp(x, y int, button walk.MouseButton) {
	if cw.mode == dragMove && cw.dragComp != nil {
		lx, ly := cw.toLogical(x, y)
		var container *models.Component
		for i := len(cw.mw.Project.Components) - 1; i >= 0; i-- {
			c := cw.mw.Project.Components[i]
			if c == cw.dragComp { continue }
			if cw.isContainer(c.Type) && lx >= c.X && lx <= c.X+c.Width && ly >= c.Y && ly <= c.Y+c.Height {
				container = c
				break
			}
		}

		if container != nil {
			cw.dragComp.ParentID = container.ID
		} else {
			cw.dragComp.ParentID = ""
		}
		cw.mw.RefreshAll()
	}

	cw.mode = dragNone
	cw.dragComp = nil
	cw.Invalidate()
}

func (cw *CanvasWidget) isContainer(typ string) bool {
	switch typ {
	case "Composite", "GroupBox", "Splitter", "TabWidget": return true
	}
	return false
}

func (cw *CanvasWidget) parentHasLayout(c *models.Component) bool {
	if c.ParentID == "" || cw.mw.Project == nil { return false }
	for _, p := range cw.mw.Project.Components {
		if p.ID == c.ParentID {
			return p.Layout != models.LayoutNone
		}
	}
	return false
}
