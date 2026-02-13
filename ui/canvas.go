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

	// Draw Background
	bgBrush, _ := walk.NewSolidColorBrush(walk.RGB(245, 245, 245))
	defer bgBrush.Dispose()
	canvas.FillRectanglePixels(bgBrush, updateBounds)

	// Draw Grid
	if cw.mw.ActionShowGrid.Checked() {
		gridColor := walk.RGB(220, 220, 220)
		pen, _ := walk.NewCosmeticPen(walk.PenSolid, gridColor)
		defer pen.Dispose()
		
		spacingPx, _ := cw.fromLogical(20, 20)
		for x := 0; x < cw.WidthPixels(); x += spacingPx {
			canvas.DrawLinePixels(pen, walk.Point{x, 0}, walk.Point{x, cw.HeightPixels()})
		}
		for y := 0; y < cw.HeightPixels(); y += spacingPx {
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

	borderPen, _ := walk.NewCosmeticPen(walk.PenSolid, walk.RGB(120, 120, 120))
	defer borderPen.Dispose()
	lightPen, _ := walk.NewCosmeticPen(walk.PenSolid, walk.RGB(220, 220, 220))
	defer lightPen.Dispose()
	bluePen, _ := walk.NewCosmeticPen(walk.PenSolid, walk.RGB(0, 120, 215))
	defer bluePen.Dispose()

	font, _ := walk.NewFont("Segoe UI", 9, 0)
	defer font.Dispose()

	displayText := c.Text
	if displayText == "" {
		displayText = c.Type
	}

	fillRect := func(r walk.Rectangle, color walk.Color) {
		b, _ := walk.NewSolidColorBrush(color)
		defer b.Dispose()
		canvas.FillRectanglePixels(b, r)
	}

	switch c.Type {
	case "PushButton":
		canvas.GradientFillRectanglePixels(walk.RGB(252, 252, 252), walk.RGB(232, 232, 232), walk.Vertical, rect)
		canvas.DrawRectanglePixels(borderPen, rect)
		canvas.DrawTextPixels(displayText, font, walk.RGB(20, 20, 20), rect, walk.TextCenter|walk.TextVCenter|walk.TextSingleLine)
	case "LineEdit", "TextEdit", "NumberEdit", "DateEdit":
		fillRect(rect, walk.RGB(255, 255, 255))
		canvas.DrawRectanglePixels(borderPen, rect)
		textRect := walk.Rectangle{X: rect.X + 6, Y: rect.Y + 4, Width: rect.Width - 12, Height: rect.Height - 8}
		canvas.DrawTextPixels(displayText, font, walk.RGB(60, 60, 60), textRect, walk.TextEditControl)
	case "CheckBox":
		boxSize, _ := cw.fromLogical(14, 14)
		box := walk.Rectangle{X: rect.X + 4, Y: rect.Y + (rect.Height-boxSize)/2, Width: boxSize, Height: boxSize}
		fillRect(box, walk.RGB(255, 255, 255))
		canvas.DrawRectanglePixels(borderPen, box)
		if c.Checked {
			canvas.DrawLinePixels(bluePen, walk.Point{X: box.X + 3, Y: box.Y + 7}, walk.Point{X: box.X + 6, Y: box.Y + 10})
			canvas.DrawLinePixels(bluePen, walk.Point{X: box.X + 6, Y: box.Y + 10}, walk.Point{X: box.X + 11, Y: box.Y + 3})
		}
		textRect := walk.Rectangle{X: box.X + box.Width + 6, Y: rect.Y, Width: rect.Width - box.Width - 10, Height: rect.Height}
		canvas.DrawTextPixels(displayText, font, walk.RGB(20, 20, 20), textRect, walk.TextVCenter|walk.TextSingleLine)
	case "RadioButton":
		circleSize, _ := cw.fromLogical(14, 14)
		circle := walk.Rectangle{X: rect.X + 4, Y: rect.Y + (rect.Height-circleSize)/2, Width: circleSize, Height: circleSize}
		fillRect(circle, walk.RGB(255, 255, 255))
		canvas.DrawEllipse(borderPen, circle)
		if c.Checked {
			dotSize, _ := cw.fromLogical(6, 6)
			dot := walk.Rectangle{X: circle.X + 4, Y: circle.Y + 4, Width: dotSize, Height: dotSize}
			fillRect(dot, walk.RGB(0, 120, 215))
		}
		textRect := walk.Rectangle{X: circle.X + circle.Width + 6, Y: rect.Y, Width: rect.Width - circle.Width - 10, Height: rect.Height}
		canvas.DrawTextPixels(displayText, font, walk.RGB(20, 20, 20), textRect, walk.TextVCenter|walk.TextSingleLine)
	case "ComboBox":
		fillRect(rect, walk.RGB(255, 255, 255))
		canvas.DrawRectanglePixels(borderPen, rect)
		btnW, _ := cw.fromLogical(20, 20)
		btn := walk.Rectangle{X: rect.X + rect.Width - btnW, Y: rect.Y + 1, Width: btnW - 1, Height: rect.Height - 2}
		canvas.GradientFillRectanglePixels(walk.RGB(248, 248, 248), walk.RGB(232, 232, 232), walk.Vertical, btn)
		canvas.DrawRectanglePixels(lightPen, btn)
		txt := ""
		if len(c.Items) > 0 {
			txt = c.Items[0]
		}
		textRect := walk.Rectangle{X: rect.X + 6, Y: rect.Y + 4, Width: rect.Width - btnW - 8, Height: rect.Height - 8}
		canvas.DrawTextPixels(txt, font, walk.RGB(40, 40, 40), textRect, walk.TextVCenter|walk.TextSingleLine)
	case "ListBox":
		fillRect(rect, walk.RGB(255, 255, 255))
		canvas.DrawRectanglePixels(borderPen, rect)
		y := rect.Y + 4
		for i, item := range c.Items {
			if i >= 10 || y > rect.Y+rect.Height-14 {
				break
			}
			lineRect := walk.Rectangle{X: rect.X + 6, Y: y, Width: rect.Width - 12, Height: 14}
			canvas.DrawTextPixels(item, font, walk.RGB(30, 30, 30), lineRect, walk.TextSingleLine)
			y += 16
		}
	case "ProgressBar":
		fillRect(rect, walk.RGB(245, 245, 245))
		canvas.DrawRectanglePixels(borderPen, rect)
		minV, maxV := c.MinValue, c.MaxValue
		if maxV <= minV { maxV = minV + 100 }
		val := c.Value
		if val < minV { val = minV }
		if val > maxV { val = maxV }
		fillW := (rect.Width - 2) * (val - minV) / (maxV - minV)
		if fillW > 0 {
			fill := walk.Rectangle{X: rect.X + 1, Y: rect.Y + 1, Width: fillW, Height: rect.Height - 2}
			canvas.GradientFillRectanglePixels(walk.RGB(138, 199, 255), walk.RGB(0, 120, 215), walk.Vertical, fill)
		}
	case "Slider":
		trackY := rect.Y + rect.Height/2
		canvas.DrawLinePixels(lightPen, walk.Point{X: rect.X + 8, Y: trackY}, walk.Point{X: rect.X + rect.Width - 8, Y: trackY})
		thumbSize, _ := cw.fromLogical(10, 16)
		thumb := walk.Rectangle{X: rect.X + rect.Width/2 - thumbSize/2, Y: trackY - thumbSize, Width: thumbSize, Height: thumbSize * 2}
		canvas.GradientFillRectanglePixels(walk.RGB(245, 245, 245), walk.RGB(220, 220, 220), walk.Vertical, thumb)
		canvas.DrawRectanglePixels(borderPen, thumb)
	case "GroupBox":
		fillRect(rect, walk.RGB(252, 252, 252))
		canvas.DrawRectanglePixels(borderPen, rect)
		titleRect := walk.Rectangle{X: rect.X + 8, Y: rect.Y, Width: rect.Width - 16, Height: 18}
		canvas.DrawTextPixels(displayText, font, walk.RGB(40, 40, 40), titleRect, walk.TextSingleLine)
	case "TabWidget":
		fillRect(rect, walk.RGB(250, 250, 250))
		canvas.DrawRectanglePixels(borderPen, rect)
		tabW, tabH := cw.fromLogical(80, 22)
		tab := walk.Rectangle{X: rect.X + 6, Y: rect.Y + 4, Width: tabW, Height: tabH}
		canvas.GradientFillRectanglePixels(walk.RGB(255, 255, 255), walk.RGB(237, 237, 237), walk.Vertical, tab)
		canvas.DrawRectanglePixels(borderPen, tab)
		canvas.DrawTextPixels("Tab 1", font, walk.RGB(30, 30, 30), tab, walk.TextCenter|walk.TextVCenter|walk.TextSingleLine)
	case "Composite", "Splitter":
		fillRect(rect, walk.RGB(247, 247, 247))
		canvas.DrawRectanglePixels(borderPen, rect)
		canvas.DrawTextPixels(c.Type, font, walk.RGB(90, 90, 90), rect, walk.TextCenter|walk.TextVCenter|walk.TextSingleLine)
	default:
		fillRect(rect, walk.RGB(240, 240, 240))
		canvas.DrawRectanglePixels(borderPen, rect)
		canvas.DrawTextPixels(fmt.Sprintf("[%s] %s", c.Type, displayText), font, walk.RGB(0, 0, 0), rect, walk.TextVCenter|walk.TextCenter)
	}
}

func (cw *CanvasWidget) drawSelectionHandles(canvas *walk.Canvas, c *models.Component) {
	handleSize := 6
	brush, _ := walk.NewSolidColorBrush(walk.RGB(0, 120, 215))
	defer brush.Dispose()

	handles := cw.getHandles(c, handleSize)
	for _, h := range handles {
		canvas.FillRectanglePixels(brush, h)
	}

	// Also draw dash border
	pen, _ := walk.NewCosmeticPen(walk.PenDash, walk.RGB(0, 120, 215))
	defer pen.Dispose()
	lx, ly := cw.fromLogical(c.X, c.Y)
	lw, lh := cw.fromLogical(c.Width, c.Height)
	canvas.DrawRectanglePixels(pen, walk.Rectangle{lx - 2, ly - 2, lw + 4, lh + 4})
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

	// Handle pending tool
	if cw.mw.PendingTool != "" {
		cw.mw.AddComponent(cw.mw.PendingTool, lx, ly)
		cw.mw.PendingTool = ""
		if cw.mw.Toolbox != nil {
			// Toolbox might be using ListBox or TreeView now
			if cw.mw.Toolbox.listBox != nil {
				cw.mw.Toolbox.listBox.SetCurrentIndex(-1)
			}
		}
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
	if cw.mw.PropertyEditor != nil {
		cw.mw.PropertyEditor.Update()
	}
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
			return p.Layout != nil && p.Layout.Type != models.LayoutNone
		}
	}
	return false
}
