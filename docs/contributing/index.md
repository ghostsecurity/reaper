---
layout: page
title: Contributing
permalink: /contributing
has_children: true
nav_order: 10
---

## Contributing

We welcome contributions to _reaper_ of all kinds, from bug reports to feature requests, and of course code
contributions!

### Bug Reports

If you find a bug, please [open an issue](https://github.com/ghostsecurity/reaper/issues/new) with as much detail as
possible, including:

- The version of _reaper_ you are using
- The operating system you are using
- The architecture you;re using (e.g. amd64, arm64 etc.)
- The steps to reproduce the bug

### Feature Requests

If you have a feature request, please [open an issue](https://github.com/ghostsecurity/reaper/issues/new) with as much
detail as possible. UI mock-ups are appreciated.

### Code Contributions

Reaper is made up of two main components: the Go backend (`/backend`), and the Vue frontend (`/frontend`). Most features
require backend and frontend changes, so please be prepared to make both. If you don't know both, don't worry, we can
help you get started. Feel free to raise a PR with changes to a single component - we can try to help support the
remaining development.

#### Hacking Locally

The easiest way to get up and running is to get Reaper running from the source on your machine.

The following steps should work across Mac, Linux and Windows:

1. Clone the repo.
2. Ensure you have Go (1.19+) and npm installed.
3. Run `make run` or `make docker` to start _Reaper_.
