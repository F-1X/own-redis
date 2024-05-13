package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type keys_expired struct {
	key     string
	expired time.Duration
}

type Memory struct {
	store map[string]string
	ke    []keys_expired
	mu    *sync.Mutex
}

var SIZE_DEF = 256

var Cache = NewMemory(SIZE_DEF)

func NewMemory(size int) *Memory {
	return &Memory{
		make(map[string]string, size),
		[]keys_expired{},
		&sync.Mutex{},
	}
}

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
			case "SET":
				if err := Cache.SETData(ret); err != nil {
					conn.Write([]byte("$-1\r\n"))
					conn.Close()
				} else {
					conn.Write([]byte("+" + "OK" + CRLF))
				}
			case "GET":
				returnGET, err := Cache.GETData(ret)
				if err != nil {
					conn.Write([]byte("$-1\r\n"))
					// conn.Close()
				} else {
					conn.Write([]byte("+" + returnGET + CRLF))
				}
			}
		}
	}
}

func (m *Memory) GETData(s []string) (string, error) {
	register_any_key := strings.ToUpper(s[1])
	fmt.Println("getGETData func", s[1], "register_any_key", register_any_key)
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.store[register_any_key]; !ok {
		fmt.Println("not ok", m.store, ok, register_any_key)
		return "empty", fmt.Errorf("not found")
	}
	fmt.Println("its exist func", m.store, m.store[register_any_key])
	return m.store[register_any_key], nil
}

func (m *Memory) SETData(s []string) error {
	//fmt.Println("setting key", s[1], "val", s[2], "mode", s[3], "after:", s[4])
	m.mu.Lock()
	register_any_key := strings.ToUpper(s[1])

	if _, ok := m.store[register_any_key]; ok {
		return fmt.Errorf("already exist")
	}
	m.store[register_any_key] = s[2]
	m.mu.Unlock()

	if len(s) > 4 {
		register_any := strings.ToUpper(s[3])
		switch register_any {
		case "PX":
			t, _ := strconv.Atoi(s[4])
			go m.addKeyExpited(register_any_key, time.Millisecond*time.Duration(t))
		case "EX":
			t, _ := strconv.Atoi(s[4])
			go m.addKeyExpited(register_any_key, time.Second*time.Duration(t))
		}
	}
	return nil
}

func (m *Memory) addKeyExpited(key string, expired time.Duration) {
	// <-time.After(expired)
	time.Sleep(expired)
	m.mu.Lock()
	delete(m.store, key)
	m.mu.Unlock()
}

type RASP struct {
	datatype string
	size     int
	raw      []byte
}

func ReadArray(d []byte) []string {
	//fmt.Println("all data:", string(d))

	buffer := bytes.NewBuffer(d)

	_, _ = buffer.ReadByte()
	SA, err := buffer.ReadByte()
	if err != nil {
		return nil
	}
	// fmt.Println("TA:", string(TA), TA)

	SAint, _ := strconv.Atoi(string(SA))
	//fmt.Println("SA:", string(SA), SA)
	_, _ = buffer.ReadByte()
	_, _ = buffer.ReadByte()
	var arrays []string
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

		switch T {
		case BULK_STRING:
			data := ""
			// читаем строку
			for i := 0; i < Sint; i++ {
				c, err := buffer.ReadByte()
				// fmt.Println("read some:", string(c))
				if err != nil {
					break
				}
				data += string(c)
			}
			// fmt.Println("%x", data)
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
