# [DEPRECATED] HiddenLake

> Current version https://github.com/number571/go-peer/tree/master/cmd/hidden_lake

> Decentralized private network. Version 1.0.7s.

### Characteristics:
1. F2F network. End to end encryption;
2. Supported hidden connections;
3. Symmetric algorithm: AES256-[CBC,OFB];
4. Asymmetric algorithm: RSA3072-OAEP;
5. Hash function: HMAC(SHA256);

### Home page:
<img src="/images/HiddenLake1.png" alt="HomePage"/>

### Abilities:
1. Private / Group chats;
2. Emails;
3. File sharing / storage;

### Chat room page:
<img src="/images/HiddenLake14.png" alt="ChatRoomPage"/>

### Used libraries/frameworks:
1. gopeer: [github.com/number571/gopeer](https://github.com/number571/gopeer);
2. vuejs: [github.com/vuejs/vue](https://github.com/vuejs/vue);
3. bootstrap: [github.com/twbs/bootstrap](https://github.com/twbs/bootstrap);
4. go-sqlite3: [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3);
5. jquery: [github.com/jquery/jquery](https://github.com/jquery/jquery);
6. popper: [github.com/popperjs/popper-core](https://github.com/popperjs/popper-core);

### Account page:
<img src="/images/HiddenLake4.png" alt="AccountPage"/>

### Modules:
1. Network (kernel): 
* connects/disconnects servers;
* sends Packages with a response to servers (TCP);
2. Intermediate (server): 
* sends API responses to the client (HTTPS, WSS);
* sends Packages with a requests to the network (TCP);
* saves information in a database;
3. Interface (client): 
* sends API requests to the server (HTTPS, WSS);
* single page application;
* native application routing;

### Email page:
<img src="/images/HiddenLake16.png" alt="EmailPage"/>

### Default configuration (config.cfg): 
> Configuration file is created when the application starts.
```json
{
	"http": {
		"ipv4": "localhost",
		"port": ":7545"
	},
	"tcp": {
		"ipv4": "",
		"port": ""
	},
	"tls": {
		"crt": "tls/cert.crt",
		"key": "tls/cert.key"
	}
}
```

### Archive page:
<img src="/images/HiddenLake6.png" alt="ArchivePage"/>

### SQL Tables (database.db):
> Database file is created when the application starts.
```sql
/* Authorization user; */
/* Username = hash(username); */
/* Hashpasw = hash(hash(password+salt)); */
CREATE TABLE IF NOT EXISTS User (
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	Username VARCHAR(44) UNIQUE,
	Salt VARCHAR(32),
	Hashpasw VARCHAR(44) UNIQUE,
	PrivateKey VARCHAR(4096) UNIQUE
);
/* User emails; */
/* Hash = hash(hash(sender_pub)+receiver+message+salt); */
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
/* User connections; */
/* Hashname = hash(public_key); */
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
/* User chat; */
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
/* User global chat; */
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
/* File information; */
/* Hash = hash(file); */
/* PathName = hash(hash(file)+random(16)); */
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
/* Connection list for F2F network */
CREATE TABLE IF NOT EXISTS Friend (
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	IdUser INTEGER,
	Hashname VARCHAR(44),
	FOREIGN KEY (IdUser) REFERENCES User (Id) ON UPDATE CASCADE,
	FOREIGN KEY (IdUser) REFERENCES User (Id) ON DELETE CASCADE
);
/* User saved state */
CREATE TABLE IF NOT EXISTS State (
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	IdUser INTEGER,
	UsedF2F BOOLEAN,
	UsedFSH BOOLEAN,
	UsedGCH BOOLEAN,
	FOREIGN KEY (IdUser) REFERENCES User (Id) ON UPDATE CASCADE,
	FOREIGN KEY (IdUser) REFERENCES User (Id) ON DELETE CASCADE
);
```

### Network page:
<img src="/images/HiddenLake9.png" alt="NetworkPage"/>
