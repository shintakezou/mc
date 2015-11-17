/*
 * Minio Client (C) 2015 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"os"
	"strings"

	"github.com/minio/mc/pkg/client"
	"github.com/minio/minio-xl/pkg/probe"
)

// ``...`` recursiveSeparator
const (
	recursiveSeparator = "..."
)

// urlJoinPath Join a path to existing URL.
func urlJoinPath(url1, url2 string) string {
	u1 := client.NewURL(url1)
	u2 := client.NewURL(url2)
	return client.JoinURLs(u1, u2).String()
}

// isURLRecursive - find out if requested url is recursive.
func isURLRecursive(urlStr string) bool {
	return strings.HasSuffix(urlStr, recursiveSeparator)
}

// stripRecursiveURL - Strip "..." from the URL if present.
func stripRecursiveURL(urlStr string) string {
	if !isURLRecursive(urlStr) {
		return urlStr
	}
	urlStr = strings.TrimSuffix(urlStr, recursiveSeparator)
	if urlStr == "" {
		urlStr = "."
	}
	return urlStr
}

// args2URLs extracts source and target URLs from command-line args.
func args2URLs(args []string) ([]string, *probe.Error) {
	config, err := getMcConfig()
	if err != nil {
		return nil, err.Trace()

	}
	// Convert arguments to URLs: expand alias, fix format...
	URLs := []string{}
	for _, arg := range args {
		URLs = append(URLs, getAliasURL(arg, config.Aliases))
	}
	return URLs, nil
}

// url2Stat returns stat info for URL.
func url2Stat(urlStr string) (client client.Client, content *client.Content, err *probe.Error) {
	client, err = url2Client(urlStr)
	if err != nil {
		return nil, nil, err.Trace(urlStr)
	}
	content, err = client.Stat()
	if err != nil {
		return nil, nil, err.Trace(urlStr)
	}
	return client, content, nil
}

// url2DirContent returns directory content info for a URL.
func url2DirContent(urlStr string) (content *client.Content, err *probe.Error) {
	clnt, err := url2Client(urlStr)
	if err != nil {
		return nil, err.Trace(urlStr)
	}
	isRecursive := false
	isIncomplete := false
	for entry := range clnt.List(isRecursive, isIncomplete) {
		if entry.Err != nil {
			return nil, entry.Err.Trace(urlStr)
		}
		if strings.HasPrefix(entry.Content.URL.Path, clnt.GetURL().Path) {
			content := new(client.Content)
			content.URL = clnt.GetURL()
			content.Type = os.ModeDir
			return content, nil
		}
	}
	content = new(client.Content)
	content.URL = clnt.GetURL()
	content.Type = os.ModeDir
	return content, nil
}
