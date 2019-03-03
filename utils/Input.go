package utils

import (
	"os"
	"bufio"
	"strings"
)

func Input() string {
    command, _ := bufio.NewReader(os.Stdin).ReadString('\n')
    return strings.Replace(command, "\n", "", -1)
}
