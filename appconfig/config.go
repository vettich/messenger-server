package appconfig

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
)

type config struct {
	Port int32    `json:"port"`
	Mode mode     `json:"mode"`
	DB   configDB `json:"db"`
}

type mode string

const (
	production  mode = "production"
	development      = "development"
	staging          = "staging"
)

type configDB struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

var configFilename = flag.String("config", "config.development.json", "JSON config filename")

func init() {
	flag.Parse()

	// считываем конфиг из файла
	dat, err := ioutil.ReadFile(*configFilename)
	if err != nil {
		log.Fatalln("Error open config file:", err)
	}
	if err := json.Unmarshal(dat, &globalConfig); err != nil {
		log.Fatalln("Error unmarshal config data:", err)
	}
}

var globalConfig config

func Config() config {
	return globalConfig
}

type MapConfig map[string]interface{}

func (m MapConfig) Get(key string) interface{} {
	v, ok := m[key]
	if ok {
		return v
	}
	return nil
}
