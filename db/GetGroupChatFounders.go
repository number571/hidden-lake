package db

import (
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func GetGroupChatFounders(user *models.User) []string {
	var (
		founders []string
		founder  string
	)
	id := GetUserId(user.Username)
	if id < 0 {
		return nil
	}
	rows, err := settings.DB.Query(
		"SELECT Founder FROM GlobalChat WHERE IdUser=$1 GROUP BY Founder ORDER BY Id",
		id,
	)
	if err != nil {
		panic("query 'getglobalchatfounders' failed")
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&founder)
		if err != nil {
			return nil
		}
		founders = append(founders, founder)
	}
	return founders
}
