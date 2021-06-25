// Copyright © 2020 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package release

import "expvar"

// NAME is the name of this application
const NAME = "gen-statsd"

var (
	// COMMIT of release in git repo
	COMMIT = "undef"
	// DATE of release
	DATE = "undef"
	// TAG of release
	TAG = "undef"
	// VERSION of the release
	VERSION = "undef"
)

// Info contains release information
type Info struct {
	Name      string
	Version   string
	Commit    string
	BuildDate string
	Tag       string
}

func init() {
	expvar.Publish("app", expvar.Func(info))
}

func info() interface{} {
	return &Info{
		Name:      NAME,
		Version:   VERSION,
		Commit:    COMMIT,
		BuildDate: DATE,
		Tag:       TAG,
	}
}
