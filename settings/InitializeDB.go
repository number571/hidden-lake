package settings

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/number571/hiddenlake/utils"
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

	// Hashpasw = sha256(sha256(username+password))
	// Hashname = sha256(pubkey)
	// Hash = sha256(file)
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
CREATE TABLE IF NOT EXISTS File (
	Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	Owner VARCHAR(44),
	Hash VARCHAR(64),
	Name VARCHAR(128),
	Path VARCHAR(64),
	Size INTEGER
);
`)
	if err != nil {
		panic("can't exec database")
	}
}
