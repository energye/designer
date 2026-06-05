package dflag

import (
	"os"
	"reflect"
	"testing"
)

func TestParsePositionalsAndSplitN(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()
	os.Args = []string{"energy", "env", "version", "-w=proxy=http://a=b"}
	cmd := New()
	var got Args
	cmd.Add(&Command{
		Name: "env",
		Run: func(args Args) {
			got = args
		},
	})
	cmd.Parse()
	if got.Get("w") != "proxy=http://a=b" {
		t.Fatalf("w got %q", got.Get("w"))
	}
	if !reflect.DeepEqual(got.Positionals(), []string{"version"}) {
		t.Fatalf("positionals got %#v", got.Positionals())
	}
}
