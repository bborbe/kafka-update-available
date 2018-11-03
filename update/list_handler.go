// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package update

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/kafka-update-available/avro"
	"github.com/bborbe/version"
	"github.com/boltdb/bolt"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

type ListHandler struct {
	DB                  *bolt.DB
	LatestBucketName    []byte
	InstalledBucketName []byte
}

func (h *ListHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	list := make([]map[string]string, 0)
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
				list = append(list, map[string]string{
					"app":              installed.App,
					"url":              installed.Url,
					"installedVersion": installed.Version,
					"latestVersion":    latestVersion.Version,
				})
			}
			return nil
		})
	})
	if err != nil {
		glog.Warningf("read versions failed: %v", err)
		http.Error(resp, "read versions failed", http.StatusInternalServerError)
	}
	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(resp).Encode(list); err != nil {
		glog.Warningf("encode json failed: %v", err)
	}
}
