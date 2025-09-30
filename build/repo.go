package build

import (
	"database/sql"

	repo "otusgoods/internal/repo"
)

func (b *Builder) NewRepo(db *sql.DB) *repo.Queries {
	repo := repo.New(db)

	return repo
}
