package xmlimport

import (
	"encoding/json"
	"testing"
)

func TestImport(t *testing.T) {
	//const testfile = "testdata/test.xml"
	const testfile = "/home/remy/Documents/banque/compta.xml"
	f, err := ImportFile(testfile)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", b)
}
