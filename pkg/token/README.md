# Token Metadata for Go

## Usage
Import this repo directly into your own go repo

## Maintain
1. update token meta submodule by `git submodule update --remote --merge`
2. run `make gen` to copy static token meta file to embed
3. develop your new util functions
4. don't forget to run `make build` to make sure no error in compiling at least
ps: better to not edit & push submodule to remote