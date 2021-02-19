# Go Generate Cloud Native Buildpack

It is enabled by `BP_GO_GENERATE` being set to `true` in your environment.

## Packaging
To package this buildpack for use with a CNB platform of choice:
```
./scripts/package.sh --version <your-version-number>
```
The resulting buildpack (.tgz) and buildpackage (.cnb) will be placed in
`build/`.

## Configuration
Use `BP_GO_GENERATE_FLAGS` and `BP_GO_GENERATE_ARGS` in your environment to
configure the behaviour of `go generate`.

For instance, building an image with
[`pack`](https://github.com/buildpacks/pack):
```
pack build my-app --buildpack gcr.io/paketo-buildpacks/go-dist \
                  --buildpack gcr.io/paketo-buildpacks/go-mod-vendor \
                  --buildpack /path/to/repo-dir/go-generate/build/buildpack.tgz \
                  --buildpack gcr.io/paketo-buildpacks/go-build \
                  --env BP_GO_GENERATE=true \
                  --env BP_GO_GENERATE_ARGS='main.go helper.go' \
                  --env BP_GO_GENERATE_FLAGS='-v -run "^//go:generate sometool"'
```
will result in the buildpack running the equivalent of:
```
go generate -v -run "^//go:generate sometool" main.go internal.go
```

See `go help generate` for information about flag and argument options for `go
generate`.
