<img src="images/hl_logo.png" alt="hl_logo.png"/>

<h2>
	<p align="center">
    	<strong>
        	Theoretically Provable Anonymous Network
   		</strong>
	</p>
	<p align="center">
        <a href="https://github.com/topics/golang">
        	<img src="https://img.shields.io/github/go-mod/go-version/number571/go-peer" alt="Go" />
		</a>
        <a href="https://github.com/number571/go-peer/releases">
        	<img src="https://img.shields.io/github/v/release/number571/go-peer.svg" alt="Release" />
		</a>
        <a href="https://github.com/number571/go-peer/blob/master/LICENSE">
        	<img src="https://img.shields.io/github/license/number571/go-peer.svg" alt="License" />
		</a>
        <a href="https://github.com/number571/hidden-lake/actions">
        	<img src="https://github.com/number571/hidden-lake/actions/workflows/build.yml/badge.svg" alt="Build" />
		</a>
		<a href="https://github.com/number571/go-peer/blob/d06ff1b7d35ceb8fa779acda2e1335896b0afdb1/cmd/Makefile#L50">
        	<img src="test/result/badge_coverage.svg" alt="Coverage" />
		</a>
        <a href="https://pkg.go.dev/github.com/number571/hidden-lake/cmd/hidden-lake?status.svg">
        	<img src="https://godoc.org/github.com/number571/go-peer?status.svg" alt="GoDoc" />
		</a>
        <a href="https://github.com/number571/go-peer">
        	<img src="https://raw.githubusercontent.com/number571/go-peer/refs/heads/master/images/go-peer_badge.svg" alt="Go-Peer" />
		</a>
	</p>
    <p align="center">
        <a href="https://goreportcard.com/report/github.com/number571/hidden-lake">
        	<img src="https://goreportcard.com/badge/github.com/number571/hidden-lake" alt="GoReportCard" />
		</a>
        <a href="https://github.com/number571/hidden-lake/pulse">
        	<img src="https://img.shields.io/github/commit-activity/m/number571/hidden-lake" alt="Activity" />
		</a>
		<a href="https://github.com/number571/hidden-lake/commits/master">
        	<img src="https://img.shields.io/github/last-commit/number571/hidden-lake.svg" alt="Commits" />
		</a>
        <a href="https://github.com/number571/hidden-lake/blob/b27339aa283eb137e680a9ca6a04391e7960510a/Makefile#L107">
        	<img src="test/result/badge_codelines.svg" alt="Code Lines" />
		</a>
        <a href="https://img.shields.io/github/languages/code-size/number571/hidden-lake.svg">
        	<img src="https://img.shields.io/github/languages/code-size/number571/hidden-lake.svg" alt="CodeSize" />
		</a>
        <a href="https://img.shields.io/github/downloads/number571/hidden-lake/total.svg">
        	<img src="https://img.shields.io/github/downloads/number571/hidden-lake/total.svg" alt="Downloads" />
		</a>
    </p>
    <p align="center">
        <a href="https://github.com/croqaz/awesome-decentralized">
        	<img src="https://awesome.re/mentioned-badge.svg" alt="Awesome-Decentralized" />
		</a>
        <a href="https://github.com/redecentralize/alternative-internet">
        	<img src="https://awesome.re/mentioned-badge.svg" alt="Alternative-Internet" />
		</a>
        <a href="https://github.com/number571/awesome-anonymity">
        	<img src="https://awesome.re/mentioned-badge.svg" alt="Awesome-Anonymity" />
		</a>
    </p>
	About project
</h2>

> [!IMPORTANT]
> The project is being actively developed, the implementation of some details may change over time. More information about the changes can be obtained from the [CHANGELOG.md](CHANGELOG.md) file.

The `Hidden Lake` is an anonymous network built on a `micro-service` architecture. At the heart of HL is the core - `HLS` (service), which generates anonymizing traffic and combines many other services (for example, `HLF` and `HLM`). Thus, Hidden Lake is not a whole and monolithic solution, but a composition of several combined services. The HL is a `friend-to-friend` (F2F) network, which means building trusted communications. Due to this approach, members of the HL network can avoid `spam` in their direction, as well as `possible attacks` if vulnerabilities are found in the code.

## Coverage map

<p align="center"><img src="test/result/coverage.svg" alt="coverage.svg"/></p>

## Releases

All cmd programs are compiled for {`amd64`, `arm64`} ARCH and {`windows`, `linux`, `darwin`} OS as pattern = `appname_arch_os`. In total, one application is compiled into six versions. The entire list of releases can be found here: [github.com/number571/hidden-lake/releases](https://github.com/number571/hidden-lake/releases "releases"). 

## Dependencies

1. Go library [github.com/number571/go-peer](https://github.com/number571/go-peer "go-peer") (used by `cmd/hls,cmd/hle`)
2. Go library [golang.org/x/net](https://golang.org/x/net "x/net") (used by `cmd/hlm`)
3. CSS/JS library [getbootstrap.com](https://getbootstrap.com "bootstrap") (used by `cmd/hlm,cmd/hlf`)

### Makefile

There are a number of dependencies that represent separate applications for providing additional information about the quality of the code. These applications are not entered into the project, but are loaded via the `make install-deps` command. The list of applications is as follows:

1. golangci-lint [github.com/golangci/golangci-lint/cmd/golangci-lintv1.60.0](https://github.com/golangci/golangci-lint/tree/v1.60.0)
2. go-cover-treemap [github.com/nikolaydubina/go-cover-treemap@v1.4.2](https://github.com/nikolaydubina/go-cover-treemap/tree/v1.4.2)

## List of applications

Basic | Applied | Helpers
:-----------------------------:|:-----------------------------:|:------------------------------:
[HL Service](cmd/hls) | [HL Messenger](cmd/hlm) | [HL Traffic](cmd/hlt)
[HL Composite](cmd/hlc) | [HL Filesharer](cmd/hlf) | [HL Loader](cmd/hll)
[HL Adapters](cmd/hla) | [HL Remoter](cmd/hlr) | [HL Encryptor](cmd/hle)

## How it works

The Hidden Lake anonymous network is based on the (queue-based) `QB-problem`, which can be described by the following list of actions:

1. Each message `m` is encrypted with the recipient's key `k`: `c = Ek(m)`,
2. Message `c` is sent during period `= T` to all network participants,
3. The period `T` of one participant is independent of the periods `T1, T2, ..., Tn` of other participants,
4. If there is no message for the period `T`, then a false message `v` is sent to the network without a recipient (with a random key `r`): `c = Er(v)`,   
5. Each participant tries to decrypt the message they received from the network: `m = Dk(c)`.

<p align="center"><img src="images/hl_qbp.png" alt="hl_qbp.png"/></p>
<p align="center">Figure 1. QB-network with three nodes {A,B,C}</p>

> More information about Hidden Lake in research paper: [hidden_lake_anonymous_network.pdf](docs/hidden_lake_anonymous_network.pdf)

## Build and run

Launching an anonymous network is primarily the launch of an anonymizing HLS service. There are two ways to run HLS: through `source code`, and through the `release version`. 

### 1. Running from source code

```bash
$ go install github.com/number571/hidden-lake/cmd/hls@latest
$ hls
```

### 2. Running from release version

```bash
$ wget https://github.com/number571/hidden-lake/releases/latest/download/hls_amd64_linux
$ chmod +x hls_amd64_linux
$ ./hls_amd64_linux
```

### Production

The HLS node is easily connected to the production environment. To do this, you just need to specify the network_key at startup. You can find them in the [networks.yml](networks.yml) file.

```bash
$ hls -network=oi4r9NW9Le7fKF9d
```

After such a launch, the hls.yml file will be created or overwritten (if it existed). The `settings` and `connections` fields will be substituted in it. When overwriting a file, only the above fields will be changed. The remaining fields of the `friends`, `services`, `address`, etc. type will not be overwritten.

<p align="center"><img src="cmd/hls/images/hls_request.gif" alt="hls_request.gif"/></p>
<p align="center">Figure 2. Example of request to echo-service</p>

> There are also examples of running HL applications in a production environment. For more information, follow the links: [echo_service](examples/anonymity/echo_service/prod_test), [anon_messenger](examples/anonymity/messenger/prod_test), [anon_filesharer](examples/anonymity/filesharer/prod_test).

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=number571/hidden-lake&type=Date)](https://star-history.com/#number571/hidden-lake&Date)

## License

Licensed under the MIT License. See [LICENSE](LICENSE) for the full license text.

**[⬆ back to top](#releases)**
