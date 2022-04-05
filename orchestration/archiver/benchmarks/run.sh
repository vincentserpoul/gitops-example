#!/bin/bash

set -eou pipefail

k6 run -d 5s -u 300 benchmarks/k6.js --no-usage-report
