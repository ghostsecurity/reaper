package bindings

import (
	"fmt"
)

func Generate() error {

	summary := examine()

	if err := generateClient(summary); err != nil {
		return fmt.Errorf("failed to generate api client: %w", err)
	}

	if err := generateTypes(summary); err != nil {
		return fmt.Errorf("failed to generate types: %w", err)
	}

	return nil
}
