// Copyright 2020 The Starship Troopers Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//tmplstorage implements filesystem templates storage with caching features
package tmplstorage

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/cache"
	"html/template"
	"path/filepath"
	"strings"
	"time"
)

var (
	DefaultCacheGCInterval = 60                                 //seconds
	DefaultTmplCacheTTL    = time.Second * time.Duration(86400) //seconds
	ERR_NOT_FOUND          = errors.New("template not found")
	//ERR_PARSE_ERROR			= errors.New("template processing error")
)

//TemplateStorage load, caching and return the templates by their id
type TemplateStorage struct {
	fsroot string
	cache  cache.Cache
}

//New create the instance of TemplateStorage with path pointed to fs root directory where templates are located
func New(path string) *TemplateStorage {
	rPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	c, err := cache.NewCache("memory", `{"interval":`+string(DefaultCacheGCInterval)+"}")
	if err != nil {
		panic(err)
	}
	return &TemplateStorage{fsroot: rPath, cache: c}

}

// Template return the template.Template instance for template id or error.
// Id is a path relative to the storage's root. If id doesn't end with '.tmpl' suffix, it will be added
// If id is empty, "default.tmpl" will be used
func (t *TemplateStorage) Template(id string) (*template.Template, error) {

	if id == "" || id == "/" {
		id = "default"
	}
	if !strings.HasSuffix(id, ".tmpl") {
		id += ".tmpl"
	}

	if hit, ok := t.cache.Get(id).(*template.Template); ok {
		return hit, nil
	}

	tPath, err := filepath.Abs(t.fsroot + string(filepath.Separator) + id)

	//preventing access outside the root folder
	if err != nil || !strings.HasPrefix(tPath, t.fsroot) {
		return nil, fmt.Errorf("%v: %s", ERR_NOT_FOUND, id)
	}

	tmpl, err := template.ParseFiles(tPath)

	if err != nil {
		return nil, fmt.Errorf("%v: %s", ERR_NOT_FOUND, id)
	}

	_ = t.cache.Put(id, tmpl, DefaultTmplCacheTTL)

	return tmpl, nil
}
