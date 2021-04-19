package analysis

import (
	"encoding/json"
	"fmt"
	"go/types"
	"strings"

	"github.com/xlab/treeprint"
	"golang.org/x/tools/go/ssa/ssautil"
)

type PkgShort struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type PkgTree struct {
	Root    string    `json:"root"`
	Self    PkgShort  `json:"self"`
	Imports []PkgTree `json:"imports"`
}

func NewPkgTree(pkg *types.Package) *PkgTree {
	tree := new(PkgTree)
	tree.Set(pkg, pkg.Path())
	return tree
}

func (tree *PkgTree) Set(pkg *types.Package, root string) {
	tree.Self = PkgShort{
		Name: pkg.Name(),
		Path: pkg.Path(),
	}
	tree.Root = root
	if !strings.HasPrefix(tree.Self.Path, root) {
		return
	}

	tree.Imports = make([]PkgTree, len(pkg.Imports()))
	for i, p := range pkg.Imports() {
		tree.Imports[i].Set(p, root)
	}
}

func (tree *PkgTree) IntoJSON() string {
	res, _ := json.MarshalIndent(tree, "", "  ")
	return string(res)
}

func (tree *PkgTree) IntoTree() treeprint.Tree {
	treePrint := treeprint.New()
	tree.addToTreePrint(treePrint)
	return treePrint
}

func (tree *PkgTree) addToTreePrint(treePrint treeprint.Tree) {
	if len(tree.Imports) > 0 {
		treePrint.AddBranch(tree.Self.Name + " (" + tree.Self.Path + ")")
	} else {
		treePrint.AddNode(tree.Self.Name + " (" + tree.Self.Path + ")")
	}

	for _, pkgTree := range tree.Imports {
		pkgTree.addToTreePrint(treePrint)
	}
}

func (opts AnalyseOpts) ImportsTree() (string, error) {
	lPkgs, err := opts.LoadMeta()
	if err != nil {
		return "", err
	}

	// Create and build SSA-form program representation.
	prog, pkgs := ssautil.AllPackages(lPkgs, 0)
	prog.Build()
	mainPkgs, err := mainPackages(pkgs)
	if err != nil {
		return "", err
	}
	// resq, err := json.MarshalIndent(mainPkgs, "", "  ")
	// fmt.Println(string(resq))

	// res, err := json.MarshalIndent(lPkgs, "", "  ")
	// fmt.Println(string(res))

	for _, pkg := range mainPkgs {
		// fmt.Println(NewPkgTree(pkg.Pkg).IntoTree().String())
		fmt.Println(NewPkgTree(pkg.Pkg).IntoJSON())
	}

	return "", err
}
