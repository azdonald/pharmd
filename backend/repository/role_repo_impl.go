package repository

import "database/sql"

type RoleRepoImpl struct {
	db *sql.DB
}

func NewRoleRepositoryImpl(db *sql.DB) RoleRepository {
	return &RoleRepoImpl{db: db}
}
