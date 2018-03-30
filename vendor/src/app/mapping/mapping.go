package mapping

import (
	"app/config"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

//GetName load map in files...
func GetNames(paths []string) (map[string]string, error) {

	var nameMapping = make(map[string]string)

	return nameMapping, config.Read(paths, func(path string) error {

		var decodeMap map[string]string

		f, err := os.Stat(path)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("Error open file %s [%s]", path, err)
		}
		defer file.Close()

		err = yaml.NewDecoder(file).Decode(&decodeMap)

		if err != nil && err != io.EOF {
			return fmt.Errorf("Error yaml decode  %s [%s]", path, err)
		}

		if err != nil && err == io.EOF {
			return nil
		}

		for ip, name := range decodeMap {
			if value, ok := nameMapping[ip]; ok {
				return fmt.Errorf("Duplicate key %s [%s] in file %s", ip, value, f.Name())
			}
			nameMapping[ip] = name
		}

		return nil
	})
}
