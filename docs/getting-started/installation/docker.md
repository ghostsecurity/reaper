---
layout: page
title: Docker
permalink: /getting-started/installation/docker
grand_parent: Getting Started
parent: Installation
nav_order: 11
---

## Run with Docker

You can get started quickly with Reaper using Docker:

```bash
docker run -v $HOME/.reaper:/.reaper -p 8080:8080 -p 31337:31337 ghcr.io/ghostsecurity/reaper
```

This will start Reaper and bind the Reaper data directory to `$HOME/.reaper` on your host machine. This means that any
workspace/settings data will be persisted between containers.

The CLI will prompt you with a GUI link shortly after starting.
