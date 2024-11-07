# How to Hack Ghostbank with Reaper
This tutorial will walk through hacking Ghostbank with Reaper.

## Steps
0. [Prerequisites](#0-prerequisites)
1. [Start Reaper](#1-start-reaper)
2. [Add a Domain](#2-add-a-domain)
3. [Capture Traffic](#3-capture-traffic)
4. [Fuzz Manually](#4-fuzz-manually)
5. [Automated Test](#5-automated-test)
6. [AI Agent-Driven Test](#6-ai-agent-driven-test)

## 0. Prerequisites


1. [Clone the Reaper repository](https://docs.github.com/en/desktop/adding-and-cloning-repositories/cloning-a-repository-from-github-to-github-desktop).
2. [Install Docker](https://docs.docker.com/engine/install/).
3. Have two web browsers ready. We recommend [Chrome](https://www.google.com/chrome/dr/download/) for using Reaper and [Firefox](https://www.mozilla.org/en-US/firefox/new/) for using Ghostbank.

## 1. Start Reaper
First things first, let's get Reaper up and running.

1. From a command line, navigate to the Reaper directory and run `docker compose up`. 
2. In a web browser (we recommend Chrome), navigate to https://127.0.0.1:8000 to reach the Reaper UI.
3. Enter a username (or use the default Reaper Admin) and click *Sign in*.

<p align="center"><img src="/docs/img/reaper_login.png" width="400" /></p>

You will land on the Scan page, where we will get started by adding a domain!

## 2. Add a Domain
Let's add a Domain and discover the available hosts.

1. On the *Scan* tab, click *Add Domain*.
2. Enter `ghostbank.net`.
3. Leave the *Auto-scan* option enabled and click *Add and scan*.

<p align="center"><img src="/docs/img/reaper_add_domain.png" width="400" /></p>

You will see *ghostbank.net* in the Domain list. When the scan completes, you can click on *ghostbank.net* to view the discovered hosts. By default, Reaper will only show live hosts.

<p align="center"><img src="/docs/img/discovered_hosts.png" width="400" /></p>


## 3. Capture Traffic
Now, let's get to the good stuff. We are going to set up a proxy in Firefox, then log in to Ghostbank and initiate a fund transfer. Reaper's proxy will capture the traffic, allowing us to inspect, replay, and modify the requests.

**Firefox Proxy Configuration**
1. Open Firefox.
2. In the URL bar, enter `about:preferences` (or open the app menu and click *Settings*).
3. Search for `proxy` on the Settings page. *Network Settings* should show up.
4. Click *Settings*.
5. Select *Manual proxy configuration*.
6. Enter `localhost` in the HTTP Proxy field and `8080` in the Port field.
7. Select *Also use this proxy for HTTPS*.
8. Click *OK*.

<p align="center"><img src="/docs/img/firefox_proxy_settings.png" width="400" /></p>

**Capturing Traffic in Reaper**
1. Return to Reaper in Chrome and switch to the *Explore* tab.
2. Make sure the proxy is on (you'll see *Proxy on*).
3. Switch back to Firefox and browse to `https://ghostbank.net`.
4. Log in to Ghostbank.
5. Initiate a funds transfer by entering $10 and clicking *Transfer*. You may need to switch the From and To accounts.
6. Switch back to Chrome. You should see a list of captured requests!

<p align="center"><img src="/docs/img/captured_requests.png" width="400" /></p>

You now have valid requests that we can replay and modify -- good work!

## 4. Fuzz Manually
*[Fuzzing](https://owasp.org/www-community/Fuzzing)* is a software testing method wherein a human or program provides invalid, unexpected, or random data as inputs to a computer program. We are going to fuzz the Transfer API in Ghostbank to see if we can transfer funds from a different customer's account into our own account. But first, let's simply replay and modify our original transfer request.

**Replaying Original and Modified Requests**
1. Go to the *Replay* tab in Reaper. You should see a variety of requests captured from our previous interaction with Ghostbank.
2. Search for `transfer` and click on the `/api/v3/transfer` endpoint.
3. You should see the *Request Headers* and *Request Body*. Notice there are three inputs in the *Request Body*: `account_to`, `account_from`, and `amount`.
4. Click the *Replay original* button.
5. In Firefox, disable the proxy (Settings > Network Settings > Settings > No Proxy > OK).
6. Refresh the Ghostbank page. You should see a new transfer from replaying the request in Reaper.
7. Return to Reaper and change the `amount` to `20` in the Request Body.
8. Click *Replay modified*.
9. Switch back to Firefox and refresh Ghostbank. You should see another transfer, this time for $20!

<p align="center"><img src="/docs/img/replay_modified.png" width="500" /></p>


## 5. Automated Test

## 6. AI Agent-Driven Test