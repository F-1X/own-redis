package main

import (
	"fmt"
	"testing"
	"time"

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

	test := "*3\r\n$3\r\nSET\r\n$6\r\norange\r\n$5\r\nmango\r\n"
	expect := "+OK\r\n"

	returnData := ReadArray([]byte(test))

	switch returnData[0] {
	case "SET":
		if err := Cache.SETData(returnData); err != nil {
			fmt.Println("some error")
		}
		assert.Equal(t, expect, "+"+"OK"+CRLF)
	}

	test2 := "*2\r\n$3\r\nGET\r\n$5\r\norange\r\n"
	// expect2 := "$5\r\ngrape\r\n"
	returnData = ReadArray([]byte(test2))

	switch returnData[0] {
	case "GET":
		returnGET, _ := Cache.GETData(returnData)
		assert.Equal(t, "mango", returnGET)
	}

}

func Test_MemoryCache(t *testing.T) {
	requsetGet := "*2\r\n$3\r\nGET\r\n$4\r\npear\r\n"
	b := ReadArray([]byte(requsetGet))
	fmt.Println("b", b)
	input := "*5\r\n$3\r\nSET\r\n$4\r\npear\r\n$6\r\nbanana\r\n$2\r\npx\r\n$3\r\n100\r\n"
	wantGet := "banana"
	a := ReadArray([]byte(input))
	fmt.Println("a", a)
	go func() {
		Cache.SETData(a)
	}()

	time.Sleep(time.Millisecond * 100)
	getActual, err := Cache.GETData(b)
	if err != nil {

		assert.Equal(t, fmt.Errorf("not found"), err)
	}

	assert.Equal(t, wantGet, getActual)

}

func Test_Jus(t *testing.T) {
	input := "*3\r\n$3\r\nSET\r\n$5\r\ngrape\r\n$5\r\nmango\r\n"
	a := ReadArray([]byte(input))
	fmt.Printf("%+v", a)
	Cache.SETData(a)
}
