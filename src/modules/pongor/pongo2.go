// Copyright 2015 ipfans
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package pongor

import (
	"github.com/labstack/echo"
	"io"
	"path/filepath"
	"sync"

	"gopkg.in/flosch/pongo2.v3"
)

type PongorOption struct {
	// Directory to load templates. Default is "templates"
	Directory string
	// Reload to reload templates everytime.
	Reload bool
}

type Renderer struct {
	PongorOption
	templates map[string]*pongo2.Template
	lock      sync.RWMutex
}

func perparOption(options []PongorOption) PongorOption {
	var opt PongorOption
	if len(options) > 0 {
		opt = options[0]
	}
	if len(opt.Directory) == 0 {
		opt.Directory = "templates"
	}
	return opt
}

func GetRenderer(opt ...PongorOption) *Renderer {
	o := perparOption(opt)
	r := &Renderer{
		PongorOption: o,
		templates:    make(map[string]*pongo2.Template),
	}
	return r
}

func getContext(templateData interface{}) pongo2.Context {
	if templateData == nil {
		return nil
	}
	contextData, isMap := templateData.(map[string]interface{})
	if isMap {
		return contextData
	}
	return nil
}

func (r *Renderer) buildTemplatesCache(name string) (t *pongo2.Template, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	t, err = pongo2.FromFile(filepath.Join(r.Directory, name))
	if err != nil {
		return
	}
	r.templates[name] = t
	return
}

func (r *Renderer) getTemplate(name string) (t *pongo2.Template, err error) {
	if r.Reload {
		return pongo2.FromFile(filepath.Join(r.Directory, name))
	}
	r.lock.RLock()
	var ok bool
	if t, ok = r.templates[name]; !ok {
		r.lock.RUnlock()
		t, err = r.buildTemplatesCache(name)
	} else {
		r.lock.RUnlock()
	}
	return
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	template, err := r.getTemplate(name)
	if err != nil {
		return err
	}
	err = template.ExecuteWriter(getContext(data), w)
	return err
}
