# Configs

By default the Hidden Lake network uses `adaptive` mode configuration. This means that the connection between the nodes is carried out using an adapter that is not typical of the default network protocol (HTTP). Currently, such an adapter is `HLA=tcp`, which represents a TCP connection between many other nodes or repeaters. The multiplicative mode allows you to use multiple adapters using a single common HTTP-based adapter.

## Modes

1. Classic: [config](../configs/classic), [example](../examples/echo_service/modes/classic)
2. Adaptive: [config](../configs/adaptive), [example](../examples/echo_service/modes/adaptive)
3. Multiplicative: [config](../configs/multiplicative), [example](../examples/echo_service/modes/multiplicative)

More information about modes in research paper: [hidden_lake_anonymous_network.pdf](../docs/hidden_lake_anonymous_network.pdf)
