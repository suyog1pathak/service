package main

import (
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/suyog1pathak/services/migration"
	"github.com/suyog1pathak/services/pkg/datastore"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func main() {

	db, _ = datastore.GetDBConnection()
	if err != nil {
		return
	}
	err := migration.RunMigrations(db)
	if err != nil {
		fmt.Println(err)
		return
	}
}
