#!/bin/bash

set -eou pipefail

k6 run -d 2s -u 30 benchmarks/k6.js
