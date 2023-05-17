---
layout: page
title: From Source
permalink: /getting-started/installation/source
grand_parent: Getting Started
parent: Installation
nav_order: 10
---

## Install From Source

First, ensure you have a recent version of Go (1.19+) and npm (9+) installed.

```bash
git clone https://github.com/ghostsecurity/reaper.git
cd reaper
make wails
wails doctor
```

Install any missing dependencies as prompted in the output. In Linux, this usually involves installing `libgtk-3-dev`
and `libwebkit2gtk-4.0-dev`, but the output will guide you.

Finally either `make run` to start _reaper_, or `make install` to install it to your `PATH` (will prompt for password).

