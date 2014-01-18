// Package dataurl parses Data URL Schemes
// according to RFC 2397
// (http://tools.ietf.org/html/rfc2397).
package dataurl

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
)

const (
	EncodingBase64 = "base64"
	EncodingASCII  = "ascii"
)

func defaultMediaType() MediaType {
	return MediaType{
		"text",
		"plain",
		map[string]string{"charset": "US-ASCII"},
	}
}

// MediaType is the combination of a media type, a media subtype
// and optional parameters.
type MediaType struct {
	Type    string
	Subtype string
	Params  map[string]string
}

func (mt *MediaType) ContentType() string {
	return fmt.Sprintf("%s/%s", mt.Type, mt.Subtype)
}

// String implements the Stringer interface.
//
// Params values are escaped with the Escape function, rather than in a quoted string.
func (mt *MediaType) String() string {
	var buf bytes.Buffer
	for k, v := range mt.Params {
		fmt.Fprintf(&buf, ";%s=%s", k, EscapeString(v))
	}
	return mt.ContentType()+(&buf).String()
}

// DataURL is the combination of a MediaType describing the type of its Data.
type DataURL struct {
	MediaType
	Encoding  string
	Data      []byte
}

type encodedDataReader func(string) ([]byte, error)

var asciiDataReader encodedDataReader = func(s string) ([]byte, error) {
	us, err := Unescape(s)
	if err != nil {
		return nil, err
	}
	return []byte(us), nil
}

var base64DataReader encodedDataReader = func(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return []byte(data), nil
}

type parser struct {
	du                  *DataURL
	l                   *lexer
	currentAttr         string
	unquoteParamVal     bool
	encodedDataReaderFn encodedDataReader
}

func (p *parser) parse() error {
	for item := range p.l.items {
		switch item.t {
		case itemError:
			return errors.New(item.String())
		case itemMediaType:
			p.du.MediaType.Type = item.val
			// Should we clear the default
			// "charset" parameter at this point?
			delete(p.du.MediaType.Params, "charset")
		case itemMediaSubType:
			p.du.MediaType.Subtype = item.val
		case itemParamAttr:
			p.currentAttr = item.val
		case itemLeftStringQuote:
			p.unquoteParamVal = true
		case itemParamVal:
			var val string = item.val
			if p.unquoteParamVal {
				p.unquoteParamVal = false
				us, err := strconv.Unquote("\"" + val + "\"")
				if err != nil {
					return err
				}
				val = us
			} else {
				us, err := UnescapeToString(val)
				if err != nil {
					return err
				}
				val = us
			}
			p.du.MediaType.Params[p.currentAttr] = val
		case itemBase64Enc:
			p.du.Encoding = EncodingBase64
			p.encodedDataReaderFn = base64DataReader
		case itemDataComma:
			if p.encodedDataReaderFn == nil {
				p.encodedDataReaderFn = asciiDataReader
			}
		case itemData:
			reader, err := p.encodedDataReaderFn(item.val)
			if err != nil {
				return err
			}
			p.du.Data = reader
		case itemEOF:
			if p.du.Data == nil {
				p.du.Data = []byte("")
			}
			return nil
		}
	}
	panic("EOF not found")
}

// DecodeString decodes a Data URL scheme string.
func DecodeString(s string) (*DataURL, error) {
	du := &DataURL{
		MediaType: defaultMediaType(),
		Encoding:  EncodingASCII,
	}

	parser := &parser{
		du: du,
		l:  lex(s),
	}
	if err := parser.parse(); err != nil {
		return nil, err
	}
	return du, nil
}

// Decode decodes a Data URL scheme from a io.Reader.
func Decode(r io.Reader) (*DataURL, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return DecodeString(string(data))
}
