package main

import (
	"log"
	"testing"
)

func TestDecodeCmd(t *testing.T) {
	log.Println(decodeCmd([]byte{73, 0, 0, 48, 57, 0, 0, 0, 101}))
	log.Println(decodeCmd([]byte{81, 0, 0, 48, 0, 0, 0, 64, 0}))
}

func TestDf(t *testing.T) {
	df := &DB{}

	log.Println(df.query(363286264, 363288164))
}

func TestRequest(t *testing.T) {
	df := &DB{}

	log.Println(handleRequest(1, df, []byte{73, 0, 0, 48, 57, 0, 0, 0, 101}))
	log.Println(handleRequest(1, df, []byte{81, 0, 0, 48, 0, 0, 0, 64, 0}))
}
