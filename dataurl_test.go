package dataurl

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

type dataURLTest struct {
	InputRawDataURL string
	ExpectedItems   []item
	ExpectedDataURL DataURL
}

func genTestTable() []dataURLTest {
	return []dataURLTest{
		dataURLTest{
			`data:;base64,aGV5YQ==`,
			[]item{
				item{itemDataPrefix, dataPrefix},
				item{itemParamSemicolon, ";"},
				item{itemBase64Enc, "base64"},
				item{itemDataComma, ","},
				item{itemData, "aGV5YQ=="},
				item{itemEOF, ""},
			},
			DataURL{
				defaultMediaType(),
				EncodingBase64,
				[]byte("heya"),
			},
		},
		dataURLTest{
			`data:text/plain;base64,aGV5YQ==`,
			[]item{
				item{itemDataPrefix, dataPrefix},
				item{itemMediaType, "text"},
				item{itemMediaSep, "/"},
				item{itemMediaSubType, "plain"},
				item{itemParamSemicolon, ";"},
				item{itemBase64Enc, "base64"},
				item{itemDataComma, ","},
				item{itemData, "aGV5YQ=="},
				item{itemEOF, ""},
			},
			DataURL{
				MediaType{
					"text",
					"plain",
					map[string]string{},
				},
				EncodingBase64,
				[]byte("heya"),
			},
		},
		dataURLTest{
			`data:text/plain;charset=utf-8;base64,aGV5YQ==`,
			[]item{
				item{itemDataPrefix, dataPrefix},
				item{itemMediaType, "text"},
				item{itemMediaSep, "/"},
				item{itemMediaSubType, "plain"},
				item{itemParamSemicolon, ";"},
				item{itemParamAttr, "charset"},
				item{itemParamEqual, "="},
				item{itemParamVal, "utf-8"},
				item{itemParamSemicolon, ";"},
				item{itemBase64Enc, "base64"},
				item{itemDataComma, ","},
				item{itemData, "aGV5YQ=="},
				item{itemEOF, ""},
			},
			DataURL{
				MediaType{
					"text",
					"plain",
					map[string]string{
						"charset": "utf-8",
					},
				},
				EncodingBase64,
				[]byte("heya"),
			},
		},
		dataURLTest{
			`data:text/plain;charset=utf-8;foo=bar;base64,aGV5YQ==`,
			[]item{
				item{itemDataPrefix, dataPrefix},
				item{itemMediaType, "text"},
				item{itemMediaSep, "/"},
				item{itemMediaSubType, "plain"},
				item{itemParamSemicolon, ";"},
				item{itemParamAttr, "charset"},
				item{itemParamEqual, "="},
				item{itemParamVal, "utf-8"},
				item{itemParamSemicolon, ";"},
				item{itemParamAttr, "foo"},
				item{itemParamEqual, "="},
				item{itemParamVal, "bar"},
				item{itemParamSemicolon, ";"},
				item{itemBase64Enc, "base64"},
				item{itemDataComma, ","},
				item{itemData, "aGV5YQ=="},
				item{itemEOF, ""},
			},
			DataURL{
				MediaType{
					"text",
					"plain",
					map[string]string{
						"charset": "utf-8",
						"foo":     "bar",
					},
				},
				EncodingBase64,
				[]byte("heya"),
			},
		},
		dataURLTest{
			`data:application/json;charset=utf-8;foo="b\"<@>\"r";style=unformatted%20json;base64,eyJtc2ciOiAiaGV5YSJ9`,
			[]item{
				item{itemDataPrefix, dataPrefix},
				item{itemMediaType, "application"},
				item{itemMediaSep, "/"},
				item{itemMediaSubType, "json"},
				item{itemParamSemicolon, ";"},
				item{itemParamAttr, "charset"},
				item{itemParamEqual, "="},
				item{itemParamVal, "utf-8"},
				item{itemParamSemicolon, ";"},
				item{itemParamAttr, "foo"},
				item{itemParamEqual, "="},
				item{itemLeftStringQuote, "\""},
				item{itemParamVal, `b\"<@>\"r`},
				item{itemRightStringQuote, "\""},
				item{itemParamSemicolon, ";"},
				item{itemParamAttr, "style"},
				item{itemParamEqual, "="},
				item{itemParamVal, "unformatted%20json"},
				item{itemParamSemicolon, ";"},
				item{itemBase64Enc, "base64"},
				item{itemDataComma, ","},
				item{itemData, "eyJtc2ciOiAiaGV5YSJ9"},
				item{itemEOF, ""},
			},
			DataURL{
				MediaType{
					"application",
					"json",
					map[string]string{
						"charset": "utf-8",
						"foo":     `b"<@>"r`,
						"style":   "unformatted json",
					},
				},
				EncodingBase64,
				[]byte(`{"msg": "heya"}`),
			},
		},
		dataURLTest{
			`data:xxx;base64,aGV5YQ==`,
			[]item{
				item{itemDataPrefix, dataPrefix},
				item{itemError, "invalid character for media type"},
			},
			DataURL{},
		},
		dataURLTest{
			`data:,`,
			[]item{
				item{itemDataPrefix, dataPrefix},
				item{itemDataComma, ","},
				item{itemEOF, ""},
			},
			DataURL{
				defaultMediaType(),
				EncodingASCII,
				[]byte(""),
			},
		},
		dataURLTest{
			`data:,A%20brief%20note`,
			[]item{
				item{itemDataPrefix, dataPrefix},
				item{itemDataComma, ","},
				item{itemData, "A%20brief%20note"},
				item{itemEOF, ""},
			},
			DataURL{
				defaultMediaType(),
				EncodingASCII,
				[]byte("A brief note"),
			},
		},
		dataURLTest{
			`data:image/svg+xml-im.a.fake;base64,cGllLXN0b2NrX1RoaXJ0eQ==`,
			[]item{
				item{itemDataPrefix, dataPrefix},
				item{itemMediaType, "image"},
				item{itemMediaSep, "/"},
				item{itemMediaSubType, "svg+xml-im.a.fake"},
				item{itemParamSemicolon, ";"},
				item{itemBase64Enc, "base64"},
				item{itemDataComma, ","},
				item{itemData, "cGllLXN0b2NrX1RoaXJ0eQ=="},
				item{itemEOF, ""},
			},
			DataURL{
				MediaType{
					"image",
					"svg+xml-im.a.fake",
					map[string]string{},
				},
				EncodingBase64,
				[]byte("pie-stock_Thirty"),
			},
		},
	}
}

func expectItems(expected, actual []item) bool {
	if len(expected) != len(actual) {
		return false
	}
	for i, _ := range expected {
		if expected[i].t != actual[i].t {
			return false
		}
		if expected[i].val != actual[i].val {
			return false
		}
	}
	return true
}

func equal(du1, du2 *DataURL) (bool, error) {
	if !reflect.DeepEqual(du1.MediaType, du2.MediaType) {
		return false, nil
	}
	if du1.Encoding != du2.Encoding {
		return false, nil
	}

	if du1.Data == nil || du2.Data == nil {
		return false, fmt.Errorf("nil Data")
	}

	if !bytes.Equal(du1.Data, du2.Data) {
		return false, nil
	}
	return true, nil
}

func TestLexDataURLs(t *testing.T) {
	for _, test := range genTestTable() {
		l := lex(test.InputRawDataURL)
		items := make([]item, 0)
		for item := range l.items {
			items = append(items, item)
		}
		if !expectItems(test.ExpectedItems, items) {
			t.Errorf("Expected %v, got %v", test.ExpectedItems, items)
		}
	}
}

func TestDataURLs(t *testing.T) {
	for _, test := range genTestTable() {
		var expectedItemError string
		for _, item := range test.ExpectedItems {
			if item.t == itemError {
				expectedItemError = item.String()
				break
			}
		}
		dataURL, err := Decode(strings.NewReader(test.InputRawDataURL))
		if expectedItemError == "" && err != nil {
			t.Error(err)
			continue
		} else if expectedItemError != "" && err == nil {
			t.Errorf("Expected error \"%s\", got nil", expectedItemError)
			continue
		} else if expectedItemError != "" && err != nil {
			if err.Error() != expectedItemError {
				t.Errorf("Expected error \"%s\", got \"%s\"", expectedItemError, err.Error())
			}
			continue
		}

		if ok, err := equal(dataURL, &test.ExpectedDataURL); err != nil {
			t.Error(err)
		} else if !ok {
			t.Errorf("Expected %v, got %v", test.ExpectedDataURL, *dataURL)
		}
	}
}

func BenchmarkLex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range genTestTable() {
			l := lex(test.InputRawDataURL)
			for _ = range l.items {
			}
		}
	}
}

const rep = `^data:(?P<mediatype>\w+/[\w\+\-\.]+)?(?P<parameter>(?:;[\w\-]+="?[\w\-\\<>@,";:%]*"?)+)?(?P<base64>;base64)?,(?P<data>.*)$`

func TestRegexp(t *testing.T) {
	re, err := regexp.Compile(rep)
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range genTestTable() {
		shouldMatch := true
		for _, item := range test.ExpectedItems {
			if item.t == itemError {
				shouldMatch = false
				break
			}
		}
		// just test it matches, do not parse
		if re.MatchString(test.InputRawDataURL) && !shouldMatch {
			t.Error("doesn't match", test.InputRawDataURL)
		} else if !re.MatchString(test.InputRawDataURL) && shouldMatch {
			t.Error("match", test.InputRawDataURL)
		}
	}
}

func BenchmarkRegexp(b *testing.B) {
	re, err := regexp.Compile(rep)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		for _, test := range genTestTable() {
			_ = re.FindStringSubmatch(test.InputRawDataURL)
		}
	}
}
