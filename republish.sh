#!/bin/bash
set -eo pipefail; [[ $PAGES_TRACE ]] && set -x

for d in */; do
  site=$(echo "$d" | tr -d /)
  if test -s "$site/ref"; then
    echo "--- $(date --utc --rfc-3339=seconds) republish $site"
    ./publish.sh "$site" "$(cat $site/ref)"
  fi
done
