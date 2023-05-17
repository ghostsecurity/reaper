---
layout: page
title: From Source
permalink: /getting-started/installation/source
grand_parent: Getting Started
parent: Installation
nav_order: 10
---

## Install From Source

First, ensure you have a recent version of Go (1.19+) and npm installed.

```bash
git clone https://github.com/ghostsecurity/reaper.git
cd reaper
make wails
wails doctor
```

Install any missing dependencies as prompted, and finally either `make run` to start _reaper_, or `make install` to
install it to your `PATH`.

