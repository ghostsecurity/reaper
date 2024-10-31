<a id="readme-top"></a>
<h1><img src="docs/img/logo-reaper-only.png" width="30px"> Reaper</h1>

Reaper by [Ghost Security](https://ghost.security) is a modern, lightweight, and extensible open-source application security testing framework built to be operated by both humans and AI Agents.  It provides several capabilities that enable manual application security workflows: target reconnaissance, request proxying, request tampering/replay, live collaboration, and active test running.  When combined with an AI Agent backed by an LLM, Reaper becomes a flexible engine to drive even more powerful application testing workflows.

<!-- LOGO AND YOUTUBE -->
<br />
<div align="center">
  <a href="https://www.youtube.com/watch?v=t0Oe1IIB9xI">
    <img src="docs/img/logo-reaper-only.png" alt="Logo" width="200px">
  </a>
</div>


> :warning:
> This project is undergoing rapid development and may change significantly in the near future.


<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#about">About</a></li>
    <li><a href="#project-goals">Project Goals</a></li>
    <li><a href="#setup">Setup</a></li>
    <li>
        <a href="#usage">Usage</a>
        <ul>
            <li><a href="#scan">Scan</a></li>
            <li><a href="#explore">Explore</a></li>
            <li><a href="#replay">Replay</a></li>
            <li><a href="#test">Tests</a></li>
            <li><a href="#ai-agent">AI Agent</a></li>
            <li><a href="#report">Reports</a></li>
      </ul>
    </li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#acknowledgments"> Acknowledgments </a></li>
  </ol>
</details>

## About

The Reaper framework was created to combine several application security workflow steps from discrete tools into one.  It aims to streamline the process of discovering targets, performing reconnaissance, tampering/replayiing requests, driving workflows via API or AI automation, and more from within the same toolset.

Existing tools (e.g. [Burp Suite](), [Zap](), and [subfinder](https://github.com/projectdiscovery/subfinder) / [katana](https://github.com/projectdiscovery/katana) / [nuclei](https://github.com/projectdiscovery/nuclei)) are able to perform individual steps of the testing lifecycle but require the end user to manually engage with each tool and export/import data between steps.

Reaper is designed to be orchestrated by humans and AI Agents (Agents) to enable almost any workflow you need to become a reality.  Agents that are backed by an LLM can act as another helpful team member and perform tasks in seconds that would take hours by analysts.  For example, it can assist with test parameter tuning, summarization of data/findings, data analysis, and even report generation.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Project Goals

- A modern, lightweight, and extensible framework for application security testing
- Usable by humans and AI Agents alike
- A platform for running autonomous workflows
- Easy to maintain and extend
- Help avoid application security engineer burn-out with helpful automation

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Installation

### Running via Docker

If you have Docker version 19.x or above, the quickest path to getting running is to clone this repo and run:

```sh
docker compose up
```

### Running via Binary

```
TODO
```


<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- Usage -->
## Usage

### Scan

The first step in reconnaissance is enumerating the available targets for a given domain/subdomain and to probe them for availability.  Click `Add Domain` and enter in a domain or subdomain that you are authorized to test.  For example `ghostbank.net` or `api.ghostbank.net`.  With the `Auto-scan` checkbox enabled, click `Add and scan` to initiate discovery of live hosts.

### Explore

To capture requests made to a target system, enable the `Proxy on` toggle at the top of the page.  From there, configure your browser or other client to route requests through the proxy at `localhost:8080` for both HTTP and HTTPS.

To install the proxy's certificate and configure your tool/browser to proxy through Reaper, [follow this guide](docs/proxy_certs.md).

### Replay

Requests/Responses that have traversed the Proxy will appear in this listing.  The filter allows filtering all requests by fuzzy match on the hostname or path.  The `All`/`APIs` toggles viewing of all or responses of content-type `application/json`.

To replay or tamper a request:

1. Select the desired request.
2. On the right pane, click `Replay original` to resend without modification.  The `Response` pane will update automatically.  In many cases, there will be no change in that field.
3. To send a modified request, live-edit either the Request headers or Request Body as desired.  Click `Replay modified` and view the response in the `Response` field.

### Tests

This workspace drives testing workflows based on endpoints and/or requests that match desired criteria.  For example, when testing for Broken Object Level Access (BOLA) / Insecure Direct Object Reference (IDOR) vulnerabilities, it typically requires capture and replay of a valid request to an endpoint while fuzzing certain parameters.  Stay tuned as we continue to develop this capability.  In the meantime, your feedback is welcomed and encouraged!

### AI Agent

The AI Agent capability is the basis for a natural language interaction with one or more Agents via a chat-like interface.  Each `session` will record all messages and actions taken by the Agent and provide human-in-the-loop confirmation for important actions as needed.  Stay tuned as we continue to develop this capability.  In the meantime, your feedback is welcomed and encouraged!

### Reports

To view reports generated and saved via the `/api/reports` `POST` endpoint, select the desired report.  Stay tuned as we continue to develop this capability.  In the meantime, your feedback is welcomed and encouraged!

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTRIBUTING -->
## Contributing

First, thank you for taking the time to check out Reaper! Our primary goal is to get as many folks using it and to drive a roadmap based on your feedback.  If you have a great idea for an enhancement or you have encountered a bug, we'd greatly appreciate a well-formed [Issue](https://github.com/ghostsecurity/reaper/issues/new) in this repo so we can triage and prioritize accordingly.

Reaper is distributed under the [Apache 2.0 License](LICENSE). All Reaper contributors and community members must adhere to the [Code of Conduct](CODE_OF_CONDUCT.md)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

Here are a list of projects we want to acknowledge:

* [ProjectDiscovery](https://github.com/projectdiscovery) - produces a suite of open source tools tailored for offensive security: security engineers, bug bounty hunters, and red teamers.  The creaters of [subfinder](https://github.com/projectdiscovery/subfinder), [katana](https://github.com/projectdiscovery/katana), [nuclei](https://github.com/projectdiscovery/nuclei), and many other great tools.

<p align="right">(<a href="#readme-top">back to top</a>)</p>
