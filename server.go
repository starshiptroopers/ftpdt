// Copyright 2020 The Starship Troopers Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// ftpdt is the Ftp server library based on goftp/server and used to do a real-time files generation from templates and exposes them as ftp downloadable files
// It was originally designed to implement the trick of escaping from Facebook's (Instagram and so-on) in-app browsers on IOS.
// Ftp is used as an intermediate gateway to download the html file with redirect to real web site.
// The main reason the ftp is using in this trick is because the ftp protocol is handled with Safary by default.
// Opening a ftp link in a webkit lead to starting the Safari. At the moment of developing this library it was the only way to escape from Facebook browser to ios default browser.
//
// You can realize your own strategy of storing templates and a data by implementing DataStorage and TemplateStorage interfaces
//
// ftpdt/tmplstorage implements a Templates storage where templates is loaded from files
// The requested ftp file is mapping to the template by its full path excluding the filename itself.
// Then template is executed and filled with data and exposed to the ftp client as a regular file.
// For example, if user download the file ftp://servername/example/redirect/abcde.txt, example/redirect.tmpl will be used as a template
// and the filename "abcde" is a UID used to invoke data from DataStorage and insert them to the template.
//
// ftpdt/datastorage implements a Data storage where data are stored in the memory and invoked by ftpdt request

package ftpdt

import (
	"github.com/starshiptroopers/ftpdt/ftp"
	"github.com/starshiptroopers/uidgenerator"
	"goftp.io/server/core"
	"io"
	"log"
	"os"
)

type Ftpdt struct {
	*core.Server
}

// Opts is a ftpdt options
type Opts struct {
	FtpOpts         *core.ServerOpts    //goftp server options
	UidGenerator    uidgenerator.UID    //uid validator used to invoke and validate uids from the ftp filepath
	TemplateStorage ftp.TemplateStorage //template storage used to invoke templates
	DataStorage     ftp.DataStorage     //data storage
	LogFtpDebug     bool                //do a verbose ftp operations logging
	LogWriter       io.Writer           //Where log will be written to (default to stdout)
}

//Create a new Ftpdt instanse
func New(opts *Opts) (server *Ftpdt) {
	if opts.UidGenerator == nil {
		panic("Uid validator isn't defined")
	}

	if opts.TemplateStorage == nil {
		panic("TemplateStorage isn't defined")
	}

	if opts.DataStorage == nil {
		panic("DataStorage isn't defined")
	}

	if opts.LogWriter == nil {
		opts.LogWriter = os.Stdout
	}

	ftpCfg := *opts.FtpOpts

	if ftpCfg.Auth == nil {
		ftpCfg.Auth = &ftp.AuthAnonymous{}
	}

	if ftpCfg.Logger == nil {
		if opts.LogFtpDebug {
			ftpCfg.Logger = ftp.NewDefaultFTPLogger(opts.LogWriter)
		} else {
			ftpCfg.Logger = &core.DiscardLogger{}
		}
	}

	if ftpCfg.Factory == nil {
		ftpCfg.Factory = ftp.NewDriverFactory(
			opts.TemplateStorage,
			opts.DataStorage,
			opts.UidGenerator,
			log.New(opts.LogWriter, "", log.LstdFlags),
		)
	}

	server = &Ftpdt{core.NewServer(&ftpCfg)}
	return
}

//ListenAndServe starts listening for ftp connection. It's blocking function
func (ftpdt *Ftpdt) ListenAndServe() error {
	return ftpdt.Server.ListenAndServe()
}
