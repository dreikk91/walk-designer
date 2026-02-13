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
			sb.WriteString(fmt.Sprintf("\tvar %s *walk.%s\n", c.Name, c.Type))
		}
	}
	sb.WriteString("\n")

	sb.WriteString("\tif _, err := (MainWindow{\n")
	sb.WriteString(fmt.Sprintf("\t\tTitle: \"%s\",\n", p.Title))
	sb.WriteString(fmt.Sprintf("\t\tMinSize: Size{Width: %d, Height: %d},\n", p.Width, p.Height))
	if p.Layout != models.LayoutNone {
		sb.WriteString(fmt.Sprintf("\t\tLayout: %s{},\n", p.Layout))
	}
	sb.WriteString("\t\tChildren: []Widget{\n")

	// Recursive children generation
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
				if c.Orientation == 2 {
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

			// Layout and Children for containers
			if isContainer(c.Type) {
				if c.Layout != models.LayoutNone {
					sb.WriteString(fmt.Sprintf("%s\tLayout: %s{},\n", indent, c.Layout))
				}

				if c.Type == "TabWidget" {
					sb.WriteString(fmt.Sprintf("%s\tPages: []TabPage{\n", indent))
					sb.WriteString(generateTabPages(p, c.ID))
					sb.WriteString(fmt.Sprintf("%s\t},\n", indent))
				} else {
					sb.WriteString(fmt.Sprintf("%s\tChildren: []Widget{\n", indent))
					sb.WriteString(generateChildren(p, c.ID))
					sb.WriteString(fmt.Sprintf("%s\t},\n", indent))
				}
			}

			// Other properties
			if c.Type == "PushButton" {
				// sb.WriteString(fmt.Sprintf("%s\tOnClicked: func() { ... },\n", indent))
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

func generateTabPages(p *models.Project, parentID string) string {
	var sb strings.Builder
	indent := "\t\t\t\t"

	for _, c := range p.Components {
		if c.ParentID == parentID {
			sb.WriteString(fmt.Sprintf("%sTabPage{\n", indent))
			sb.WriteString(fmt.Sprintf("%s\tTitle: \"%s\",\n", indent, c.Text))
			if c.Layout != models.LayoutNone {
				sb.WriteString(fmt.Sprintf("%s\tLayout: %s{},\n", indent, c.Layout))
			}
			sb.WriteString(fmt.Sprintf("%s\tChildren: []Widget{\n", indent))
			sb.WriteString(generateChildren(p, c.ID))
			sb.WriteString(fmt.Sprintf("%s\t},\n", indent))
			sb.WriteString(fmt.Sprintf("%s},\n", indent))
		}
	}
	return sb.String()
}
