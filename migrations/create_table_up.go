package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"os"
)

// Функция, создающая таблицу в БД (накатывающая миграция)
func Up(ctx context.Context, db *sql.DB) error {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ()", os.Getenv("TABLE_NAME"))
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
