// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package update

import (
	"bytes"

	"github.com/Shopify/sarama"
	"github.com/bborbe/kafka-update-available/avro"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"github.com/seibert-media/go-kafka/schema"
)

type InstalledMessageHandler struct {
	BucketName []byte
}

func (a *InstalledMessageHandler) HandleMessage(tx *bolt.Tx, msg *sarama.ConsumerMessage) error {
	buf := bytes.NewBuffer(msg.Value)
	if err := schema.RemoveMagicHeader(buf); err != nil {
		return errors.Wrap(err, "remove magic headers failed")
	}
	newVersion, err := avro.DeserializeApplicationVersionInstalled(buf)
	if err != nil {
		return errors.Wrap(err, "deserialize version failed")
	}
	versionRegistry := InstalledRegistry{
		Tx:         tx,
		BucketName: a.BucketName,
	}
	if err := versionRegistry.Set(*newVersion); err != nil {
		return errors.Wrap(err, "save version failed")
	}
	return nil
}
