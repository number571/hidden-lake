# HLA

> Hidden Lake Adapters

<img src="images/hla_logo.png" alt="hla_logo.png"/>

The `Hidden Lake Adapters` are the main way to exchange data between multiple HLS processes using network protocols. 

> More information about HLA in the [habr.com/ru/post/720544](https://habr.com/ru/post/720544/ "Habr HLA")

## How it works

By default, [HLS](https://github.com/number571/hidden-lake/tree/master/cmd/hls) contains `HLA=http`, which allows you to produce and consume ciphertexts over the HTTP protocol. This method works ideally in a local microservice environment, where the main way of communication between services is via HTTP. However, in a global environment, HTTP is not a good fit, because it requires all participants to have a public IP address. Because of this, `HLA=tcp`, based on the TCP protocol, is becoming the most popular on the Hidden Lake network.

<p align="center"><img src="images/hla_arch.jpg" alt="hla_arch.jpg"/></p>
<p align="center">Figure 1. Architecture of HLA.</p>

## Example 

Adapters are based on two functions: Consume and Produce. Due to this, at the interface level, users do not care about the nature of communication: where ciphertexts are read from and where they are written. Due to this property, as well as the properties of QB networks to preserve anonymity in any communication environment, it becomes possible to write adapters not only for network protocols, but also for centralized services, thereby creating secret communication channels.
