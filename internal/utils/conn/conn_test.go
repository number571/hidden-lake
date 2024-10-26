package conn

import (
	"testing"

	hiddenlake "github.com/number571/hidden-lake"
)

func TestConn(t *testing.T) {
	t.Parallel()

	if IsAmI(testNewAddr(""), hiddenlake.SConnection{}) != false {
		t.Error(`IsAmI(testNewAddr(""), hiddenlake.SConnection{}) != false`)
		return
	}
	if IsAmI(testNewAddr("localhost:9581"), hiddenlake.SConnection{FHost: "127.0.0.1", FPort: 9582}) != false {
		t.Error(`IsAmI(testNewAddr("localhost:9581"), hiddenlake.SConnection{FHost: "127.0.0.1", FPort: 9582}) != false`)
		return
	}
	if IsAmI(testNewAddr("localhost:9581"), hiddenlake.SConnection{FHost: "127.0.0.1", FPort: 9581}) != true {
		t.Error(`IsAmI(testNewAddr("localhost:9581"), hiddenlake.SConnection{FHost: "127.0.0.1", FPort: 9581}) != true`)
		return
	}
}

var (
	_ IAddress = &sAddress{}
)

type sAddress struct {
	a string
}

func testNewAddr(a string) IAddress {
	return &sAddress{a: a}
}

func (p *sAddress) GetTCP() string {
	return p.a
}
