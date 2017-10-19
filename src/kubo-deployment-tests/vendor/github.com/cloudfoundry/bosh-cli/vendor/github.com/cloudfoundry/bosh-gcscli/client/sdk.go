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

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/option"

	"net/http"

	"cloud.google.com/go/storage"
	"github.com/cloudfoundry/bosh-gcscli/config"
	"golang.org/x/oauth2/jwt"
)

const uaString = "bosh-gcscli"

// NewSDK returns context and client necessary to instantiate a client
// based off of the provided configuration.
func NewSDK(c config.GCSCli) (context.Context, *storage.Client, error) {
	ctx := context.Background()

	var client *storage.Client
	var err error
	ua := option.WithUserAgent(uaString)
	var opt option.ClientOption
	switch c.CredentialsSource {
	case config.ApplicationDefaultCredentialsSource:
		var tokenSource oauth2.TokenSource
		tokenSource, err = google.DefaultTokenSource(ctx,
			storage.ScopeFullControl)
		if err == nil {
			opt = option.WithTokenSource(tokenSource)
		}
	case config.NoneCredentialsSource:
		opt = option.WithHTTPClient(http.DefaultClient)
	case config.ServiceAccountFileCredentialsSource:
		var token *jwt.Config
		token, err = google.JWTConfigFromJSON([]byte(c.ServiceAccountFile),
			storage.ScopeFullControl)
		if err == nil {
			tokenSource := token.TokenSource(ctx)
			opt = option.WithTokenSource(tokenSource)
		}
	default:
		err = errors.New("unknown credentials_source in configuration")
	}
	if err != nil {
		return ctx, client, err
	}

	client, err = storage.NewClient(ctx, ua, opt)
	return ctx, client, err
}
