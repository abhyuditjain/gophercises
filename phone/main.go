package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
)

const (
	host     = "0.0.0.0"
	port     = 5432
	user     = "abhyuditjain"
	password = "25111992"
	dbname   = "gophercises"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		host,
		port,
		user,
		password,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = resetDB(db, dbname)
	if err != nil {
		panic(err)
	}
	_ = db.Close()

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()

	must(db.Ping())
	must(createPhoneNumbersTable(db))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
}

func createPhoneNumbersTable(db *sql.DB) error {
	statement := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS phone_numbers (
			id SERIAL,
			value VARCHAR(255)
		)
	`)

	_, err := db.Exec(statement)
	return err
}

func normalize(phone string) string {
	var normalized strings.Builder
	for _, b := range phone {
		if b >= '0' && b <= '9' {
			normalized.WriteRune(b)
		}
	}
	return normalized.String()
}

//func normalize(phone string) string {
//	re := regexp.MustCompile("\\D")
//	return re.ReplaceAllString(phone, "")
//}
