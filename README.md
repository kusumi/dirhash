dirhash ([v0.4.5](https://github.com/kusumi/dirhash/releases/tag/v0.4.5))
========

## About

Recursively walk directory trees and print message digest of regular files.

## Supported platforms

Unix-likes in general

## Requirements

go 1.18 or above

## Build

    $ make

## Usage

    $ ./dirhash
    usage: dirhash: [<options>] <paths>
      -abs
            Print file paths in absolute path
      -debug
            Enable debug print
      -follow_symlink
            Follow symbolic links unless directory
      -h    Print usage and exit
      -hash_algo string
            Hash algorithm to use (default "sha256")
      -hash_only
            Do not print file paths
      -hash_verify string
            Message digest to verify in hex string
      -ignore_dot
            Ignore entries start with .
      -ignore_dot_dir
            Ignore directories start with .
      -ignore_dot_file
            Ignore files start with .
      -ignore_symlink
            Ignore symbolic links
      -sort
            Print sorted file paths
      -squash
            Print squashed message digest instead of per file
      -swap
            Print file path first in each line
      -v    Print version and exit
      -verbose
            Enable verbose print
