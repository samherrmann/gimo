package models

import "github.com/samherrmann/gimo"

type Book struct {
	gimo.DocumentBase `json:",inline"    bson:",inline"`
	Title             string `json:"title"`
	Author            string `json:"author"`
	Publisher         string `json:"publisher"`
}

func (p *Book) New() gimo.Document {
	return &Book{}
}

func (p *Book) Slice() interface{} {
	return &[]Book{}
}
