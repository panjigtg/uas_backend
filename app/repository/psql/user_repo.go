package psql

import (
	"uas/app/repository"
	"database/sql"
	"uas/app/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) repository.UserRepository {
	return &UserRepository{DB: db}
}


func (r *UserRepository) GetAll() ([]models.UserWithRole, error) {
	query := `
	SELECT 
		u.id, u.username, u.email, u.full_name,
		r.name AS role, u.is_active
	FROM users u
	JOIN roles r ON r.id = u.role_id
	ORDER BY u.created_at ASC;
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserWithRole

	for rows.Next() {
		var u models.UserWithRole
		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.FullName,
			&u.RoleName,
			&u.IsActive,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepository) GetByID(id string) (*models.UserWithRole, error) {
	query := `
	SELECT 
		u.id, u.username, u.email, u.full_name,
		r.name AS role, u.is_active
	FROM users u
	JOIN roles r ON r.id = u.role_id
	WHERE u.id = $1;
	`

	var u models.UserWithRole

	err := r.DB.QueryRow(query, id).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.FullName,
		&u.RoleName,
		&u.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) Update(id string, data models.UserUpdateRequest) error {
	query := `
	UPDATE users 
	SET username=$1, email=$2, full_name=$3, role_id=$4, updated_at=NOW()
	WHERE id=$5;
	`

	_, err := r.DB.Exec(query,
		data.Username,
		data.Email,
		data.FullName,
		data.RoleID,
		id,
	)

	return err
}

func (r *UserRepository) Create(user *models.Users) error {
	query := `
	INSERT INTO users (username, email, password_hash, full_name, role_id)
	VALUES ($1, $2, $3, $4, $5);
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


func (r *UserRepository) Delete(id string) error {
	_, err := r.DB.Exec(`DELETE FROM users WHERE id=$1`, id)
	return err
}

func (r *UserRepository) UpdateRole(id string, roleID string) error {
	_, err := r.DB.Exec(`UPDATE users SET role_id=$1 WHERE id=$2`, roleID, id)
	return err
}

func (r *UserRepository) GetIDByIndex(idx int) (string, error) {
    query := `
        SELECT id
        FROM users
        ORDER BY created_at ASC
        LIMIT 1 OFFSET $1
    `

    var id string
    err := r.DB.QueryRow(query, idx).Scan(&id)
    if err != nil {
        return "", err
    }

    return id, nil
}

