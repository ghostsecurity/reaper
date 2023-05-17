---
layout: page
title: Chrome (MacOS)
parent: Browser Setup
grand_parent: Getting Started
permalink: /getting-started/browser-setup/chrome-macos
nav_order: 2
---

## Chrome (MacOS) Setup

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
8. Select `System` and then `All Items`.
9. Drag and drop the certificate you exported from Reaper into the list.
10. Double-click on the certificate you just added in the list.
11. Expand the `Trust` section and select `Always Trust` for `When using this certificate`.
12. Close the certificate window and enter your password to confirm the change.
13. Restart Chrome.

![img_2.png](../../images/browsers/chrome/macos/img_2.png)

### Configure Chrome to use Reaper

It is recommended to use a Chrome extension to manage your proxy use, such
as [Proxy Switcher](https://chrome.google.com/webstore/detail/proxy-switcher/iejkjpdckomcjdhmkemlfdapjodcpgih). Tools
like this allow you to quickly
switch between proxies and configure particular proxies for particular URLs. This means you can configure Chrome to
only use Reaper for a target web application, and route all other traffic as usual. You can use the proxy
address `127.0.0.1:8081` to send traffic through Reaper.

Alternatively, you can configure Chrome to use Reaper directly (not recommended).

Finally, you should [test](test) your setup.
