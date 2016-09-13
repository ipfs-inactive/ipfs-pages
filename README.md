# ipfs-pages

[![](https://img.shields.io/badge/made%20by-Protocol%20Labs-blue.svg?style=flat-square)](http://ipn.io)
[![](https://img.shields.io/badge/project-ipfs-blue.svg?style=flat-square)](http://github.com/ipfs/ipfs)
[![](https://img.shields.io/badge/freenode-%23ipfs-blue.svg?style=flat-square)](http://webchat.freenode.net/?channels=%23ipfs)

> Easy publishing of static web content on IPFS

This is purely about *publishing* and *persisting* the stuff.
Later on this might also:

- *build* the stuff
- *listen* for changes
- *verify* dnslinks, maybe even *create* dnslinks
- *measure* update propagation

## Usage

```sh
> ssh pages@pages.ipfs.team ./create.sh examplesite
> ssh pages@pages.ipfs.team ./publish.sh examplesite Qmfoobar
> ssh pages@pages.ipfs.team tail -f republish.log

## Setup

```sh
# create user
> groupadd pages
> useradd -g pages -d /opt/pages -s /bin/bash pages
> mkdir -p /opt/pages/.ssh
> vim /opt/pages/.ssh/authorized_keys
> chown pages:pages /opt/pages

# setup cronjob for periodic republishing
> echo "@hourly pages cd && ./republish.sh >> republish.log"
```

## Debugging

Enable the `PAGES_TRACE` environment variable.

```sh
> ssh pages@pages.ipfs.team env PAGES_TRACE=1 ./create.sh examplesite
```

## Maintainers

Captain: [@lgierth](https://github.com/lgierth).

## Contributing

Please see [contribute.md](contribute.md)!

This repository falls under the IPFS [Code of Conduct](https://github.com/ipfs/community/blob/master/code-of-conduct.md).

### Want to hack on IPFS?

[![](https://cdn.rawgit.com/jbenet/contribute-ipfs-gif/master/img/contribute.gif)](https://github.com/ipfs/community/blob/master/contributing.md)

## License

MIT
