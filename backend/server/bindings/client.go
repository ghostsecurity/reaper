package bindings

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/ghostsecurity/reaper/backend/server/api"
)

var markerImportsStart = []byte("//%IMPORTS:START%\n")
var markerImportsEnd = []byte("//%IMPORTS:END%\n")
var markerMethodsStart = []byte("//%METHODS:START%\n")
var markerMethodsEnd = []byte("//%METHODS:END%\n")

const clientPath = "./frontend/src/lib/api/Client.ts"

func generateClient() ([]PackageType, error) {

	raw, err := os.ReadFile(clientPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read client file: %w", err)
	}

	beforeImports := bytes.Split(raw, markerImportsStart)[0]
	afterImports := bytes.Split(raw, markerImportsEnd)[1]
	beforeMethods := bytes.Split(afterImports, markerMethodsStart)[0]
	afterMethods := bytes.Split(afterImports, markerMethodsEnd)[1]

	imports, methods, types, err := generateClientMethods()
	if err != nil {
		return nil, fmt.Errorf("failed to generate client methods: %w", err)
	}

	generated := bytes.Join([][]byte{
		beforeImports,
		markerImportsStart,
		imports,
		markerImportsEnd,
		beforeMethods,
		markerMethodsStart,
		methods,
		markerMethodsEnd,
		afterMethods,
	}, []byte{})
	return types, os.WriteFile(clientPath, generated, 0644)
}

func generateClientMethods() ([]byte, []byte, []PackageType, error) {

	buffer := bytes.NewBuffer(nil)
	importsBuffer := bytes.NewBuffer(nil)

	var types []PackageType
	importMap := make(map[string]PackageType)

	conv := &converter{}

	apiVal := reflect.ValueOf(&api.API{})
	apiType := apiVal.Type()
	for i := 0; i < apiType.NumMethod(); i++ {
		method := apiType.Method(i)
		var argStrs []string
		var callArgs []string
		for j := 1; j < method.Type.NumIn(); j++ {
			arg := method.Type.In(j)
			baseType, err := conv.convertType(arg)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("failed to convert argument type: %w", err)
			}
			importMap[baseType.Name] = *baseType
			argName := 'a' + j - 1
			argStrs = append(argStrs, fmt.Sprintf("%c: %s", argName, baseType.Alias))
			callArgs = append(callArgs, fmt.Sprintf("%c", argName))
		}
		var numOut int
		for i := 0; i < method.Type.NumOut(); i++ {
			if method.Type.Out(i).String() == "error" {
				continue
			}
			numOut++
		}
		if numOut == 0 {
			_, _ = fmt.Fprintf(buffer, `    %s(%s): void {
        this.callMethod("%[1]s", [%[3]s], (args: string[]) => void {}, (reason?: any) => void {});
    }

`,
				method.Name,
				strings.Join(argStrs, ", "),
				strings.Join(callArgs, ", "),
			)
		} else if numOut == 1 {
			baseType, err := conv.convertType(method.Type.Out(0))
			if err != nil {
				return nil, nil, nil, fmt.Errorf("failed to convert return type: %w", err)
			}
			importMap[baseType.Name] = *baseType

			_, _ = fmt.Fprintf(buffer, `    %s(%s): Promise<%s> {
        return new Promise<%[3]s>((resolve, reject) => {
            const receive = (args: string[]) => {
                let output: %[3]s = JSON.parse(args[0]);
                resolve(output);
            }
            this.callMethod("%[1]s", [%[4]s], receive, reject);
        })
    }

`,
				method.Name,
				strings.Join(argStrs, ", "),
				baseType.Alias,
				strings.Join(callArgs, ", "),
			)
		} else {
			return nil, nil, nil, fmt.Errorf("method %s has too many return values", method.Name)
		}

	}

	// flatten imports
	for _, pType := range importMap {
		if pType.PackagePath == "" {
			continue
		}
		_, _ = fmt.Fprintf(importsBuffer, `import {%s} from "./%s";
`, pType.Simplified, pType.PackageName)
	}

	for _, pType := range flattenDeps(importMap) {
		if pType.PackagePath == "" {
			continue
		}
		types = append(types, pType)
	}

	return importsBuffer.Bytes(), buffer.Bytes(), types, nil
}

func flattenDeps(importMap map[string]PackageType) map[string]PackageType {
	flat := make(map[string]PackageType)
	for _, pType := range importMap {
		flat[pType.Name] = pType
		if len(pType.Deps) == 0 {
			continue
		}
		for _, e := range pType.Deps {
			flat[e.Name] = e
			if len(e.Deps) == 0 {
				continue
			}
			sub := flattenDeps(map[string]PackageType{
				e.Name: e,
			})
			for _, s := range sub {
				flat[s.Name] = s
			}
		}
	}
	return flat
}

type PackageType struct {
	Name        string // e.g. []int
	Alias       string // e.g. number[]
	Simplified  string // e.g. number
	PackageName string
	PackagePath string
	Go          reflect.Type
	Deps        []PackageType
}

type converter struct {
	cache    map[reflect.Type]PackageType
	depCache map[reflect.Type]struct{}
}

func (c *converter) convertType(t reflect.Type) (*PackageType, error) {

	if c.cache == nil {
		c.cache = make(map[reflect.Type]PackageType)
	}
	if pt, ok := c.cache[t]; ok {
		return &pt, nil
	}

	var suffix string

	if t.Kind() == reflect.Slice {
		suffix = "[]"
	} else if t.Kind() == reflect.Ptr {
		suffix = "|null"
	}

	if t.Kind() == reflect.Map {
		if t.Key().String() != "string" {
			return nil, fmt.Errorf("only string keys are supported for maps")
		}
		parts := strings.Split(t.String(), ".")
		typeName := parts[len(parts)-1]
		sub, err := c.convertType(t.Elem())
		if err != nil {
			return nil, fmt.Errorf("failed to convert map type: %w", err)
		}
		pt := PackageType{
			Name:       typeName,
			Alias:      fmt.Sprintf("{[key: string]: %s}", sub.Alias),
			Simplified: sub.Simplified,
			Go:         t,
			Deps:       []PackageType{*sub},
		}
		return &pt, nil
	}

	t = simplifyType(t)
	parts := strings.Split(t.String(), ".")
	typeName := parts[len(parts)-1]

	pt := PackageType{
		Name: typeName,
		Go:   t,
	}
	if t.PkgPath() != "" {
		pt.PackageName = parts[len(parts)-2]
		pt.PackagePath = t.PkgPath()
		pt.Alias = typeName + suffix
		pt.Simplified = typeName
		pt.Deps = c.findDeps(t)
		return &pt, nil
	}
	switch typeName {
	case "interface {}":
		pt.Simplified = "any"
	case "string":
		pt.Simplified = "string"
	case "int", "int64", "int32", "int16", "int8", "uint", "uint64", "uint32", "uint16", "uint8", "float32", "float64":
		pt.Simplified = "number"
	case "bool":
		pt.Simplified = "boolean"
	case "Regexp":
		pt.Simplified = "string"
	default:
		return nil, fmt.Errorf("unsupported Go type '%s'", t)
	}
	pt.Alias = pt.Simplified + suffix

	c.cache[t] = pt
	return &pt, nil
}

func simplifyType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		return simplifyType(t)
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
		return simplifyType(t)
	}
	return t
}

func (c *converter) findDeps(t reflect.Type) []PackageType {

	if c.depCache == nil {
		c.depCache = make(map[reflect.Type]struct{})
	}
	if _, ok := c.depCache[t]; ok {
		return nil
	}
	c.depCache[t] = struct{}{}

	t = simplifyType(t)
	if t.Kind() != reflect.Struct {
		return nil
	}
	var deps []PackageType
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		ft, err := c.convertType(field.Type)
		if err != nil {
			continue
		}
		deps = append(deps, *ft)
		deps = append(deps, c.findDeps(field.Type)...)
	}
	return deps
}
