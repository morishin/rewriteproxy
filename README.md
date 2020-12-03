# rewriteproxy

Emulator of Firebase Hosting [_rewrites_](https://firebase.google.com/docs/hosting/full-config#rewrites) for your development environment.

## Examples

### Rewrite path

![rewrite_path](https://user-images.githubusercontent.com/1413408/101050973-4dcddb00-35c8-11eb-9ad2-07a5e4713f82.png)

### Cloud Functions

It emulates this: https://firebase.google.com/docs/hosting/full-config#rewrite-functions

![cloud_functions](https://user-images.githubusercontent.com/1413408/101050971-4dcddb00-35c8-11eb-936c-5deb8e2604a4.png)

### `cleanUrls: true`

It emulates this: https://firebase.google.com/docs/hosting/full-config#rewrites

![clean_urls1](https://user-images.githubusercontent.com/1413408/101051746-29263300-35c9-11eb-8693-89014592d179.png)

![clean_urls2](https://user-images.githubusercontent.com/1413408/101051730-275c6f80-35c9-11eb-8426-8ad7cb086c32.png)

### Root path

![root_path](https://user-images.githubusercontent.com/1413408/101050970-4d354480-35c8-11eb-9b96-58d83d03a8f5.png)

## Installation

```sh
$ go get github.com/morishin/rewriteproxy
```

## Usage

```sh
$ rewriteproxy \
  --port=3000 \
  --firebase-json=/path/to/firebase.json \
  --web-app-url=http://localhost:1234 \
  --cloud-function-base-url=http://localhost:5001/your-project-id/us-central1
```

## Development

### Run

```sh
$ go run . \
  --port=3000 \
  --firebase-json=firebase.json.example \
  --web-app-url=http://localhost:1234 \
  --cloud-function-base-url=http://localhost:5001/your-project-id/us-central1
```
