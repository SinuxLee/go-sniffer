[![GitHub license](https://img.shields.io/github/license/40t/go-sniffer.svg?style=popout-square)](https://github.com/sinuxlee/go-sniffer/blob/master/LICENSE)
# go-sniffer
The go-sniffer can capture packets from the network and intercept database requests in the project and parse them into the corresponding statements, 
such as parsing mysql protocol into SQL statements, facilitating daily analysis and debugging. 
No need to modify the project source code, just sniff data requests in the project directly.

## Support List
go-sniffer can run on common operating systems, including macOS, Linux, Unix, and Windows. 
I have provided the relevant executable file in a compressed package, which can be found on the repository's release page. 
Currently supported protocols for capturing include `HTTP` `Redis` `MySQL` `MongoDB` `Kafka` and `MSSQL`
If the protocol you want is not supported, feel free to submit an issue or pull request.

![image](https://github.com/sinuxlee/go-sniffer/raw/master/images/demo.gif)

## Install:
Before downloading and installing, you need to pre-install the `libpcap` library.

Centos
``` bash
$ yum -y install libpcap-devel
```
Ubuntu
``` bash
$ apt-get install libpcap-dev
```

MacOs
``` bash
$ brew install libpcap
```

Windows
``` bash
$ choco install nmap -y
```

You can now proceed with the official installation.  
If you encounter problems in the `go install` process, try upgrading the go version(1.18+)
``` bash
$ go install -v -u github.com/sinuxlee/go-sniffer@latest
$ cp -rf $(go env GOPATH)/bin/go-sniffer /usr/local/bin
$ go-sniffer
```

## Example:
``` bash
$ go-sniffer lo0 mysql 
$ go-sniffer en0 redis 
$ go-sniffer eth0 http -p 8080
$ go-sniffer eth1 mongodb
$ go-sniffer eth0 mssql 
$ go-sniffer eth0 kafka 
```

## TODO
- [x] Introduction on GitHub Page
- [x] Release product which building by GitHub Action
- [ ] FAQ in WIKI
- [ ] Choose device and plugin via TUI
- [ ] Unit test
- [ ] print db's user and database name

## License:
[MIT](http://opensource.org/licenses/MIT)
