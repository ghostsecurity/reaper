# Reaper

> :warning:
> This is a work in progress. It's an experimental PoC and will likely change almost entirely over time.

<img width="75" align="right" src="frontend/src/assets/images/logo.png">

Reaper is a reconnaissance and attack proxy, built to be a modern, lightweight, and efficient equivalent to Burp
Suite/ZAP etc. This is an attack proxy with a heavy focus on automation, collaboration, and building universally
distributable workflows.

![Reaper Screenshot](screenshot.png)

## Documentation

For further documentation on installation, configuration and usage, check out
the [docs](https://ghostsecurity.github.io/reaper).

## Building and Hacking Locally

The following steps should work across Mac, Linux and Windows:

1. Clone the repo.
2. Ensure you have Go (1.19+) and npm installed.
3. Run `make wails` to install Wails v2.
4. Run `wails doctor` to ensure your environment is configured correctly. Install any missing dependencies as prompted.
5. Run `make run` to start _reaper_.
