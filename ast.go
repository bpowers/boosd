package main

import (

)

type ObjKind int
const (
	ObjModel ObjKind = iota
	ObjInterface
	ObjKynd
	ObjFlow
	ObjStock
	ObjTime
	ObjAux
	ObjString
)



type Object struct {
	Name string
	Kind ObjKind
	Decl interface{}
	Data interface{}
	Type string
	Unit string
}
