// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/newrelic/go-agent/v3/stress-test/benchmark"
)

func main() {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("ApplicationLogging Stress Test Golang"),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		newrelic.ConfigInfoLogger(os.Stdout),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tests := []benchmark.Benchmark{
		benchmark.Zerolog(10, 10),
		benchmark.Zerolog(100, 10),
		benchmark.Zerolog(1000, 10),

		benchmark.NRZerolog(10, 10),
		benchmark.NRZerolog(100, 10),
		benchmark.NRZerolog(1000, 10),

		benchmark.CustomEvent(10, 10),
		benchmark.CustomEvent(100, 10),
		benchmark.CustomEvent(1000, 10),
	}

	// Wait for the application to connect.
	if err := app.WaitForConnection(5 * time.Second); nil != err {
		fmt.Println(err)
	}

	for _, test := range tests {
		test.Benchmark(app)
	}

	// Make sure the metrics get sent
	time.Sleep(60 * time.Second)
	// Shut down the application to flush data to New Relic.
	app.Shutdown(10 * time.Second)

	var metrics string
	for _, test := range tests {
		metrics += test.Sprint()
	}

	fmt.Println(metrics)
}
