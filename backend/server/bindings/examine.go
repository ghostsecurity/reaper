package bindings

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ghostsecurity/reaper/backend/server/api"
)

type Method struct {
	Name     string
	InTypes  []Type
	OutTypes []Type
}

type Type struct {
	TSName    string
	TSPackage string
	Go        reflect.Type
	IsBase    bool
	Fields    []Type
	FieldName string
	Base      *Type
}

func (t Type) TSProp() string {
	return goToTS(t.Go, false, true)
}

func (t Type) TSDefinition() string {
	for t.Go.Kind() == reflect.Ptr {
		t.Go = t.Go.Elem()
	}
	switch t.Go.Kind() {
	case reflect.Struct:
		output := fmt.Sprintf(`export interface %s {`, t.TSName)
		for _, field := range t.Fields {
			output += fmt.Sprintf("\n  %s: %s;", field.FieldName, field.TSProp())
		}
		output += "\n}"
		return output
	default:
		return fmt.Sprintf("export type %s = %s;", t.TSName, goToTS(t.Go, false, false))
	}
}

func goToTS(t reflect.Type, ignorePtrs, allowAlias bool) string {
	switch t.Kind() {
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return "Uint8Array"
		}
		return goToTS(t.Elem(), true, allowAlias) + "[]"
	case reflect.Ptr:
		if ignorePtrs {
			return goToTS(t.Elem(), ignorePtrs, allowAlias)
		}
		return goToTS(t.Elem(), true, allowAlias) + "|null"
	case reflect.Map:
		if allowAlias && t.PkgPath() != "" {
			return t.Name()
		}
		return fmt.Sprintf("{[key: string]: %s}", goToTS(t.Elem(), true, allowAlias))
	case reflect.String:
		if allowAlias && t.PkgPath() != "" {
			return t.Name()
		}
		return "string"
	case reflect.Uint8, reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Int8, reflect.Float32, reflect.Float64:
		if allowAlias && t.PkgPath() != "" {
			return t.Name()
		}
		return "number"
	case reflect.Bool:
		if allowAlias && t.PkgPath() != "" {
			return t.Name()
		}
		return "boolean"
	case reflect.Interface:
		return "any"
	case reflect.Struct:
		return t.Name()
	default:
		return "any"
	}
}

type Summary struct {
	Methods []Method
	Types   []Type
}

func examine() Summary {
	var s Summary
	s.scan()
	return s
}

func (s *Summary) scan() {
	apiVal := reflect.ValueOf(&api.API{})
	apiType := apiVal.Type()

	for i := 0; i < apiType.NumMethod(); i++ {
		mt := apiType.Method(i)
		method := Method{
			Name: mt.Name,
		}

		for j := 1; j < mt.Type.NumIn(); j++ {
			arg := mt.Type.In(j)
			method.InTypes = append(method.InTypes, s.examineType(arg, nil))
		}
		for j := 0; j < mt.Type.NumOut(); j++ {
			out := mt.Type.Out(j)
			if out.String() == "error" {
				continue
			}
			method.OutTypes = append(method.OutTypes, s.examineType(out, nil))
		}

		s.Methods = append(s.Methods, method)
	}

}

func (s *Summary) examineType(t reflect.Type, chain []reflect.Type) Type {

	var result Type
	result.Go = t

	parts := strings.Split(t.PkgPath(), "/")
	result.TSPackage = parts[len(parts)-1]
	if n := t.Name(); n != "" && n != goToTS(t, false, true) {
		result.TSName = n
	} else {
		result.TSName = goToTS(t, false, true)
	}

	if result.TSName == "" {
		panic(fmt.Sprintf("empty type name for %s", t.String()))
	}

	// recurse
	base := t
	for base.Kind() == reflect.Slice || base.Kind() == reflect.Map || base.Kind() == reflect.Ptr {
		base = base.Elem()
	}

	result.IsBase = base == t

	for _, c := range chain {
		if c.String() == base.String() {
			return Type{
				Go: t,
			}
		}
	}

	if !result.IsBase {
		b := s.examineType(base, append(chain, t))
		result.Base = &b
	}

	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			j := field.Tag.Get("json")
			if j == "-" || j == "" {
				continue
			}
			j = strings.Split(j, ",")[0]
			sub := s.examineType(field.Type, append(chain, t))
			sub.FieldName = j
			result.Fields = append(result.Fields, sub)
		}
	}

	var found bool
	for _, existing := range s.Types {
		if existing.TSName == result.TSName && existing.TSPackage == result.TSPackage {
			found = true
			break
		}
	}
	if !found {
		s.Types = append(s.Types, result)
	}

	return result
}
