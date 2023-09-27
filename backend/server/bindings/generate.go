package bindings

import (
	"fmt"
)

func Generate() error {

	types, err := generateClient()
	if err != nil {
		return fmt.Errorf("failed to generate api client: %w", err)
	}

	if err := generateTypes(types); err != nil {
		return fmt.Errorf("failed to generate types: %w", err)
	}

	return nil
}
