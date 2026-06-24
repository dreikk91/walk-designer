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
	// For robustness with lxn/walk, simply expanding all after reset is often easier,
	// but let's try to restore the selected item.

	ol.model.PublishItemsReset(nil)

	// Expand the root by default
	if ol.model.RootCount() > 0 {
		root := ol.model.RootAt(0)
		ol.TreeView.SetExpanded(root, true)

		// Expand all recursively for better UX in small hierarchies
		var expandAll func(item walk.TreeItem)
		expandAll = func(item walk.TreeItem) {
			if item == nil { return }
			ol.TreeView.SetExpanded(item, true)
			count := ol.model.ChildCount(item)
			for i := 0; i < count; i++ {
				expandAll(ol.model.ChildAt(item, i))
			}
		}
		expandAll(root)

		// Select the currently selected component
		selComp := ol.mw.SelectedComponent
		if selComp != nil {
			var findItem func(item walk.TreeItem) walk.TreeItem
			findItem = func(item walk.TreeItem) walk.TreeItem {
				if item == nil { return nil }
				objItem, ok := item.(*ObjectItem)
				if ok && objItem.Component != nil && objItem.Component.ID == selComp.ID {
					return item
				}
				count := ol.model.ChildCount(item)
				for i := 0; i < count; i++ {
					if found := findItem(ol.model.ChildAt(item, i)); found != nil {
						return found
					}
				}
				return nil
			}
			found := findItem(root)
			if found != nil {
				ol.TreeView.SetCurrentItem(found)
			}
		}
	}
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

// ChildCount and ChildAt for the model
func (m *ObjectModel) ChildCount(parent walk.TreeItem) int {
	if obj, ok := parent.(*ObjectItem); ok {
		return obj.ChildCount()
	}
	return 0
}

func (m *ObjectModel) ChildAt(parent walk.TreeItem, index int) walk.TreeItem {
	if obj, ok := parent.(*ObjectItem); ok {
		return obj.ChildAt(index)
	}
	return nil
}

type ObjectModel struct {
	walk.TreeModelBase
	mw *DesignerWindow
}

func (m *ObjectModel) RootCount() int {
	if m.mw == nil || m.mw.Project == nil {
		return 0
	}
	return 1 // The Project root
}

func (m *ObjectModel) RootAt(index int) walk.TreeItem {
	if m.mw == nil || m.mw.Project == nil {
		return nil
	}
	root := &ObjectItem{Component: nil}
	m.fillChildren(root)
	return root
}

func (m *ObjectModel) fillChildren(item *ObjectItem) {
	if item == nil || m.mw == nil || m.mw.Project == nil {
		return
	}
	parentID := ""
	if item.Component != nil {
		parentID = item.Component.ID
	}

	for _, c := range m.mw.Project.Components {
		if c.ParentID == parentID {
			child := &ObjectItem{Component: c, parentItem: item}
			item.Children = append(item.Children, child)
			m.fillChildren(child) // Recursively fetch children to build the hierarchy correctly
		}
	}
}
