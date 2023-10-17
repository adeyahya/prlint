package main

import (
	"os"
	"testing"
)

func TestConfigGetErrorString(t *testing.T) {
	type TestCase struct {
		Name  string
		Envar string
		Want  string
		Config
	}
	testCaseList := []TestCase{
		{
			Name:  "Null Envar",
			Envar: "",
			Want:  "Undefined Envar",
			Config: Config{
				Description: "Null Envar",
				Rules:       []string{},
				Envar:       "NO_ENV",
			},
		},
		{
			Name:  "Regex Postive",
			Envar: "chore: setup architecture",
			Want:  "",
			Config: Config{
				Description: "Conventional Commit",
				Rules:       []string{"^(CHORE|docs|feat):"},
				Envar:       "TITLE",
			},
		},
		{
			Name:  "Regex Negative",
			Envar: "CHORES: setup architecture",
			Want:  "^(CHORE|docs|feat):",
			Config: Config{
				Description: "Conventional Commit",
				Rules:       []string{"^(CHORE|docs|feat):"},
				Envar:       "TITLE",
			},
		},
	}

	for _, test := range testCaseList {
		t.Log(test.Name)
		if test.Envar != "" {
			os.Setenv(test.Config.Envar, test.Envar)
		}
		errString := test.Config.GetErrorString()
		if errString != test.Want {
			t.Errorf("want: %s\ngot: %s\n", test.Want, errString)
		}
	}
}
