# Certificate Installation and Proxy Configuration

In order to browse websites protected by TLS/SSL through the Reaper proxy, you must download and install the certificate authority to your browser or operating system's trust store.

## Operating System

First, either clone this repo or download the contents of the `tls` directory.

### OSX Keychain

To install in OSX's keychain so that all applications can trust the proxy CA:

1. Open `Utilities` > `Keychain Access`.
2. Under the `Default Keychain` > `login`, select the `Certificates` tab
3. Drag the `tls/ca.pem` file into the listing.  It should appear as `reaper.ghostsecurity.com` with a red `x` icon.
4. Double click on the certificate to bring up the info dialog.
5. Under `Trust`, choose `Always Trust` for the `Secure Sockets Layer (SSL)` option.
6. Click the window's `x`. This should ask for your password to save changes. Do so and click `Update Settings`.
7. The certificate should now show with a blue icon.
8. If your browser is configured to trust CA certs from the operating system (the default), you are all set.

### Linux and Windows

Instructions coming soon.  PRs are welcomed.

## Tools

Some tools allow specifying the CA pem file on demand. First, either clone this repo or download the contents of the `tls` directory.

### curl

Use `-x` and `--cacert` to proxy a request through Reaper without TLS warnings, respectively:

```sh
curl -x localhost:8080 --cacert path/to/reaper/tls/ca.pem https://icanhazip.com
<ip response here>
```

## Browsers

First, either clone this repo or download the contents of the `tls` directory.

### Firefox

To install the certificate in Firefox's Trusted Authority store:

1. Navigate to [about:preferences](about:preferences)
2. Search for "cert" in the `Find in Settings` input.
3. Click `View Certificates...`.
4. Under the `Authorities` tab, click `Import...`.
5. Navigate to the `tls` directory and select `ca.pem`.
6. In the dialog that pops up, check `Trust this CA to identify websites` and click `Ok`.
7. You should see the `reaper.ghostsecurity.com` entry.  Click `Ok`.

To configure Firefox to proxy through Reaper:

1. Navigate to [about:preferences](about:preferences)
2. Search for "proxy" in the `Find in Settings` input.
3. Click `Settings...`.
4. Check `Manual proxy configuration`, enter `localhost` and port `8080` for the `HTTP Proxy` fields.
5. Check `Also use this proxy for HTTPS`.
6. Click `Ok`.
7. Now navigate to the target and observe requests being proxied in Reaper.

Disable proxying by switching the proxy settings dialog back to `No proxy`.  Remove the certificate from the Certificate Manager by selecting the `reaper.ghostsecurity.com` CA in the `Authorities` list and clicking on `Delete or Distrust...`.

### Chrome

To install the certificate for Chrome to use, follow the steps to install for your <a href="#operating-system">Operating system</a>.

To configure Chrome to proxy through Reaper, you must configure your OS proxy settings.  Chrome doesn't have native proxy configuration without help from a [proxy switcher extension](https://chromewebstore.google.com/search/proxy%20switcher):

1. Either use one of the extensions or navigate to [chrome://settings/?search=proxy](chrome://settings/?search=proxy)
2. Click `Open your computer's proxy settings`
3. Enable the `HTTP` and `HTTPS` proxy with a host:`localhost` and port: `8080`

Disable proxying by reverting the proxy settings to disabled.