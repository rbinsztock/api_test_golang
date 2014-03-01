package main

import (
	"database/sql"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
)

var db *sql.DB

func main() {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/api_test_development")

	if err!= nil{
		fmt.Println(err)
	}
	defer db.Close()

	m := martini.Classic()
	m.Use(render.Renderer())

	m.Use(func(res http.ResponseWriter, req *http.Request, r render.Render) {
		api_token := ""
		api_token = req.URL.Query().Get("key")
		if api_token == "" {
			r.JSON(404, map[string]interface{}{"status": "Fail", "error_message": "Need api token"})
			return
		}
		current_user, user_id := GetUser(db, api_token)
		if user_id < 0 {
	      r.JSON(404, map[string]interface{}{"status": "Fail", "error_message": "Bad api key"})
	      return
	    }
	    m.Map(current_user)
	})
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Get("/campaigns", func(current_user User, r render.Render) {
  		campaigns := GetCampaign(db, current_user.Id)
  		r.JSON(200, map[string]interface{}{"status": "Success", "data": campaigns})
	})
	m.Get("/campaigns/:id", func(current_user User, params martini.Params, r render.Render) {
		paramId, err := strconv.Atoi(params["id"])
		if err != nil {
			r.JSON(404, map[string]interface{}{"status": "Fail", "error_message": err.Error()})
			return
		}
		campaign, id := GetCampaign(db, current_user.Id, paramId)
		if id > 0 {
			 r.JSON(200, map[string]interface{}{"status": "Success", "data": campaign})
			 return
		}
		r.JSON(404, map[string]interface{}{"status": "Fail", "error_message": "Campaign not found"})
		
	})
	m.Run()
}