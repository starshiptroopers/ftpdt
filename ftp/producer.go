// Copyright 2020 The Starship Troopers Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Producer implements Driver for goftp server framework
//It does a real-time content generation from templates an exposes them to the FTP as downloadable files
package ftp

import (
	"bytes"
	"errors"
	"goftp.io/server/core"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ERR_NOT_SUPPORTED = errors.New("operation isn't supported")
	ERR_WRONG_PATH    = errors.New("wrong path")
	LOG_PREFIX        = "FTPDT "
)

type TemplateStorage interface {
	Template(id string) (*template.Template, error)
}

type DataStorage interface {
	Get(uid string) (payload interface{}, createdAt time.Time, err error)
}

//UID validator
type UID interface {
	//searching the UID in the string
	Validate(string) (string, error)
}

// Driver implements Driver for goftp server framework
type Driver struct {
	ts           TemplateStorage
	ps           DataStorage
	uidGenerator UID
	fileCache    map[string]string
	logger       *log.Logger
}

// Stat return FileInfo for entity located at path
func (d Driver) Stat(filename string) (core.FileInfo, error) {

	p, err := d.produce(filename)

	if err == ERR_WRONG_PATH {
		/*
			we do not pass this error,
			because some ftp clients instead invoking the file directly by its full path
			trying to enter to the each filepath's directory element and Stat() is called for each directory too
			returning the error on this stage can break a ftp clients workflow
		*/
		return &file{
			fullname: filename,
			body:     []byte{},
			created:  time.Now(),
		}, nil
	} else if err != nil {
		d.logger.Printf("%sWARN %s %v", LOG_PREFIX, filename, err)
		return nil, errors.New("file unavailable")
	}

	return p, nil
}

//implements a dummy readerCloser required by goftp driver interface
type readCloser struct {
	f *file
	l *log.Logger
	io.Reader
}

func (rc readCloser) Close() error {
	rc.l.Printf("%sGET %s", LOG_PREFIX, rc.f.Name())
	return nil
}

// GetFile expose the content of a filename as an io.ReadCloser interface
// returns size, io.ReadCloser interface and error on errors
func (d Driver) GetFile(filename string, offset int64) (int64, io.ReadCloser, error) {

	p, err := d.produce(filename)

	if err != nil {
		d.logger.Printf("%sWARN %s %v", LOG_PREFIX, filename, err)
		return 0, nil, errors.New("file unavailable")
	}

	length := p.Size()

	if offset < 0 || offset > length {
		return 0, nil, io.EOF
	}

	rc := readCloser{p, d.logger, bytes.NewReader(p.body[offset:])}
	return length - offset, &rc, nil
}

//parse the file path and invoke template and data ids
func (d Driver) parsePath(path string) (uid string, templateId string, err error) {
	paths := strings.Split(path, string(filepath.Separator))

	//relative paths isn't supported
	if paths[0] == "." || paths[0] == ".." {
		if len(paths) == 1 {
			return "", "", ERR_WRONG_PATH
		}
		paths = paths[1:]
	}

	filename := paths[len(paths)-1]
	if filename == "" {
		return "", "", ERR_WRONG_PATH
	}

	uid, err = d.uidGenerator.Validate(filename)
	if err != nil {
		return "", "", ERR_WRONG_PATH
	}

	templateId = filepath.Join(paths[:len(paths)-1]...)
	return
}

//invoke template and data ids from filepath and generate the file content
func (d Driver) produce(filepath string) (*file, error) {

	uid, templateId, err := d.parsePath(filepath)
	if err != nil {
		return nil, err
	}

	t, err := d.ts.Template(templateId)
	if err != nil {
		return nil, err
	}

	payload, createdAt, err := d.ps.Get(uid)
	if err != nil {
		return nil, err
	}

	//fill the template
	var b bytes.Buffer
	if err = t.Execute(&b, payload); err != nil {
		return nil, err
	}

	return &file{
		fullname: filepath,
		body:     b.Bytes(),
		created:  createdAt,
	}, nil
}

// ListDir defined to satisfy goftp driver interface, but not implemented and always returns an error
func (d Driver) ListDir(string, func(core.FileInfo) error) error {
	return ERR_NOT_SUPPORTED
}

// DeleteDir defined to satisfy goftp driver interface, but not implemented and always returns an error
func (d Driver) DeleteDir(string) error {
	return ERR_NOT_SUPPORTED
}

// DeleteFile defined to satisfy goftp driver interface, but not implemented and always returns an error
func (d Driver) DeleteFile(string) error {
	return ERR_NOT_SUPPORTED
}

// Rename defined to satisfy goftp driver interface, but not implemented and always returns an error
func (d Driver) Rename(string, string) error {
	return ERR_NOT_SUPPORTED
}

// MakeDir defined to satisfy goftp driver interface, but not implemented and always returns an error
func (d Driver) MakeDir(string) error {
	return ERR_NOT_SUPPORTED
}

// PutFile defined to satisfy goftp driver interface, but not implemented and always returns an error
func (d Driver) PutFile(string, io.Reader, bool) (int64, error) {
	return 0, ERR_NOT_SUPPORTED
}

// Implements DriverFactory which creates Driver instance for each ftp client connection
type DriverFactory struct {
	ts           TemplateStorage
	ps           DataStorage
	uidGenerator UID
	logger       *log.Logger
}

// Create Driver instance for each ftp client connection
func (factory *DriverFactory) NewDriver() (core.Driver, error) {
	return &Driver{
		factory.ts,
		factory.ps,
		factory.uidGenerator,
		make(map[string]string),
		factory.logger,
	}, nil
}

//NewDriverFactory create the instance of DriverFactory
func NewDriverFactory(ts TemplateStorage, ps DataStorage, uidGenerator UID, logger *log.Logger) *DriverFactory {

	if ts == nil {
		panic("templateStorage isn't defined")
	}
	if ps == nil {
		panic("payloadStorage isn't defined")
	}

	if uidGenerator == nil {
		panic("uidGenerator isn't defined")
	}

	if logger == nil {
		logger = log.New(os.Stderr, "", log.LstdFlags)
	}
	return &DriverFactory{ts, ps, uidGenerator, logger}
}
