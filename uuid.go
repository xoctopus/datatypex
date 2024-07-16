package datatypex

// uuid wraps uuid v7(time-ordered) to improve pg insertion and indexing performance
// ref: https://maciejwalkowiak.com/blog/postgres-uuid-primary-key/

import "github.com/google/uuid"

func NewUUID() (UUID, error) {
	id, err := uuid.NewV7()
	return UUID{id}, err
}

type UUID struct {
	uuid.UUID
}

func (UUID) DataType(driver string) string {
	if driver == "postgres" {
		return "uuid"
	}
	return "text"
}
