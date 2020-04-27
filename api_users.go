// +build !wasm

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type usersAPI struct {
	data *DataModel
}

func NewUsersAPI(data *DataModel) *usersAPI {
	u := &usersAPI{data: data}
	config := make(map[string]interface{})
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		// fallback to default install
		if _, err := toml.DecodeFile("/usr/local/etc/poker.toml", &config); err != nil {
			fmt.Println(err)
			return u
		}
	}
	for _, v := range config["users"].([]interface{}) {
		u.data.AppendUser(&User{
			Name: v.(string),
		})
	}
	return u
}

func (u *usersAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		println(time.Now().String(), r.RemoteAddr, "GET", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(u.data.users)
	case http.MethodPost:
		println(time.Now().String(), r.RemoteAddr, "POST", r.URL.Path)
		switch {
		case strings.HasSuffix(r.URL.Path, "/login"):
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			user := r.FormValue("user")
			println("login", user)
			u.data.Login(user)
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(fmt.Sprintf(`{"%v": "logged in"}`, user)))
		case strings.HasSuffix(r.URL.Path, "/logout"):
			println(time.Now().String(), r.RemoteAddr, "POST", r.URL.Path)
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			user := r.FormValue("user")
			println("logout", user)
			u.data.Logout(user)
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(fmt.Sprintf(`{"%v": "logged out"}`, user)))
		case strings.HasSuffix(r.URL.Path, "/add"):
			println(time.Now().String(), r.RemoteAddr, "POST", r.URL.Path)
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			user := r.FormValue("user")
			u.data.AppendUser(&User{
				Name: user,
			})
			println("adding user", user)
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(fmt.Sprintf(`{"%v": "added"}`, user)))
		}
	}
}

func (u *usersAPI) RequestMatch(r *http.Request) bool {
	return strings.HasPrefix(r.URL.String(), "/api/v1/users")
}
