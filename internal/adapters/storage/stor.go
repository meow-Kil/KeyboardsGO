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
	// Подключение к PostgreSQL
	// Измените параметры подключения при необходимости
	connStr := "user=postgres password=111 dbname=KeyboardGO sslmode=disable host=localhost port=5432"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Проверка подключения
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to PostgreSQL database")
	
	// Создание таблицы, если не существует
	if err := createTable(db); err != nil {
		return nil, err
	}
	
	log.Println("Table 'keyboards' created or already exists")
	
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

// ... остальные методы MemoryStorage