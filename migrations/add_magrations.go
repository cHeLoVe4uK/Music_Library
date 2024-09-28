package migrations

import "github.com/pressly/goose/v3"

// Функция, регистрирующая доступные миграции в приложении
func AddMigration() {
	goose.AddMigrationNoTxContext(Up, Down)
}
