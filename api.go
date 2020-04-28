// +build !wasm

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (d *DataModel) RequestMatch(r *http.Request) bool {
	return strings.HasPrefix(r.URL.Path, "/api/v1")
}

func (d *DataModel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		println(time.Now().String(), r.RemoteAddr, "GET", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		users := r.URL.Query()["user"]
		user := ""
		if len(users) > 0 {
			user = users[0]
		}
		json.NewEncoder(w).Encode(d.UserView(user))
	case http.MethodPost:
		println(time.Now().String(), r.RemoteAddr, "POST", r.URL.Path)
		switch {
		case strings.HasSuffix(r.URL.Path, "/users/login"):
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			user := r.FormValue("user")
			println("login", user)
			d.Login(user)
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(fmt.Sprintf(`{"%v": "logged in"}`, user)))
		case strings.HasSuffix(r.URL.Path, "/users/logout"):
			println(time.Now().String(), r.RemoteAddr, "POST", r.URL.Path)
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			user := r.FormValue("user")
			println("logout", user)
			d.Logout(user)
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(fmt.Sprintf(`{"%v": "logged out"}`, user)))
		case strings.HasSuffix(r.URL.Path, "/users/add"):
			println(time.Now().String(), r.RemoteAddr, "POST", r.URL.Path)
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			user := r.FormValue("user")
			d.AppendUser(&User{
				Name: user,
			})
			println("adding user", user)
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(fmt.Sprintf(`{"%v": "added"}`, user)))
		case strings.HasSuffix(r.URL.Path, "/start/topic"):
			println(time.Now().String(), r.RemoteAddr, "POST", r.URL.Path)
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			user := r.FormValue("user")
			state := r.FormValue("state")
			println("set state", user, state)
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(fmt.Sprintf(`{"%v": "%v"}`, state, user)))
			d.SetState(StateSetTopic)
		case strings.HasSuffix(r.URL.Path, "/topic"):
			println(time.Now().String(), r.RemoteAddr, "POST", r.URL.Path)
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			user := r.FormValue("user")
			topic := r.FormValue("topic")
			println("set topic", user, topic)
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(fmt.Sprintf(`{"%v": "%v"}`, topic, user)))
			d.SetTopic(topic)
		case strings.HasSuffix(r.URL.Path, "/start/vote"):
			println(time.Now().String(), r.RemoteAddr, "POST", r.URL.Path)
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			user := r.FormValue("user")
			println("start vote", user)
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(fmt.Sprintf(`{"%v": "%v"}`, StateVote, user)))
			d.SetState(StateVote)
		case strings.HasSuffix(r.URL.Path, "/vote"):
			println(time.Now().String(), r.RemoteAddr, "POST", r.URL.Path)
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			user := r.FormValue("user")
			vote := r.FormValue("vote")
			println("cast vote", user, vote)
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(fmt.Sprintf(`{"%v": "%v"}`, user, vote)))
			d.CastVote(user, vote)
		}
	}
}

func (d *DataModel) UserView(username string) *DataModel {
	newd := &DataModel{
		State: d.State,
		Topic: d.Topic,
		Users: make([]*User, len(d.Users)),
	}
	anyBlank := false
	for k, v := range d.Users {
		newd.Users[k] = &User{
			Name:   v.Name,
			Status: v.Status,
			Vote:   v.Vote,
		}
		if v.Vote == "" {
			anyBlank = true
		}
	}
	// if any are blank, then obfuscate the other votes
	if anyBlank {
		for _, v := range newd.Users {
			if v.Name != username {
				v.Vote = "?"
			}
		}
	}
	return newd
}
