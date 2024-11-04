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

Use `-x` and `-k` to proxy a request through Reaper without TLS warnings, respectively:

```sh
curl -x localhost:8080 -k https://ipinfo.io

{
  "ip": "98.76.54.32",
  "hostname": "ip98-76-54-32.isp.net",
  "city": "San Diego",
  "region": "California",
  "country": "US",
  "loc": "38.8462,-77.3064",
  "org": "AS1337 Ghost Communications Inc.",
  "postal": "90001",
  "timezone": "America/Los_Angeles"
}   
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

Launch a new Chrome instance using the local Reaper proxy (without altering your system settings).

#### macOS

```sh
/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome \
  --proxy-server="127.0.0.1:8080" \
  --ignore-certificate-errors \
  --user-data-dir="$HOME/.chrome-reaper"
```


#### Linux

```sh
google-chrome \
  --proxy-server="127.0.0.1:8080" \
  --ignore-certificate-errors \
  --user-data-dir="$HOME/.chrome-reaper"
```