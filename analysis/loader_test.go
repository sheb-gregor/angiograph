package analysis

//
// func TestLoad(t *testing.T) {
//
// 	mode := packages.NeedName |
// 		packages.NeedFiles |
// 		packages.NeedCompiledGoFiles |
// 		packages.NeedImports |
// 		packages.NeedDeps |
// 		packages.NeedExportsFile |
// 		packages.NeedTypes |
// 		packages.NeedTypesSizes |
// 		packages.NeedSyntax |
// 		packages.NeedTypesInfo |
// 		packages.NeedModule
//
// 	// Load, parse, and type-check the packages named on the command line.
// 	cfg := &packages.Config{
// 		Mode:  mode,
// 		Tests: false,
// 		// BuildFlags: app.BuildFlags,
// 	}
//
// 	lpkgs, err := packages.Load(cfg)
// 	if err != nil {
// 		t.Error(err)
// 		t.FailNow()
// 	}
//
// 	// -deps: print dependencies too.
//
// 	// We can't use packages.All because
// 	// we need an ordered traversal.
// 	var all []*packages.Package // postorder
// 	seen := make(map[*packages.Package]bool)
// 	var visit func(*packages.Package)
// 	visit = func(lpkg *packages.Package) {
// 		if !seen[lpkg] {
// 			seen[lpkg] = true
//
// 			// visit imports
// 			var importPaths []string
// 			for path := range lpkg.Imports {
// 				importPaths = append(importPaths, path)
// 			}
// 			sort.Strings(importPaths) // for determinism
// 			for _, path := range importPaths {
// 				visit(lpkg.Imports[path])
// 			}
//
// 			all = append(all, lpkg)
// 		}
// 	}
// 	for _, lpkg := range lpkgs {
// 		visit(lpkg)
// 	}
// 	lpkgs = all
//
// 	for _, lpkg := range lpkgs {
// 		app.print(lpkg)
// 	}
// }
