// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package lfs

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"gitea.dev/modules/util"
)

// DownloadCallback gets called for every requested LFS object to process its content
type DownloadCallback func(p Pointer, content io.ReadCloser, objectError error) error

// UploadCallback gets called for every requested LFS object to provide its content
type UploadCallback func(p Pointer, objectError error) (io.ReadCloser, error)

// Client is used to communicate with a LFS source
type Client interface {
	BatchSize() int
	Download(ctx context.Context, objects []Pointer, callback DownloadCallback) error
	Upload(ctx context.Context, objects []Pointer, callback UploadCallback) error
}

// newClient creates a LFS client
func newClient(endpoint *url.URL, httpTransport *http.Transport) Client {
	return newClientWithHeaders(endpoint, httpTransport, nil)
}

func newClientWithHeaders(endpoint *url.URL, httpTransport *http.Transport, headers map[string]string) Client {
	if endpoint.Scheme == "file" {
		return newFilesystemClient(endpoint)
	}
	return newHTTPClientWithHeaders(endpoint, httpTransport, headers)
}

// NewClientFromEndpoint creates a LFS client after resolving its endpoint.
func NewClientFromEndpoint(cloneurl, lfsurl string, httpTransport *http.Transport) (Client, error) {
	return NewClientFromEndpointWithHeaders(cloneurl, lfsurl, httpTransport, nil)
}

// NewClientFromEndpointWithHeaders creates a LFS client and applies headers to requests.
func NewClientFromEndpointWithHeaders(cloneurl, lfsurl string, httpTransport *http.Transport, headers map[string]string) (Client, error) {
	endpoint := DetermineEndpoint(cloneurl, lfsurl)
	if endpoint == nil {
		source := cloneurl
		if lfsurl != "" {
			source = lfsurl
		}
		return nil, fmt.Errorf("unable to determine LFS endpoint from %q", util.SanitizeCredentialURLs(source))
	}
	return newClientWithHeaders(endpoint, httpTransport, headers), nil
}
