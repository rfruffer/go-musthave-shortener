package main

import (
	"go/ast"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// OsExitInMainAnalyzer запрещает прямой вызов os.Exit в функции main пакета main
// внутри нашего проекта (папка /cmd/...), чтобы игнорировать синтетические main-пакеты go run.
var OsExitInMainAnalyzer = &analysis.Analyzer{
	Name: "noosexitinmain",
	Doc:  "Запрещает прямой вызов os.Exit в функции main пакета main (только внутри проекта, под /cmd/)",
	Run:  runNoOsExit,
}

func runNoOsExit(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg == nil || pass.Pkg.Name() != "main" {
		return nil, nil
	}

	// Проверяем, что хотя бы один файл пакета лежит внутри нашего репозитория и под /cmd/
	// (чтобы игнорировать синтетические/внешние main-пакеты при go run)
	inProjectCmd := false
	for _, f := range pass.Files {
		pos := pass.Fset.PositionFor(f.Pos(), false)
		// подстрой под реальный корень; достаточно эвристики по имени модуля и /cmd/
		if strings.Contains(pos.Filename, string(filepath.Separator)+"go-musthave-shortener"+string(filepath.Separator)) &&
			strings.Contains(pos.Filename, string(filepath.Separator)+"cmd"+string(filepath.Separator)) {
			inProjectCmd = true
			break
		}
	}
	if !inProjectCmd {
		return nil, nil
	}

	for _, f := range pass.Files {
		if f.Name == nil || f.Name.Name != "main" {
			continue
		}
		ast.Inspect(f, func(n ast.Node) bool {
			fd, ok := n.(*ast.FuncDecl)
			if !ok || fd.Name == nil || fd.Name.Name != "main" || fd.Body == nil {
				return true
			}
			ast.Inspect(fd.Body, func(n ast.Node) bool {
				if call, ok := n.(*ast.CallExpr); ok {
					if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
						if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "os" && sel.Sel != nil && sel.Sel.Name == "Exit" {
							pass.Reportf(call.Pos(), "запрещён прямой вызов os.Exit в функции main пакета main")
						}
					}
				}
				return true
			})
			return false
		})
	}
	return nil, nil
}
