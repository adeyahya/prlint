package main

import (
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"
)

type ConfigMap map[string]Config

func (cm *ConfigMap) Parse() {
	configFileList := []string{
		"./.prlint",
		"./.prlint.yml",
		"./.prlint.yaml",
	}
	for idx, configFile := range configFileList {
		configRaw, err := os.ReadFile(configFile)
		if err != nil && idx+1 == len(configFileList) {
			PrintRed("Could not find config file in current directory\n")
			os.Exit(1)
		}

		err = yaml.Unmarshal(configRaw, cm)
		if err != nil {
			panic(err)
		}
	}
}

type Config struct {
	Description string    `yaml:"description"`
	Files       *[]string `yaml:"files"`
	Rules       []string  `yaml:"rules"`
	Envar       string    `yaml:"envar"`
}

func (c *Config) GetErrorString() string {
	env, exist := os.LookupEnv(c.Envar)
	if !exist {
		return "Undefined Envar"
	}

	for _, rule := range c.Rules {
		matched, err := regexp.MatchString("(?i)"+rule, env)
		if err != nil {
			panic(err)
		}
		if !matched {
			return rule
		}
	}

	return ""
}

func (c *Config) IsMatch(files []string) bool {
	if c.Files == nil {
		return true
	}

	for _, pattern := range *c.Files {
		for _, f := range files {
			matched, err := filepath.Match(pattern, f)
			if err != nil {
				panic(err)
			}

			if matched {
				return true
			}
		}
	}

	return false
}
