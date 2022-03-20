#!/bin/bash

set -eou pipefail

k6 run -d 5s -u 30 benchmarks/k6.js --no-usage-report
