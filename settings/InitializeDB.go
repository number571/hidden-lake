package settings

import (
	"../utils"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitializeDB(dbname string) {
	var err error

	if !utils.FileIsExist(dbname) {
		err = utils.CreateFile(dbname)
		if err != nil {
			panic("can't create database")
		}
	}

	DB, err = sql.Open("sqlite3", dbname)
	if err != nil {
		panic("can't open database")
	}

	// Hashpasw = hash(hash(username+password))
	// Hashname = hash(pubkey)
	_, err = DB.Exec(`
CREATE TABLE IF NOT EXISTS User (
    Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    Hashpasw VARCHAR(44) UNIQUE,
    Key VARCHAR(4096) UNIQUE
);
CREATE TABLE IF NOT EXISTS Client (
    Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    Contributor VARCHAR(44),
    Hashname VARCHAR(44),
    Address VARCHAR(64),
    Public VARCHAR(2048)
);
CREATE TABLE IF NOT EXISTS Chat (
    Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    Hashname VARCHAR(44),
    Companion VARCHAR(44),
    Name VARCHAR(44),
    Text TEXT,
    Time VARCHAR(128)
);
`)
	if err != nil {
		panic("can't exec database")
	}
}
