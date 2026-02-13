package logic

import (
	"fmt"
	"strings"
	"walk-gui-designer-pro/models"
)

func GenerateGoCode(p *models.Project) string {
	var sb strings.Builder

	sb.WriteString("package main\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"log\"\n")
	sb.WriteString("\t\"github.com/lxn/walk\"\n")
	sb.WriteString("\t. \"github.com/lxn/walk/declarative\"\n")
	sb.WriteString(")\n\n")

	sb.WriteString("func main() {\n")

	// Generate variables for AssignTo
	for _, c := range p.Components {
		if c.Name != "" {
			walkType := c.Type
			sb.WriteString(fmt.Sprintf("\tvar %s *walk.%s\n", c.Name, walkType))
		}
	}
	sb.WriteString("\n")

	sb.WriteString("\tif _, err := (MainWindow{\n")
	sb.WriteString(fmt.Sprintf("\t\tTitle: \"%s\",\n", p.Title))
	sb.WriteString(fmt.Sprintf("\t\tSize: Size{Width: %d, Height: %d},\n", p.Width, p.Height))
	if p.Layout != nil {
		sb.WriteString(generateLayoutCode(p.Layout, 2))
	}
	sb.WriteString("\t\tChildren: []Widget{\n")

	sb.WriteString(generateChildren(p, ""))

	sb.WriteString("\t\t},\n")
	sb.WriteString("\t}.Run()); err != nil {\n")
	sb.WriteString("\t\tlog.Fatal(err)\n")
	sb.WriteString("\t}\n")
	sb.WriteString("}\n")

	return sb.String()
}

func generateChildren(p *models.Project, parentID string) string {
	var sb strings.Builder
	indent := "\t\t\t"

	for _, c := range p.Components {
		if c.ParentID == parentID {
			typ := c.Type
			if typ == "Splitter" {
				if c.Orientation == "Vertical" {
					typ = "VSplitter"
				} else {
					typ = "HSplitter"
				}
			}

			sb.WriteString(fmt.Sprintf("%s%s{\n", indent, typ))
			if c.Name != "" {
				sb.WriteString(fmt.Sprintf("%s\tAssignTo: &%s,\n", indent, c.Name))
			}
			if c.Text != "" {
				sb.WriteString(fmt.Sprintf("%s\tText: \"%s\",\n", indent, c.Text))
			}
			if c.Enabled == false {
				sb.WriteString(fmt.Sprintf("%s\tEnabled: false,\n", indent))
			}
			if c.Visible == false {
				sb.WriteString(fmt.Sprintf("%s\tVisible: false,\n", indent))
			}

			// Specific properties
			switch c.Type {
			case "LineEdit":
				if c.ReadOnly { sb.WriteString(fmt.Sprintf("%s\tReadOnly: true,\n", indent)) }
				if c.IsPassword { sb.WriteString(fmt.Sprintf("%s\tPasswordMode: true,\n", indent)) }
			case "CheckBox", "RadioButton":
				if c.Checked { sb.WriteString(fmt.Sprintf("%s\tChecked: true,\n", indent)) }
			case "NumberEdit":
				sb.WriteString(fmt.Sprintf("%s\tMinValue: %d,\n", indent, c.MinValue))
				sb.WriteString(fmt.Sprintf("%s\tMaxValue: %d,\n", indent, c.MaxValue))
				sb.WriteString(fmt.Sprintf("%s\tValue: %d,\n", indent, c.Value))
			case "Slider", "ProgressBar":
				sb.WriteString(fmt.Sprintf("%s\tMinValue: %d,\n", indent, c.MinValue))
				sb.WriteString(fmt.Sprintf("%s\tMaxValue: %d,\n", indent, c.MaxValue))
				sb.WriteString(fmt.Sprintf("%s\tValue: %d,\n", indent, c.Value))
				if c.Type == "Slider" && c.Orientation == "Vertical" {
					sb.WriteString(fmt.Sprintf("%s\tOrientation: Vertical,\n", indent))
				}
			case "ComboBox", "ListBox":
				if len(c.Items) > 0 {
					sb.WriteString(fmt.Sprintf("%s\tModel: []string{%s},\n", indent, formatItems(c.Items)))
				}
			}

			// Layout and Children for containers
			if isContainer(c.Type) {
				if c.Layout != nil {
					sb.WriteString(generateLayoutCode(c.Layout, indentCount(indent)+1))
				}

				if c.Type == "TabWidget" {
					sb.WriteString(fmt.Sprintf("%s\tPages: []TabPage{\n", indent))
					sb.WriteString(generateTabPages(p, c.ID, c.TabPages))
					sb.WriteString(fmt.Sprintf("%s\t},\n", indent))
				} else {
					sb.WriteString(fmt.Sprintf("%s\tChildren: []Widget{\n", indent))
					sb.WriteString(generateChildren(p, c.ID))
					sb.WriteString(fmt.Sprintf("%s\t},\n", indent))
				}
			}

			sb.WriteString(fmt.Sprintf("%s},\n", indent))
		}
	}
	return sb.String()
}

func isContainer(typ string) bool {
	switch typ {
	case "Composite", "GroupBox", "Splitter", "TabWidget":
		return true
	}
	return false
}

func generateTabPages(p *models.Project, parentID string, configs []models.TabPageConfig) string {
	var sb strings.Builder
	indent := "\t\t\t\t"

	if len(configs) > 0 {
		for _, tp := range configs {
			sb.WriteString(fmt.Sprintf("%sTabPage{\n", indent))
			sb.WriteString(fmt.Sprintf("%s\tTitle: \"%s\",\n", indent, tp.Title))
			if tp.Layout != nil {
				sb.WriteString(generateLayoutCode(tp.Layout, indentCount(indent)+1))
			}
			sb.WriteString(fmt.Sprintf("%s\tChildren: []Widget{\n", indent))
			sb.WriteString(fmt.Sprintf("%s\t},\n", indent))
			sb.WriteString(fmt.Sprintf("%s},\n", indent))
		}
	} else {
		for _, c := range p.Components {
			if c.ParentID == parentID {
				sb.WriteString(fmt.Sprintf("%sTabPage{\n", indent))
				sb.WriteString(fmt.Sprintf("%s\tTitle: \"%s\",\n", indent, c.Text))
				if c.Layout != nil {
					sb.WriteString(generateLayoutCode(c.Layout, indentCount(indent)+1))
				}
				sb.WriteString(fmt.Sprintf("%s\tChildren: []Widget{\n", indent))
				sb.WriteString(generateChildren(p, c.ID))
				sb.WriteString(fmt.Sprintf("%s\t},\n", indent))
				sb.WriteString(fmt.Sprintf("%s},\n", indent))
			}
		}
	}
	return sb.String()
}

func generateLayoutCode(cfg *models.LayoutConfig, indent int) string {
	if cfg == nil || cfg.Type == models.LayoutNone { return "" }
	indentStr := strings.Repeat("\t", indent)

	var params []string
	if cfg.MarginsZero { params = append(params, "MarginsZero: true") }
	if cfg.SpacingZero { params = append(params, "SpacingZero: true") }
	if cfg.Type == models.LayoutGrid { params = append(params, fmt.Sprintf("Columns: %d", cfg.Columns)) }

	return fmt.Sprintf("%sLayout: %s{%s},\n", indentStr, cfg.Type, strings.Join(params, ", "))
}

func formatItems(items []string) string {
	var formatted []string
	for _, it := range items {
		formatted = append(formatted, fmt.Sprintf("\"%s\"", it))
	}
	return strings.Join(formatted, ", ")
}

func indentCount(indent string) int {
	return len(indent)
}
