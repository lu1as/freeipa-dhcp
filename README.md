# freeipa-dhcp

A tool for generating DHCP hosts files from FreeIPA DNS zones.

Currently supported DHCP servers are: `dnsmasq`, `dhcpd`

## Description

The program uses the FreeIPA API for listing all hosts. After a successful request each host will be split into hostname and domain.
The domain is compared with the given `dhcpZone`, if it's equal another API request is started to find out the host IP.
Make sure that all FreeIPA host entries have a MAC address, otherwise they will be skipped.

## Usage

```
freeipa-dhcp [parameter] ...

parameter:
  -debug
        set log level to debug
  -dhcpHostsFile string
        dhcp hostsfile path (default "freeipa-dhcp.hosts")
  -dhcpServer string
        dhcp server (options: dnsmasq, dhcpd) (default "dnsmasq")
  -dhcpZone string
        dns zone (default "example.de")
  -ipaInsecure
        allow invalid ssl certificate
  -ipaPassword string
        freeipa user password (default "secret")
  -ipaServer string
        freeipa server address (default "https://master.ipa.example.de")
  -ipaUser string
        freeipa user name (default "admin")
  -updateInterval duration
        update interval for new hosts (default 1m0s)
```
