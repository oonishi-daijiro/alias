package main

import (
	"config"
	"errors"
	"fmt"
	"os"
	"path"

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
	exeDir := path.Dir(exepath)
	return config.Init[string](exeDir + "/config-alias.json")
}

func getAliasValue(aliasName string) (string, error) {
	c, configErr := getConfig()
	if configErr != nil {
		return "", configErr
	}
	value, getErr := c.Get(aliasName)
	if getErr != nil {
		return "", getErr
	}
	if value == "" {
		return "", errors.New("No such as alias")
	}
	return value, nil
}

func setAliasValue(aliasName string, value string) error {
	c, configErr := getConfig()
	if configErr != nil {
		return configErr
	}
	return c.Set(aliasName, value)
}

func removeAlias(aliasName string) error {
	c, configErr := getConfig()
	if configErr != nil {
		return configErr
	}
	return c.Delete(aliasName)
}

func copyAliasValueToClipboard(aliasName string) error {
	value, aliasError := getAliasValue(aliasName)
	if aliasError != nil {
		return aliasError
	}
	clipboardError := clipboard.WriteAll(value)
	if clipboardError != nil {
		return clipboardError
	}

	return nil
}

func listAllAliases() error {
	c, err := getConfig()
	if err != nil {
		return err
	}
	keys, getKeyErr := c.GetKeyList()
	if getKeyErr != nil {
		return getKeyErr
	}
	for _, key := range keys {
		v, _ := c.Get(key)
		fmt.Println(key, ":", v)
	}
	return nil
}

func operateByCmdArg() error {
	if len(os.Args) == 1 {
		return errors.New("Please set argument")
	}
	switch os.Args[1] {
	case "add":
		if len(os.Args) > 3 {
			if os.Args[2] == "add" || os.Args[2] == "remove" || os.Args[2] == "list" {
				return errors.New("Cannot use this alias name")
			}
			return setAliasValue(os.Args[2], os.Args[3])
		}
	case "remove":
		if len(os.Args) > 2 {
			return removeAlias(os.Args[2])
		}
	case "list":
		if len(os.Args) > 1 {
			return listAllAliases()
		}
	default:
		if len(os.Args) > 1 {
			return copyAliasValueToClipboard(os.Args[1])
		}
	}
	return nil
}
