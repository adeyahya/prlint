package main

import (
	"github.com/fatih/color"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func PrintYellow(format string, a ...interface{}) (n int, err error) {
	c := color.New(color.FgYellow)
	return c.Printf(format, a...)
}

func PrintGreen(format string, a ...interface{}) (n int, err error) {
	c := color.New(color.FgGreen)
	return c.Printf(format, a...)
}

func PrintRed(format string, a ...interface{}) (n int, err error) {
	c := color.New(color.FgRed)
	return c.Printf(format, a...)
}

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
		// ignore deletion
		if change.From.Name == "" {
			continue
		}
		result = append(result, change.From.Name)
	}
	return result
}
