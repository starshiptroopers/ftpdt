package ftpdt

import (
	"ftpdt/src/ftpdt/ftp"
	"ftpdt/src/ftpdt/uid"
	"goftp.io/server/core"
	"io"
	"log"
	"os"
)

type Ftpdt struct {
	*core.Server
}

type Opts struct {
	FtpOpts         *core.ServerOpts        //goftp server options
	UidGenerator    uid.UID                 //uid validator used to invoke and validate uids from the ftp filepath
	TemplateStorage ftp.TemplateStorage 	//template storage used to invoke templates
	DataStorage     ftp.DataStorage			//data storage
	LogFtpDebug		bool					//do a verbose ftp operations logging
	LogWriter		io.Writer				//default to stdout
}


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