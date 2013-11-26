package riaken_core

import (
	"testing"
)

import (
	"github.com/riaken/riaken-core/rpb"
)

func TestListKeys(t *testing.T) {
	client := dial()
	defer client.Close()
	session, err := client.Session()
	if err != nil {
		t.Error(err.Error())
	}
	defer session.Close()

	bucket := session.GetBucket("b2")
	object := bucket.Object("o2")
	if _, err := object.Store([]byte("o2-data")); err != nil {
		t.Error(err.Error())
	}

	var keys [][]byte
	// Loop until done is received from Riak
	for out, err := bucket.ListKeys(); !out.GetDone(); out, err = bucket.ListKeys() {
		if err != nil {
			t.Error(err.Error())
			break
		}
		keys = append(keys, out.GetKeys()...)
	}
	if len(keys) > 0 {
		if string(keys[0]) != "o2" {
			t.Errorf("expected: o2, got: %s", keys[0])
		}
	}

	if _, err := object.Delete(); err != nil {
		t.Error(err.Error())
	}
}

func TestSetGetBucketProps(t *testing.T) {
	client := dial()
	defer client.Close()
	session, err := client.Session()
	if err != nil {
		t.Error(err.Error())
	}
	defer session.Close()

	bucket := session.GetBucket("b2")
	tb := true
	ti := uint32(1)
	props := &rpb.RpbBucketProps{
		NVal:      &ti,
		AllowMult: &tb,
	}
	if ok, err := bucket.SetBucketProps(props); !ok {
		t.Error("could not set bucket props")
	} else if err != nil {
		t.Error(err.Error())
	}

	out, err := bucket.GetBucketProps()
	if err != nil {
		t.Error(err.Error())
	}
	if out.GetProps().GetAllowMult() != true {
		t.Errorf("expected: true, got: %t", out.GetProps().GetAllowMult())
	}
}
