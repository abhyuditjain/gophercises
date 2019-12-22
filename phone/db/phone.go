package db

import "database/sql"

type DB struct {
	db *sql.DB
}

func Open(driver, dataSource string) (*DB, error) {
	if db, err := sql.Open(driver, dataSource); err != nil {
		return nil, err
	} else {
		return &DB{db}, err
	}
}

func (db *DB) Close() error {
	return db.db.Close()
}

// Phone represents the phone_numbers table
type Phone struct {
	ID     int
	Number string
}

func (db *DB) Seed() error {
	data := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}

	for _, number := range data {
		if _, err := insertPhone(db.db, number); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) AllPhones() ([]Phone, error) {
	return allPhones(db.db)
}

func (db *DB) FindPhone(number string) (*Phone, error) {
	return findPhone(db.db, number)
}

func (db *DB) UpdatePhone(p *Phone) error {
	return updatePhone(db.db, p)
}

func (db *DB) DeletePhone(id int) error {
	return deletePhone(db.db, id)
}

func Reset(driver, dataSource, dbName string) error {
	db, err := sql.Open(driver, dataSource)
	if err != nil {
		panic(err)
	}
	err = resetDB(db, dbName)
	if err != nil {
		return err
	}
	return db.Close()
}

func Migrate(driver, dataSource string) error {
	db, err := sql.Open(driver, dataSource)
	if err != nil {
		panic(err)
	}
	if err := createPhoneNumbersTable(db); err != nil {
		return err
	}
	return db.Close()
}

func allPhones(db *sql.DB) ([]Phone, error) {
	rows, err := db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []Phone

	for rows.Next() {
		var p Phone
		if err := rows.Scan(&p.ID, &p.Number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func insertPhone(db *sql.DB, phone string) (int, error) {
	statement := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`
	var id int
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func findPhone(db *sql.DB, number string) (*Phone, error) {
	var p Phone

	err := db.QueryRow("SELECT * FROM phone_numbers WHERE value=$1", number).Scan(&p.ID, &p.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func updatePhone(db *sql.DB, p *Phone) error {
	statement := `UPDATE phone_numbers SET value=$2 WHERE id=$1`
	_, err := db.Exec(statement, p.ID, p.Number)
	return err
}

func deletePhone(db *sql.DB, id int) error {
	statement := `DELETE FROM phone_numbers WHERE id=$1`
	_, err := db.Exec(statement, id)
	return err
}

func createPhoneNumbersTable(db *sql.DB) error {
	statement := `
		CREATE TABLE IF NOT EXISTS phone_numbers (
			id SERIAL,
			value VARCHAR(255)
		)`

	_, err := db.Exec(statement)
	return err
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
