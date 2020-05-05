package mgosession

import (
	"testing"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func TestNewMgoSession_Set(t *testing.T) {
	var err error

	mgo.SetDebug(true)

	m := NewMgoSession()
	err = m.Init("localhost:27017")
	if err != nil {
		t.Fatal(err)
	}
	err = m.DialWithInfo()
	if err != nil {
		t.Fatal(err)
	}
	m.Set(func() {
		m.SetMode(mgo.Strong, true)
	})
	t.Logf("%#v", m.Session)

	m.Set(nil)
	t.Logf("%#v", m.Session) // consistency:2

	t.Log("done")
}

func TestNewMgoSession_MultiSession(t *testing.T) {
	var err error
	m1 := NewMgoSession()
	err = m1.Dial("localhost:27017")
	if err != nil {
		t.Fatal(err)
	}
	err = m1.Cmd("test1", func(c *mgo.Collection) error {
		return c.Insert(bson.M{"1": 1})
	})
	if err != nil {
		t.Fatal(err)
	}

	m2 := NewMgoSession()
	err = m2.Dial("localhost:27017")
	if err != nil {
		t.Fatal(err)
	}
	err = m2.Cmd("test1", func(c *mgo.Collection) error {
		return c.Insert(bson.M{"1": 2})
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("done")
}

func TestMgoSession_DropTable(t *testing.T) {
	m := NewMgoSession()
	err := m.DialWithTimeout("localhost:27017", 10*time.Second)
	if err != nil {
		t.Fatal(err)
	}

}
