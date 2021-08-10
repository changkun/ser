# ser

a simple http server as replacement of python -m http.server

## Install

```
go install changkun.de/x/ser@latest
```

## Usage

```
$ ser --help
ser is a simple http server.

Command line usage:

$ ser [--help] [-addr <addr>] [-p <port>] [<dir>]

options:
  --help
        print this message
  -addr string
        address for listening (default "localhost")
  -p string
        port for listening (default "8080")

examples:
ser
ser .
        serve . directory using port 8080
ser -p 8088
        serve . directory using port 8088
ser ..
        serve .. directory using port 8080
ser -p 8088 ..
        serve .. directory using port 8088
ser -addr 0.0.0.0 -p 9999
        server . directory using address 0.0.0.0:9999
```

## License

MIT &copy; 2021 [Changkun Ou](https://changkun.de)