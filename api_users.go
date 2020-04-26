// build !wasm

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/BurntSushi/toml"
)

type usersAPI struct {
	Users []*User `json:"users"`
}

func newUsersAPI() *usersAPI {
	u := &usersAPI{}
	config := make(map[string]interface{})
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
		return u
	}
	for _, v := range config["users"].([]interface{}) {
		u.Users = append(u.Users, &User{
			Name: v.(string),
		})
	}
	return u
}

func (u *usersAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u.Users)
}

func (u *usersAPI) RequestMatch(r *http.Request) bool {
	return strings.HasPrefix(r.URL.String(), "/api/v1/users")
}
