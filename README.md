dirhash ([v0.3.2](https://github.com/kusumi/dirhash/releases/tag/v0.3.2))
========

## About

Recursively walk directory trees and print message digest of regular files.

## Supported platforms

Unix-likes in general

## Requirements

go 1.18 or above

## Build

    $ make

or

    $ gmake

## [Usage](doc/dirhash.1.txt)

    $ ./dirhash
    usage: dirhash: [<options>] <paths>
      -abs
            Print file path in absolute path
      -debug
            Enable debug print
      -h    Print usage and exit
      -hash_algo string
            Hash algorithm to use (default "sha256")
      -hash_only
            Do not print file path
      -hash_verify string
            Message digest to verify in hex string
      -ignore_dot
            Ignore entry starts with .
      -ignore_dot_dir
            Ignore directory starts with .
      -ignore_dot_file
            Ignore file starts with .
      -ignore_symlink
            Ignore symbolic link
      -lstat
            Do not resolve symbolic link
      -squash
            Print squashed message digest instead of per file
      -v    Print version and exit
      -verbose
            Enable verbose print

## Resource

[https://github.com/kusumi/dirhash/](https://github.com/kusumi/dirhash/)
