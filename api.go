package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Api string `json:"api_token"`
	AccountId int `json:"account_id"`
}

type Campaign struct {
	Id  int    `json:"id"`
	Name string `json:"name"`
	AccountId int `json:"account_id"`
}

func GetUser(db *sql.DB, key string) (User, int) {
	var user_id int
	var api_token string
	err := db.QueryRow("select id, api from users where token_api = ? limit 1", key).Scan(&user_id, &api_token)
	if err == sql.ErrNoRows {
		return User{}, 0
	}
	if err != nil {
		fmt.Println(err)
		return User{}, -1
	}
	return User{user_id, api_token}, user_id
}

func GetCampaign(db *sql.DB, user_id int, campaign_id int) (Campaign, int) {
	var (
		id   int
		name string
	)
	err := db.QueryRow("select id, name from campaigns where id = ? and user_id = ? limit 1", campaign_id, user_id).Scan(&id, &name)
	if err == sql.ErrNoRows {
		return Campaign{}, 0
	}
	if err != nil {
		fmt.Println(err)
		return Campaign{}, -1
	}
	return Campaign{id, name}, id
}

func GetCampaigns(db *sql.DB, UserId int) []Campaign {
	campaigns, err := db.Query("select id, name from campaigns where account_id = ?", UserId)
	if err != nil {
		fmt.Println(err)
	}

	var (
		id   int
		name string
	)

	p := make([]Campaign, 0)
	defer campaigns.Close()
	for campaigns.Next() {
		err := campaigns.Scan(&id, &name)
		if err != nil {
			fmt.Println(err)
		} else {
			p = append(p, Campaign{id, name})
		}
	}
	return p
}