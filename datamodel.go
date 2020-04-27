// +build !wasm

package main

type DataModel struct {
	users     []*User `json:"users"`
	listeners []chan bool
}

func NewDataModel() *DataModel {
	return &DataModel{}
}

func (d *DataModel) Login(user string) {
	for _, v := range d.users {
		if v.Name == user {
			v.Status = true
			d.Update("login")
		}
	}
}

func (d *DataModel) Logout(user string) {
	for _, v := range d.users {
		if v.Name == user {
			v.Status = false
			d.Update("login")
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
	d.users = append(d.users, user)
	d.Update("append user")
}
