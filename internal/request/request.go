package request

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/red7-c/httpfromtcp/internal/headers"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	State       parserState
	Body        string
}

type parserState string

const (
	StateInit   parserState = "init"
	StateDone   parserState = "done"
	StateError  parserState = "error"
	StateHeader parserState = "headers"
	StateBody   parserState = "Body"
)

var ErrorMalformedRequestLine = fmt.Errorf("malformed request line")
var ErrUnsupportedHttpVersion = fmt.Errorf("unsupported http version")
var ErrorRequestInErrorState = fmt.Errorf("request in error state")
var SP = []byte("\r\n")

func newRequest() *Request {
	return &Request{
		State:   StateInit,
		Headers: headers.NewHeaders(),
		Body: "",
	}
}

func (r *Request) Done() bool {
	return r.State == StateDone || r.State == StateError
}

func getInt(headers *headers.Headers, name string, defaultValue int) int {
	valuestr, exist := headers.Get(name)
	if !exist {
		return defaultValue
	}
	value, err := strconv.Atoi(valuestr)
	if err != nil {
		return defaultValue
	}
	return value
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {

	idx := bytes.Index(b, SP)

	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	read := idx + len(SP)

	parts := bytes.Split(startLine, []byte(" "))

	if len(parts) != 3 {
		return nil, 0, ErrorMalformedRequestLine
	}

	httpParts := bytes.Split(parts[2], []byte("/"))

	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ErrUnsupportedHttpVersion
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpParts[1]),
	}

	return rl, read, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		currentData := data[read:]
		switch r.State {

		case StateError:
			return 0, ErrorRequestInErrorState
		case StateInit:
			rl, n, err := parseRequestLine(currentData)

			if err != nil {
				r.State = StateError
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n
			r.State = StateHeader
		case StateHeader:
			n, done, err := r.Headers.Parse(currentData)

			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}

			read += n

			if done {
				r.State = StateBody
			}
		case StateBody:
			length := getInt(r.Headers, "content-length", 0)
			if length == 0 {
				r.State = StateDone
			}


	
		case StateDone:
			break outer
		default:
			panic("well well well...")
		}
	}
	return read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	//NOTE: buf size could be an issue, a request that exceeds 1k would be an issue...
	buf := make([]byte, 1024)
	bufLen := 0

	for !request.Done() {
		n, err := reader.Read(buf[bufLen:])
		//TODO: Check for the right error handling here.
		if err != nil {
			return nil, err
		}
		bufLen += n
		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}
	return request, nil
}
