package shared

import "io"

type MetaResources struct {
	MetaRow string `json:"MetaRow"`
	Count   int    `json:"Count"`
}

type Talents struct {
	RowName string `json:"RowName"`
	Rank    int    `json:"Rank"`
}

type Metadata struct {
	FileName string `json:"-"`
	Path     string `json:"-"`
}

type PlayerData interface {
	Read() error
	ReadF(io.Reader) error
	Write() error
	WriteF(io.Writer) error
}
