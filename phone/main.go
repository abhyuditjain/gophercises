package main

import (
	"database/sql"
	"fmt"
	"github.com/abhyuditjain/gophercices/phone/db"
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
	must(db.Reset("postgres", psqlInfo, dbname))

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	must(db.Migrate("postgres", psqlInfo))

	phoneDb, err := db.Open("postgres", psqlInfo)
	must(err)
	defer phoneDb.Close()

	if err := phoneDb.Seed(); err != nil {
		panic(err)
	}

	phones, err := phoneDb.AllPhones()
	must(err)
	for _, p := range phones {
		fmt.Printf("Working on... %+v\n", p)
		number := normalize(p.Number)
		if number != p.Number {
			fmt.Println("Updating or removing...", number)
			existing, err := phoneDb.FindPhone(number)
			must(err)
			if existing != nil {
				// delete this number
				must(phoneDb.DeletePhone(p.ID))
			} else {
				// update this number
				p.Number = number
				must(phoneDb.UpdatePhone(&p))
			}
		} else {
			fmt.Println("No changes required")
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func getPhone(db *sql.DB, id int) (string, error) {
	var number string

	err := db.QueryRow("SELECT * FROM phone_numbers WHERE id=$1", id).Scan(&id, &number)
	if err != nil {
		return "", err
	}
	return number, nil
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
