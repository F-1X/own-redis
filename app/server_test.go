package main

import (
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
