package config

import (
	"encoding/json"
	"os"
)

type JsonTypes interface {
	int | bool | string | []string | []int | []bool
}

type Config[T JsonTypes] struct {
	path    string
	rawJson map[string]T
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func Init[T JsonTypes](filePath string) (Config[T], error) {
	if !isExist(filePath) {
		if initErr := initFile(filePath); initErr != nil {
			return Config[T]{}, initErr
		}
		return Config[T]{}, nil
	}

	c, err := os.ReadFile(filePath)
	if err != nil {
		return Config[T]{}, err
	}

	config := Config[T]{}
	config.path = filePath
	config.rawJson = make(map[string]T)
	json.Unmarshal(c, &config.rawJson)
	return config, nil
}

func initFile(filePath string) error {
	_, err := os.Create(filePath)
	if err != nil {
		return err
	}
	return nil
}

func (p *Config[T]) Get(key string) (T, error) {
	c, err := os.ReadFile(p.path)
	if err != nil {
		var i T
		return i, err
	}
	json.Unmarshal(c, &p.rawJson)
	return p.rawJson[key], nil
}

func (p *Config[T]) Set(key string, value T) error {
	p.rawJson[key] = value
	c, err := json.Marshal(&p.rawJson)
	if err != nil {
		return err
	}
	if err := os.WriteFile(p.path, c, 0664); err != nil {
		return err
	}
	return nil
}

func (p *Config[T]) Delete(key string) error {
	delete(p.rawJson, key)
	c, err := json.Marshal(&p.rawJson)
	if err != nil {
		return err
	}
	if err := os.WriteFile(p.path, c, 0664); err != nil {
		return err
	}
	return nil
}

func (p *Config[T]) GetKeyList() ([]string, error) {
	c, err := os.ReadFile(p.path)
	if err != nil {
		return []string{}, err
	}
	json.Unmarshal(c, &p.rawJson)
	list := make([]string, 0)
	for key := range p.rawJson {
		list = append(list, key)
	}
	return list, nil
}
