// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/TheNhatAT/e2e"
	e2edb "github.com/TheNhatAT/e2e/db"
	e2einteractive "github.com/TheNhatAT/e2e/interactive"
	e2emonitoring "github.com/TheNhatAT/e2e/monitoring"
	"github.com/efficientgo/core/testutil"
	"github.com/efficientgo/examples/pkg/sum/sumtestutil"
	"github.com/go-kit/log"
	"github.com/thanos-io/objstore/client"
	"github.com/thanos-io/objstore/providers/s3"
	"gopkg.in/yaml.v3"
)

const bktName = "test"

func marshal(t testing.TB, i interface{}) string {
	t.Helper()

	b, err := yaml.Marshal(i)
	testutil.Ok(t, err)

	return string(b)
}

// TestLabeler_LabelObject runs interactive macro benchmark for `labeler` program.
// Prerequisites:
// * `docker` CLI and docker engine installed.
// * Run `make docker` from root project to build `labeler:latest` docker image. TODO: don't forget to del the image in docker before running this test.
// * Run `go test . -v -run TestLabeler_LabelObject` from `pkg/benchmark/macro/labeler` directory to run this test.
// Read more in "Efficient Go"; Example 8-19, 8-20,
func TestLabeler_LabelObject(t *testing.T) {
	//t.Skip("Comment this line if you want to run it - it's interactive test. Won't be useful in CI")

	e, err := e2e.NewDockerEnvironment("labeler")
	testutil.Ok(t, err)
	t.Cleanup(e.Close)

	// Start monitoring.
	mon, err := e2emonitoring.Start(e)
	testutil.Ok(t, err)
	testutil.Ok(t, mon.OpenUserInterfaceInBrowser())

	// Start storage.
	minio := e2edb.NewMinio(e, "object-storage", bktName)
	testutil.Ok(t, e2e.StartAndWaitReady(minio))

	// Run program we want to test and benchmark.
	labeler := e2emonitoring.AsInstrumented(e.Runnable("labeler").
		WithPorts(map[string]int{"http": 8080}).
		Init(e2e.StartOptions{
			Image:     "labeler:test", // `make docker` for building this image.
			LimitCPUs: 4.0,            // limit to 4 CPUs to prevent saturating the local machine.
			Command: e2e.NewCommand(
				"/labeler",
				"-listen-address=:8080",
				"-objstore.config="+marshal(t, client.BucketConfig{
					Type: client.S3,
					Config: s3.Config{
						Bucket:    bktName,
						AccessKey: e2edb.MinioAccessKey,
						SecretKey: e2edb.MinioSecretKey,
						Endpoint:  minio.InternalEndpoint(e2edb.AccessPortName),
						Insecure:  true,
					},
				}),
			),
		}), "http")
	testutil.Ok(t, e2e.StartAndWaitReady(labeler))

	// Add test file.
	testutil.Ok(t, uploadTestInput(minio, "object1.txt", 2e6))

	// Start continuous profiling (not present in examples 8-19, 8-20).
	parca := e2emonitoring.AsInstrumented(e.Runnable("parca").
		WithPorts(map[string]int{"http": 7070}).
		Init(e2e.StartOptions{
			Image: "ghcr.io/parca-dev/parca:main-4e20a666",
			Command: e2e.NewCommand("/bin/sh", "-c",
				`mkdir -p /tmp/shared/data && \
  cat << EOF > /tmp/shared/data/config.yml && \
  /parca --config-path=/tmp/shared/data/config.yml
object_storage:
  bucket:
    type: "FILESYSTEM"
    config:
      directory: "./data"
scrape_configs:
- job_name: "labeler"
  scrape_interval: "15s"
  static_configs:
    - targets: [ '`+labeler.InternalEndpoint("http")+`' ]
  profiling_config:
    pprof_config:
      fgprof:
        enabled: true
        path: /debug/fgprof/profile
        delta: true
EOF
	`),
			User:      strconv.Itoa(os.Getuid()),
			Readiness: e2e.NewTCPReadinessProbe("http"),
		}), "http")
	testutil.Ok(t, e2e.StartAndWaitReady(parca))
	testutil.Ok(t, e2einteractive.OpenInBrowser("http://"+parca.Endpoint("http")))

	// Load test labeler from 1 clients with k6 and export result to Prometheus.
	k6 := e.Runnable("k6").Init(e2e.StartOptions{
		Command: e2e.NewCommandRunUntilStop(),
		Image:   "grafana/k6:0.39.0",
	})
	testutil.Ok(t, e2e.StartAndWaitReady(k6))

	url := fmt.Sprintf(
		"http://%s/label_object?object_id=object1.txt",
		labeler.InternalEndpoint("http"),
	)
	testutil.Ok(t, k6.Exec(e2e.NewCommand( // k6 batch jobs with config: 1 virtual user (-u 1), 5 minutes duration (-d 5m).
		"/bin/sh", "-c",
		`cat << EOF | k6 run -u 1 -d 5m -
import http from 'k6/http';
import { check, sleep } from 'k6';

export default function () {
	const res = http.get('`+url+`');
	check(res, {
		'is status 200': (r) => r.status === 200,
		'response': (r) =>
			r.body.includes('{"object_id":"object1.txt","sum":6221600000,"checksum":"SUUr'),
	});
	sleep(0.5)
}
EOF`)))

	// Once done, wait for user input so user can explore the results in Prometheus UI and logs.
	testutil.Ok(t, e2einteractive.RunUntilEndpointHit())
}

func uploadTestInput(m e2e.Runnable, objID string, numLen int) error {
	bkt, err := s3.NewBucketWithConfig(log.NewNopLogger(), s3.Config{
		Bucket:    bktName,
		AccessKey: e2edb.MinioAccessKey,
		SecretKey: e2edb.MinioSecretKey,
		Endpoint:  m.Endpoint(e2edb.AccessPortName),
		Insecure:  true,
	}, "test")
	if err != nil {
		return err
	}

	b := bytes.Buffer{}
	if _, err := sumtestutil.CreateTestInputWithExpectedResult(&b, numLen); err != nil {
		return err
	}

	return bkt.Upload(context.Background(), objID, &b)
}

/** result of TestLabeler_LabelObject run:
21:10:03 k6-exec: checks.........................: 100.00% ✓ 996      ✗ 0
21:10:03 k6-exec: data_received..................: 113 kB  374 B/s
21:10:03 k6-exec: data_sent......................: 60 kB   199 B/s
21:10:03 k6-exec: http_req_blocked...............: avg=54.33µs  min=5.04µs   med=26.6µs   max=11.35ms  p(90)=48.31µs  p(95)=64.02µs
21:10:03 k6-exec: http_req_connecting............: avg=1.15µs   min=0s       med=0s       max=577.25µs p(90)=0s       p(95)=0s
21:10:03 k6-exec: http_req_duration..............: avg=101.3ms  min=84.13ms  med=100.99ms max=251.1ms  p(90)=107.8ms  p(95)=113.73ms	||=> TODO:NOTE: track latency of the total HTTP request
21:10:03 k6-exec: { expected_response:true }...: avg=101.3ms  min=84.13ms  med=100.99ms max=251.1ms  p(90)=107.8ms  p(95)=113.73ms
21:10:03 k6-exec: http_req_failed................: 0.00%   ✓ 0        ✗ 498
21:10:03 k6-exec: http_req_receiving.............: avg=130.12µs min=75.91µs  med=111.25µs max=3.97ms   p(90)=157.18µs p(95)=176.3µs
21:10:03 k6-exec: http_req_sending...............: avg=141.04µs min=22.37µs  med=108.52µs max=4.93ms   p(90)=183.56µs p(95)=345.44µs
21:10:03 k6-exec: http_req_tls_handshaking.......: avg=0s       min=0s       med=0s       max=0s       p(90)=0s       p(95)=0s
21:10:03 k6-exec: http_req_waiting...............: avg=101.03ms min=83.99ms  med=100.73ms max=250.83ms p(90)=107.57ms p(95)=113.5ms
21:10:03 k6-exec: http_reqs......................: 498     1.656632/s																	||=> TODO:NOTE: track total HTTP requests
21:10:03 k6-exec: iteration_duration.............: avg=603.41ms min=585.44ms med=602.97ms max=753.84ms p(90)=610.25ms p(95)=615.31ms
21:10:03 k6-exec: iterations.....................: 498     1.656632/s																	||=> TODO:NOTE: track total iterations -> more iterations, more reliable
21:10:03 k6-exec: vus............................: 1       min=1      max=1
21:10:03 k6-exec: vus_max........................: 1       min=1      max=1
*/
