# bnet-id

What's my BNet ID?

## Building

```sh
docker build -t bnet-id:latest .
```

## Running

```sh
docker run --name="bnet" --rm -p 8080:8080 -v ./secrets:/secrets bnet-id:latest
```
