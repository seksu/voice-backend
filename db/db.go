package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {

	dbUser := "doadmin"
	dbPass := "d2ui8cjl2u640wiw"
	dbName := "postgres"
	dbURL := "db-postgresql-sgp1-54259-do-user-8097391-0.b.db.ondigitalocean.com"
	dbPort := 25060

	// var connectDB string = "" + dbUser + ":" + dbPass + "@tcp(" + dbURL + ")/" + dbName + ""

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", dbURL, dbPort, dbUser, dbPass, dbName)

	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("DB Connect Init::Err::", err)
	}

	conn.SetConnMaxIdleTime(0)

	return conn
}
