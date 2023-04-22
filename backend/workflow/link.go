package workflow

import "github.com/google/uuid"

type Link struct {
	From       LinkDirection
	To         LinkDirection
	Annotation string
}

type LinkDirection struct {
	Node      uuid.UUID
	Connector uuid.UUID
}
