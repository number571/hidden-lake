# HLA

> Hidden Lake Adapters

<img src="images/hla_logo.png" alt="hla_logo.png"/>

`Hidden Lake Adapters` are based on two functions: Consume and Produce. Due to this, at the interface level, users do not care about the nature of communication: where ciphertexts are read from and where they are written. Due to this property, as well as the properties of QB networks to preserve anonymity in any communication environment, it becomes possible to write adapters not only for network protocols, but also for centralized services, thereby creating secret communication channels.

<p align="center"><img src="images/hla_arch.png" alt="hla_arch.png"/></p>
<p align="center">Figure 1. architecture of HLA using the example of HLA=tcp</p>

## List of adapters

1. [HLA=tcp](hla_tcp) - adapts HL traffic to a custom TCP connection
