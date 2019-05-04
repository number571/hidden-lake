package utils

import (
	"fmt"
)

func PrintHelp() {
	fmt.Print(`
| I. In the start client:

| Commands with parameters:
| 1. [--login, -l] = set login (first run is signup)
| 2. [--password, -p] = set password (first run is signup)
| 3. [--address, -a] = set address ipv4:port

| Commands without parameters:
| 1. [--interface, -i] = run GUI interface in browser on port 7545
| 2. [--delete, -d] = delete Archive, Config and database files with multiple overwriting
| 3. [--help, -h] = get information about client
| 4. [--f2f, -f] = run F2F connection

| II. In the run client:

| Commands in CLI client for all users:
| 1. [:exit] = exit from client
| 2. [:help] = get information about client
| 3. [:interface] = on/off GUI interface

| Commands in CLI client if not authorized:
| 1. [:login] = set login (first run is signup)
| 2. [:password] = set password (first run in signup)
| 3. [:enter] = authorization from the entered login and password
| 4. [:address] = set address ipv4:port

| Commands in CLI client if authorized:
| 1.  [:whoami] = get hashname
| 2.  [:logout] = logout from authorized user
| 3.  [:network] = get list of connections
| 4.  [:send] = send local message to another user
| 5.  [:email] = read or write email to another user
| 6.  [:archive] = get list or download files from archive another user
| 7.  [:history] = get local/global messages or delete messages
| 8.  [:connect] = connect to another user
| 9.  [:disconnect] = disconnect from user
| 10. [] = send global message to another users

`)
}
