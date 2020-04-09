package settings

import (
	"github.com/number571/hiddenlake/utils"
	"time"
)

func ClearTempEmails(checkTime time.Duration) {
	for {
		time.Sleep(checkTime)
		checkLifetimeEmails()
	}
}

func checkLifetimeEmails() {
	var (
		id       uint64
		lasttime string
		currTime = utils.ParseTime(utils.CurrentTime())
	)
	rows, err := DB.Query("SELECT Id, LastTime FROM Email WHERE Temporary=1")
	if err != nil {
		panic("query 'checklifetimeemails' failed")
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &lasttime)
		if err != nil {
			break
		}
		checktime := utils.ParseTime(lasttime)
		if checktime.Add(LIFETIME).Before(currTime) {
			DB.Exec("DELETE FROM Email WHERE Id=$1", id)
		}
	}
}
