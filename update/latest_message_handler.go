package update

import (
	"bytes"

	"github.com/Shopify/sarama"
	"github.com/bborbe/kafka-update-available/avro"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"github.com/seibert-media/go-kafka/schema"
)

type LatestMessageHandler struct {
	BucketName []byte
}

func (a *LatestMessageHandler) HandleMessage(tx *bolt.Tx, msg *sarama.ConsumerMessage) error {
	buf := bytes.NewBuffer(msg.Value)
	if err := schema.RemoveMagicHeader(buf); err != nil {
		return errors.Wrap(err, "remove magic headers failed")
	}
	newVersion, err := avro.DeserializeApplicationVersionAvailable(buf)
	if err != nil {
		return errors.Wrap(err, "deserialize version failed")
	}
	versionRegistry := LatestRegistry{
		Tx:         tx,
		BucketName: a.BucketName,
	}
	if err := versionRegistry.Set(*newVersion); err != nil {
		return errors.Wrap(err, "save version failed")
	}
	return nil
}
