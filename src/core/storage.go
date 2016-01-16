package main

import "gopkg.in/mgo.v2"

var session *mgo.Session

type Storage struct {
	Session *mgo.Session
	DB      *mgo.Database
}

func NewStorage() *Storage {
	storage := new(Storage)
	storage.Session = session
	return storage
}
