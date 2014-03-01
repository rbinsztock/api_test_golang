package main

import (
	"database/sql"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/martini-contrib/auth"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var m *martini.Martini

func init() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:5234)/api_test_development")

	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	m = martini.New()
	// Setup middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	m.Use(MapEncoder)
	// Setup routes
	r := martini.NewRouter()
	r.Get("/", func() string {
		return "Hello world!"
	})

	r.Get("/campaigns", func(current_user User, r render.Render) {
		campaigns := GetCampaigns(db, current_user.Id)
		r.JSON(200, map[string]interface{}{"status": "Success", "data": campaigns})
	})
	r.Get("/campaigns/:id", func(current_user User, params martini.Params, r render.Render) {
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
	// Add the router action
	m.Action(r.Handle)
}

var rxExt = regexp.MustCompile(`(\.(?:json))\/?$`)

func MapEncoder(c martini.Context, w http.ResponseWriter, r *http.Request) {
	// Get the format extension
	matches := rxExt.FindStringSubmatch(r.URL.Path)
	ft := ".json"
	if len(matches) > 1 {
		// Rewrite the URL without the format extension
		l := len(r.URL.Path) - len(matches[1])
		if strings.HasSuffix(r.URL.Path, "/") {
			l--
		}
		r.URL.Path = r.URL.Path[:l]
		ft = matches[1]
	}
	if ft == ".json" {
		// Inject the requested encoder
		c.MapTo(jsonEncoder{}, (*Encoder)(nil))
		w.Header().Set("Content-Type", "application/json")
	}
}

func main() {
	//m.Use(auth.Basic("username", "secretpassword"))

	m.Use(func(current_user User, res http.ResponseWriter, req *http.Request, r render.Render) {
		//current_user, id := GetUser(db, current_user.Id)
		auth.Basic(current_user.Email, current_user.Api)
		// if user.api_token == "" {
		// 	r.JSON(404, map[string]interface{}{"status": "Fail", "error_message": "Need api token"})
		// 	return
		// }
		// if user < 0 {
		// 	r.JSON(404, map[string]interface{}{"status": "Fail", "error_message": "Bad api key"})
		// 	return
		// }
		m.Map(current_user)
	})
	m.Run()
}
