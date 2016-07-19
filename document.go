package gimo

type Document interface {
	GetID() string
	SetID(id string)
	New() Document
	Slice() interface{}
}

type DocumentBase struct {
	ID string `json:"id" bson:"_id"`
}

func (d *DocumentBase) SetID(id string) {
	d.ID = id
}

func (d *DocumentBase) GetID() string {
	return d.ID
}
