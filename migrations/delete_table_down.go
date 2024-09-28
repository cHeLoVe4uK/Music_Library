package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"os"
)

// Функция, удаляющая таблицу из БД (откатывающая миграция)
func Down(ctx context.Context, db *sql.DB) error {
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", os.Getenv("TABLE_NAME"))

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
