# netlink-api

## What is this?

This is (going to be) a feature-complete http wrapper to the vishvananda/netlink library.

## Why vishvananda/netlink ?

It's the library used by docker and similar projects.

## How do I use this ?

For now, run `go build . && ./netlink-api` (as root if you want to create/modify interfaces)

In the future this will be a containerized blob that will need to run with the `NET_ADMIN` cap or something.

## Is this safe to start using now?

Nope. Expect breaking changes with every commit. Once we're somewhat close to feature parity we'll tag a release and that'll be safe to use.

## Is there any documentation?

Not yet! We're trying to match the api provided by the upstream library, so check those docs first.

## What is this for anyway?

PlonkFW is (going to be) a linux-based firewall applicance along the lines of pf/opnsense, but built with modern technologies. This project here will be one of many components.

## How can I help?

Fork and start writing endpoints! I'm not a "real" dev, so if you have any structure/design recommendations those are very welcome.
