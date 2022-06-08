# GoCV Docker
GitHub action for automating build for armv7 gocv


## Build
```sh
docker build -t ghcr.io/jonathongardner/gocv:b4.5.5 -f Dockerfile.base .
docker build -t ghcr.io/jonathongardner/gocv:v4.5.5 .
```

## Run
```sh
docker run --rm -it --device /dev/video0 -p 3000:3000 ghcr.io/jonathongardner/gocv:4.5.5-0.0.0 0 0.0.0.0:3000
```
