# Escada

A self-hosted alternative to [12ft.io](https://12ft.io).


## Usage

```
escada [-addr=<domain or ip address>] [-port=<port>] [-help]

addr: address to bind to
port: port to bind to
```

By default the program listens on `127.0.0.1:9982`.


## Bookmarklet

The following bookmarklet will redirect the current paywalled page to the `escada` version.

```
javascript:(function(){window.location.href='http://<ADDR>:<PORT>/pages/'+encodeURIComponent(window.location.href);})();
```

Change `ADDR` and `PORT` to the address and port you set for `escada`.

## Create a Service

### Linux/SystemD

```
[Unit]
Description=escada server
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
ExecStart=/usr/local/bin/escada -port=9982

[Install]
WantedBy=default.target
```

### Windows

In an Admin cmd prompt:


```
sc.exe create Escada binpath= "PATH TO SERVICE"
```

## Similar Projects

`escada` is a personal project to learn a bit about using Go as a web service. It was inspired by the following projects:

- https://github.com/wasi-master/13ft
- https://github.com/everywall/ladder


> **Disclaimer:** This project is intended for educational purposes only. The author does not endorse its use for illegal activities.
