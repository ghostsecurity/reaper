---
layout: page
title: Other Browsers/Tools
parent: Browser Setup
grand_parent: Getting Started
permalink: /getting-started/browser-setup/other
nav_order: 5
---

## Other Browsers

### Install the Reaper CA Certificate

1. First, ensure you have [installed](../installation) Reaper.
2. Launch Reaper and select or create a new workspace.
3. Hit the cog icon in the bottom left.

   ![img.png](../../images/browsers/reaper/img.png)
4. Select `Certificates` and hit the `Export CA Certificate` button. Export the certificate somewhere convenient,
   perhaps your home directory or desktop - it only needs to be there temporarily.

   ![img_1.png](../../images/browsers/reaper/img_1.png)
5. Import and trust the certificate in your browser settings.

### Configure Browser to use Reaper

It is recommended to use a browser extension to manage your proxy use. Tools like this allow you to quickly switch
between proxies and configure particular proxies for particular URLs. This means you can configure your browser to only
use Reaper for a target web application, and route all other traffic as usual. You can use the proxy
address `127.0.0.1:8081` to send traffic through Reaper. Finally, you should [test](test) your setup.

### Configuring CLI Tools to use Reaper

Many CLI tools support proxying through an HTTP proxy. You can use the proxy address `127.0.0.1:8081` to send traffic
through Reaper.

Many CLI tools respect the `http_proxy` and `https_proxy` environment variables (most often lower case, but it's worth
trying both permutations). You can set these variables to the proxy address `127.0.0.1:8081` to send traffic through
Reaper.

Remember to [test](test) your setup.
