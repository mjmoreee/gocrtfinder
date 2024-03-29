# GoCrtFinder
Simple golang tools to extract all subdomain from crt.sh

This tool is inspired by https://github.com/eslam3kl/crtfinder. I simply converted it from Python to Golang due to the need for automating the recon process during bug hunting, and it's used for a private framework I employ for hunting.

## Installation
```Bash
git clone https://github.com/mjmoreee/gocrtfinder
cd gocrtfinder && go build gocrtfinder.go
```

## Usage
```Bash
gocrtfinder -u target.com
```
