# multik

A fast, parallelized multi-cluster Kubernetes CLI written in Go.
Query multiple clusters at once, filter results, and output in human-friendly tables or machine readble JSON.

## v1 Goals

- `get pods` across multiple clusters in parallel
- Filters: namespace, labels, age, status
- Output: table (default) or JSON
- Per-cluster timeouts, partial results
- Human + CI-friendly UX

## Install

```bash
git clone https://github.com/josiarod/multik.git
cd multik
make build
./bin/multik
```