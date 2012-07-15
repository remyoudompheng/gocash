package xmlimport

import (
	"encoding/json"
	"testing"
)

func TestRead(t *testing.T) {
	const testfile = "testdata/test.xml"
	f, err := ReadFile(testfile)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", b)
}

func TestImport(t *testing.T) {
	const testfile = "testdata/test.xml"
	book, err := ImportFile(testfile)
	if err != nil {
		t.Fatal(err)
	}
	js, err := json.MarshalIndent(book, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", js)
}
