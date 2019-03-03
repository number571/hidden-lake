package utils

import (
	"fmt"
)

func PrintHelp() {
	fmt.Print(`| Arguments:
| - -a, --address = set ipv4:port;
| Example:
| - ./main -a 127.0.0.1:8080
| Program arguments:
| 1 :exit = exit from client
| 2 :send = send to one node message:
| - | :send addr message
| 3 :help = get info about client
| 4 :history = read or delete history:
| - | ~ without arguments = read global history
| - | del, delete = delete global or local history:
| - | - | :history delete = delete global history
| - | - | :history delete addrs = delete local history
| - | loc, local = read local history:
| - | - | :history local addrs
| 5 :network = get info about connections
| 6 :connect = connect to node[s]:
| - | :connect addrs
| 7 :archive = check/download files
| - | ~ without arguments = read your archive
| - | list = read list of files from node:
| - | - | :archive list addr
| - | download = download file from node:
| - | - | :archive download addr file
| 8 :disconnect = disconnect from node[s]
| - | :disconnect addrs
`)
}
