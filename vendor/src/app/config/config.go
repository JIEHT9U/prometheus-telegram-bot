package config

import (
	"path/filepath"
)

func Read(paths []string, f func(string) error) error {
	for _, tp := range paths {
		p, err := filepath.Glob(tp)
		if err != nil {
			return err
		}
		if len(p) > 0 {
			for _, path := range p {
				if err := f(path); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
