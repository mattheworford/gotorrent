# GoTorrent

[![Go Report Card](https://goreportcard.com/badge/github.com/mattheworford/gotorrent)](https://goreportcard.com/report/github.com/mattheworford/gotorrent)
[![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/mattheworford/gotorrent)
[![LICENSE](https://img.shields.io/github/license/mattheworford/gotorrent.svg?style=flat-square)](https://github.com/mattheworford/gotorrent/LICENSE)

GoTorrent is a lightweight BitTorrent client implemented in Go. It allows users to download and share files using the BitTorrent protocol. The client is designed to be fast, efficient, and easy to use.

## Features

- Supports downloading and uploading of files using the BitTorrent protocol.
- Lightweight and resource-efficient design.
- Cross-platform compatibility: works on Windows, macOS, and Linux.
- Simple and intuitive command-line interface.
- Built-in support for magnet links and torrent files.
- Configurable settings for maximum download/upload speed, port number, etc.
- Uses concurrency and parallelism to maximize performance.

## Installation

To install GoTorrent, you need to have Go installed on your system. Once you have Go installed, you can install GoTorrent using the following command:

```bash
go get -u github.com/mattheworford/gotorrent
