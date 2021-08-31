package analysis

import (
	"encoding/json"
	"go/types"
	"strings"
)

type BiDirectPkgTree struct {
	Root string `json:"root"`
	Name string `json:"name"`
	Path string `json:"path"`
	// UsedBy  PkgIndex `json:"used_by"`
	// Imports PkgIndex `json:"imports"`
	UsedBy  map[string]struct{} `json:"used_by"`
	Imports map[string]struct{} `json:"imports"`
}

func (BiDirectPkgTree) New(name, root, path string) *BiDirectPkgTree {
	return &BiDirectPkgTree{
		Root: root,
		Name: name,
		Path: path,
		// UsedBy:  map[string]*BiDirectPkgTree{},
		// Imports: map[string]*BiDirectPkgTree{},
		UsedBy:  map[string]struct{}{},
		Imports: map[string]struct{}{},
	}

}

func (tree *BiDirectPkgTree) AddUsedBy(newTree *BiDirectPkgTree) {
	// prev, ok := tree.UsedBy[newTree.Name]
	// if !ok {
	// 	prev = newTree
	// } else {
	// 	prev.Merge(newTree)
	// }
	// tree.UsedBy[newTree.Name] = prev
	tree.UsedBy[newTree.Name] = struct{}{}
}

func (tree *BiDirectPkgTree) AddImport(newTree *BiDirectPkgTree) {
	// prev, ok := tree.Imports[newTree.Name]
	// if !ok {
	// 	prev = newTree
	// } else {
	// 	prev.Merge(newTree)
	// }
	// tree.Imports[newTree.Name] = prev
	tree.Imports[newTree.Name] = struct{}{}
}

func (tree *BiDirectPkgTree) Merge(newTree *BiDirectPkgTree) {
	for s := range newTree.Imports {
		// childTree, ok := tree.Imports[s]
		// if !ok {
		// 	childTree = newTree
		// } else {
		// 	childTree.Merge(newTree.Imports[s])
		// }
		tree.Imports[s] = struct{}{}
	}

	for s := range newTree.UsedBy {
		// childTree, ok := tree.UsedBy[s]
		// if !ok {
		// 	childTree = newTree
		// } else {
		// 	childTree.Merge(newTree.UsedBy[s])
		// }
		// tree.UsedBy[s] = childTree
		tree.UsedBy[s] = struct{}{}
	}
}

type PkgIndex map[string]*BiDirectPkgTree

func (index PkgIndex) Insert(tree *BiDirectPkgTree) {
	prev, ok := index[tree.Name]
	if !ok {
		prev = tree
	} else {
		prev.Merge(tree)
	}

	index[tree.Name] = prev

	for s := range tree.Imports {
		prev = index[s]
		if !ok {
			prev = BiDirectPkgTree{}.New(s, tree.Root, "")
		}
		prev.AddUsedBy(tree)
		index[s] = prev
		// imp := tree.Imports[s]
		// imp.AddUsedBy(tree)
		// index.Insert(imp)
	}

	for s := range tree.UsedBy {
		prev = index[s]
		if !ok {
			prev = BiDirectPkgTree{}.New(s, tree.Root, "")
		}
		prev.AddImport(tree)
		index[s] = prev
		// imp := tree.Imports[s]
		// imp.AddImport(tree)
		// index.Insert(imp)
	}
}

func NewPkgIndex(pkg *types.Package) PkgIndex {
	tree := PkgIndex{}
	tree.Set(pkg, pkg.Path())
	return tree
}

func (index PkgIndex) Set(pkg *types.Package, root string) {
	name := pkg.Path()

	_, ok := stdLib[strings.Split(pkg.Path(), "/")[0]]
	if ok {
		return
	}

	if !strings.Contains(pkg.Path(), root) {
		return
	}

	node, ok := index[name]
	if !ok {
		node = BiDirectPkgTree{}.New(name, root, pkg.Path())
		index[name] = node
	}

	imports := pkg.Imports()
	for i := range imports {
		child, ok := index[imports[i].Path()]
		if !ok {
			child = BiDirectPkgTree{}.New(imports[i].Path(), root, imports[i].Path())
			index[imports[i].Path()] = child
		}
		child.AddUsedBy(node)
		node.AddImport(child)

		index.Set(imports[i], root)
	}

}

func (index PkgIndex) IntoJSON() string {
	type treeNode struct {
		Root    string   `json:"root"`
		Name    string   `json:"name"`
		UsedBy  []string `json:"used_by"`
		Imports []string `json:"imports"`
	}

	resArr := make([]treeNode, 0, len(index))
	for _, val := range index {
		n := treeNode{
			Root:    val.Root,
			Name:    val.Name,
			UsedBy:  make([]string, 0, len(val.UsedBy)),
			Imports: make([]string, 0, len(val.Imports)),
		}
		for s := range val.UsedBy {
			n.UsedBy = append(n.UsedBy, s)
		}
		for s := range val.Imports {
			n.Imports = append(n.Imports, s)
		}
		resArr = append(resArr, n)
	}
	res, _ := json.MarshalIndent(resArr, "", "  ")
	return string(res)
}
