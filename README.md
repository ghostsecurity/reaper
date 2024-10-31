# Reaper

A modern web app attack and testing framework.

## What does it do?

- intercepting web proxy
- enumerate and discover hosts & subdomains (using pd subfinder)
- probe discovered hosts (using pd httpx)
- TODO: crawl discovered hosts (using pd katana)
- TODO: attack target hosts (using pd nuclei)
- TODO: fuzz inputs against specific host endpoints
- team collaboration (multiple users shared workspace)

## Running it

Run locally via Docker compose.

```
docker compose up
```

Locally via make

```
make run
```

Browse to [http://localhost:8000](http://localhost:8000)


