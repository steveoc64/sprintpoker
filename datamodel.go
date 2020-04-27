// +build !wasm

package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type DataModel struct {
	State     AppState    `json:"state"`
	Topic     string      `json:"topic"`
	Users     []*User     `json:"users"`
	listeners []chan bool `json:"-"`
}

func NewDataModel() *DataModel {
	d := &DataModel{}
	config := make(map[string]interface{})
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		// fallback to default install
		if _, err := toml.DecodeFile("/usr/local/etc/poker.toml", &config); err != nil {
			fmt.Println(err)
			return d
		}
	}
	for _, v := range config["users"].([]interface{}) {
		d.Users = append(d.Users, &User{
			Name: v.(string),
		})
	}
	return d
}

func (d *DataModel) Login(user string) {
	for _, v := range d.Users {
		if v.Name == user {
			v.Status = true
			d.Update("login " + user)
		}
	}
}

func (d *DataModel) Logout(user string) {
	for _, v := range d.Users {
		if v.Name == user {
			v.Status = false
			d.Update("logout " + user)
		}
	}
}

func (d *DataModel) Update(why string) {
	println("updating because", why, len(d.listeners))
	for _, v := range d.listeners {
		v <- true
	}
}

func (d *DataModel) Listen() chan bool {
	newChan := make(chan bool, 100)
	d.listeners = append(d.listeners, newChan)
	return newChan
}

func (d *DataModel) AppendUser(user *User) {
	d.Users = append(d.Users, user)
	d.Update("append user " + user.Name)
}

func (d *DataModel) SetState(state AppState) {
	d.State = state
	println("setting state", state)
	d.Update("state")
}

func (d *DataModel) SetTopic(topic string) {
	d.Topic = topic
	println("setting topic", topic)
	d.Update("topic set")
}
