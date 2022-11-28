<p align="center">
<img width="400" src="frontend/src/assets/images/logo.png">
</p>

# Reaper

> :dragon: HERE BE DRAGONS!
> This is a work in progress. It's an experimental PoC and will likely change almost entirely over time. I'm currently working on it as a side-project to test out some ideas.

Reaper is a reconnaissance and attack proxy, built to be a lightweight, API-focused equivalent to Burp Suite/ZAP etc.

## Installation

For now you can clone the repo and use `make run` to try it out. Ensure you have `go` and `npm` installed and in your path first.

## Usage

One day we should put some docs here.

## Extensibility

One day we should add plugin support.

## Live Development

To run in live development mode, run `make run` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `make build`.

## TODO: For MVPoC

- [x] Basic proxy functionality
- [x] Request history
- [x] Request viewer
- [x] Response viewer
- [x] Request editor
- [x] Response editor
- [x] Request interception
- [ ] Response interception
- [ ] Rewrite proxy dependency (the current module has dodgy TLS handling and lots of globals)
- [ ] Fuzzer (Intruder style)
- [ ] Scope editing (ignore all non-scope requests/responses)
- [ ] Target map
- [ ] Websockets
- [x] Save/load settings, restart proxy as required
- [x] Dark mode (most important)
- [x] Ghost branded UI theme
- [x] Ghost branded non-proxy page
- [x] Ghost branded syntax highlighting theme
- [ ] Proxy status indicator (toggle switch in top right?)
- [ ] History descending - sortable?
- [ ] History filtering
- [ ] Intercept filtering
- [ ] JWT viewer/editor
- [ ] Protobuf support
