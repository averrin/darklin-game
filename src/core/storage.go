package main

import "gopkg.in/mgo.v2"

var session *mgo.Session

//Storage - some db logic
type Storage struct {
	Session *mgo.Session
	DB      *mgo.Database
}

//NewStorage constructor
func NewStorage() *Storage {
	storage := new(Storage)
	storage.Session = session
	return storage
}
