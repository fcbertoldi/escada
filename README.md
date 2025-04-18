# Escada

A self-hosted alternative to [12ft.io](https://12ft.io).

## How it Works

`escada` works as a proxy server. It pretends to be GoogleBot web crawler by changing the `User-Agent` HTTP header of the request to GoogleBot's. It may not work with some paywalled sites that uses other forms of identifying non-bot requests, such as IP ranges.


## Installation

The following instruction will install the executable and create a http service available at http://localhost:9982

### Linux

1. The following shell command will build install the executable into `usr/local/bin`
```sh
make build && sudo make install
```

2. Configure the  [SystemD service](#linuxsystemd)
```sh
mv escada.service <systemd_unit_path>
systemctl enable escada.service
systemctl start escada.service
```

### Windows

1. Open a Windows command prompt with admin priviledges and cd into the repo. The following will build and install the executable into %PROGRAMFILES%Escada, and

```bat
build.cmd
install.cmd
```

2. [Create a service](#windows-service)

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

### Windows Service

In an Admin cmd prompt:


```
sc.exe create Escada binpath= "PATH TO SERVICE"
```

## Similar Projects

`escada` is a personal project to learn a bit about using Go as a web service. It was inspired by the following projects:

- https://github.com/wasi-master/13ft
- https://github.com/everywall/ladder


> **Disclaimer:** This project is intended for educational purposes only. The author does not endorse its use for illegal activities.
