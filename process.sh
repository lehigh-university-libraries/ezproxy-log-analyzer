#!/usr/bin/env bash

set -eou pipefail

cd logs
if [ -f aggregated.log ]; then
  rm aggregated.log
fi

find . -type d -mindepth 1 | sort | while read dir; do
    find "$dir" -maxdepth 1 -type f | sort -V -r | while read LOG; do
        grep -v ezproxy.lib.lehigh.edu "$LOG" >> aggregated.log
    done
done

curl -v -XPOST \
  --proxy "" \
  --header "Accept: application/json" \
  --no-buffer \
  --data-binary @aggregated.log \
  -o ezpaarse.json \
  http://localhost:59599
