# How to Hack Ghostbank with Reaper
This tutorial will walk through hacking Ghostbank with Reaper.

## Steps
0. [Prerequisites](#prerequisites)
1. Start Reaper
2. Add a Domain
3. Capture Traffic
4. Manual Test
5. Automated Test
6. AI Agent-Driven Test

## Prerequisites


1. [Clone the Reaper repository](https://docs.github.com/en/desktop/adding-and-cloning-repositories/cloning-a-repository-from-github-to-github-desktop).
2. [Install Docker](https://docs.docker.com/engine/install/).
3. Have two web browsers ready. We recommend [Chrome](https://www.google.com/chrome/dr/download/) for using Reaper and [Firefox](https://www.mozilla.org/en-US/firefox/new/) for using Ghostbank.

## 1. Start Reaper
First things first, let's get Reaper up and running.

1. From a command line, navigate to the Reaper directory and run `docker compose up`. 
2. In a web browser (we recommend Chrome), navigate to https://127.0.0.1:8000 to reach the Reaper UI.
3. Enter a username (or use the default Reaper Admin) and click *Sign in*.

<p align="center"><img src="/docs/img/reaper_login.png" width="400" /></p>

## 2. Add a Domain
Now, let's add a Domain and discover the available hosts.

1. On the *Scan* tab, click *Add Domain*.
2. Enter `ghostbank.net`.
3. Leave the *Auto-scan* option enabled and click *Add and scan*.

<p align="center"><img src="/docs/img/reaper_add_domain.png" width="400" /></p>
