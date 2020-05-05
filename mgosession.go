package mgosession

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

type MgoSession struct {
	*mgo.Session
	*mgo.DialInfo
	*sync.Once
}

func NewMgoSession() *MgoSession {
	return &MgoSession{
		Once: &sync.Once{},
	}
}

func (m *MgoSession) Init(url string) (err error) {
	var info *mgo.DialInfo

	m.Do(func() {
		info, err = mgo.ParseURL(url)
		if err != nil {
			m.Once = new(sync.Once)
			log.Printf("MgoSession ParseURL %s error: %s", url, err)
			return
		}
		//m.DialInfo.Database // 如果数据库名为空, mgo驱动将默认使用 "test" 库
		m.DialInfo = info
	})
	return
}

func (m *MgoSession) Set(f func()) {
	if f == nil {
		m.Session.SetMode(mgo.Monotonic, true)
	} else {
		f()
	}
}

func (m *MgoSession) Dial(url string) error {
	var err error
	if err = m.Init(url); err != nil {
		return err
	}
	m.Session, err = mgo.Dial(url)
	if err != nil {
		return err
	}
	m.Set(nil)
	return nil
}

func (m *MgoSession) DialWithTimeout(url string, timeout time.Duration) error {
	var err error
	if err = m.Init(url); err != nil {
		return err
	}
	m.Session, err = mgo.DialWithTimeout(url, timeout)
	if err != nil {
		return err
	}
	m.Set(nil)
	return nil
}

// 调 Init 后调 DialWithInfo
func (m *MgoSession) DialWithInfo() (err error) {
	if m.DialInfo == nil {
		panic("Init must be called first")
	}
	m.Session, err = mgo.DialWithInfo(m.DialInfo)
	if err != nil {
		return err
	}
	m.Set(nil)
	return nil
}

func (m *MgoSession) Close() {
	m.Session.Close()
}

func (m *MgoSession) conn() *mgo.Session {
	return m.Session.Clone()
}

func (m *MgoSession) Cmd(collection string, f func(*mgo.Collection) error) (err error) {
	session := m.conn()
	defer func() {
		session.Close()
		if _err := recover(); _err != nil {
			err = fmt.Errorf("%v", _err)
			return
		}
	}()
	c := session.DB(m.Database).C(collection)
	return f(c)
}
