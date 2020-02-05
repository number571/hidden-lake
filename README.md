# HiddenLake

> Decentralized network. Version 1.1.0s.

### Characteristics:
1. F2F network. End to end encryption;
2. Package transfer in blockchain;
3. Symmetric algorithm: AES256-CBC;
4. Asymmetric algorithm: RSA2048-OAEP;
5. Hash function: HMAC(SHA256);

### Home page:
<img src="/images/HiddenLake1.png" alt="HomePage"/>

### Used libraries/frameworks:
1. gopeer: [github.com/number571/gopeer](https://github.com/number571/gopeer);
2. vuejs: [github.com/vuejs/vue](https://github.com/vuejs/vue);
3. bootstrap: [github.com/twbs/bootstrap](https://github.com/twbs/bootstrap);
4. go-sqlite3: [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3);
5. websocket: [golang.org/x/net/websocket](https://golang.org/x/net/websocket);
6. jquery: [github.com/jquery/jquery](https://github.com/jquery/jquery);
7. popper: [github.com/popperjs/popper-core](https://github.com/popperjs/popper-core);

### Chat room page:
<img src="/images/HiddenLake8.png" alt="ChatRoomPage"/>

### Modules:
1. Network (kernel): 
* connects/disconnects servers;
* creates blockchain transfer;
* sends Packages with a response to servers (TCP);
2. Intermediate (server): 
* sends API responses to the client (HTTPS, WSS);
* sends Packages with a requests to the network (TCP);
* saves information in a database;
3. Interface (client): 
* sends API requests to the server (HTTPS, WSS);
* single page application;
* native application routing;

### Account page:
<img src="/images/HiddenLake4.png" alt="AccountPage"/>

### Default configuration (config.cfg): 
> Configuration file is created when the application starts.
```json
{
	"host": {
		"http": {
			"ipv4": "localhost",
			"port": ":7545",
			"tls": {
				"crt": "tls/cert.crt",
				"key": "tls/cert.key"
			}
		},
		"tcp": {
			"ipv4": "localhost",
			"port": ":8080"
		}
	}
}
```

### Settings page:
<img src="/images/HiddenLake5.png" alt="SettingsPage"/>

### SQL Tables (database.db):
> Database file is created when the application starts.
```sql
/* Authorization user; */
/* Hashpasw = hash(hash(username+password)); */
CREATE TABLE IF NOT EXISTS User (
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	Hashpasw VARCHAR(44) UNIQUE,
	PrivateKey VARCHAR(4096) UNIQUE
);
/* User connections; */
/* Hashname = hash(public_key); */
CREATE TABLE IF NOT EXISTS Client (
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	IdUser INTEGER,
	Hashname VARCHAR(44),
	Address VARCHAR(64),
	PublicKey VARCHAR(2048),
	FOREIGN KEY (IdUser)  REFERENCES User (Id)
);
/* User chat; */
CREATE TABLE IF NOT EXISTS Chat (
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	IdUser INTEGER,
	Companion VARCHAR(44),
	Name VARCHAR(44),
	Message TEXT,
	LastTime VARCHAR(128),
	FOREIGN KEY (IdUser)  REFERENCES User (Id)
);
/* File information; */
/* Hash = hash(file); */
CREATE TABLE IF NOT EXISTS File (
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	IdUser INTEGER,
	Hash VARCHAR(44),
	Name VARCHAR(128),
	Path VARCHAR(44),
	Size INTEGER,
	FOREIGN KEY (IdUser)  REFERENCES User (Id)
);
```

### Archive page:
<img src="/images/HiddenLake6.png" alt="ArchivePage"/>
