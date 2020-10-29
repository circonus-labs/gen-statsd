// Copyright © 2020 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package main

import (
	"fmt"

	"github.com/circonus-labs/gen-statsd/internal/release"
)

func main() {

	//Generate the config
	conf := genConfig()

	//Check and see if version is being called
	if conf.version {
		fmt.Printf("%s v%s - commit: %s, date: %s, tag: %s\n", release.NAME, release.VERSION, release.COMMIT, release.DATE, release.TAG)
		return
	}

	//Start the agent controller
	agentController := NewAgentController(conf.agents)
	agentController.Start(conf)

}
