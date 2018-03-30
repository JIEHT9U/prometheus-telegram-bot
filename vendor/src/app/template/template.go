// Copyright 2015 Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package template

import (
	"app/config"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tmplhtml "html/template"

	"github.com/pkg/errors"
)

const (
	DEFAULT_TEMPLATE      = "default.tmpl"
	DEFAULT_TEMPLATE_NAME = "default"
)

type ErrDefaultTempatestruct struct {
	message string
}

func NewDefaultTempate(message string) *ErrDefaultTempatestruct {
	return &ErrDefaultTempatestruct{
		message: message,
	}
}

func (e *ErrDefaultTempatestruct) Error() string {
	return e.message
}

type Tempate struct {
	Templates map[string]*tmplhtml.Template
	Default   *tmplhtml.Template
}

func Load(mapsInstance map[string]string, paths []string) (map[string]*tmplhtml.Template, error) {

	var err error
	var tmps = map[string]*tmplhtml.Template{}

	if tmps[trimExtension(DEFAULT_TEMPLATE)], err = getDefaultTemplate(DEFAULT_TEMPLATE, mapsInstance); err != nil {
		return tmps, fmt.Errorf("Error load default.tmpl ( embedded in code ) [ %s ]", err)
	}

	return tmps, config.Read(paths, func(path string) error {
		stat, err := os.Stat(path)
		if err != nil {
			return errors.Wrap(err, "Not found:"+path)
		}

		templateName := stat.Name()

		if _, ok := tmps[trimExtension(templateName)]; ok {
			return fmt.Errorf("Duplicate template name %s ", templateName)
		}

		if tmps[trimExtension(templateName)], err = tmplhtml.New(templateName).Option("missingkey=zero").Funcs(initFuncMap(mapsInstance)).ParseFiles(path); err != nil {
			return fmt.Errorf("Err create Templates from %s [%s]", path, err)
		}
		return nil

	})
}

func AlignmentPath(paths []string) string {
	var result string
	for _, p := range paths {
		result += "[ " + p + "] "
	}
	return result
}

func trimExtension(tmpName string) string {
	return strings.TrimSuffix(tmpName, filepath.Ext(tmpName))
}

// ExecuteTextString needs a meaningful doc comment (TODO(fabxc)).
func ExecuteTextString(tmpl *tmplhtml.Template, data interface{}) (string, error) {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	return buf.String(), err
}

func Find(tmps map[string]*tmplhtml.Template, templateName string) (*tmplhtml.Template, error) {
	if template, ok := tmps[templateName]; ok {
		return template, nil
	}
	if template, ok := tmps[trimExtension(DEFAULT_TEMPLATE_NAME)]; ok {
		return template, NewDefaultTempate(fmt.Sprintf("Not found template %s will be use default template instedd", templateName))
	}

	return nil, errors.New("Not found template")
}

/*
// FromGlobs calls ParseGlob on all path globs provided and returns the
// resulting Template.
func FromGlobs(paths ...string) (*Template, error) {
	t := &Template{
		text: tmpltext.New("").Option("missingkey=zero"),
		html: tmplhtml.New("").Option("missingkey=zero"),
	}
	var err error

	t.text = t.text.Funcs(tmpltext.FuncMap(DefaultFuncs))
	t.html = t.html.Funcs(tmplhtml.FuncMap(DefaultFuncs))

	b, err := Asset("template/default.tmpl")
	if err != nil {
		return nil, err
	}
	if t.text, err = t.text.Parse(string(b)); err != nil {
		return nil, err
	}
	if t.html, err = t.html.Parse(string(b)); err != nil {
		return nil, err
	}
	for _, tp := range paths {
		// ParseGlob in the template packages errors if not at least one file is
		// matched. We want to allow empty matches that may be populated later on.
		p, err := filepath.Glob(tp)
		if err != nil {
			return nil, err
		}

		log.Println("p:", p)
		if len(p) > 0 {
			if t.text, err = t.text.ParseGlob(tp); err != nil {
				return nil, err
			}
			if t.html, err = t.html.ParseGlob(tp); err != nil {
				return nil, err
			}
		}
	}
	return t, nil
}




// ExecuteHTMLString needs a meaningful doc comment (TODO(fabxc)).
func (t *Template) ExecuteHTMLString(html string, data interface{}) (string, error) {
	if html == "" {
		return "", nil
	}
	tmpl, err := t.html.Clone()
	if err != nil {
		return "", err
	}
	tmpl, err = tmpl.New("").Option("missingkey=zero").Parse(html)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	return buf.String(), err
}
*/
