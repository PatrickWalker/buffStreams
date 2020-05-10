package migrations

import (
	"database/sql"

	"fmt"

	helpers "github.com/buffup/api/helpers"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
)

//Migrate will run in a go routine to perform db migrations as needed
func Migrate(conf helpers.Config) error {

	//this opens the SQL connection using the mysql driver
	db, err := sql.Open("mysql", conf.DB.ConnectionString)
	if err != nil {
		fmt.Printf("sql open %v", err)
		return err
	}
	// this returns a driver which can be used with the migrate object
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		fmt.Printf("sql instance %v", err)
		return err
	}
	//Again the migrations path could be configurable but I think that's overkill right now
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/scripts",
		"mysql", driver)
	if err != nil {
		fmt.Println("Unable to create DB Migration Instance")
		return err
	}
	defer m.Close()
	//This runs up on the DB as much as needed. Going by current schema migration value
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		fmt.Println("Unable to run DB Migration")
		return err
	}
	return nil
}
