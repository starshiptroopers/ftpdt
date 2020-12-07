// Copyright 2020 The Starship Troopers Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ftp

import (
	"os"
	"time"
)

//implements goftp.FileInfo and os.FileInfo interfaces
type file struct {
	fullname string //full filename
	body     []byte //file file
	created  time.Time
}

func (i *file) Name() string {
	return i.fullname
}

func (i *file) Size() int64 {
	return int64(len(i.body))
}

func (i *file) Mode() os.FileMode {
	return os.ModePerm
}

func (i *file) ModTime() time.Time {
	return i.created
}

func (i *file) IsDir() bool {
	return i.Size() == 0
}

func (i *file) Sys() interface{} {
	return nil
}

func (i *file) Owner() string {
	return "tmpl"
}

func (i *file) Group() string {
	return "tmpl"
}
