package bindings

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

var markerImportsStart = []byte("// %IMPORTS:START%\n")
var markerImportsEnd = []byte("// %IMPORTS:END%\n")
var markerMethodsStart = []byte("// %METHODS:START%\n")
var markerMethodsEnd = []byte("// %METHODS:END%\n")

const clientPath = "./frontend/src/lib/api/Client.ts"

func generateClient(summary Summary) error {

	raw, err := os.ReadFile(clientPath)
	if err != nil {
		return fmt.Errorf("failed to read client file: %w", err)
	}

	beforeImports := bytes.Split(raw, markerImportsStart)[0]
	afterImports := bytes.Split(raw, markerImportsEnd)[1]
	beforeMethods := bytes.Split(afterImports, markerMethodsStart)[0]
	afterMethods := bytes.Split(afterImports, markerMethodsEnd)[1]

	imports, methods, err := generateClientMethods(summary)
	if err != nil {
		return fmt.Errorf("failed to generate client methods: %w", err)
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

	return os.WriteFile(clientPath, generated, 0644)
}

func generateClientMethods(summary Summary) ([]byte, []byte, error) {

	buffer := bytes.NewBuffer(nil)
	importsBuffer := bytes.NewBuffer(nil)

	for _, method := range summary.Methods {
		var argStrs []string
		var callArgs []string
		for i, inType := range method.InTypes {
			argName := fmt.Sprintf("a%d", i)
			argStrs = append(argStrs, fmt.Sprintf("%s: %s", argName, inType.TSProp()))
			callArgs = append(callArgs, argName)
		}
		if len(method.OutTypes) == 0 {

			_, _ = fmt.Fprintf(buffer, `    %s(%s): Promise<void> {
        return new Promise<void>((resolve, reject) => {
            const receive = () => {
                resolve();
            }
            this.callMethod("%[1]s", [%[3]s], receive, reject);
        })
    }

`,
				method.Name,
				strings.Join(argStrs, ", "),
				strings.Join(callArgs, ", "),
			)
		} else {

			var promiseTypes []string
			var outputs string
			var outputList []string
			for i, outType := range method.OutTypes {
				if outType.TSName == "" {
					panic(fmt.Sprintf("method %s has empty output type", method.Name))
				}
				promiseTypes = append(promiseTypes, outType.TSProp())
				outputs += fmt.Sprintf("\n                let output%d: %s = JSON.parse(args[%d])", i, outType.TSProp(), i)
				outputList = append(outputList, fmt.Sprintf("output%d", i))
			}
			promiseType := strings.Join(promiseTypes, ", ")

			_, _ = fmt.Fprintf(buffer, `    %s(%s): Promise<%s> {
        return new Promise<%[3]s>((resolve, reject) => {
            const receive = (args: string[]) => {
				%s
                resolve(%s);
            }
            this.callMethod("%[1]s", [%[6]s], receive, reject);
        })
    }

`,
				method.Name,
				strings.Join(argStrs, ", "),
				promiseType,
				outputs,
				strings.Join(outputList, ", "),
				strings.Join(callArgs, ", "),
			)
		}

	}

	// look for unique types/packages to import
	// map -> type -> package
	importMap := make(map[string]map[string]struct{})
	for _, method := range summary.Methods {
		for _, typ := range append(method.InTypes, method.OutTypes...) {
			if !typ.IsBase && typ.Base != nil {
				typ = *typ.Base
			}
			if typ.TSPackage == "" {
				continue
			}
			if _, ok := importMap[typ.TSPackage]; !ok {
				importMap[typ.TSPackage] = make(map[string]struct{})
			}
			importMap[typ.TSPackage][typ.TSName] = struct{}{}
		}
	}

	// and import them
	for pkg, types := range importMap {
		var typeNames []string
		for typ := range types {
			typeNames = append(typeNames, typ)
		}
		_, _ = fmt.Fprintf(importsBuffer, "import { %s } from './%s'\n", strings.Join(typeNames, ", "), pkg)
	}

	return importsBuffer.Bytes(), buffer.Bytes(), nil
}
