package gimo

import "gopkg.in/mgo.v2"

// dialDB establishes a new session to the mongoDB cluster
// identified by info.
func dialDB(info *mgo.DialInfo) *mgo.Session {
	s, err := mgo.DialWithInfo(info)
	if err != nil {
		panic("Failed to establish session with database: " + err.Error())
	}
	return s
}
