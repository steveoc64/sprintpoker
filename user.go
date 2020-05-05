package main

type User struct {
	Name   string `json:"name"`
	Admin  bool   `json:"admin"`
	Status bool   `json:"status"`
	Vote   string `json:"vote"`
}
