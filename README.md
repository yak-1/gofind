# gofind

A file find utility modeled after the unix `find` written in Go.

## Why

This project is used as an intro to Go exersize. Implementing `find` is great because it allows you to utilize many concepts that you may be already
familiar with from other languages in a quick and concise project to get a feel for how to do those things in Go. In specific, 
this project calls for the use of recursion, tree traversal, CLI, file IO, structs, error handling, closures, lifetimes, regex, and testing.

## Usage
```
% gofind -h
gofind searches the file system for files that match your filters

Usage:
  gofind [flags] PATH

Flags:
  -c, --contains string        string that the file name contains
      --ends-with string       string that the file name ends with
  -e, --ext string             file extension
  -h, --help                   help for gofind
  -i, --icontains string       case insensitive string that the file name contains
      --inot-contains string   case insensitive string that the file name doesn't contain
  -d, --max-depth int          max recursion depth to search
      --not-contains string    string that the file name doesn't contain
  -r, --regex string           regex pattern that the file name matches
      --size-at-least int      min file size
      --size-at-most int       max file size
      --starts-with string     string that the file name starts with
```
