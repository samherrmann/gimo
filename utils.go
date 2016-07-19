package gimo

import (
	"reflect"

	"gopkg.in/mgo.v2"
)

// makeSlice returns a slice of i as an interface{}.
// Noteworthy resource: http://stackoverflow.com/a/25386460/3634032
func makeSlice(i interface{}) interface{} {
	iType := reflect.TypeOf(i)
	sliceType := reflect.SliceOf(iType)
	goValue := reflect.New(sliceType)
	goValue.Elem().Set(reflect.MakeSlice(sliceType, 0, 0))
	return goValue.Interface()
}

// dialDB establishes a new session to the mongoDB cluster
// identified by info.
func dialDB(info *mgo.DialInfo) *mgo.Session {
	s, err := mgo.DialWithInfo(info)
	if err != nil {
		panic("Failed to establish session with database: " + err.Error())
	}
	return s
}
