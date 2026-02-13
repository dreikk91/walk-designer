package logic

import (
	"encoding/json"
	"os"
	"walk-gui-designer-pro/models"
)

func SaveProject(p *models.Project, filename string) error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func LoadProject(filename string) (*models.Project, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var p models.Project
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}

	// Normalization
	normalizeProject(&p)

	return &p, nil
}

func normalizeProject(p *models.Project) {
	// If project has children in components (legacy), flatten them
	flat := []*models.Component{}
	var flatten func([]*models.Component, string)
	flatten = func(comps []*models.Component, parentID string) {
		for _, c := range comps {
			if c.ParentID == "" {
				c.ParentID = parentID
			}
			// We can't easily detect missing fields in JSON with bools 
			// unless we use pointers, but we can assume true for new components.
			flat = append(flat, c)
			if len(c.Children) > 0 {
				flatten(c.Children, c.ID)
				c.Children = nil // Clear after flattening
			}
		}
	}

	flatten(p.Components, "")
	p.Components = flat
}
