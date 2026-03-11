package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	modulePath         = "github.com/Hughzu/trackstack/apps/server"
	maxAdapterBranches = 8
)

type violation struct {
	File   string
	Symbol string
	Rule   string
}

func main() {
	files, err := filepath.Glob("internal/modules/*/adapters/db/*.go")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "lintguard: list adapter files: %v\n", err)
		os.Exit(1)
	}

	violations := make([]violation, 0)
	for _, filePath := range files {
		if strings.HasSuffix(filePath, "_test.go") {
			continue
		}

		src, err := os.ReadFile(filePath)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "lintguard: read %s: %v\n", filePath, err)
			os.Exit(1)
		}

		fileViolations, err := checkSource(filepath.ToSlash(filePath), string(src))
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "lintguard: %v\n", err)
			os.Exit(1)
		}
		violations = append(violations, fileViolations...)
	}

	sort.Slice(violations, func(i, j int) bool {
		if violations[i].File != violations[j].File {
			return violations[i].File < violations[j].File
		}
		if violations[i].Symbol != violations[j].Symbol {
			return violations[i].Symbol < violations[j].Symbol
		}
		return violations[i].Rule < violations[j].Rule
	})

	if len(violations) > 0 {
		for _, v := range violations {
			if v.Symbol == "" {
				_, _ = fmt.Fprintf(os.Stderr, "lintguard: %s (%s)\n", v.File, v.Rule)
				continue
			}
			_, _ = fmt.Fprintf(os.Stderr, "lintguard: %s [%s] (%s)\n", v.File, v.Symbol, v.Rule)
		}
		os.Exit(1)
	}

	fmt.Printf("lintguard: ok (%d adapter files checked)\n", len(files))
}

func checkSource(relPath string, src string) ([]violation, error) {
	moduleName, err := adapterModuleName(relPath)
	if err != nil {
		return nil, err
	}

	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, relPath, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", relPath, err)
	}

	violations := make([]violation, 0)
	importAliases := make(map[string]string)
	for _, spec := range file.Imports {
		importPath, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid import in %s: %w", relPath, err)
		}

		alias := importAlias(spec, importPath)
		if alias != "" {
			importAliases[alias] = importPath
		}

		if rule := validateAdapterImport(moduleName, importPath); rule != "" {
			violations = append(violations, violation{File: relPath, Rule: rule})
		}
	}

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}

		symbol := fn.Name.Name
		if fn.Recv != nil && len(fn.Recv.List) > 0 {
			symbol = recvName(fn.Recv.List[0].Type) + "." + fn.Name.Name
		}

		branches := branchCount(fn.Body)
		if branches > maxAdapterBranches {
			violations = append(violations, violation{
				File:   relPath,
				Symbol: symbol,
				Rule:   fmt.Sprintf("db adapter branch count %d exceeds max %d", branches, maxAdapterBranches),
			})
		}

		for _, rule := range forbiddenCalls(fn.Body, importAliases) {
			violations = append(violations, violation{File: relPath, Symbol: symbol, Rule: rule})
		}
	}

	return violations, nil
}

func adapterModuleName(relPath string) (string, error) {
	parts := strings.Split(path.Clean(relPath), "/")
	if len(parts) < 5 || parts[0] != "internal" || parts[1] != "modules" || parts[3] != "adapters" || parts[4] != "db" {
		return "", fmt.Errorf("unexpected adapter path: %s", relPath)
	}
	return parts[2], nil
}

func validateAdapterImport(moduleName string, importPath string) string {
	ownModuleImport := modulePath + "/internal/modules/" + moduleName
	if importPath == ownModuleImport {
		return ""
	}

	if strings.HasPrefix(importPath, modulePath+"/") {
		return "db adapters may only import their own module package from the application"
	}

	if isStdlibImport(importPath) {
		return ""
	}

	return "db adapters may not import third-party packages"
}

func isStdlibImport(importPath string) bool {
	first := importPath
	if idx := strings.Index(importPath, "/"); idx >= 0 {
		first = importPath[:idx]
	}
	return !strings.Contains(first, ".")
}

func importAlias(spec *ast.ImportSpec, importPath string) string {
	if spec.Name != nil {
		if spec.Name.Name == "_" || spec.Name.Name == "." {
			return ""
		}
		return spec.Name.Name
	}
	return path.Base(importPath)
}

func recvName(expr ast.Expr) string {
	switch value := expr.(type) {
	case *ast.Ident:
		return value.Name
	case *ast.StarExpr:
		return recvName(value.X)
	default:
		return "receiver"
	}
}

func branchCount(node ast.Node) int {
	count := 0
	ast.Inspect(node, func(current ast.Node) bool {
		switch current.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
			count++
		}
		return true
	})
	return count
}

func forbiddenCalls(body *ast.BlockStmt, importAliases map[string]string) []string {
	rules := make([]string, 0)
	ast.Inspect(body, func(current ast.Node) bool {
		call, ok := current.(*ast.CallExpr)
		if !ok {
			return true
		}

		selector, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		pkgIdent, ok := selector.X.(*ast.Ident)
		if !ok {
			return true
		}

		importPath := importAliases[pkgIdent.Name]
		switch {
		case importPath == "time" && selector.Sel.Name == "Now":
			rules = append(rules, "db adapters must not generate timestamps; receive them from module services")
		case importPath == "github.com/google/uuid" && (selector.Sel.Name == "New" || selector.Sel.Name == "NewString"):
			rules = append(rules, "db adapters must not generate IDs; receive them from module services")
		case importPath == "strings" && (selector.Sel.Name == "TrimSpace" || selector.Sel.Name == "ToLower" || selector.Sel.Name == "ToUpper"):
			rules = append(rules, "db adapters must not normalize domain input; validate before persistence")
		}

		return true
	})

	return rules
}
