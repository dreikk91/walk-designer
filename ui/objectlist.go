package ui

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"walk-gui-designer-pro/models"
)

type ObjectListWidget struct {
	*walk.TreeView
	mw    *DesignerWindow
	model *ObjectModel
}

func CreateObjectList(mw *DesignerWindow) Widget {
	ol := &ObjectListWidget{mw: mw}
	mw.ObjectList = ol
	ol.model = &ObjectModel{mw: mw}

	return TreeView{
		AssignTo: &ol.TreeView,
		Model:    ol.model,
		OnCurrentItemChanged: func() {
			if item := ol.CurrentItem(); item != nil {
				mw.SelectComponent(item.(*ObjectItem).Component)
			}
		},
	}
}

func (ol *ObjectListWidget) Refresh() {
	ol.model.PublishItemsReset(nil)
}

type ObjectItem struct {
	walk.TreeItem
	Component  *models.Component
	parentItem *ObjectItem
	Children   []*ObjectItem
}

func (i *ObjectItem) Text() string {
	if i.Component == nil {
		return "Project"
	}
	return i.Component.Name + " [" + i.Component.Type + "]"
}

func (i *ObjectItem) Parent() walk.TreeItem {
	if i.parentItem == nil {
		return nil
	}
	return i.parentItem
}

func (i *ObjectItem) ChildAt(index int) walk.TreeItem {
	return i.Children[index]
}

func (i *ObjectItem) ChildCount() int {
	return len(i.Children)
}

func (i *ObjectItem) Image() interface{} {
	return nil
}

type ObjectModel struct {
	walk.TreeModelBase
	mw *DesignerWindow
}

func (m *ObjectModel) RootCount() int {
	return 1 // The Project root
}

func (m *ObjectModel) RootAt(index int) walk.TreeItem {
	root := &ObjectItem{Component: nil}
	m.fillChildren(root)
	return root
}

func (m *ObjectModel) fillChildren(item *ObjectItem) {
	parentID := ""
	if item.Component != nil {
		parentID = item.Component.ID
	}

	for _, c := range m.mw.Project.Components {
		if c.ParentID == parentID {
			child := &ObjectItem{Component: c, parentItem: item}
			item.Children = append(item.Children, child)
			m.fillChildren(child)
		}
	}
}
