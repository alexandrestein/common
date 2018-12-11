package common_test

import (
	"net"
	"reflect"
	"testing"

	"github.com/alexandrestein/common"
)

func TestAddr(t *testing.T) {
	addr, err := common.NewAddr(2466)
	if err != nil {
		t.Fatal(err)
	}

	addrMust := common.MustNewAddr(2466)
	if !reflect.DeepEqual(addr, addrMust) {
		t.Fatalf("addr and must are not equal: %v and %v", addr, addrMust)
	}

	if broadcast := addr.ForListenerBroadcast(); broadcast != ":2466" {
		t.Fatalf("broadcast address is %q but must be %q", broadcast, ":2466")
	}

	addr.AddAddrAndSwitch("127.0.0.1")
	if main := addr.String(); main != "127.0.0.1:2466" {
		t.Fatalf("main address is %q but must be %q", main, "127.0.0.1:2466")
	}

	if notExist := addr.SwitchMain(9999999); notExist != "" {
		t.Fatalf("try to switch to address %d which is expected to not exist but apparently it does: %s", 9999999, notExist)
	}
	if main := addr.String(); main != "127.0.0.1:2466" {
		t.Fatalf("main address is %q but must be %q", main, "127.0.0.1:2466")
	}

	if !addr.IP().Equal(net.ParseIP("127.0.0.1")) {
		t.Fatalf("IP %q is not equal to IPv4 localhost", addr.IP().String())
	}

	if testing.Verbose() {
		t.Logf("%s address is: %q", addr.Network(), addr.String())
	}
}
