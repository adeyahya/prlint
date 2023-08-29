package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Description string   `yaml:"description"`
	Rules       []string `yaml:"rules"`
	Files       []string `yaml:"files"`
	Envar       string   `yaml:"envar"`
}
type Param struct {
	Config `yaml:",inline"`
}

func main() {
	args := os.Args
	if len(args) < 4 {
		panic(`should have args with format prlint [repo_path] [source_commit_hash] [dest_commit_hash/"HEAD"]`)
	}

	// why not .yaml?
	configFile, err := os.ReadFile("./.prlint")
	if err != nil {
		panic(err)
	}
	confs := make(map[string]Param)
	err = yaml.Unmarshal(configFile, &confs)
	if err != nil {
		panic(err)
	}

	diffFiles := []string{}
	configToValid := map[string]bool{}
	for confKey, conf := range confs {
		configToValid[confKey] = false
		hasFileMatch := false
		if len(conf.Config.Files) > 0 && len(diffFiles) == 0 {
			diffFiles := getDiff(args[1], args[2], args[3])
			for _, confFile := range conf.Config.Files {
				if !hasFileMatch {
					for _, diffFile := range diffFiles {
						match, err := filepath.Match(confFile, diffFile)
						if err != nil {
							panic(err)
						}
						if match {
							hasFileMatch = true
							break
						}
					}
				}

			}
		}

		for _, rule := range conf.Config.Rules {
			match, err := regexp.MatchString(rule, conf.Config.Envar)
			if err != nil {
				panic(err)
			}
			if match {
				configToValid[confKey] = true
			}
		}

	}
	// all params valid

	hasInvalidConfig := false
	for conf, valid := range configToValid {
		if !valid {
			hasInvalidConfig = true
			fmt.Printf("%s - is not valid\n", conf)
		}
	}
	if hasInvalidConfig {
		os.Exit(1)
	}

	fmt.Println("all config is valid")
	os.Exit(0)
}

// git util
func getDiff(repoPath, source, destination string) []string {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		panic(err)
	}
	var destinationHash plumbing.Hash
	var sourceHash plumbing.Hash

	headRef, _ := repo.Head()
	if source == string(plumbing.HEAD) {
		sourceHash = headRef.Hash()
	} else {
		sourceHash = plumbing.NewHash(source)
	}

	if destination == string(plumbing.HEAD) {
		destinationHash = headRef.Hash()
	} else {
		destinationHash = plumbing.NewHash(destination)
	}

	currentCommit, err := repo.CommitObject(destinationHash)
	if err != nil {
		panic(err)
	}
	previousCommit, err := repo.CommitObject(sourceHash)
	if err != nil {
		panic(err)
	}

	currentTree, err := currentCommit.Tree()
	if err != nil {
		panic(err)
	}

	previousTree, err := previousCommit.Tree()
	if err != nil {
		panic(err)
	}
	changes, err := currentTree.Diff(previousTree)
	if err != nil {
		panic(err)
	}

	result := []string{}
	for _, change := range changes {
		result = append(result, change.From.Name)
	}
	return result
}
