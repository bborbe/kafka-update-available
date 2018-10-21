// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package update

import (
	"bytes"
	"fmt"

	"github.com/bborbe/kafka-update-available/avro"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

type InstalledRegistry struct {
	Tx         *bolt.Tx
	BucketName []byte
}

func (v *InstalledRegistry) Get(url string) (*avro.ApplicationVersionInstalled, error) {
	bucket := v.Tx.Bucket(v.BucketName)
	buf := bytes.NewBuffer(bucket.Get([]byte(url)))
	if buf == nil {
		return nil, fmt.Errorf("not found")
	}
	version, err := avro.DeserializeApplicationVersionInstalled(buf)
	return version, errors.Wrap(err, "deserialize version failed")
}

func (v *InstalledRegistry) Set(version avro.ApplicationVersionInstalled) error {
	key := []byte(version.Url)
	if len(key) == 0 {
		return errors.New("key empty")
	}
	bucket := v.Tx.Bucket(v.BucketName)
	buf := &bytes.Buffer{}
	if err := version.Serialize(buf); err != nil {
		return errors.Wrap(err, "serialize version failed")
	}
	err := bucket.Put(key, buf.Bytes())
	return errors.Wrap(err, "put version failed")
}

func (v *InstalledRegistry) ForEach(fn func(avro.ApplicationVersionInstalled) error) error {
	bucket := v.Tx.Bucket(v.BucketName)
	return bucket.ForEach(func(k, v []byte) error {
		version, err := avro.DeserializeApplicationVersionInstalled(bytes.NewBuffer(v))
		if err != nil {
			return errors.Wrap(err, "deserialize version failed")
		}
		return fn(*version)
	})
}
