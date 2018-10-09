package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	DestNets []string `json:"destnets" yaml:"destnets"`
	Nic      string   `json:"nic" yaml:"nic"`
	Ports    []string `json:"ports" yaml:"ports"`
	LogName  string   `json:"logname" yaml:"logname"`
}

func getConf() *Config {
	var c Config
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return &c
}
