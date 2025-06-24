# motel

An example go CLI to demo and learn new Go tooling!

## Project status (2025-06-13)

I use `motel` to test CI/CD, so it's not really useful to anyone else.

## Convert to a new project

See [Go Project Notes](https://www.bbkane.com/blog/go-project-notes/#creating-a-new-go-project).

## Use

![./demo.gif](./demo.gif)

```bash
motel hello
```

## Install

- [Homebrew](https://brew.sh/): `brew install bbkane/tap/motel`
- [Scoop](https://scoop.sh/):

```
scoop bucket add bbkane https://github.com/bbkane/scoop-bucket
scoop install bbkane/motel
```

- Download Mac/Linux/Windows executable: [GitHub releases](https://github.com/bbkane/motel/releases)
- Go: `go install go.bbkane.com/motel@latest`
- Build with [goreleaser](https://goreleaser.com/) after cloning: `goreleaser release --snapshot --clean`

## Notes

See [Go Project Notes](https://www.bbkane.com/blog/go-project-notes/) for notes on development tooling.
