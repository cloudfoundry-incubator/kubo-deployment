/*
 * Copyright 2017 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"encoding/json"
	"errors"
	"io"
)

// GCSCli represents the configuration for the gcscli
type GCSCli struct {
	// BucketName is the GCS bucket operations will use.
	BucketName string `json:"bucket_name"`
	// CredentialsSource is the location of a Service Account File.
	// If left empty, Application Default Credentials will be used.
	// If equal to 'none', read-only scope will be used.
	// If equal to 'static', json_key will be used.
	CredentialsSource string `json:"credentials_source"`
	// ServiceAccountFile is the contents of a JSON Service Account File.
	// Required if credentials_source is 'static', otherwise ignored.
	ServiceAccountFile string `json:"json_key"`
	// StorageClass is the type of storage used for objects added to the bucket
	// https://cloud.google.com/storage/docs/storage-classes
	StorageClass string `json:"storage_class"`
	// EncryptionKey is a Customer-Supplied encryption key used to
	// encrypt objects added to the bucket.
	// If left empty, no explicit encryption key will be used;
	// GCS transparently encrypts data using server-side encryption keys.
	// https://cloud.google.com/storage/docs/encryption
	EncryptionKey []byte `json:"encryption_key"`
}

const (
	defaultRegionalLocation          = "US-EAST1"
	defaultMultiRegionalLocation     = "US"
	defaultRegionalStorageClass      = "REGIONAL"
	defaultMultiRegionalStorageClass = "MULTI_REGIONAL"
)

// ApplicationDefaultCredentialsSource specifies that
// Application Default Credentials should be used for authentication.
const ApplicationDefaultCredentialsSource = ""

// NoneCredentialsSource specifies that credentials are explicitly empty
// and that the client should be restricted to a read-only scope.
const NoneCredentialsSource = "none"

// ServiceAccountFileCredentialsSource specifies that a service account file
// included in json_key should be used for authentication.
const ServiceAccountFileCredentialsSource = "static"

// ErrEmptyBucketName is returned when a bucket_name in the config is empty
var ErrEmptyBucketName = errors.New("bucket_name must be set")

// ErrEmptyServiceAccountFile is returned when json_key in the
// config is empty when StaticCredentialsSource is explicitly requested.
var ErrEmptyServiceAccountFile = errors.New("json_key must be set")

// ErrWrongLengthEncryptionKey is returned when a non-nil encryption_key
// in the config is not exactly 32 bytes.
var ErrWrongLengthEncryptionKey = errors.New("encryption_key not 32 bytes")

// getDefaultStorageClass returns the default StorageClass for a given location.
// This takes into account regional/multi-regional incompatibility.
//
// Empty string is returned if the location cannot be matched.
func getDefaultStorageClass(location string) (string, error) {
	if _, ok := GCSMultiRegionalLocations[location]; ok {
		return defaultMultiRegionalStorageClass, nil
	}
	if _, ok := GCSRegionalLocations[location]; ok {
		return defaultRegionalStorageClass, nil
	}
	return "", ErrUnknownLocation
}

// NewFromReader returns the new gcscli configuration struct from the
// contents of the reader.
//
// reader.Read() is expected to return valid JSON.
func NewFromReader(reader io.Reader) (GCSCli, error) {

	dec := json.NewDecoder(reader)
	var c GCSCli
	if err := dec.Decode(&c); err != nil {
		return GCSCli{}, err
	}

	if c.BucketName == "" {
		return GCSCli{}, ErrEmptyBucketName
	}

	if c.CredentialsSource == ServiceAccountFileCredentialsSource &&
		c.ServiceAccountFile == "" {
		return GCSCli{}, ErrEmptyServiceAccountFile
	}

	if len(c.EncryptionKey) != 32 && c.EncryptionKey != nil {
		return GCSCli{}, ErrWrongLengthEncryptionKey
	}

	return c, nil
}

// FitCompatibleLocation returns whether a provided Location
// can have c.StorageClass objects written to it.
//
// When c.StorageClass is empty, a compatible default is filled in.
//
// nil return value when compatible, otherwise a non-nil explanation.
func (c *GCSCli) FitCompatibleLocation(loc string) error {
	if c.StorageClass == "" {
		var err error
		if c.StorageClass, err = getDefaultStorageClass(loc); err != nil {
			return err
		}
	}

	_, regional := GCSRegionalLocations[loc]
	_, multiRegional := GCSMultiRegionalLocations[loc]
	if !(regional || multiRegional) {
		return ErrUnknownLocation
	}

	if _, ok := GCSStorageClass[c.StorageClass]; !ok {
		return ErrUnknownStorageClass
	}

	return validLocationStorageClass(loc, c.StorageClass)
}

// IsReadOnly returns if the configuration is limited to only read operations.
func (c *GCSCli) IsReadOnly() bool {
	return c.CredentialsSource == NoneCredentialsSource
}
