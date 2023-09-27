package bindings

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
)

const basePath = "./frontend/src/lib/api/"

type TSPkg struct {
	code    []byte
	imports map[string]map[string]struct{} // path -> types
}

func generateTypes(types []PackageType) error {

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

	packageMap := make(map[string]TSPkg)
	cache := make(map[string]struct{})

	for _, typ := range types {
		if typ.PackagePath == "" {
			continue
		}
		data, err := generateType(typ, cache)
		if err != nil {
			return fmt.Errorf("failed to generate type '%s': %w", typ.Go, err)
		}
		fmt.Println(typ.Simplified)
		pkg := packageMap[typ.PackageName]
		pkg.code = append(pkg.code, data...)
		if pkg.imports == nil {
			pkg.imports = make(map[string]map[string]struct{})
		}
		for _, dep := range typ.Deps {
			if dep.PackagePath == "" {
				continue
			}
			existing := pkg.imports[dep.PackageName]
			if existing == nil {
				existing = make(map[string]struct{})
			}
			existing[dep.Simplified] = struct{}{}
			pkg.imports[dep.PackageName] = existing
		}

		packageMap[typ.PackageName] = pkg
	}

	for dir, pkg := range packageMap {

		if len(pkg.code) == 0 {
			continue
		}

		if err := os.MkdirAll(basePath+dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory '%s': %w", dir, err)
		}

		file := bytes.NewBuffer(nil)

		for path, names := range pkg.imports {
			if path == dir {
				continue
			}
			var typeList []string
			for name := range names {
				typeList = append(typeList, name)
			}
			_, _ = fmt.Fprintf(file, "import {%s} from \"../%s\";\n", strings.Join(typeList, ", "), path)
		}

		if _, err := file.Write(pkg.code); err != nil {
			return fmt.Errorf("failed to write file '%s': %w", dir+"/index.ts", err)
		}

		if err := os.WriteFile(basePath+dir+"/index.ts", file.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write file '%s': %w", dir+"/index.ts", err)
		}
	}

	return nil
}

func generateType(typ PackageType, cache map[string]struct{}) ([]byte, error) {
	if _, ok := cache[typ.Go.PkgPath()+":"+typ.Go.String()]; ok {
		return nil, nil
	}
	cache[typ.Go.PkgPath()+":"+typ.Go.String()] = struct{}{}
	switch typ.Go.Kind() {
	case reflect.Struct:
		return generateStruct(typ, cache)
	default:
		return nil, nil
		//return nil, nil, fmt.Errorf("unsupported kind '%s'", typ.Go.Kind())
	}
}

func generateStruct(typ PackageType, cache map[string]struct{}) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	_, _ = fmt.Fprintf(buffer, "\nexport interface %s {\n", typ.Simplified)
	for f := 0; f < typ.Go.NumField(); f++ {
		field := typ.Go.Field(f)
		if !field.IsExported() {
			continue
		}
		tsType, err := (&converter{}).convertType(field.Type)
		if err != nil {
			return nil, fmt.Errorf("failed to convert type '%s': %w", field.Type.Name(), err)
		}
		_, _ = fmt.Fprintf(buffer, "%s: %s\n", field.Name, tsType.Alias)
	}
	_, _ = fmt.Fprintf(buffer, "}\n\n")
	return buffer.Bytes(), nil
}
