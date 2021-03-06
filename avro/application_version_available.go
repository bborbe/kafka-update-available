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

type ApplicationVersionAvailable struct {
	App     string
	Version string
}

func DeserializeApplicationVersionAvailable(r io.Reader) (*ApplicationVersionAvailable, error) {
	return readApplicationVersionAvailable(r)
}

func NewApplicationVersionAvailable() *ApplicationVersionAvailable {
	v := &ApplicationVersionAvailable{}

	return v
}

func (r *ApplicationVersionAvailable) Schema() string {
	return "{\"fields\":[{\"name\":\"App\",\"type\":\"string\"},{\"name\":\"Version\",\"type\":\"string\"}],\"name\":\"ApplicationVersionAvailable\",\"type\":\"record\"}"
}

func (r *ApplicationVersionAvailable) Serialize(w io.Writer) error {
	return writeApplicationVersionAvailable(r, w)
}
