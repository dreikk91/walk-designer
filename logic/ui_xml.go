package logic

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"walk-gui-designer-pro/models"
)

// Simplified Qt .ui XML structures

type UIFile struct {
	XMLName xml.Name  `xml:"ui"`
	Version string    `xml:"version,attr"`
	Class   string    `xml:"class"`
	Widget  UIWidget  `xml:"widget"`
}

type UIWidget struct {
	Class      string        `xml:"class,attr"`
	Name       string        `xml:"name,attr"`
	Properties []UIProperty  `xml:"property"`
	Widgets    []UIWidget    `xml:"widget"`
}

type UIProperty struct {
	Name   string  `xml:"name,attr"`
	Rect   *UIRect `xml:"rect,omitempty"`
	String string  `xml:"string,omitempty"`
	Bool   string  `xml:"bool,omitempty"`
}

type UIRect struct {
	X      int `xml:"x"`
	Y      int `xml:"y"`
	Width  int `xml:"width"`
	Height int `xml:"height"`
}

// Map walk/models types to Qt types roughly
func mapTypeToQt(t string) string {
	switch t {
	case "PushButton": return "QPushButton"
	case "Label": return "QLabel"
	case "LineEdit": return "QLineEdit"
	case "TextEdit": return "QTextEdit"
	case "CheckBox": return "QCheckBox"
	case "RadioButton": return "QRadioButton"
	case "ComboBox": return "QComboBox"
	case "ListBox": return "QListWidget"
	case "ProgressBar": return "QProgressBar"
	case "Slider": return "QSlider"
	case "NumberEdit": return "QSpinBox"
	case "Composite", "GroupBox": return "QWidget"
	case "TabWidget": return "QTabWidget"
	case "TableView": return "QTableView"
	case "TreeView": return "QTreeView"
	}
	return "QWidget"
}

func mapQtToType(qtClass string) string {
	switch qtClass {
	case "QPushButton": return "PushButton"
	case "QLabel": return "Label"
	case "QLineEdit": return "LineEdit"
	case "QTextEdit": return "TextEdit"
	case "QCheckBox": return "CheckBox"
	case "QRadioButton": return "RadioButton"
	case "QComboBox": return "ComboBox"
	case "QListWidget": return "ListBox"
	case "QProgressBar": return "ProgressBar"
	case "QSlider": return "Slider"
	case "QSpinBox": return "NumberEdit"
	case "QTabWidget": return "TabWidget"
	case "QTableView": return "TableView"
	case "QTreeView": return "TreeView"
	case "QMainWindow", "QDialog", "QWidget": return "Composite"
	}
	return "Composite"
}

// Export Project to Qt .ui XML
func ExportToUI(p *models.Project, filename string) error {
	ui := UIFile{
		Version: "4.0",
		Class:   p.Title,
		Widget: UIWidget{
			Class: "QWidget",
			Name:  "MainWindow",
			Properties: []UIProperty{
				{Name: "geometry", Rect: &UIRect{0, 0, p.Width, p.Height}},
				{Name: "windowTitle", String: p.Title},
			},
		},
	}

	// Build a map of children to nest them recursively
	childMap := make(map[string][]*models.Component)
	for _, c := range p.Components {
		childMap[c.ParentID] = append(childMap[c.ParentID], c)
	}

	var buildWidget func(*models.Component) UIWidget
	buildWidget = func(c *models.Component) UIWidget {
		w := UIWidget{
			Class: mapTypeToQt(c.Type),
			Name:  c.Name,
		}

		w.Properties = append(w.Properties, UIProperty{
			Name: "geometry",
			Rect: &UIRect{c.X, c.Y, c.Width, c.Height},
		})

		if c.Text != "" {
			w.Properties = append(w.Properties, UIProperty{
				Name: "text",
				String: c.Text,
			})
		}
		if c.ToolTip != "" {
			w.Properties = append(w.Properties, UIProperty{
				Name: "toolTip",
				String: c.ToolTip,
			})
		}

		for _, child := range childMap[c.ID] {
			w.Widgets = append(w.Widgets, buildWidget(child))
		}

		return w
	}

	for _, c := range childMap[""] {
		ui.Widget.Widgets = append(ui.Widget.Widgets, buildWidget(c))
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	enc := xml.NewEncoder(f)
	enc.Indent("", "  ")
	return enc.Encode(ui)
}

// Import Project from Qt .ui XML
func ImportFromUI(filename string) (*models.Project, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var ui UIFile
	if err := xml.Unmarshal(data, &ui); err != nil {
		return nil, err
	}

	p := &models.Project{
		Title:  ui.Class,
		Width:  800,
		Height: 600,
	}

	for _, prop := range ui.Widget.Properties {
		if prop.Name == "geometry" && prop.Rect != nil {
			p.Width = prop.Rect.Width
			p.Height = prop.Rect.Height
		}
		if prop.Name == "windowTitle" && prop.String != "" {
			p.Title = prop.String
		}
	}

	var parseWidget func(UIWidget, string)
	parseWidget = func(w UIWidget, parentID string) {
		id := w.Name
		if id == "" {
			id = fmt.Sprintf("widget_%d", len(p.Components)+1)
		}

		c := &models.Component{
			ID:       id,
			Name:     w.Name,
			Type:     mapQtToType(w.Class),
			ParentID: parentID,
			Enabled:  true,
			Visible:  true,
		}

		for _, prop := range w.Properties {
			switch prop.Name {
			case "geometry":
				if prop.Rect != nil {
					c.X = prop.Rect.X
					c.Y = prop.Rect.Y
					c.Width = prop.Rect.Width
					c.Height = prop.Rect.Height
				}
			case "text":
				c.Text = prop.String
			case "toolTip":
				c.ToolTip = prop.String
			}
		}

		p.Components = append(p.Components, c)

		for _, child := range w.Widgets {
			parseWidget(child, c.ID)
		}
	}

	for _, w := range ui.Widget.Widgets {
		parseWidget(w, "")
	}

	return p, nil
}
