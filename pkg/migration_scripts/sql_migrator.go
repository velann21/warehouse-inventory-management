package migration_scripts

import (
	"database/sql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	"github.com/golang-migrate/migrate/source/file"
	"github.com/velann21/warehouse-inventory-management/pkg/helpers"
)

func MigrateDb(stepCount uint) error {
	helper := helpers.Helper{}
	fsrc, err := (&file.File{}).Open("file://")
	if err != nil {
		return err
	}

	db, err := sql.Open("mysql", helper.ReadEnv(helpers.MYSQLCONNECTIONSTRING)+"multiStatements=true")
	if err != nil {
		return err
	}
	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	m, err := migrate.NewWithInstance("", fsrc, helpers.DATABASENAME, driver)
	defer func(m *migrate.Migrate) {
		_, _ = m.Close()
	}(m)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil {
		return err
	}
	return nil
}
