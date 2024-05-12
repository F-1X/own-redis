package main

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert"
)

func Test_ReadArray(t *testing.T) {

	test := "*1\r\n$4\r\nPING\r\n"
	expect := "+PONG\r\n"

	returnData := ReadArray([]byte(test))

	// fmt.Println("return:", returnData, "expect", expect)
	switch returnData[0] {
	case "PING":
		assert.Equal(t, expect, "+"+"PONG"+CRLF)
	}

}

func Test_SET(t *testing.T) {

	test := "*3\r\n$3\r\nSET\r\n$5\r\napple\r\n$5\r\ngrape\r\n"
	expect := "+OK\r\n"

	returnData := ReadArray([]byte(test))

	switch returnData[0] {
	case "SET":
		if err := Cache.SETData(returnData); err != nil {
			fmt.Println("some error")
		}
		assert.Equal(t, expect, "+"+"OK"+CRLF)
	}

	test2 := "*2\r\n$3\r\nGET\r\n$5\r\napple\r\n"
	// expect2 := "$5\r\ngrape\r\n"
	returnData = ReadArray([]byte(test2))

	switch returnData[0] {
	case "GET":
		returnGET, _ := Cache.GETData(returnData)
		assert.Equal(t, "grape", returnGET)
	}

}



func Test_SETGET(t *testing.T) {

	test := "*3\r\n$3\r\nSET\r\n$5\r\napple\r\n$9\r\nblueberry\r\n"
	expect := "+OK\r\n"

	returnData := ReadArray([]byte(test))

	switch returnData[0] {
	case "SET":
		if err := Cache.SETData(returnData); err != nil {
			fmt.Println("some error")
		}
		assert.Equal(t, expect, "+"+"OK"+CRLF)
	}

	test2 := "*2\r\n$3\r\nGET\r\n$5\r\napple\r\n"
	// expect2 := "$5\r\ngrape\r\n"
	returnData = ReadArray([]byte(test2))

	switch returnData[0] {
	case "GET":
		returnGET, _ := Cache.GETData(returnData)
		assert.Equal(t, "grape", returnGET)
	}

}
