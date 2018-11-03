// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package update

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/bborbe/run"
	"github.com/boltdb/bolt"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/seibert-media/go-kafka/persistent"
)

type App struct {
	DataDir                    string
	KafkaInstalledVersionTopic string
	KafkaBrokers               string
	KafkaLatestVersionTopic    string
	KafkaSchemaRegistryUrl     string
	Port                       int
}

func (a *App) Validate() error {
	if a.KafkaBrokers == "" {
		return errors.New("KafkaBrokers missing")
	}
	if a.KafkaInstalledVersionTopic == "" {
		return errors.New("KafkaInstalledVersionTopic missing")
	}
	if a.KafkaLatestVersionTopic == "" {
		return errors.New("KafkaLatestVersionTopic missing")
	}
	if a.KafkaSchemaRegistryUrl == "" {
		return errors.New("SchemaRegistryUrl missing")
	}
	if a.DataDir == "" {
		return errors.New("DataDir missing")
	}
	if a.Port <= 0 {
		return errors.New("Port invalid")
	}
	return nil
}

func (a *App) Run(ctx context.Context) error {
	db, err := bolt.Open(path.Join(a.DataDir, "kafka-update-available.db"), 0600, nil)
	if err != nil {
		return errors.Wrap(err, "open bolt db failed")
	}
	defer db.Close()

	installedVersionsConsumer, err := a.createInstalledVersionsConsumer(db)
	if err != nil {
		return errors.Wrap(err, "create consumer failed")
	}
	availableVersionsConsumer, err := a.createLatestVersionConsumer(db)
	if err != nil {
		return errors.Wrap(err, "create consumer failed")
	}

	runHttpServer := func(ctx context.Context) error {
		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", a.Port),
			Handler: a.createHttpHandler(db),
		}
		go func() {
			select {
			case <-ctx.Done():
				if err := server.Shutdown(ctx); err != nil {
					glog.Warningf("shutdown failed: %v", err)
				}
			}
		}()
		return server.ListenAndServe()
	}

	return run.CancelOnFirstFinish(ctx, installedVersionsConsumer.Consume, availableVersionsConsumer.Consume, runHttpServer)
}

func (a *App) createHttpHandler(db *bolt.DB) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/healthz", a.check)
	router.HandleFunc("/readiness", a.check)
	router.Handle("/metrics", promhttp.Handler())
	router.Handle("/updates", &ListHandler{
		DB:                  db,
		LatestBucketName:    []byte("version_latest"),
		InstalledBucketName: []byte("version_installed"),
	})
	router.Handle("/", &IndexHandler{})
	return router
}

func (a *App) check(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusOK)
}

func (a *App) createInstalledVersionsConsumer(db *bolt.DB) (*persistent.Consumer, error) {
	offsetBucketName := []byte("offset_installed")
	installedBucketName := []byte("version_installed")

	err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(installedBucketName); err != nil {
			return fmt.Errorf("create bucket failed: %s", err)
		}
		glog.V(2).Infof("bucket version created")
		if _, err := tx.CreateBucketIfNotExists(offsetBucketName); err != nil {
			return fmt.Errorf("create bucket failed: %s", err)
		}
		glog.V(2).Infof("bucket offset created")
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "create default bolt buckets failed")
	}

	return &persistent.Consumer{
		KafkaTopic:   a.KafkaInstalledVersionTopic,
		KafkaBrokers: a.KafkaBrokers,
		OffsetManager: &persistent.MessageHandler{
			DB: db,
			MessageHandler: &InstalledMessageHandler{
				BucketName: installedBucketName,
			},
			OffsetBucketName: offsetBucketName,
		},
	}, nil
}

func (a *App) createLatestVersionConsumer(db *bolt.DB) (*persistent.Consumer, error) {
	offsetBucketName := []byte("offset_available")
	latestBucketName := []byte("version_latest")

	err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(latestBucketName); err != nil {
			return fmt.Errorf("create bucket failed: %s", err)
		}
		glog.V(2).Infof("bucket version created")
		if _, err := tx.CreateBucketIfNotExists(offsetBucketName); err != nil {
			return fmt.Errorf("create bucket failed: %s", err)
		}
		glog.V(2).Infof("bucket offset created")
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "create default bolt buckets failed")
	}

	return &persistent.Consumer{
		KafkaTopic:   a.KafkaLatestVersionTopic,
		KafkaBrokers: a.KafkaBrokers,
		OffsetManager: &persistent.MessageHandler{
			DB: db,
			MessageHandler: &LatestMessageHandler{
				BucketName: latestBucketName,
			},
			OffsetBucketName: offsetBucketName,
		},
	}, nil
}
