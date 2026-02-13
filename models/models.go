package models

type LayoutType string

const (
	LayoutNone LayoutType = ""
	LayoutHBox LayoutType = "HBox"
	LayoutVBox LayoutType = "VBox"
	LayoutGrid LayoutType = "Grid"
)

type Component struct {
	ID          string       `json:"id"`
	Type        string       `json:"type"`
	Name        string       `json:"name"`
	Text        string       `json:"text"`
	ParentID    string       `json:"parent_id"`
	X           int          `json:"x"`
	Y           int          `json:"y"`
	Width       int          `json:"width"`
	Height      int          `json:"height"`
	Enabled     bool         `json:"enabled"`
	Visible     bool         `json:"visible"`
	Layout      LayoutType   `json:"layout,omitempty"`
	Columns     int          `json:"columns,omitempty"`
	MarginsZero bool         `json:"margins_zero,omitempty"`
	SpacingZero bool         `json:"spacing_zero,omitempty"`
	Items       []string     `json:"items,omitempty"`
	Min         int          `json:"min,omitempty"`
	Max         int          `json:"max,omitempty"`
	Value       int          `json:"value,omitempty"`
	Checked     bool         `json:"checked,omitempty"`
	ReadOnly    bool         `json:"read_only,omitempty"`
	Orientation int          `json:"orientation,omitempty"` // 1 for Horizontal, 2 for Vertical (lxn/walk constants)
	Children    []*Component `json:"children,omitempty"`    // For legacy support and tree representation
}

type Project struct {
	Title      string       `json:"window_title"`
	Width      int          `json:"window_width"`
	Height     int          `json:"window_height"`
	Layout     LayoutType   `json:"layout"`
	Components []*Component `json:"components"` // We'll keep it as a flat list for the core, but can be nested in JSON
}
