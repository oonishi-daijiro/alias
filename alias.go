package main

import (
	"config"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
)

func main() {
	err := operateByCmdArg()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getConfig() (config.Config[string], error) {
	exepath, exePathErr := os.Executable()
	if exePathErr != nil {
		var i config.Config[string]
		return i, exePathErr
	}
	exeDir := filepath.Dir(exepath)
	return config.Init[string](exeDir + "/config-alias.json")
}

func getAliasValue(c *config.Config[string], aliasName string) (string, error) {
	value, getErr := c.Get(aliasName)
	if getErr != nil {
		return "", getErr
	}
	if value == "" {
		return "", errors.New("No such as alias")
	}
	return value, nil
}

func setAliasValue(c *config.Config[string], aliasName string, value string) error {
	return c.Set(aliasName, value)
}

func removeAlias(aliasName string) error {
	c, configErr := getConfig()
	if configErr != nil {
		return configErr
	}
	return c.Delete(aliasName)
}

func copyAliasValueToClipboard(c *config.Config[string], aliasName string) error {
	value, aliasError := getAliasValue(c, aliasName)
	if aliasError != nil {
		return aliasError
	}
	clipboardError := clipboard.WriteAll(value)
	if clipboardError != nil {
		return clipboardError
	}

	return nil
}

func listAllAliases(c *config.Config[string]) error {
	keys, getKeyErr := c.GetKeyList()
	if getKeyErr != nil {
		return getKeyErr
	}

	maxKeyLength := 0

	for _, key := range keys {
		if len(key) > maxKeyLength {
			maxKeyLength = len(key)
		}
	}

	for _, key := range keys {
		margin := strings.Repeat(" ", maxKeyLength-len(key))
		if v, err := c.Get(key); err == nil {
			fmt.Println(margin, key, ":", v)
		}
	}
	return nil
}

func operateByCmdArg() error {
	c, configErr := getConfig()
	if configErr != nil {
		return configErr
	}
	if len(os.Args) == 1 {
		return errors.New("Please set argument")
	}
	switch os.Args[1] {
	case "add":
		if len(os.Args) > 3 {
			if os.Args[2] == "add" || os.Args[2] == "remove" || os.Args[2] == "list" {
				return errors.New("Cannot use this alias name")
			}
			return setAliasValue(&c, os.Args[2], os.Args[3])
		}
	case "remove":
		if len(os.Args) > 2 {
			return removeAlias(os.Args[2])
		}
	case "list":
		if len(os.Args) > 1 {
			return listAllAliases(&c)
		}
	default:
		if len(os.Args) > 1 {
			return copyAliasValueToClipboard(&c, os.Args[1])
		}
	}
	return nil
}
