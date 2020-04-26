package main

type User struct {
	Name   string `json:"name"`
	Status bool   `json:"status"`
	Vote   string `json:"vote"`
}
