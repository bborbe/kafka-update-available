// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package update

import (
	"fmt"
	"net/http"

	"github.com/bborbe/kafka-update-available/avro"
	"github.com/bborbe/version"
	"github.com/boltdb/bolt"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

type Handler struct {
	DB                  *bolt.DB
	LatestBucketName    []byte
	InstalledBucketName []byte
}

func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/plain")
	resp.WriteHeader(http.StatusOK)
	err := h.DB.View(func(tx *bolt.Tx) error {
		latestRegistry := LatestRegistry{
			Tx:         tx,
			BucketName: h.LatestBucketName,
		}
		installedRegistry := InstalledRegistry{
			Tx:         tx,
			BucketName: h.InstalledBucketName,
		}
		return installedRegistry.ForEach(func(installed avro.ApplicationVersionInstalled) error {
			latestVersion, err := latestRegistry.Get(installed.App)
			if err != nil {
				return errors.Wrap(err, "get latest version failed")
			}
			if version.Version(installed.Version).Less(version.Version(latestVersion.Version)) {
				fmt.Fprintf(resp, "%s at %s need update from %s to %s\n", installed.App, installed.Url, installed.Version, latestVersion.Version)
			}
			return nil
		})
	})
	if err != nil {
		glog.Warningf("read versions failed: %v", err)
	}
}
