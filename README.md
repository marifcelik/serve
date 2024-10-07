# Serve - Simple HTTP File Server

A minimalist HTTP file server written in Go that serves files from a specified directory. It provides functionality similar to `python -m http.server`.

## Usage Examples

- Serve current directory: `./serve`
- Specify port and directory: `./serve -p 3000 /path/to/directory`
- Serve on all interfaces: `./serve -h 0.0.0.0 -p 8000 /path/to/directory`

## Installation

To install Serve, run:

```sh
go install github.com/marifcelik/serve@latest
```