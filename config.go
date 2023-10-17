package main

import (
	"os"
	"regexp"

	"github.com/gobwas/glob"
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
		if err != nil {
			if idx+1 == len(configFileList) {
				PrintRed("Could not find config file in current directory\n%s", err.Error())
				os.Exit(1)
			}
		} else {
			err = yaml.Unmarshal(configRaw, cm)
			if err != nil {
				panic(err)
			}
			break
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
			g, err := glob.Compile(pattern)
			if err != nil {
				PrintRed("invalid file pattern %s, please use valid glob pattern. see https://en.wikipedia.org/wiki/Glob_(programming)", pattern)
				os.Exit(1)
			}
			matched := g.Match(f)

			if matched {
				return true
			}
		}
	}

	return false
}
