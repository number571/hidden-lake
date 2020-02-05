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
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	Hashpasw VARCHAR(44) UNIQUE,
	PrivateKey VARCHAR(4096) UNIQUE
);
CREATE TABLE IF NOT EXISTS Client (
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	IdUser INTEGER,
	Hashname VARCHAR(44),
	Address VARCHAR(64),
	PublicKey VARCHAR(2048),
	FOREIGN KEY (IdUser)  REFERENCES User (Id)
);
CREATE TABLE IF NOT EXISTS Chat (
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	IdUser INTEGER,
	Companion VARCHAR(44),
	Name VARCHAR(44),
	Message TEXT,
	LastTime VARCHAR(128),
	FOREIGN KEY (IdUser)  REFERENCES User (Id)
);
CREATE TABLE IF NOT EXISTS File (
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	IdUser INTEGER,
	Hash VARCHAR(44),
	Name VARCHAR(128),
	Path VARCHAR(44),
	Size INTEGER,
	FOREIGN KEY (IdUser)  REFERENCES User (Id)
);
`)
	if err != nil {
		panic("can't exec database")
	}
}
