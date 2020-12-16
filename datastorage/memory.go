// Copyright 2020 The Starship Troopers Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//dataStorage implements the storage ( wow :) )
package datastorage

import (
	"errors"
	"github.com/astaxie/beego/cache"
	"time"
)

var (
	DefaultCacheGCInterval = 60                                 //seconds
	DefaultCacheTTL        = time.Second * time.Duration(86400) //seconds
	ErrNFound              = errors.New("template not found")
)

type MemoryDataStorage struct {
	cache                  cache.Cache
	DefaultCacheGCInterval uint //seconds
	DefaultCacheTTL        time.Duration
}

type dataRecord struct {
	created time.Time
	ttl     time.Duration
	payload interface{}
}

func NewMemoryDataStorage() *MemoryDataStorage {
	c, err := cache.NewCache("memory", `{"interval":`+string(DefaultCacheGCInterval)+"}")
	if err != nil {
		panic(err)
	}
	return &MemoryDataStorage{
		cache:           c,
		DefaultCacheTTL: DefaultCacheTTL,
	}
}

func (t *MemoryDataStorage) Get(uid string) (payload interface{}, createdAt time.Time, ttl time.Duration, err error) {

	r, ok := t.cache.Get(uid).(*dataRecord)
	if !ok {
		err = ErrNFound
		return
	}

	return r.payload, r.created, r.ttl, nil
}

func (t *MemoryDataStorage) Put(uid string, payload interface{}, ttl *time.Duration) error {

	if ttl == nil {
		ttl = &t.DefaultCacheTTL
	}

	return t.cache.Put(uid, &dataRecord{
		created: time.Now(),
		ttl:     *ttl,
		payload: payload,
	}, *ttl)
}
