dirhash (v0.1.0)
========

## About

Recursively walk directory trees and print message digest of regular files.

## Build

Run make(1) or gmake(1).

    $ make

## Usage

    $ ./dirhash
    usage: ./dirhash: [<options>] <paths>
      -abs
            Use absolute path in output
      -debug
            Enable debug print
      -h    Print usage and exit
      -hash_algo string
            Hash algorithm to use (default "sha256")
      -hash_verify string
            Message digest to verify in hex string
      -v    Print version and exit
      -verbose
            Enable verbose print

## Resource

[https://github.com/kusumi/dirhash/](https://github.com/kusumi/dirhash/)
