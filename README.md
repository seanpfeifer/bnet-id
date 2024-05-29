# bnet-id

What's my BNet ID?

## Building

Via Docker

```sh
docker build -t bnet-id:latest .
```

## Setup

Ensure you have a `secrets` directory with `secret.toml` in it with the following format:

```toml
clientID = "CLIENT_ID_HERE"
clientSecret = "CLIENT_SECRET_HERE"
redirectURL = "https://YOUR_SITE_HERE/oauthCallback"
```

Your clientID and clientSecret can be found at [Battle.net's API Access portal](https://develop.battle.net/access/clients/). You will also need to configure the associated client to have the redirectURL that you will be using in your application.

## Running

Locally:

```sh
go run ./cmd/server
```

Via Docker:

```sh
docker run --name="bnet" --rm -p 8080:8080 -v ./secrets:/secrets bnet-id:latest
```
