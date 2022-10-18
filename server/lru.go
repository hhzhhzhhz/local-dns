package server

import glrn "github.com/hashicorp/golang-lru"


var lru *glrn.Cache

func NewLru() (err error) {
	lru, err = glrn.New(1024)
	return err
}

func Lru() *glrn.Cache {
	return lru
}

func dump()  {

}