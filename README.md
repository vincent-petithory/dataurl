# Data URL Schemes for Go [![GoDoc](https://godoc.org/github.com/vincent-petithory/dataurl?status.png)](https://godoc.org/github.com/vincent-petithory/dataurl)

This package parses and generates Data URL Schemes for the Go language, according to [RFC 2397](http://tools.ietf.org/html/rfc2397).

Data URLs are small chunks of data commonly used in browsers to display inline data,
typically like small images, or when you use the FileReader API of the browser.

Common use-cases:

 * generate a data URL out of a `string`, `[]byte`, `io.Reader` for inclusion in HTML templates,
 * parse a data URL sent by a browser in a http.Handler, and do something with the data (save to disk, etc.)
 * ...

Install the package with:
~~~
go get github.com/vincent-petithory/dataurl
~~~

## Usage

~~~ go
package main

import (
	"github.com/vincent-petithory/dataurl"
)

func main() {
	dataURL, err := dataurl.DecodeString(`data:text/plain;charset=utf-8;base64,aGV5YQ==`)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("content type: %s, data: %s\n", dataURL.MediaType.ContentType(), string(dataURL.Data))
	// Output: content type: text/plain, data: heya
}
~~~

From a `http.Handler`:

~~~ go
func handleDataURLUpload(w http.ResponseWriter, r *http.Request) {
	dataURL, err := dataurl.Decode(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.BadRequest)
		return
	}
	if dataURL.ContentType() == "image/png" {
		ioutil.WriteFile("image.png", dataURL.Data, 0644)
	} else {
		http.Error(w, err.Error(), http.BadRequest)
	}

}
~~~

## Contributing

Feel to file an issue/make a pull request if you find any bug, or want to suggest enhancements.
