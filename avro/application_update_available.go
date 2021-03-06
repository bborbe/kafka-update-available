// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
/*
 * SOURCES:
 *     application_update_available.avsc
 *     application_version_available.avsc
 *     application_version_installed.avsc
 */

package avro

import (
	"io"
)

type ApplicationUpdateAvailable struct {
	App     string
	Version string
	Url     string
}

func DeserializeApplicationUpdateAvailable(r io.Reader) (*ApplicationUpdateAvailable, error) {
	return readApplicationUpdateAvailable(r)
}

func NewApplicationUpdateAvailable() *ApplicationUpdateAvailable {
	v := &ApplicationUpdateAvailable{}

	return v
}

func (r *ApplicationUpdateAvailable) Schema() string {
	return "{\"fields\":[{\"name\":\"App\",\"type\":\"string\"},{\"name\":\"Version\",\"type\":\"string\"},{\"name\":\"Url\",\"type\":\"string\"}],\"name\":\"ApplicationUpdateAvailable\",\"type\":\"record\"}"
}

func (r *ApplicationUpdateAvailable) Serialize(w io.Writer) error {
	return writeApplicationUpdateAvailable(r, w)
}
