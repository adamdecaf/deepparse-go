#!/usr/bin/env bash
# Install golangci-lint into ./bin with curl.
# moov-io/infra lint-project.sh uses wget -q, which exits 8 on GitHub HTTP
# errors and prints nothing. Pre-placing the binary lets that script reuse it.
set -euo pipefail

version="${GOLANGCI_LINT_VERSION:-v2.15.0}"
if [[ "$version" == "latest" ]]; then
  version=$(curl -fsSL -o /dev/null -w '%{url_effective}' \
    https://github.com/golangci/golangci-lint/releases/latest \
    | grep -Eo 'v[0-9][0-9.]*$')
  if [[ -z "$version" ]]; then
    echo "failed to resolve golangci-lint latest tag" >&2
    exit 1
  fi
fi
tag="v${version#v}"
ver="${tag#v}"

if [[ -x ./bin/golangci-lint ]] &&
  ./bin/golangci-lint version 2>/dev/null | grep -qE "version v?${ver}($|[^0-9.])"; then
  echo "Reusing existing ./bin/golangci-lint for ${tag}"
  exit 0
fi

os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)
case "$arch" in
  x86_64) arch=amd64 ;;
  aarch64 | arm64) arch=arm64 ;;
esac

name="golangci-lint-${ver}-${os}-${arch}"
base="https://github.com/golangci/golangci-lint/releases/download/${tag}"

mkdir -p ./bin
tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT

echo "Downloading ${name}.tar.gz"
curl -fsSL --retry 5 --retry-all-errors --retry-delay 2 \
  -o "${tmpdir}/${name}.tar.gz" "${base}/${name}.tar.gz"
curl -fsSL --retry 5 --retry-all-errors --retry-delay 2 \
  -o "${tmpdir}/checksums.txt" "${base}/golangci-lint-${ver}-checksums.txt"

sha_cmd="sha256sum"
command -v sha256sum >/dev/null 2>&1 || sha_cmd="shasum -a 256"
(cd "$tmpdir" && grep " ${name}.tar.gz\$" checksums.txt | $sha_cmd -c -)

tar -xzf "${tmpdir}/${name}.tar.gz" -C ./bin --strip-components=1 "${name}/golangci-lint"
chmod +x ./bin/golangci-lint
./bin/golangci-lint version
