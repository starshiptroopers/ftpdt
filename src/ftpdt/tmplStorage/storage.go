// Copyright 2020 The Starship Troopers Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//tmplStorage implements filesystem templates storage with caching features

package tmplStorage

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/cache"
	"html/template"
	"path/filepath"
	"time"
)

var (
	DefaultCacheGCInterval  = 60     //seconds
 	DefaultTmplCacheTTL		= time.Second * time.Duration(86400)  //seconds
 	ERR_NOT_FOUND			= errors.New("template not found")
 	ERR_PARSE_ERROR			= errors.New("template processing error")
)

type TemplateStorage struct {
	fsroot string
	cache	cache.Cache
}

func New(path string) *TemplateStorage {
	rPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	cache, err := cache.NewCache("memory", `{"interval":` + string(DefaultCacheGCInterval) + "}")
	if err != nil {
		panic(err)
	}
	return &TemplateStorage{fsroot: rPath, cache: cache}

}

func (t *TemplateStorage) Template(id string) (*template.Template, error) {

	if id == "" {
		id = "default"
	}
	id += ".tmpl"

	if hit, ok := t.cache.Get(id).(*template.Template); ok {
		return hit, nil
	}

	tPath, err := filepath.Abs(t.fsroot + string(filepath.Separator) + id)
	if err != nil {
		return nil, fmt.Errorf( "%v: %s", ERR_NOT_FOUND, id)
	}

	tmpl, err := template.ParseFiles(tPath)

	if err != nil {
		return nil, fmt.Errorf( "%v: %s", ERR_NOT_FOUND, id)
	}

	t.cache.Put(id, tmpl, DefaultTmplCacheTTL)

	return tmpl, nil
}