[中文](./README_zh.md)

Boast
=========

"I want track all requests, and replay it easily."

## Usage

Install the boast package:

`go get https://github.com/dcb9/boast`

After installing, modify your server file

```diff
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/dcb9/boast"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	}))

-	// the old way
-	http.ListenAndServe(":8080", mux)
+	// the boast way
+	server := httptest.NewServer(mux)
+	addr, debugAddr := ":8080", ":8079"
+	boast.Serve(server, addr, debugAddr)
}
```

Then run your server:

`go run server.go`

First request your server as usual, then you can visit debug panel ( http://localhost:8079 )

## Standalone

### Install

Download the latest binary file from the [Releases](https://github.com/dcb9/boast/releases) page.

(darwin_amd64.tar.gz is for Mac OS X users)

### Usage

```
$ cat .boast.json
{
	"debug_addr": ":8079",
	"list": [
		{ "url": "https://www.baidu.com/", "addr": ":8080" },
		{ "url": "https://github.com/", "addr": ":8081" }
	]
}

$ boast -c .boast.json

$ boast --help
Usage of boast:
  -c string
       config file path (default ".boast.json")
```

## Sketch

```
HTTP Client                   Boast                       WebServer
| GET http://localhost:8080/   | Record and Reverse Proxy    | Response 200 OK
| ---------------------------> | --------------------------> | ------┐
|                              |                             |       |
|                              |     Record and Forward      |  <----┘
| <--------------------------- | <-------------------------- |

┌----------------------------------------------------------------------------┐
| url: http://localhost:8081                                                 |
| ---------------------------------------------------------------------------|
| All Transactions         ┌ - - - - - - - - - - - - - - - - - - - - - - - ┐ |
| ----------------------   | time: 10 hours ago  Client: 127.0.0.1         | |
| |GET / 200 OK 100 ms |   |                                               | |
| ----------------------   | Request                      [ Replay ]       | |
|                          | -   -   -   -   -   -   -   -   -   -   -   - | |
|                          | GET http://localhost/ HTTP/1.1                | |
|                          | User-Agent: curl/7.51.0                       | |
|                          | Accept: */*                                   | |
|                          |                                               | |
|                          | Response                                      | |
|                          | -   -   -   -   -   -   -   -   -   -   -   - | |
|                          | HTTP/1.1 200 OK                               | |
|                          | X-Server: HTTPLab                             | |
|                          | Date: Thu, 02 Mar 2017 02:25:27 GMT           | |
|                          | Content-Length: 13                            | |
|                          | Content-Type: text/plain; charset=utf-8       | |
|                          |                                               | |
|                          | Hello, World                                  | |
|                          └ - - - - - - - - - - - - - - - - - - - - - - - ┘ |
|                                                                            |
└----------------------------------------------------------------------------┘
```

## Warning

DO NOT USE ON PRODUCTION!

Boast is heavily inspired by [ngrok](https://github.com/inconshreveable/ngrok/).
