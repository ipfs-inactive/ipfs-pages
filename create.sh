#!/bin/bash
set -eo pipefail; [[ $PAGES_TRACE ]] && set -x

site="$1"
initRef="QmSAGVJed32haEhK7LJPL57PFjA6yPHnxbsRbdQPVSgWQu"
keySize=4096

mkdir -p "$site/"
chmod 700 "$site/"
echo "created $site/"

if ! test -s "$site/key"; then
  ipfs-key -bitsize=$keySize 2>/dev/null > "$site/key"
  echo "created $site/key"
else
  echo "skipped $site/key (exists)"
fi

if ! test -s "$site/ref"; then
  ./publish.sh "$site" "$initRef"
else
  echo "skipped initial publish"
fi
