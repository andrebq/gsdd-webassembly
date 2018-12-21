package server

import (
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/andrebq/gsdd-webassembly/hchan"
)

func TestReadWrite(t *testing.T) {
	srv := httptest.NewServer(MakeRWHandler())

	type randomData struct {
		Name string
	}

	t.Logf("Server address: %v", srv.URL)
	in := randomData{Name: "something"}

	err := hchan.Write(srv.URL+"/write/super-unique-channel-name", in)
	if err != nil {
		t.Fatal(err)
	}

	var out randomData
	err = hchan.Read(&out, srv.URL+"/read/super-unique-channel-name")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(in, out) {
		t.Fatalf("Should get %v got %v", in, out)
	}

}
