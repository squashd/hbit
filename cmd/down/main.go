package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/feat"
	"github.com/SQUASHD/hbit/rpg"
	"github.com/SQUASHD/hbit/task"
	"github.com/SQUASHD/hbit/user"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	var rpgDownErr, userDownErr, featDownErr, taskDownErr error
	// RPG DB DOWN
	connectionStr := os.Getenv("RPG_DB_URL")
	db, err := hbit.NewDatabase(hbit.NewDbParams{
		ConnectionStr: connectionStr,
		Driver:        hbit.DbDriverLibsql,
	})
	if err != nil {
		log.Fatalf("cannot connect to rpg database: %s", err)
	}
	rpgDownErr = hbit.DBMigratedDownTo(db, 0, hbit.MigrationData{
		FS:      rpg.Migrations,
		Dialect: "sqlite",
		Dir:     "schemas",
	})

	// TASK DB DOWN
	connectionStr = os.Getenv("TASK_DB_URL")
	db, err = hbit.NewDatabase(hbit.NewDbParams{
		ConnectionStr: connectionStr,
		Driver:        hbit.DbDriverLibsql,
	})
	if err != nil {
		log.Fatalf("cannot connect to task database: %s", err)
	}
	taskDownErr = hbit.DBMigratedDownTo(db, 0, hbit.MigrationData{
		FS:      task.Migrations,
		Dialect: "sqlite",
		Dir:     "schemas",
	})
	// USER DB DOWN
	connectionStr = os.Getenv("USER_DB_URL")
	db, err = hbit.NewDatabase(hbit.NewDbParams{
		ConnectionStr: connectionStr,
		Driver:        hbit.DbDriverLibsql,
	})
	if err != nil {
		log.Fatalf("cannot connect to user database: %s", err)
	}
	userDownErr = hbit.DBMigratedDownTo(db, 0, hbit.MigrationData{
		FS:      user.Migrations,
		Dialect: "sqlite",
		Dir:     "schemas",
	})

	// FEAT DB DOWN
	connectionStr = os.Getenv("ACH_DB_URL")
	db, err = hbit.NewDatabase(hbit.NewDbParams{
		ConnectionStr: connectionStr,
		Driver:        hbit.DbDriverLibsql,
	})
	if err != nil {
		log.Fatalf("cannot connect to feat database: %s", err)
	}
	featDownErr = hbit.DBMigratedDownTo(db, 0, hbit.MigrationData{
		FS:      feat.Migrations,
		Dialect: "sqlite",
		Dir:     "schemas",
	})

	var errs []error
	for _, err := range []error{rpgDownErr, userDownErr, featDownErr, taskDownErr} {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		fmt.Println(errors.Join(errs...))
	} else {
		fmt.Println("All databases down")
	}

}
