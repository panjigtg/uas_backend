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


func (r *UserRepository) Update(tx *sql.Tx, userID string, req models.UserUpdateRequest) error {
	query := `
	UPDATE users 
	SET username=$1, email=$2, full_name=$3, updated_at=NOW()
	WHERE id=$4;
	`

	_, err := tx.Exec(query,
		req.Username,
		req.Email,
		req.FullName,
		userID,
	)

	return err
}


func (r *UserRepository) Create(tx *sql.Tx, user *models.Users) (string, error) {
	query := `
	INSERT INTO users (username, email, password_hash, full_name, role_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id;
	`

	var newID string
	err := tx.QueryRow(query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.RoleID,
	).Scan(&newID)

	return newID, err
}


func (r *UserRepository) Delete(tx *sql.Tx, id string) error {
	_, err := tx.Exec(`DELETE FROM users WHERE id=$1`, id)
	return err
}

func (r *UserRepository) UpdateRole(tx *sql.Tx, userID string, roleID string) error {
	_, err := tx.Exec(`
		UPDATE users SET role_id=$1, updated_at=NOW() WHERE id=$2
	`, roleID, userID)
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

