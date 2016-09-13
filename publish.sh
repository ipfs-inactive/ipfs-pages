#!/bin/bash
set -eo pipefail; [[ $PAGES_TRACE ]] && set -x

site="$1"
ref="$2"

test -s "$site/key" || (echo "run this first: ./create.sh $site" && exit 1)

echo "pinning (a minority of api mismatch errors is okay)"
for t in `cat targets`; do
  echo "pinning on $t ..."
  out=$(ipfs --api="$t" pin add "$ref" || true)
  printf %s\\n "$out" | grep pinned >>"$site/publish.log" || printf %s\\n "$out"
done

if test -s "$site/ref"; then
  cat "$site/ref" >> "$site/history"
fi
echo "$ref" > "$site/ref"

echo "publishing ..."
out=$(ipns-pub -key="$site/key" "$ref" || true)
printf %s\\n "$out" >>"$site/publish.log"

peerid=$(printf %s\\n "$out" | grep "Local peer ID" | cut -d':' -f2 | xargs)
echo "published https://ipfs.io/ipns/$peerid"
