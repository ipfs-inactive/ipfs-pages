# ipfs.li

- https://ipfs.li
- entrypoint for project knowledge+activity
- indexes github, pinbot, irc, mailing list, twitter, meetup, sprint, ipfs-community calendar, waffle, incidents, etc. activity
  - imports files, attachments, gists, hackmd, zoom
  - first steps: github issues, irc lines
  - latency is hidden in the indexer, instead of presented to the UI
- ultrafast results
  - ipfs.li/#searchterm, and redirct from searchterm.ipfs.li
  - finger tree
    - great because we can very quickly provide recent results, and then load more older results
    - how fast is finger tree on ipfs generally?
    - https://github.com/krl/aolog
  - fast autofocus on search form
    - don't render search with js
    - just plain onload=focusInput, onsubmit=donothing, then when ipfs.js is loaded give feedback and do search
  - dweb.link and cache-control: immutable
  - windows.ipfs and polyfill
  - time to first byte
  - load js-ipfs node on first edit, until then just use dweb.link gateway (or local if rewritten by addon)
- daily/weekly digest posting or email
- browse by date, type, author, top-k words

# docs.ipfs.io

- Home
- The IPFS stack
  - IPFS
    - cli, api, fs-repo, datastore, bitswap, dag-pb / merkledag, unixfs, dnslink, fs path / dweb link, ipns, gateway
  - IPLD
  - libp2p
  - Multiformats
- Implementations
  - go-ipfs
  - js-ipfs
- Usage
  - Installing, updating
  - HTTP-to-IPFS gateway
  - Pinning
  - In web browsers
  - With large datasets
  - Naming / IPNS
    - multiple names
    - republishing
  - Pubsub
  - Running in production
    - Repo migrations
    - Garbage collection
    - Monitoring
    - Failure modes
    - Analyzing performance metrics
    - Debugging
    - Containers
    - Abuse handling
    - ...
  - ...
- Reference
  - IPFS Core API
  - HTTP API
  - go-ipfs packages
    - go-ipfs-api
    - go-ipfs-gateway
    - ...
  - js-ipfs packages
    - js-ipfs-api
    - js-ipfs-browser
    - ...
- Tools
  - ipfs-pack
  - ipfs-cluster
  - ipget
  - gx
  - iptb
  - fs-repo-migrations
- Development
  - Testing
  - ...


# IPFS Pages

Inputs:
- Github org or repo webhooks
  - react on commit-push
- Commit tree .webhook.yml

Outputs:
- IPFS DAG with rendered page, and log output
- Status API
- Pinning
- DNS record updates
- Trigger rebuild of targets which mount the current target (this is why we scan for repos)

```yml
---
# libp2p/docs/.webhook.yml
commands:
  master:
    - make build
  docs-ng:
    - make build
mounts:
  - src: github.com/libp2p/js-libp2p-*
    dest: js/*
  - src: github.com/libp2p/go-libp2p-*
    dest: go/*
  - src: github.com/libp2p/go-*-transport
    dest: go/go-*-transport
  - src: github.com/openbazaar
    dest: go/go-onion-transport
dnslinks:
  master: /ipns/docs.libp2p.io
  docs-ng: /ipns/docs-ng.libp2p.io

---
# libp2p/website/.webhook.yml

commands:
  master:
    - make build
dnslinks:
  master: /ipns/libp2p.io
```

- mounts:
  - these also act as triggers
  - when a build of the referenced repo succeeds, we update the mount
    - if we can receive webhooks for it
  - foo-* means all repos starting in foo- with .webhook.yml in their default branch
  - we scan for repos at startup, then handle repo.created/deleted webhooks
- dnslinks:
  - branch uses the current repo's named branch, and the same-named branches (or default branches) of mounted repos



---

pagesworkspace/
  builds/
    COMMIT_ID_1/
      tree/
        ...
      build.log
      cid
    COMMIT_ID_N/
      ...
  repos/
    github.com/
      multiformats/
        website/
          .git

> pagebuild -q github.com/multiformats/website COMMIT_ID
CID
