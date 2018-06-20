package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
)

type RvsPxy struct {
	URL  string `json:"url"`
	Addr string `json:"addr"`
}

var DefaultDBPath = "boast.db"

type JsonConfig struct {
	DebugAddr string   `json:"debug_addr"`
	List      []RvsPxy `json:"list"`
	// it'll persist to DB, if DBPath is not empty
	DBPath    string `json:"db_path"`
}

var Config JsonConfig

func CmdInit() {
	filePath := flag.String("c", ".boast.json", "config file path")
	flag.Parse()

	if _, err := os.Stat(*filePath); os.IsNotExist(err) {
		log.Fatal("config file: ", *filePath, " is not exist.")
	}

	if b, err := ioutil.ReadFile(*filePath); err != nil {
		log.Fatal("Read config error: ", err)
	} else {
		err := json.Unmarshal(b, &Config)
		if err != nil {
			log.Fatal("Parse json config error: ", err)
		}
	}
}

func Init(s *httptest.Server, addr, debugAddr string) {
	Config.DebugAddr = debugAddr
	Config.List = []RvsPxy{
			RvsPxy{
				s.URL, addr,
			},
	}
}

func SetDBPath(path string) {
	Config.DBPath = path
}
