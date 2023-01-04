

# Reaper


> :dragon: HERE BE DRAGONS!
> This is a work in progress. It's an experimental PoC and will likely change almost entirely over time. I'm currently working on it as a side-project to test out some ideas.

<img width="75" align="right" src="frontend/src/assets/images/logo.png">

Reaper is a reconnaissance and attack proxy, built to be a lightweight, API-focused equivalent to Burp Suite/ZAP etc. Imagine if your favourite attack proxy had a baby with your favourite API testing tool - that's our goal.

![Reaper Screenshot](screenshot.png)

## Installation

Eventually you can grab a binary from the releases page, or use your favourite package manager. For now, you'll need to build from source as described below...

## Building and Hacking Locally

The following steps should work across Mac, Linux and Windows:

1. Clone the repo.
2. Ensure you have Go (1.19+) and npm installed.
3. Run `make wails` to install Wails v2.
4. Run `wails doctor` to ensure your environment is configured correctly. Install any missing dependencies as prompted.
5. Run `make run` to start _reaper_.

In order to build cross-platform, production binaries, creating and pushing a git tag will (eventually) trigger GitHub actions to publish versioned binaries as release artifacts.

## TODO

### To lock down main branch

- [x] Run Go tests GitHub action (Linux, Mac, Windows)
- [x] Run JS tests GitHub action
- [ ] Run Go linter GitHub action (blocked by action whitelisting)
- [ ] Raise the below list of issues to GitHub issues
- [ ] Squash commits on main branch
- [ ] Protect main and require all build steps to pass etc.

### Ghostbank milestone

- [ ] Workspace scope/name editor
- [ ] Workspace switch button
- [ ] Add modal for intercepts?
- [ ] Add collection/saved requests
- [ ] Rework request/response popout/editor
- [ ] Contextual tags (request has auth)
- [ ] Add context menu to request parts (fuzz this shiz)

### Post lockdown (MVPoC milestone)

- [ ] In-memory db of saved requests, sequences etc.
- [ ] Add pre/post scripts + variables (a la postman)
- [ ] Add sequences + sequence editor
- [ ] Import saved request from `curl` command
- [ ] Try editors like https://codemirror.net/ + https://shiki.matsu.io/
- [ ] Export HAR files/curl
- [ ] Investigate TINES style request editor (raw vs. smart)
- [ ] Add notes to requests
- [ ] Rewrite proxy dependency (the current module has dodgy TLS handling and lots of globals)

### Future Development

- [ ] JWT viewer/editor
- [ ] Protobuf support
- [ ] Add omnibox for EVERYTHING - https://tailwindui.com/components/application-ui/navigation/command-palettes
