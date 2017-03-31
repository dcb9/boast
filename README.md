[中文](./README_zh.md)

Boast
=========

"I want track all requests, and replay it easily."

## Install

Download the latest binary file from the [Releases](https://github.com/dcb9/boast/releases) page.

(darwin_amd64.tar.gz is for Mac OS X users)

## Usage

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
