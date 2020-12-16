package datastorage

import (
	"testing"
	"time"
)

func TestMemoryDataStorage(t *testing.T) {
	s := NewMemoryDataStorage()
	if s.cache == nil {
		t.Error("Cache isn't defined")
	}

	if s.DefaultCacheTTL != DefaultCacheTTL {
		t.Error("Wrong Cache TTL value")
	}

	p, _, _, e := s.Get("something")
	if e != ErrNFound {
		t.Error("ErrNFound is expected for non-existing element")
	}

	if p != nil {
		t.Error("Element should be a null for non-existing element")
	}

	type datarecord struct {
		body string
	}

	d := datarecord{"somebody"}

	e = s.Put("ELEMENT1", &d, nil)
	if e != nil {
		t.Errorf("Error on Put: %v", e)
	}

	p, c, _, e := s.Get("ELEMENT1")
	if e != nil {
		t.Errorf("A Stored element not found: %v", e)
	}

	dt := c.Sub(time.Now())
	if dt == 0 || dt > time.Second {
		t.Error("Creation timestamp is wrong")
	}

	if p == nil {
		t.Error("Get returns a null element")
		return
	}

	if v, ok := p.(*datarecord); !ok {
		t.Error("The returned element type is wrong")
		return
	} else {
		if v.body != "somebody" {
			t.Error("Get returns wrong data")
		}
	}
}
