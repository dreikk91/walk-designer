# Walk GUI Designer Pro

A visual interface designer for the `github.com/lxn/walk` Go library.

## Features
- **Toolbox**: Drag & drop or double-click to add components.
- **Canvas**: Visual editing with hit-testing, selection, moving, and resizing.
- **Property Editor**: Context-aware property editing with data binding.
- **Object List**: Hierarchical view of all components.
- **Layout Management**: Support for VBox, HBox, and Grid layouts with live preview on canvas.
- **Live Preview**: Launch a real Walk window to test your design.
- **Code Generation**: Generate Go declarative code for your project.
- **Storage**: Save and load projects in JSON format.

## Installation

Ensure you have Go installed. This application is for Windows only.

```bash
go get github.com/lxn/walk
go build
```

## How to Use
1. Select a component from the **Toolbox** and double-click to add it to the canvas.
2. Drag components on the **Canvas** to move them.
3. Use the handles around a selected component to resize it.
4. Dropping a component onto another container (e.g., GroupBox, Composite) will make it a child of that container.
5. Change properties in the **Property Editor** on the right and click **Apply**.
6. Set the **Layout** property on a container to automatically arrange its children.
7. Use **File -> Save** to save your project as JSON.
8. Use **Design -> Preview** (Ctrl+P) to see the live result.
9. Use **Design -> Generate Code** (Ctrl+G) to get the Go code.

## Supported Components
- **Basic**: PushButton, Label, LineEdit, TextEdit, CheckBox, RadioButton, ComboBox, ListBox
- **Numeric**: NumberEdit, DateEdit, Slider, ProgressBar
- **Advanced**: TableView, TreeView, TabWidget
- **Containers**: Composite, GroupBox, Splitter
- **Spacers**: HSpacer, VSpacer

## Sample Project
A `sample_project.json` is provided to get you started.
