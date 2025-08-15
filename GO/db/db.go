package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

const DateFormat = "20060102"

type Store struct {
	DB *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{DB: db}
}

func ConnectDB(dbFile string) (*sql.DB, error) {
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	if install {
		fmt.Println("Файл БД не найден, создаём и применяем миграции...")

		if err := goose.SetDialect("sqlite3"); err != nil {
			return nil, fmt.Errorf("не удалось установить диалект goose: %w", err)
		}

		migrationsDir := "./GO/db/migrations"

		if err := goose.Up(db, migrationsDir); err != nil {
			return nil, fmt.Errorf("ошибка при применении миграций: %w", err)
		}
	}

	fmt.Println("Соединение с БД успешно")
	return db, nil
}
