package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"os"
)

// Функция, создающая таблицу в БД (накатывающая миграция)
func Up(ctx context.Context, db *sql.DB) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s("group" text, song text, releaseDate text, text text[], link text)`, os.Getenv("TABLE_NAME"))
	_, err := db.Exec(query)
	return err
}
