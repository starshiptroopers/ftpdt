// Copyright 2020 The Starship Troopers Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//tmplStorage implements templates storage with caching features

package example

import (
	"errors"
	"html/template"
)

type DummyTemplateStorage struct {}

func NewDummyTemplateStorage() *DummyTemplateStorage {
	return &DummyTemplateStorage{}

}

func (t *DummyTemplateStorage) Template(id string) (*template.Template, error) {

	if id != "" {
		return nil, errors.New("not found")
	}

	return template.New("default").Parse(`
<!DOCTYPE html>
<!-- This is an example template -->
<html>
<head>
    <title>{{.Title}}</title>
</head>
<body>
    <h1>{{.Caption}}</h1>
<script>
    window.location.href = "{{.Url}}"
</script>
</body>
</html>
`)
}