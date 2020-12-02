// Copyright 2020 The Starship Troopers Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//dummyStorage - implements the data storage ( wow :) )
package example

import (
	"time"
)

type DummyDataStorage struct {}

func NewDummyDataStorage() *DummyDataStorage {
	return &DummyDataStorage{}
}

func (t *DummyDataStorage) Get(uid string) (payload interface{}, createdAt time.Time, err error) {
	return &struct {
		Title string
		Caption string
		Url string}{"Title", "Caption", "https://starshiptroopers.dev"}, time.Now(), nil
}

func (t *DummyDataStorage) Put(uid string, payload interface{}, ttl *time.Duration) error {
	return nil
}
