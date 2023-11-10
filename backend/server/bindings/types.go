package bindings

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

const basePath = "./frontend/src/lib/api/"

type TSPkg struct {
	code    []byte
	imports map[string]map[string]struct{} // path -> types
}

func generateTypes(summary Summary) error {

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return fmt.Errorf("failed to read directory '%s': %w", basePath, err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			if err := os.RemoveAll(basePath + entry.Name()); err != nil {
				return fmt.Errorf("failed to remove directory '%s': %w", entry.Name(), err)
			}
		}
	}

	// organise types by package
	packageMap := make(map[string][]Type)
	for _, t := range summary.Types {
		if !t.IsBase || t.TSPackage == "" {
			continue
		}
		packageMap[t.TSPackage] = append(packageMap[t.TSPackage], t)
	}

	/// and write them
	for dir, types := range packageMap {

		if err := os.MkdirAll(basePath+dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory '%s': %w", dir, err)
		}

		file := bytes.NewBuffer(nil)

		// add imports
		importMap := make(map[string]map[string]struct{})
		for _, t := range types {
			for _, field := range t.Fields {
				if !field.IsBase && field.Base != nil {
					field = *field.Base
				}
				if field.TSPackage == "" || field.TSPackage == dir {
					continue
				}
				if _, ok := importMap[field.TSPackage]; !ok {
					importMap[field.TSPackage] = make(map[string]struct{})
				}
				importMap[field.TSPackage][field.TSName] = struct{}{}
			}
		}
		// and import them
		for pkg, types := range importMap {
			var typeNames []string
			for typ := range types {
				typeNames = append(typeNames, typ)
			}
			_, _ = fmt.Fprintf(file, "import { %s } from '../%s'\n", strings.Join(typeNames, ", "), pkg)
		}

		if len(importMap) > 0 {
			_, _ = fmt.Fprint(file, "\n")
		}

		for _, t := range types {
			_, _ = fmt.Fprint(file, t.TSDefinition()+"\n\n")
		}

		if err := os.WriteFile(basePath+dir+"/index.ts", file.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write file '%s': %w", dir+"/index.ts", err)
		}
	}

	return nil
}
