package storage

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" 
	"github.com/meow-Kil/KeyboardsGO/internal/core/domain"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	connStr := "user=postgres password=111 dbname=KeyboardGO sslmode=disable host=localhost port=5432"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}


	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to PostgreSQL database")
	

	if err := createTable(db); err != nil {
		return nil, err
	}
	
	log.Println("Table 'keyboards' created or already exists")

	if err := createUsersTable(db); err != nil {
		return nil, err
	}
	
	log.Println("Table 'users' created or already exists")

	if err := createKeycapTypesTable(db); err != nil {
    return nil, err
	}
	
	
	return &PostgresStorage{db: db}, nil
}

func createTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS keyboards (
		id SERIAL PRIMARY KEY,
		keycap_type VARCHAR(100) NOT NULL,
		base_type VARCHAR(100) NOT NULL,
		switch_type VARCHAR(100) NOT NULL,
		color VARCHAR(50) NOT NULL
	)`
	
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error creating table: %v", err)
	}
	return err
}

func createUsersTable(db *sql.DB) error {
    query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        login VARCHAR(100) NOT NULL UNIQUE,
        password VARCHAR(100) NOT NULL,
        is_admin BOOLEAN DEFAULT FALSE
    )`
    _, err := db.Exec(query)
    return err
}

func createKeycapTypesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS keycap_types (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL UNIQUE
	)`
	_, err := db.Exec(query)
	return err
}


func (s *PostgresStorage) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *PostgresStorage) Add(k domain.Keyboard) domain.Keyboard {
	query := `
	INSERT INTO keyboards (keycap_type, base_type, switch_type, color) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id`
	
	var id uint
	err := s.db.QueryRow(query, k.KeycapType, k.BaseType, k.SwitchType, k.Color).Scan(&id)
	if err != nil {
		log.Printf("Error adding keyboard: %v", err)
		return k
	}
	
	k.Id = id
	return k
}

func (s *PostgresStorage) Get() []domain.Keyboard {
	query := `SELECT id, keycap_type, base_type, switch_type, color FROM keyboards`
	rows, err := s.db.Query(query)
	if err != nil {
		log.Printf("Error getting keyboards: %v", err)
		return []domain.Keyboard{}
	}
	defer rows.Close()

	var keyboards []domain.Keyboard
	for rows.Next() {
		var k domain.Keyboard
		err := rows.Scan(&k.Id, &k.KeycapType, &k.BaseType, &k.SwitchType, &k.Color)
		if err != nil {
			log.Printf("Error scanning keyboard: %v", err)
			continue
		}
		keyboards = append(keyboards, k)
	}
	
	return keyboards
}

func (s *PostgresStorage) GetById(id uint) *domain.Keyboard {
	query := `SELECT id, keycap_type, base_type, switch_type, color FROM keyboards WHERE id = $1`
	row := s.db.QueryRow(query, id)
	
	var k domain.Keyboard
	err := row.Scan(&k.Id, &k.KeycapType, &k.BaseType, &k.SwitchType, &k.Color)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Printf("Error getting keyboard by id: %v", err)
		return nil
	}
	
	return &k
}

func (s *PostgresStorage) Remove(id uint) {
	query := `DELETE FROM keyboards WHERE id = $1`
	_, err := s.db.Exec(query, id)
	if err != nil {
		log.Printf("Error removing keyboard: %v", err)
	}
}

func (s *PostgresStorage) Update(id uint, keyboard domain.Keyboard) *domain.Keyboard {
	query := `
	UPDATE keyboards 
	SET keycap_type = $1, base_type = $2, switch_type = $3, color = $4 
	WHERE id = $5 
	RETURNING id, keycap_type, base_type, switch_type, color`
	
	var updatedKeyboard domain.Keyboard
	err := s.db.QueryRow(query, 
		keyboard.KeycapType, 
		keyboard.BaseType, 
		keyboard.SwitchType, 
		keyboard.Color, 
		id,
	).Scan(
		&updatedKeyboard.Id,
		&updatedKeyboard.KeycapType,
		&updatedKeyboard.BaseType,
		&updatedKeyboard.SwitchType,
		&updatedKeyboard.Color,
	)
	
	if err != nil {
		log.Printf("Error updating keyboard: %v", err)
		return nil
	}
	
	return &updatedKeyboard
}

func (s *PostgresStorage) AddUser(login, password string, isAdmin bool) (*domain.User, error) {
    query := `INSERT INTO users (login, password, is_admin) VALUES ($1, $2, $3) RETURNING id`
    var id uint
    err := s.db.QueryRow(query, login, password, isAdmin).Scan(&id)
    if err != nil {
        return nil, err
    }
    return &domain.User{ID: id, Login: login, Password: password, IsAdmin: isAdmin}, nil
}

func (s *PostgresStorage) GetUserByLogin(login string) (*domain.User, error) {
    query := `SELECT id, login, password, is_admin FROM users WHERE login = $1`
    row := s.db.QueryRow(query, login)
    var user domain.User
    err := row.Scan(&user.ID, &user.Login, &user.Password, &user.IsAdmin)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (s *PostgresStorage) AddKeycapType(kt domain.KeycapType) domain.KeycapType {
	query := `INSERT INTO keycap_types (name) VALUES ($1) RETURNING id`
	var id uint
	err := s.db.QueryRow(query, kt.Name).Scan(&id)
	if err != nil {
		log.Printf("Error adding keycap type: %v", err)
		return kt
	}
	kt.ID = id
	return kt
}

func (s *PostgresStorage) GetKeycapTypes() []domain.KeycapType {
	rows, err := s.db.Query(`SELECT id, name FROM keycap_types ORDER BY id`)
	if err != nil {
		log.Printf("Error getting keycap types: %v", err)
		return []domain.KeycapType{}
	}
	defer rows.Close()

	var types []domain.KeycapType
	for rows.Next() {
		var kt domain.KeycapType
		if err := rows.Scan(&kt.ID, &kt.Name); err != nil {
			log.Printf("Error scanning keycap type: %v", err)
			continue
		}
		types = append(types, kt)
	}
	return types
}

func (s *PostgresStorage) GetKeycapTypeByID(id uint) *domain.KeycapType {
	var kt domain.KeycapType
	err := s.db.QueryRow(`SELECT id, name FROM keycap_types WHERE id = $1`, id).
		Scan(&kt.ID, &kt.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Printf("Error getting keycap type by id: %v", err)
		return nil
	}
	return &kt
}

func (s *PostgresStorage) UpdateKeycapType(id uint, kt domain.KeycapType) *domain.KeycapType {
	query := `UPDATE keycap_types SET name = $1 WHERE id = $2 RETURNING id, name`
	var updated domain.KeycapType
	err := s.db.QueryRow(query, kt.Name, id).Scan(&updated.ID, &updated.Name)
	if err != nil {
		log.Printf("Error updating keycap type: %v", err)
		return nil
	}
	return &updated
}

func (s *PostgresStorage) DeleteKeycapType(id uint) {
	_, err := s.db.Exec(`DELETE FROM keycap_types WHERE id = $1`, id)
	if err != nil {
		log.Printf("Error deleting keycap type: %v", err)
	}
}

