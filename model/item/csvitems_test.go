package item

import (
	"encoding/csv"
	"io"
	"strings"
	"testing"

	"github.com/jaguilar/testify/assert"
)

func TestCsvItems(t *testing.T) {
	data := "Foo,Bar,Blaz\n1,2,3"
	it, err := newCsvMapReader(csv.NewReader(strings.NewReader(data)))
	if !assert.Nil(t, err) {
		return
	}
	for it.next() {
		assert.Equal(t, "1", it.get("Foo"))
		assert.Equal(t, "2", it.get("Bar"))
		assert.Equal(t, "3", it.get("Blaz"))
	}
	assert.Equal(t, io.EOF, it.err)
}
