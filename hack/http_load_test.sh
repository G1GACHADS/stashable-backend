#!/bin/bash

set -ex

echo "GET http://localhost:5000/warehouses?limit=100" | vegeta attack -duration=10s -connections=10000 -rate=2000/1s | tee tmp/results.bin | vegeta report
  vegeta report -type=json tmp/results.bin > tmp/metrics.json
  cat tmp/results.bin | vegeta plot > tmp/plot.html
  cat tmp/results.bin | vegeta report -type="hist[0,100ms,200ms,300ms]"