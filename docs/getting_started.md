# 💿 Installation

Reaper runs as either a Docker container (recommended) or as a binary and is controlled by humans using the local web UI.

## Running via Docker Run or Docker Compose

If you have Docker version 19.x or above, the quickest path to getting running is to run:

```sh
docker run -t --rm  \
  -e HOST=0.0.0.0 \
  -e PORT=8000 \
  -e PROXY_PORT=8080 \
  -e OPENAI_API_KEY=sk-your-key-here \
  -p 8000:8000 \
  -p 8080:8080 \
  ghcr.io/ghostsecurity/reaper:latest
```

Once the container is up and running, navigate to [http://localhost:8000](http://localhost:8000) to activate the Reaper UI.

Or, you can use `docker compose`:

 1. Clone this repo locally,
 2. Modify the docker-compose.yml environment variables to add your OpenAI API key if desired, and
 3. Run from the command line:

```sh
docker compose up
```

Again, once the container is up and running, navigate to [http://localhost:8000](http://localhost:8000) to activate the Reaper UI.

# Usage

## Scan

The first step in reconnaissance is enumerating the available targets for a given domain/subdomain and to probe them for availability. Click `Add Domain` and enter in a domain or subdomain that you are authorized to test. With the `Auto-scan` checkbox enabled, click `Add and scan` to initiate discovery of live hosts.

## Explore

✨ *Learn how to set up Reaper's proxy [using this guide](proxy_certs.md).*

Capture traffic by following these steps:
1. The proxy is on by default. It can be stopped and restarted using the `Proxy on/off` toggle at the top of the page Explore page.
2. Configure your browser or other client to route requests through the proxy at `localhost:8080` for both HTTP and HTTPS.
3. Browse through your target app, targeting any interactions or workflows you want to test for vulnerabilities.
4. The Explore page in Reaper will show an inventory of hosts and endpoints captured by the proxy.

*If you're having issues with the proxy, check the [guide](proxy_certs.md). If the issue persists, please let us know what's haunting you at reaper@ghost.security.*

## Replay

Requests/Responses that have traversed the Proxy will appear in this listing.  The filter allows filtering all requests by fuzzy match on the hostname or path.  The `All`/`APIs` toggles viewing of all or responses of content-type `application/json`.

To replay or tamper a request:

1. Select the desired request.
2. On the right pane, click `Replay original` to resend without modification.  The `Response` pane will update automatically.  In many cases, there will be no change in that field.
3. To send a modified request, edit either the Request Headers or Request Body as desired.  Click `Replay modified` and view the response in the `Response` field.

## Tests

This workspace drives testing workflows based on endpoints and/or requests that match desired criteria. For example, when testing for Broken Object Level Access (BOLA) / Insecure Direct Object Reference (IDOR) vulnerabilities, it typically requires capture and replay of a valid request to an endpoint while fuzzing certain parameters.  Stay tuned as we continue to develop this capability.  In the meantime, your feedback is welcomed and encouraged!

## AI Agent

The AI Agent capability is the basis for a natural language interaction with one or more Agents via a chat-like interface.  Each `session` will record all messages and actions taken by the Agent and provide human-in-the-loop confirmation for important actions as needed.  Stay tuned as we continue to develop this capability.  In the meantime, your feedback is welcomed and encouraged!

## Reports

To view reports generated and saved via the `/api/reports` `POST` endpoint, select the desired report.  Stay tuned as we continue to develop this capability.  In the meantime, your feedback is welcomed and encouraged!
