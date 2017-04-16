# Tasks

> docker run -d -v `pwd`/queues:/queues -p 4242:80 -e WEBHOOK_SECRET=... -e GITHUB_TOKEN=... ipfs-pages github -d
> docker run -d -v `pwd`/queues:/queues -e IPFS_API=... ipfs-pages build -d
> docker run -d -v `pwd`/queues:/queues -e DNSIMPLE_TOKENS='[...]' ipfs-pages dnsimple -d

- GHI maintains targets list because that's interaction with GH
- Builder maintains mount table because that's build logic
- GHI pushes master build job at startup for every target
  - Builder gets a chance to refresh if needed

- GH-Interface
  - [ ] Scans orgs for targets at startup (reason: want builder to see all repos)
    - [ ] Scans concurrently
  - [ ] Receives webhook
  - [ ] Processes webhook if it's a push event for a target
  - [ ] Updates targets list if push event indicates change of .ipfspages.yml
  - [ ] Writes mounts to mount table
  - [ ] Skips dnslinks which aren't allowlisted for the target
  - [ ] Pushes job to build queue
  - [ ] Shifts job from status queue
  - [ ] Writes to Status API
  - [ ] Shifts job from announce queue
  - [ ] Announces to IRC #ipfs-pages
  - [ ] Announces to logfile
  - [ ] Exposes Prometheus metrics

- Builder
  - [ ] Shifts job from build queue
  - [ ] Clones Git repo at commit
  - [ ] Runs hugo build container
  - [ ] Adds hugo result or error to ipfs
  - [ ] Writes result hash to mount table
  - [ ] Pushes job2 to build queue if job is a mount in another target
  - [ ] Pushes job to deploy queue
  - [ ] Pushes job to status queue

- Deployer
  - [ ] Shifts job from deploy queue
  - [ ] Checks if DNS record needs to be updated
  - [ ] Acquires dnsimple token for respective domain name
  - [ ] Updates the DNS record

- Data structures
  - Build job: commit => (repo)
  - Deploy job: domain => (txtrr)
  - Status job: commit => (status, repo, hash)
  - Mount table: repo => (hash, mountedIn+, dnslink)


# .ipfspages.yml

```yml
---
# libp2p/docs/.webhook.yml
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
  - /ipns/docs.libp2p.io
  - /ipns/docs-ng.libp2p.io

---
# libp2p/website/.webhook.yml

dnslinks:
  - /ipns/libp2p.io
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
    - Repo migrations
  - HTTP-to-IPFS gateway
  - Pinning & GC
  - In web browsers
  - With large datasets
    - add --nocopy
    - Datastore.NoSync
  - Naming / IPNS
    - multiple names
    - republishing
  - Pubsub
  - Running in production
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
