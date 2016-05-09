// +build !linux,!darwin,!openbsd,!freebsd,!netbsd,!dragonfly

/*
 * Minio Cloud Storage, (C) 2016 Minio, Inc.
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
	"io"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
)

// Return all the entries at the directory dirPath.
func readDir(dirPath string) (entries []string, err error) {
	d, err := os.Open(dirPath)
	if err != nil {
		log.WithFields(logrus.Fields{
			"dirPath": dirPath,
		}).Debugf("Open failed with %s", err)

		// File is really not found.
		if os.IsNotExist(err) {
			return nil, errFileNotFound
		}

		// File path cannot be verified since one of the parents is a file.
		if strings.Contains(err.Error(), "not a directory") {
			return nil, errFileNotFound
		}
		return nil, err
	}
	defer d.Close()

	for {
		fis, err := d.Readdir(1000)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		for _, fi := range fis {
			// Skip special files.
			if hasReservedPrefix(fi.Name()) || hasReservedSuffix(fi.Name()) {
				continue
			}
			if fi.Mode().IsDir() {
				// append "/" instead of "\" so that sorting is done as expected.
				entries = append(entries, fi.Name()+slashSeparator)
			} else if fi.Mode().IsRegular() {
				entries = append(entries, fi.Name())
			}
		}
	}
	return
}
