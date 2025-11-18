package datatypex

// uuid wraps uuid v7(time-ordered) to improve pg insertion and indexing performance
// ref: https://maciejwalkowiak.com/blog/postgres-uuid-primary-key/

import (
	"github.com/google/uuid"
)

func NewUUID() (UUID, error) {
	v, err := uuid.NewV7()
	return UUID{v}, err
}

type UUID struct {
	uuid.UUID
}

func (UUID) DBType(driver string) string {
	switch driver {
	case "postgres":
		return "uuid"
	case "mysql", "sqlite", "sqlite3":
		return "varchar(36)" // NOTE: sqlite does not strictly enforce the length of data.
	default:
		return "text"
	}
}
