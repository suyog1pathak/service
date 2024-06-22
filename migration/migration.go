package migration

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"os"
	// this is required to read migrations file
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

// RunMigrations runs the migrations up phase

func RunMigrations(db *gorm.DB) error {
	fmt.Println("Running db migrations..")
	dbSQL, _ := db.DB()
	driver, err := mysql.WithInstance(dbSQL, &mysql.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	dbName := "services"
	migrator, err := migrate.NewWithDatabaseInstance("file://migration/sql", dbName, driver)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Run migrations
	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil

}
