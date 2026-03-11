package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type listedPackage struct {
	ImportPath string
	Imports    []string
	Module     *listedModule
}

type listedModule struct {
	Path string
}

type packageKind struct {
	Layer     string
	Module    string
	IsModule  bool
	IsAdapter bool
}

type violation struct {
	From string
	To   string
	Rule string
}

func main() {
	packages, err := listPackages()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "archguard: %v\n", err)
		os.Exit(1)
	}

	violations := checkPackages(packages)
	if len(violations) > 0 {
		for _, v := range violations {
			_, _ = fmt.Fprintf(os.Stderr, "archguard: %s imports %s (%s)\n", v.From, v.To, v.Rule)
		}
		os.Exit(1)
	}

	fmt.Printf("archguard: ok (%d packages checked)\n", len(packages))
}

func listPackages() ([]listedPackage, error) {
	cmd := exec.Command("go", "list", "-json", "./...")
	cmd.Dir = "."

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			return nil, fmt.Errorf("go list: %w", err)
		}
		return nil, fmt.Errorf("go list: %w: %s", err, message)
	}

	decoder := json.NewDecoder(&stdout)
	packages := make([]listedPackage, 0)
	for {
		var pkg listedPackage
		if err := decoder.Decode(&pkg); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("decode go list output: %w", err)
		}
		if pkg.Module == nil || pkg.Module.Path == "" {
			continue
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}

func checkPackages(packages []listedPackage) []violation {
	violations := make([]violation, 0)

	for _, pkg := range packages {
		fromRel, ok := localImportPath(pkg.ImportPath, pkg.Module.Path)
		if !ok {
			continue
		}

		fromKind := classifyPackage(fromRel)
		for _, imported := range pkg.Imports {
			toRel, local := localImportPath(imported, pkg.Module.Path)
			if !local {
				continue
			}

			toKind := classifyPackage(toRel)

			if fromKind.IsModule {
				switch toKind.Layer {
				case "transport":
					violations = append(violations, violation{From: fromRel, To: toRel, Rule: "modules must not depend on transport"})
				case "wiring":
					violations = append(violations, violation{From: fromRel, To: toRel, Rule: "modules must not depend on wiring"})
				case "cmd":
					violations = append(violations, violation{From: fromRel, To: toRel, Rule: "modules must not depend on commands"})
				case "core":
					violations = append(violations, violation{From: fromRel, To: toRel, Rule: "modules must not depend on core runtime packages"})
				}

				if toKind.IsModule && toKind.Module != fromKind.Module {
					violations = append(violations, violation{From: fromRel, To: toRel, Rule: "modules must not import other modules directly"})
				}
			}

			if toKind.IsAdapter && fromKind.Layer != "wiring" {
				violations = append(violations, violation{From: fromRel, To: toRel, Rule: "only wiring packages may import concrete adapters"})
			}
		}
	}

	sort.Slice(violations, func(i, j int) bool {
		if violations[i].From != violations[j].From {
			return violations[i].From < violations[j].From
		}
		if violations[i].To != violations[j].To {
			return violations[i].To < violations[j].To
		}
		return violations[i].Rule < violations[j].Rule
	})

	return violations
}

func localImportPath(importPath string, modulePath string) (string, bool) {
	if importPath == modulePath {
		return "", false
	}

	prefix := modulePath + "/"
	if !strings.HasPrefix(importPath, prefix) {
		return "", false
	}

	return strings.TrimPrefix(importPath, prefix), true
}

func classifyPackage(rel string) packageKind {
	if strings.HasPrefix(rel, "cmd/") {
		return packageKind{Layer: "cmd"}
	}
	if strings.HasPrefix(rel, "internal/transport/") {
		return packageKind{Layer: "transport"}
	}
	if strings.HasPrefix(rel, "internal/wiring/") {
		return packageKind{Layer: "wiring"}
	}
	if strings.HasPrefix(rel, "internal/core/") {
		return packageKind{Layer: "core"}
	}

	rest, ok := strings.CutPrefix(rel, "internal/modules/")
	if !ok {
		return packageKind{Layer: "other"}
	}

	parts := strings.Split(rest, "/")
	if len(parts) == 0 || parts[0] == "" {
		return packageKind{Layer: "other"}
	}

	kind := packageKind{
		Layer:    "module",
		Module:   parts[0],
		IsModule: true,
	}

	if len(parts) > 1 && parts[1] == "adapters" {
		kind.IsAdapter = true
	}

	return kind
}
