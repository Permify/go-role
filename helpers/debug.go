package helpers

import (
	"os"

	"github.com/davecgh/go-spew/spew"
)

// Pre exit running project.
// @param interface{}
// @param ...interface{}
func Pre(x interface{}, y ...interface{}) {
	spew.Dump(x)
	os.Exit(1)
}
