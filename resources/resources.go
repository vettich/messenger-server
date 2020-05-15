// Use `go generate` to pack all *.graphql files under this directory (and sub-directories) into
// a binary format.
//
//go:generate go-bindata -ignore=\.go -pkg=resources -o=bindata.go ./...
package resources

import (
	"bytes"
)

const schemaDir = "schema"

func GetSchema() (string, error) {
	assets, err := AssetDir(schemaDir)
	if err != nil {
		return "", err
	}

	buf := bytes.Buffer{}
	for _, name := range assets {
		b, err := Asset(schemaDir + "/" + name)
		if err != nil {
			return "", err
		}

		buf.Write(b)

		// add a newline if file does not end in a newline
		if len(b) > 0 && b[len(b)-1] != '\n' {
			buf.WriteByte('\n')
		}
	}

	return buf.String(), nil
}
