package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const tiledir = "Kiev_Adm"
const dbpath = "Kiev_Adm.db"
const scheme = "zxy"

func main() {
	database, _ := sql.Open("sqlite3", dbpath)

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS tiles (id INTEGER PRIMARY KEY, z INT, x INT, y INT, data BLOB)")
	statement.Exec()
	defer database.Close()

	err := filepath.Walk(tiledir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			a := strings.Split(path, "\\")[1:]

			if len(a) > 2 {
				z, _ := strconv.ParseInt(a[0], 0, 64)
				x, _ := strconv.ParseInt(a[1], 0, 64)
				y, _ := strconv.ParseInt(strings.Split(a[2], ".")[0], 0, 64)

				blob, _ := ioutil.ReadFile(path)
				statement, _ = database.Prepare("INSERT INTO tiles (z, x, y, data) VALUES (?, ?, ?, ?)")
				if scheme == "zyx" {
					statement.Exec(z, y, x, blob)
					fmt.Printf("\nINSERT INTO tiles  (z, x, y) VALUES (%d, %d, %d)", z, y, x)
				} else {
					statement.Exec(z, x, y, blob)
					fmt.Printf("\nINSERT INTO tiles  (z, x, y) VALUES (%d, %d, %d)", z, x, y)
				}
			}

			return nil
		})

	if err != nil {
		log.Println(err)
	}
}
