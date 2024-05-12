package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

const (
	SIMPLE_STRINGS   = '+'
	SIMPLE_ERRORS    = '-'
	INTEGERS         = ':'
	BULK_STRING      = '$'
	ARRAYS           = '*'
	NULLS            = '_'
	BOOLEANS         = '#'
	DOUBLES          = ','
	BIG_NUMBERS      = '('
	BULK_ERRORS      = '!'
	VERBATIM_STRINGS = '='
	MAPS             = '%'
	SETS             = '~'
	PUSHES           = '>'

	CR   = '\r'
	LF   = '\n'
	CRLF = "\r\n"

	EX = 0
	PX = 0
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go acceptLoop(conn)
	}

}

func acceptLoop(conn net.Conn) {
	tmp := make([]byte, 4096)
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		buf := make([]byte, n)
		copy(buf, tmp)

		switch buf[0] {
		case ARRAYS:
			ret := ReadArray(buf)
			fmt.Println("ret", ret, len(ret))
			switch ret[0] {
			case "ECHO":
				fmt.Println("d", ret, len(ret))
				conn.Write([]byte("+" + ret[1] + CRLF))
				conn.Close()
			case "PING":
				conn.Write([]byte("+" + "PONG" + CRLF))
				// conn.Close()
			}
		}

		switch string(buf[0:3]) {
		case "GET":

		case "SET":
			
		}
	}
}

type RASP struct {
	datatype string
	size     int
	raw      []byte
}

func ReadArray(d []byte) []string {
	fmt.Println("all data:", string(d))
	arrays := make([]string, 0, len(d))
	buffer := bytes.NewBuffer(d)

	TA, err := buffer.ReadByte()
	SA, err := buffer.ReadByte()
	if err != nil {
		return nil
	}
	fmt.Println("TA:", string(TA), TA)

	SAint, _ := strconv.Atoi(string(SA))
	fmt.Println("SA:", string(SA), SA)
	_, _ = buffer.ReadByte()
	_, _ = buffer.ReadByte()

	for e := 0; e < SAint; e++ {
		// считываем ТИП и РАЗМЕР
		T, err := buffer.ReadByte()
		if err != nil {
			break
		}

		sizeString := ""
		for {
			sByte, err := buffer.ReadByte()
			if err != nil {
				break
			}
			if sByte == '\r' {
				sByte, err := buffer.ReadByte()
				if err != nil {
					break
				}
				if sByte == '\n' {
					break
				}
			}
			sizeString += string(sByte)
		}

		Sint, _ := strconv.Atoi(string(sizeString))
		// пропускаем CRLF

		fmt.Println("Sint:", Sint, string(sizeString))
		switch T {
		case BULK_STRING:
			data := ""
			// читаем строку
			for i := 0; i < Sint; i++ {
				c, err := buffer.ReadByte()
				fmt.Println("read some:", string(c))
				if err != nil {
					break
				}
				data += string(c)
			}
			arrays = append(arrays, data)

		default:
			return []string{"unknown type"}
		}
		_, _ = buffer.ReadByte()
		_, _ = buffer.ReadByte()

	}

	return arrays
}

func readBulkString(s []byte) string {
	lenght := int(s[0])

	str := string(s[3:lenght])

	return str

}
func readSimpleString(s []byte) []byte {
	return s[1 : len(s)-2]
}

func readSimpleErrors(s []byte) []byte {
	// 	-ERR unknown command 'asdf'
	// -WRONGTYPE Operation against a key holding the wrong kind of value
	return s[1 : len(s)-2]
}

func readIntegers(s []byte) int64 {
	var numberString string
	var i int
	if s[0] == '+' || s[0] == '-' {
		i++
	}
	for ; i < len(s)-1; i++ {
		if s[i] >= '0' || s[i] <= '9' {
			numberString += string(s[i])
		} else {
			panic("wtf number")
		}
	}

	number, err := strconv.ParseInt(numberString, 10, 64)
	if err != nil {
		panic(err)
	}

	return number
}

var Types map[byte]string

func initTypes() {
	Types = make(map[byte]string, 14)

	Types['+'] = "SIMPLE_STRINGS"
	Types['$'] = "BULK_STRING"
}

func parseRESP(s []byte) {
	switch s[0] {
	case SIMPLE_STRINGS:
	case INTEGERS:
	case BULK_STRING:
	case ARRAYS:
	case NULLS:
	case BOOLEANS:
	case DOUBLES:
	case BIG_NUMBERS:
	case BULK_ERRORS:
	case VERBATIM_STRINGS:
	case MAPS:
	case SETS:
	case PUSHES:
	}
}
