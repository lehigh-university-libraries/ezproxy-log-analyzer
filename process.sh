#!/usr/bin/env bash

set -eou pipefail

cd logs
if [ -f aggregated.log ]; then
  rm aggregated.log
fi

for LOG in $(find . -name "ezproxy.log*" | sort -t. -k3 -rn); do
    grep -v ezproxy.lib.lehigh.edu "$LOG" >> aggregated.log
done

curl -v -XPOST \
  --proxy "" \
  --header "Accept: application/json" \
  --no-buffer \
  --data-binary @aggregated.log \
  -o ezpaarse.json \
  http://localhost:59599
