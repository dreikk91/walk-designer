package models

type LayoutType string

const (
	LayoutNone LayoutType = "None"
	LayoutVBox LayoutType = "VBox"
	LayoutHBox LayoutType = "HBox"
	LayoutGrid LayoutType = "Grid"
)

type LayoutConfig struct {
	Type        LayoutType `json:"type"`
	Margins     int        `json:"margins"`
	Spacing     int        `json:"spacing"`
	Columns     int        `json:"columns,omitempty"`
	MarginsZero bool       `json:"margins_zero,omitempty"`
	SpacingZero bool       `json:"spacing_zero,omitempty"`
}

type TabPageConfig struct {
	Title    string        `json:"title"`
	Layout   *LayoutConfig `json:"layout,omitempty"`
	Children []Component   `json:"children,omitempty"`
}

type FontConfig struct {
	Family    string `json:"family,omitempty"`
	PointSize int    `json:"point_size,omitempty"`
	Bold      bool   `json:"bold,omitempty"`
	Italic    bool   `json:"italic,omitempty"`
	Underline bool   `json:"underline,omitempty"`
	StrikeOut bool   `json:"strike_out,omitempty"`
}

type TableColumnConfig struct {
	Name      string `json:"name,omitempty"`
	Title     string `json:"title,omitempty"`
	Width     int    `json:"width,omitempty"`
	Alignment string `json:"alignment,omitempty"` // e.g. "Near", "Center", "Far"
}

type Component struct {
	ID          string         `json:"id"`
	Type        string         `json:"type"`
	Name        string         `json:"name"`
	Text        string         `json:"text"`
	ParentID    string         `json:"parent_id,omitempty"`
	X           int            `json:"x"`
	Y           int            `json:"y"`
	Width       int            `json:"width"`
	Height      int            `json:"height"`
	Enabled     bool           `json:"enabled"`
	Visible     bool           `json:"visible"`

	// Layout
	Layout      *LayoutConfig  `json:"layout,omitempty"`

	// Extra properties
	MinValue    int             `json:"min_value,omitempty"`
	MaxValue    int             `json:"max_value,omitempty"`
	Value       int             `json:"value,omitempty"`
	Checked     bool            `json:"checked,omitempty"`
	ReadOnly    bool            `json:"read_only,omitempty"`
	Items       []string        `json:"items,omitempty"`
	TableCols   []TableColumnConfig `json:"table_cols,omitempty"`
	ColumnSpan  int             `json:"column_span,omitempty"`
	RowSpan     int             `json:"row_span,omitempty"`
	Alignment   string          `json:"alignment,omitempty"`
	TabPages    []TabPageConfig `json:"tab_pages,omitempty"`
	Orientation string          `json:"orientation,omitempty"`
	IsPassword  bool            `json:"is_password,omitempty"`
	ToolTip     string          `json:"tool_tip,omitempty"`
	ImagePath   string          `json:"image_path,omitempty"`

	// Advanced / Appearance
	Font        *FontConfig     `json:"font,omitempty"`
	BgColor     string          `json:"bg_color,omitempty"`   // e.g. "#FFFFFF" or "255,255,255"
	TextColor   string          `json:"text_color,omitempty"`
	MaxLength   int             `json:"max_length,omitempty"`
	Marquee     bool            `json:"marquee,omitempty"`
	ImageMode   string          `json:"image_mode,omitempty"` // e.g. "Ideal", "Stretch", "Zoom"

	Children    []*Component    `json:"children,omitempty"` // For tree representation
}

type Project struct {
	Title      string       `json:"window_title"`
	Width      int          `json:"window_width"`
	Height     int          `json:"window_height"`
	Layout     *LayoutConfig `json:"layout,omitempty"`
	Components []*Component `json:"components"`
}
