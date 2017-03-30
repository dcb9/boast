[English](./README.md)

Boast
=========

“我想跟踪所有访问我 Web 服务器的请求及返回数据。”

## 安装

下载最新的[二进制](http://blog.phpor.me/boast)文件

## 使用

```
$ cat .boast.json
{
	"debug_addr": ":8079",
	"list": [
		{ "url": "https://www.baidu.com/", "addr": ":8080" }
		{ "url": "https://github.com/", "addr": ":8081" }
	]
}

$ boast -c .boast.json

$ boast --help
Usage of boast:
  -c string
       config file path (default ".boast.json")
```

## 整体架构

```
HTTP 客户端                   Boast                       Web 服务器
| GET http://localhost:8080/   | 记录请求并进行反向代理      | Response 200 OK
| ---------------------------> | --------------------------> | ------┐
|                              |                             |       |
|                              | 记录返回信息并转发给客户端  |  <----┘
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

## 警告

此产品严禁用于线上产品，只适用于开发、测试环境。

Boast 灵感源于 [ngrok](https://github.com/inconshreveable/ngrok/)。
