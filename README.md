MDDNS: Matt's Dynamic DNS
=========================

A lightweight and simpler alternative to DDClient.

## Features

- **Lightweight**: less than 1MB binary
- **Runs anywhere**: runs in every architecture with no dependencies
- **Easy to configure**: simple text file configuration using uris
- **Excellent support**: our goal is to cover every possible provider in existence
- **Modular Architecture**: implement your custom provider
- **Flexible**: use as a program or a library

## Supported Providers

| Name                     | Status      | Package Path                                           | Example Config                    |
|--------------------------|-------------|--------------------------------------------------------|-----------------------------------|
| [Google Domains][google] | Implemented | `github.com/manavarrocarter/mddns/provider/google`     | `google://user:pass@hostname.com` |
| DynDNS                   | Planned     | `github.com/manavarrocarter/mddns/provider/dyndns`     | `dyndns://`                       |
| Zoneedit                 | Planned     | `github.com/manavarrocarter/mddns/provider/zoneedit`   | `zoneedit://`                     |
| EasyDNS                  | Planned     | `github.com/manavarrocarter/mddns/provider/easydns`    | `easydns://`                      |
| NameCheap                | Planned     | `github.com/manavarrocarter/mddns/provider/namecheap`  | `namecheap://`                    |
| Noip                     | Planned     | `github.com/manavarrocarter/mddns/provider/noip`       | `noip://`                         |
| Freedns                  | Planned     | `github.com/manavarrocarter/mddns/provider/freedns`    | `freedns://`                      |
| CloudFlare               | Planned     | `github.com/manavarrocarter/mddns/provider/cloudflare` | `cloudflare://`                   |
| GoDaddy                  | Planned     | `github.com/manavarrocarter/mddns/provider/godaddy`    | `godaddy://`                      |
| DuckDNS                  | Planned     | `github.com/manavarrocarter/mddns/provider/duckdns`    | `duckdns://`                      |

[google]: https://support.google.com/domains/answer/6147083?authuser=0&hl=en-GB#zippy=%2Cuse-the-api-to-update-your-dynamic-dns-record