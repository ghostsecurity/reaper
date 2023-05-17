---
layout: page
title: Chrome (Windows)
parent: Browser Setup
grand_parent: Getting Started
permalink: /getting-started/browser-setup/chrome-windows
nav_order: 3
---

## Chrome (Windows) Setup

### Install the Reaper CA Certificate

1. First, ensure you have [installed](../installation) Reaper.
2. Launch Reaper and select or create a new workspace.
3. Hit the cog icon in the bottom left.

   ![img.png](../../images/browsers/reaper/img.png)
4. Select `Certificates` and hit the `Export CA Certificate` button. Export the certificate somewhere convenient,
   perhaps your home directory or desktop - it only needs to be there temporarily.

   ![img_1.png](../../images/browsers/reaper/img_1.png)
5. Launch Chrome and open the Settings menu.

   ![img.png](../../images/browsers/chrome/img.png)
6. Navigate to `Privacy and security` and select `Security`.

   ![img_1.png](../../images/browsers/chrome/img_1.png)
7. Select `Manage device certificates`.
8. Select `Trusted Root Certification Authorities` and click `Import`.
9. Import the certificate you exported from Reaper earlier.
10. Select the `Trusted Root Certification Authorities` certificate store and click Next.

    ![img_3.png](../../images/browsers/chrome/windows/img_3.png)
11. Click `Finish` and restart Chrome.

### Configure Chrome to use Reaper

It is recommended to use a Chrome extension to manage your proxy use, such
as [Proxy Switcher](https://chrome.google.com/webstore/detail/proxy-switcher/iejkjpdckomcjdhmkemlfdapjodcpgih). Tools
like this allow you to quickly
switch between proxies and configure particular proxies for particular URLs. This means you can configure Chrome to
only use Reaper for a target web application, and route all other traffic as usual. You can use the proxy
address `127.0.0.1:8081` to send traffic through Reaper.

Alternatively, you can configure Chrome to use Reaper directly (not recommended).

Finally, you should [test](test) your setup.
