package repository

import (
	"database/sql"
	"uas/app/models"
)

type AuthRepository interface {
	Register(user *models.Users) error
	GetUserByEmail(email string) (*models.UserWithRole, error)
	GetUserByID(userID string) (*models.UserWithRole, error)
	GetPermissionsByUserID(userID string) ([]string, error)
	GetRoleIDByName(name string) (string, error) 
}

type authRepository struct {
	DB *sql.DB
}

func NewAuthRepo(db *sql.DB) AuthRepository {
	return &authRepository{DB: db}
}

func (r *authRepository) Register(user *models.Users) error {
	query := `
	INSERT INTO users (username, email, password_hash, full_name, role_id)
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.DB.Exec(query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.RoleID,
	)
	
	return err
}

func (r *authRepository) GetUserByEmail(email string) (*models.UserWithRole, error) {
	query := `
	SELECT 
		u.id,
		u.username,
		u.email,
		u.password_hash,
		u.full_name,
		r.name AS role,
		u.is_active
	FROM users u
	JOIN roles r ON u.role_id = r.id
	WHERE u.email = $1
	`

	u := &models.UserWithRole{}
	err := r.DB.QueryRow(query, email).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.PasswordHash,
		&u.FullName,
		&u.RoleName,
		&u.IsActive,
	)

	return u, err
}


func (r *authRepository) GetUserByID(userID string) (*models.UserWithRole, error) {
	query := `
	SELECT 
		u.id,
		u.username,
		u.email,
		u.password_hash,
		u.full_name,
		r.name AS role,
		u.is_active
	FROM users u
	JOIN roles r ON u.role_id = r.id
	WHERE u.id = $1
	`

	u := &models.UserWithRole{}
	err := r.DB.QueryRow(query, userID).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.PasswordHash,
		&u.FullName,
		&u.RoleName,
		&u.IsActive,
	)

	return u, err
}


func (r *authRepository) GetPermissionsByUserID(userID string) ([]string, error) {
	query := `
	SELECT p.name
	FROM permissions p
	JOIN role_permissions rp ON p.id = rp.permission_id
	JOIN users u ON rp.role_id = u.role_id
	WHERE u.id = $1
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var p string
		rows.Scan(&p)
		perms = append(perms, p)
	}

	return perms, nil
}

func (r *authRepository) GetRoleIDByName(name string) (string, error) {
	var roleID string
	query := `SELECT id FROM roles WHERE name = $1`
	err := r.DB.QueryRow(query, name).Scan(&roleID)
	return roleID, err
}
