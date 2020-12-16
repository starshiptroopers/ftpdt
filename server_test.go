package ftpdt

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jlaffaye/ftp"
	"github.com/phayes/freeport"
	"github.com/starshiptroopers/uidgenerator"
	"goftp.io/server/core"
	"html/template"
	"io/ioutil"
	"strconv"
	"testing"
	"time"
)

type DummyDataStorage struct{}

func NewDummyDataStorage() *DummyDataStorage {
	return &DummyDataStorage{}
}

func (t *DummyDataStorage) Get(uid string) (payload interface{}, createdAt time.Time, ttl time.Duration, err error) {
	return &struct {
		Title   string
		Caption string
		Url     string
	}{"Title", "Caption", "https://starshiptroopers.dev"}, time.Now(), 0, nil
}

func (t *DummyDataStorage) Put(uid string, payload interface{}, ttl *time.Duration) error {
	return nil
}

type DummyTemplateStorage struct{}

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

func downloadFile(server string, path string) ([]byte, error) {
	c, err := ftp.Dial(server, ftp.DialWithTimeout(time.Second))
	if err != nil {
		return nil, fmt.Errorf("Can't connect ftp server: %v, ", err)
	}

	err = c.Login("anonymous", "anonymous")
	if err != nil {
		return nil, fmt.Errorf("Can't login ftp server as anonymous: %v, ", err)
	}

	r, err := c.Retr(path)
	if err != nil {
		return nil, fmt.Errorf("Can't open the remote file for download %s: %v, ", path, err)
	}
	defer func() { _ = r.Close() }()

	buf, err := ioutil.ReadAll(r)

	if err != nil {
		return nil, fmt.Errorf("Can't download the file %s: %v, ", path, err)
	}

	return buf, nil
}

func TestServer(t *testing.T) {

	port, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal("Can't get a free tcp port")
	}
	host := "127.0.0.1"
	path := ""
	uidGenerator := uidgenerator.New(
		&uidgenerator.Cfg{
			Alfa:      "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
			Format:    "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			Validator: "[0-9a-zA-Z]{32}",
		},
	)
	uid := uidGenerator.New()
	filename := path + uid + ".html"

	dStorage := NewDummyDataStorage()
	tStorage := NewDummyTemplateStorage()

	tmpl, err := tStorage.Template(path)
	if tmpl == nil || err != nil {
		t.Fatalf("Null template has been returned for path %s, %v", path, err)
	}

	data, _, _, err := dStorage.Get(uid)

	buff := bytes.NewBuffer(make([]byte, 0))
	err = tmpl.Execute(buff, data)
	if err != nil {
		t.Fatalf("Template execute error: %v", err)
	}
	generated := buff.Bytes()

	configuration := &Opts{
		FtpOpts:         &core.ServerOpts{Port: port, Hostname: host},
		TemplateStorage: tStorage,
		DataStorage:     dStorage,
		UidGenerator:    uidGenerator,
		LogFtpDebug:     true,
	}
	ftpd := New(configuration)

	closeCh := make(chan error)

	go func() {
		err = ftpd.ListenAndServe()
		if err != nil {
			closeCh <- err
		}
		close(closeCh)
	}()

	//wait for server became ready
	select {
	case err := <-closeCh:
		t.Fatalf("Can't start ftp server: %v", err)
	//it's quite unreliable to use timeouts, but, it is simple and reasonable in this case
	case <-time.After(time.Second):

	}

	defer func() {
		_ = ftpd.Shutdown()
		//wait for server shutdown ready
		select {
		case <-closeCh:

		case <-time.After(time.Second):
			t.Fatalf("Server didn't stop after waiting timeout: %v", err)
		}
	}()

	downloaded, err := downloadFile(host+":"+strconv.Itoa(port), filename)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(downloaded, generated) {
		t.Fatal("FTP server returns a wrong file content")
		return
	}
}
