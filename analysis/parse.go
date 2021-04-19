package analysis

import (
	"encoding/json"
	"fmt"
	"go/types"
	"os"
	"sort"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/types/typeutil"
)

const mode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedCompiledGoFiles |
	packages.NeedImports |
	packages.NeedDeps |
	packages.NeedExportsFile |
	packages.NeedTypes |
	packages.NeedTypesSizes |
	packages.NeedSyntax |
	packages.NeedTypesInfo |
	packages.NeedModule

type AnalyseOpts struct {
	Path      string
	Deps      bool
	Tests     bool
	Private   bool
	PrintJSON bool
	Patterns  []string
}

func (opts *AnalyseOpts) LoadMeta() ([]*packages.Package, error) {
	// Load, parse, and type-check the packages named on the command line.
	cfg := &packages.Config{
		Mode:  mode,
		Tests: opts.Tests,
		Dir:   opts.Path,
		// BuildFlags: opts.BuildFlags,
	}

	lpkgs, err := packages.Load(cfg, opts.Patterns...)
	if err != nil {
		return nil, err
	}

	// -deps: print dependencies too.
	if opts.Deps {
		opts.loadDeps(lpkgs)
	}

	return lpkgs, nil
}

func (opts *AnalyseOpts) loadDeps(lpkgs []*packages.Package) {
	// We can't use packages.All because
	// we need an ordered traversal.
	var all []*packages.Package // postorder
	seen := make(map[*packages.Package]bool)
	var visit func(*packages.Package)
	visit = func(lpkg *packages.Package) {
		if !seen[lpkg] {
			seen[lpkg] = true

			// visit imports
			var importPaths []string
			for path := range lpkg.Imports {
				importPaths = append(importPaths, path)
			}
			sort.Strings(importPaths) // for determinism
			for _, path := range importPaths {
				visit(lpkg.Imports[path])
			}

			all = append(all, lpkg)
		}
	}
	for _, lpkg := range lpkgs {
		visit(lpkg)
	}
	lpkgs = all
}

// mainPackages returns the main packages to analyze.
// Each resulting package is named "main" and has a main function.
func mainPackages(pkgs []*ssa.Package) ([]*ssa.Package, error) {
	var mains []*ssa.Package
	for _, p := range pkgs {
		if p != nil && p.Pkg.Name() == "main" && p.Func("main") != nil {
			mains = append(mains, p)
		}
	}
	if len(mains) == 0 {
		return nil, fmt.Errorf("no main packages")
	}
	return mains, nil
}

// Run takes the args after flag processing and performs the specified query.
func (opts *AnalyseOpts) Run() error {
	lpkgs, err := opts.LoadMeta()
	if err != nil {
		return err
	}
	for _, lpkg := range lpkgs {
		opts.print(lpkg)
	}
	return nil
}

func (opts *AnalyseOpts) print(lpkg *packages.Package) {
	if opts.PrintJSON {
		data, _ := json.MarshalIndent(lpkg, "", "\t")
		os.Stdout.Write(data)
		return
	}
	// title
	var kind string
	if lpkg.Name == "main" {
		kind += "command"
	} else {
		kind += "package"
	}
	fmt.Printf("Go %s %q:\n", kind, lpkg.ID) // unique ID
	fmt.Printf("\tpackage %s\n", lpkg.Name)

	// characterize type info
	if lpkg.Types == nil {
		fmt.Printf("\thas no exported type info\n")
	} else if !lpkg.Types.Complete() {
		fmt.Printf("\thas incomplete exported type info\n")
	} else if len(lpkg.Syntax) == 0 {
		fmt.Printf("\thas complete exported type info\n")
	} else {
		fmt.Printf("\thas complete exported type info and typed ASTs\n")
	}
	if lpkg.Types != nil && lpkg.IllTyped && len(lpkg.Errors) == 0 {
		fmt.Printf("\thas an error among its dependencies\n")
	}

	// source files
	for _, src := range lpkg.GoFiles {
		fmt.Printf("\tfile %s\n", src)
	}

	// imports
	var lines []string
	for importPath, imp := range lpkg.Imports {
		var line string
		if imp.ID == importPath {
			line = fmt.Sprintf("\timport %q", importPath)
		} else {
			line = fmt.Sprintf("\timport %q => %q", importPath, imp.ID)
		}
		lines = append(lines, line)
	}
	sort.Strings(lines)
	for _, line := range lines {
		fmt.Println(line)
	}

	// errors
	for _, err := range lpkg.Errors {
		fmt.Printf("\t%s\n", err)
	}

	// package members (TypeCheck or WholeProgram mode)
	if lpkg.Types != nil {
		qual := types.RelativeTo(lpkg.Types)
		scope := lpkg.Types.Scope()
		for _, name := range scope.Names() {
			obj := scope.Lookup(name)
			if !obj.Exported() && !opts.Private {
				continue // skip unexported names
			}

			fmt.Printf("\t%s\n", types.ObjectString(obj, qual))
			if _, ok := obj.(*types.TypeName); ok {
				for _, meth := range typeutil.IntuitiveMethodSet(obj.Type(), nil) {
					if !meth.Obj().Exported() && !opts.Private {
						continue // skip unexported names
					}
					fmt.Printf("\t%s\n", types.SelectionString(meth, qual))
				}
			}
		}
	}

	fmt.Println()
}
