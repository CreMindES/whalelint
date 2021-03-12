# <img width="22px" src="https://user-images.githubusercontent.com/5306361/110181582-6c807f80-7e0c-11eb-81c8-36d6a9c0db0b.png"> WhaleLint <img align="right" style="position: relative; top: 10px;" src="https://github.com/cremindes/whalelint/workflows/build/badge.svg" />

Dockerfile linter written in Go.

It provides static analysis for Dockerfiles, identifying common mistakes and promotes best practices.

<p align="center">
  <img width="500px" src="https://user-images.githubusercontent.com/5306361/110991142-870aa980-8374-11eb-8855-9f3ce400049e.png"/> 
</p>

## Usage

```bash
docker pull cremindes/whalelint:[tag]
docker run --rm -v $(pwd)/Dockerfile:/Dockerfile cremindes/whalelint:[tag] Dockerfile
```

## Sample output

<p align="center">
  <img width="750px" src="https://user-images.githubusercontent.com/5306361/110198673-775f0280-7e54-11eb-8e4e-ab6350fb4e7d.png"/>
</p>

## GitHub

[Official repository](https://github.com/CreMindES/whalelint) with further information.
