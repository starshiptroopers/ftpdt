package example

import (
	"fmt"
	"ftpdt/src/ftpdt"
	"ftpdt/src/ftpdt/uid"
	"goftp.io/server/core"
)

func main() {
	configuration := &ftpdt.Opts{
		FtpOpts: &core.ServerOpts{
			Port: 2000,
			Hostname: "127.0.0.1",
		},
		TemplateStorage: NewDummyTemplateStorage(),
		DataStorage: NewDummyDataStorage(),
		UidGenerator: uid.NewGenerator(
			&uid.Cfg{
				Alfa:      "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
				Format:    "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				Validator: "[0-9a-zA-Z]{32}",
			},
		),
		LogFtpDebug: true,
	}
	ftpd := ftpdt.New( configuration )
	err := ftpd.ListenAndServe()
	if err != nil {
		panic(fmt.Errorf( "хуйня какая-то вышла: %v", err))
	}
}