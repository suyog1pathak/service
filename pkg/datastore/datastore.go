package datastore

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/suyog1pathak/services/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sync"
)

var db *gorm.DB
var err error
var once sync.Once

// GetMysqlGORM creates db connection, fails if something went wrong.
func getMysqlGORM() (*gorm.DB, error) {
	c := config.GetConfig()
	dsn := BuildDsn(
		c.Db.User,
		c.Db.Password,
		c.Db.Host,
		c.Db.Port,
		c.Db.Name)

	db, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Fatal("Server Shutdown: unable to connect to DB", err)
	}
	return db, nil
}

// BuildDsn builds connection string for mysql
func BuildDsn(user string, password string, host string, port int, name string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s",
		user,
		password,
		host,
		port,
		name,
		"?parseTime=true") //parseTime=true is added to ensure that the time output provided is parsed into time.Time and not []byte/string
}

// GetDBConnection returns already created mysql based orm db connection
// singleton implementation of mysql gorm in GO
func GetDBConnection() (*gorm.DB, error) {
	once.Do(func() {
		db, err = getMysqlGORM()
	})
	return db, err
}
