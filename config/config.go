package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

type RvsPxy struct {
	URL  string `json:"url"`
	Addr string `json:"addr"`
}
type JsonConfig struct {
	DebugAddr string   `json:"debug_addr"`
	List      []RvsPxy `json:"list"`
}

var Config JsonConfig

func init() {
	// http://stackoverflow.com/questions/14249217/how-do-i-know-im-running-within-go-test
	if flag.Lookup("test.v") == nil {
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
}
