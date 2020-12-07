package tmplstorage

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestTemplateStorage(t *testing.T) {

	dir, err := ioutil.TempDir("", "testing")
	if err != nil {
		t.Fatalf("Can't create temporay directory for testing, %v", err)
		return
	}
	defer func() { _ = os.Remove(dir) }()

	f, err := ioutil.TempFile(dir, "ftpdt_*.tmpl")
	if err != nil {
		t.Fatalf("Can't create temporay file for testing, %v", err)
		return
	}
	defer func() { _ = os.Remove(f.Name()) }()

	storage := New(dir)
	if storage.cache == nil {
		t.Error("Cache isn't defined")
	}

	s := "test__test"
	_, err = f.WriteString(s)
	if err != nil {
		t.Fatalf("Can't write to temporary file %s: %v", f.Name(), err)
		return
	}
	_ = f.Close()

	tmpl, e := storage.Template(strings.TrimPrefix(f.Name(), dir))
	if e != nil {
		t.Errorf("TemplateStorage.Template return error: %v", e)
		return
	}

	buf := bytes.NewBuffer(make([]byte, 0, len([]byte(s))))

	buf.Reset()

	if e := tmpl.Execute(buf, nil); e != nil {
		t.Errorf("Temporary template processing error: %v", e)
		return
	}

	if !bytes.Equal([]byte(s), buf.Bytes()) {
		t.Error("Wrong template processing result")
		return
	}
}
