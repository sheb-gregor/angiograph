package analysis

import (
	"encoding/json"
	"fmt"
	"go/types"
	"strings"

	"github.com/sheb-gregor/angiograph/data"
	"github.com/xlab/treeprint"
	"golang.org/x/tools/go/ssa/ssautil"
)

type PkgShort struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func (s PkgShort) String() string {
	return s.Name + " (" + s.Path + ")"
}

type PkgTree struct {
	Root    string    `json:"-"`
	Self    PkgShort  `json:"-"`
	ID      string    `json:"id"`
	Imports []PkgTree `json:"nodes"`
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
	tree.ID = tree.Self.String()
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
	return tree.addToTreePrint(treePrint)
}

func (tree *PkgTree) addToTreePrint(treePrint treeprint.Tree) treeprint.Tree {
	if len(tree.Imports) > 0 {
		treePrint = treePrint.AddBranch(tree.Self.Name + " (" + tree.Self.Path + ")")
	} else {
		treePrint = treePrint.AddNode(tree.Self.Name + " (" + tree.Self.Path + ")")
	}

	for _, pkgTree := range tree.Imports {
		pkgTree.addToTreePrint(treePrint)
	}
	return treePrint
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

	for _, pkg := range mainPkgs {
		// print("", pkg.Pkg, pkg.Pkg.Path())

		// fmt.Println(NewPkgTree(pkg.Pkg).IntoJSON())
		// fmt.Println(NewPkgTree(pkg.Pkg).IntoTree().String())
		index := data.NewPkgIndex(pkg.Pkg)

		fmt.Println(index.IntoPUML(pkg.Pkg.Path()))
		// fmt.Println(NewPkgTree(pkg.Pkg).IntoJSON())
	}

	return "", err
}

func print(prefix string, pkg *types.Package, root string) {

	_, ok := data.STDLib[strings.Split(pkg.Path(), "/")[0]]
	if ok {
		return
	}

	if !strings.Contains(pkg.Path(), root) {
		return
	}

	for _, spkg := range pkg.Imports() {
		fmt.Println(prefix+"-->", spkg.Path())
		print(prefix+"----", spkg, root)
	}
}
