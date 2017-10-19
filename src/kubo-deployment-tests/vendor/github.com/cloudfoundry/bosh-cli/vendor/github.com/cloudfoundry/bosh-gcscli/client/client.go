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

package client

import (
	"context"
	"errors"
	"fmt"
	"io"

	"log"

	"cloud.google.com/go/storage"
	"github.com/cloudfoundry/bosh-gcscli/config"
)

// ErrInvalidROWriteOperation is returned when credentials associated with the
// client disallow an attempted write operation.
var ErrInvalidROWriteOperation = errors.New("the client operates in read only mode. Change 'credentials_source' parameter value ")

// GCSBlobstore encapsulates interaction with the GCS blobstore
type GCSBlobstore struct {
	// gcsClient is a pre-configured storage.Client.
	client *storage.Client
	// gcscliConfig is the configuration for interactions with the blobstore
	config *config.GCSCli
}

// checkLocation determines if the configured StorageClass of the
// GCSBlobstore is compatible with the configured bucket's location.
func (client *GCSBlobstore) checkLocation() error {
	bucket := client.client.Bucket(client.config.BucketName)
	attrs, err := bucket.Attrs(context.Background())
	if err != nil {
		return err
	}
	return client.config.FitCompatibleLocation(attrs.Location)
}

// validateRemoteConfig determines if the configuration of the client matches
// against the remote configuration.
//
// If operating in read-only mode, no mutations can be performed
// so the remote bucket location is always compatible.
func (client *GCSBlobstore) validateRemoteConfig() error {
	if client.config.IsReadOnly() {
		return nil
	}
	return client.checkLocation()
}

// getObjectHandle returns a handle to an object at src.
func (client GCSBlobstore) getObjectHandle(src string) *storage.ObjectHandle {
	handle := client.client.Bucket(client.config.BucketName).Object(src)
	if client.config.EncryptionKey != nil {
		handle = handle.Key(client.config.EncryptionKey)
	}
	return handle
}

// New returns a BlobstoreClient configured to operate using the given config
// and client.
//
// non-nil error is returned on invalid client or config. If the configuration
// is incompatible with the GCS bucket, a non-nil error is also returned.
func New(ctx context.Context, gcsClient *storage.Client,
	gcscliConfig *config.GCSCli) (GCSBlobstore, error) {
	if gcsClient == nil {
		return GCSBlobstore{},
			errors.New("nil client causes invalid blobstore")
	}
	if gcscliConfig == nil {
		return GCSBlobstore{},
			errors.New("nil config causes invalid blobstore")
	}
	blobstore := GCSBlobstore{gcsClient, gcscliConfig}
	return blobstore, blobstore.checkLocation()
}

// Get fetches a blob from the GCS blobstore.
// Destination will be overwritten if it already exists.
func (client GCSBlobstore) Get(src string, dest io.Writer) error {
	remoteReader, err := client.getObjectHandle(src).NewReader(context.Background())
	if err != nil {
		return err
	}
	_, err = io.Copy(dest, remoteReader)
	return err
}

// Put uploads a blob to the GCS blobstore.
// Destination will be overwritten if it already exists.
//
// Put does not retry if upload fails. This is a change from s3cli/client
// which does retry an upload multiple times.
// TODO: implement retry
func (client GCSBlobstore) Put(src io.ReadSeeker, dest string) error {
	if client.config.IsReadOnly() {
		return ErrInvalidROWriteOperation
	}

	remoteWriter := client.getObjectHandle(dest).NewWriter(context.Background())
	remoteWriter.ObjectAttrs.StorageClass = client.config.StorageClass
	if _, err := io.Copy(remoteWriter, src); err != nil {
		log.Println("Upload failed", err.Error())
		return fmt.Errorf("upload failure: %s", err.Error())
	}
	return remoteWriter.Close()
}

// Delete removes a blob from from the GCS blobstore.
//
// If the object does not exist, Delete returns a nil error.
func (client GCSBlobstore) Delete(dest string) error {
	if client.config.IsReadOnly() {
		return ErrInvalidROWriteOperation
	}

	err := client.getObjectHandle(dest).Delete(context.Background())
	if err == storage.ErrObjectNotExist {
		return nil
	}
	return err
}

// Exists checks if a blob exists in the GCS blobstore.
func (client GCSBlobstore) Exists(dest string) (bool, error) {
	_, err := client.getObjectHandle(dest).Attrs(context.Background())
	if err == nil {
		log.Printf("File '%s' exists in bucket '%s'\n",
			dest, client.config.BucketName)
		return true, nil
	} else if err == storage.ErrObjectNotExist {
		log.Printf("File '%s' does not exist in bucket '%s'\n",
			dest, client.config.BucketName)
		return false, nil
	}
	return false, err
}
