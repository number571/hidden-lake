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

	// Username = sha256(username)
	// Hashpasw = sha256(sha256(password+salt))
	// Hashname = sha256(pubkey)
	// Hash = sha256(file)
	_, err = DB.Exec(`
CREATE TABLE IF NOT EXISTS User (
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	Username VARCHAR(44) UNIQUE,
	Salt VARCHAR(32),
	Hashpasw VARCHAR(44) UNIQUE,
	PrivateKey VARCHAR(4096) UNIQUE
);
CREATE TABLE IF NOT EXISTS Email (
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	IdUser INTEGER,
	Incoming BOOLEAN,
	Temporary BOOLEAN,
	LastTime VARCHAR(128),
	SenderPub VARCHAR(1024),
	Receiver VARCHAR(44),
	Session VARCHAR(128) NULL,
	Title VARCHAR(128),
	Message VARCHAR(2048),
	Salt VARCHAR(32),
	Hash VARCHAR(44),
	Sign VARCHAR(512),
	Nonce INTEGER,
	FOREIGN KEY (IdUser) REFERENCES User (Id) ON UPDATE CASCADE,
	FOREIGN KEY (IdUser) REFERENCES User (Id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS Client (
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	IdUser INTEGER,
	Hashname VARCHAR(44),
	Address VARCHAR(64),
	PublicKey VARCHAR(1024),
	ThrowClient VARCHAR(1024),
	Certificate VARCHAR(3072),
	FOREIGN KEY (IdUser) REFERENCES User (Id) ON UPDATE CASCADE,
	FOREIGN KEY (IdUser) REFERENCES User (Id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS Chat (
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	IdUser INTEGER,
	IdClient INTEGER,
	Name VARCHAR(44),
	Message VARCHAR(1024),
	LastTime VARCHAR(128),
	FOREIGN KEY (IdClient) REFERENCES Client (Id) ON UPDATE CASCADE,
	FOREIGN KEY (IdClient) REFERENCES Client (Id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS GlobalChat (
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	IdUser INTEGER,
	Founder VARCHAR(44),
	Name VARCHAR(44),
	Message VARCHAR(1024),
	LastTime VARCHAR(128),
	FOREIGN KEY (IdUser) REFERENCES User (Id) ON UPDATE CASCADE,
	FOREIGN KEY (IdUser) REFERENCES User (Id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS File (
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	IdUser INTEGER,
	Hash VARCHAR(44),
	Name VARCHAR(128),
	PathName VARCHAR(44),
	Size INTEGER,
	Encr BOOLEAN,
	FOREIGN KEY (IdUser) REFERENCES User (Id) ON UPDATE CASCADE,
	FOREIGN KEY (IdUser) REFERENCES User (Id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS Friend (
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	IdUser INTEGER,
	Hashname VARCHAR(44),
	FOREIGN KEY (IdUser) REFERENCES User (Id) ON UPDATE CASCADE,
	FOREIGN KEY (IdUser) REFERENCES User (Id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS State (
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	IdUser INTEGER,
	UsedF2F BOOLEAN,
	UsedFSH BOOLEAN,
	UsedGCH BOOLEAN,
	FOREIGN KEY (IdUser) REFERENCES User (Id) ON UPDATE CASCADE,
	FOREIGN KEY (IdUser) REFERENCES User (Id) ON DELETE CASCADE
);
`)
	if err != nil {
		panic("can't exec database")
	}
}
