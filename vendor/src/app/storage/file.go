package storage

import (
	"crypto/rand"
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type ConfigFile struct {
	AuthBotToken string  `yaml:"auth_bot_token"`
	ChatID       []int64 `yaml:"chat_id"`
}

//FileStorage implements interface Storage
type FileStorage struct {
	mutex      *sync.Mutex
	file       *os.File
	configFile ConfigFile
}

func (fs FileStorage) AddChatId(chatId int64) error {

	return nil
}

func (fs FileStorage) RemoveChatId(chatId int64) error {

	return nil
}

func (fs FileStorage) LoadAllChatId() ([]int64, error) {

	return nil, nil
}

func (fs FileStorage) GetAuthToken() (string, error) {
	return "", nil
}

func New(path string) (FileStorage, error) {
	var err error
	var fs FileStorage

	if fs.file, fs.configFile, err = openFile(path); err == nil {
		return fs, nil
	}

	return fs, err
}

func randToken() string {
	b := make([]byte, 10)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func defaultConfigFile() ConfigFile {
	return ConfigFile{
		AuthBotToken: randToken(),
	}
}

func openFile(path string) (*os.File, ConfigFile, error) {
	var cf ConfigFile
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return nil, cf, err
		}
		dConfig := defaultConfigFile()

		if err := yaml.NewEncoder(file).Encode(dConfig); err != nil {
			return nil, cf, fmt.Errorf("Error create default yaml config %s", err)
		}

		return file, dConfig, nil
	}

	if err != nil {
		return nil, cf, fmt.Errorf("Error create file %v [%v]", path, err)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, cf, err
	}
	if err := yaml.NewDecoder(file).Decode(&cf); err != nil {
		return nil, cf, fmt.Errorf("Error decode file: %s [%s]", path, err)
	}
	return file, cf, nil
}
