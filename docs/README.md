# Quick start

Run HLC (composite of HL services). The following services are started by default: HLK, HLS=messenger, HLA=tcp.

```bash
$ go install github.com/number571/hidden-lake/cmd/hlc@latest
$ hlc --network oi4r9NW9Le7fKF9d
```

The list of networks can be found here: [hidden-lake/build/networks.yml](https://github.com/number571/hidden-lake/blob/master/build/networks.yml).

After launching, open the browser and go to `localhost:9591` (HLS=messenger).

To start communicating with someone on the network, you need to follow the following list of actions:
1. Go to `Settings` (with HLS=messenger) and click on the `Copy key` button. Your public key will be copied to the clipboard,
2. Send your public key to the person you want to contact, as well as receive from him, already his, public key. That is, you just need to `exchange keys` with each other,
3. Log in to `Friends`, enter any nickname in the `Alias` field (this will be the naming of the contact) and insert the public key in the `Key` field. Next, click the ◀ button to add a friend.

As a result, when a friend does the same list of actions, you can start chatting. In order to check if a friend is online, you can send a `ping` request to HLS=messenger.

You can also commit your public key here: [number571/hidden-public-keys](https://github.com/number571/hidden-public-keys) to make it easier for people to contact you.

## For developers

The Hidden Lake anonymous network can be supplemented with new functions and features in several ways:
1. `High-level` method. Write applications and adapters through a micro-service architecture. With this method, you can write services in any convenient programming language or technology - you just need to adhere to the `HLK API`. [Examples](../examples).
2. `Middle-level` method. Use the `pkg/network` package, which is located inside the Hidden Lake project. With this method, you can write applications without using the micro-service architecture, but at the same time you gain a dependency on the Go programming language. [Examples](../pkg/network/examples).
3. `Low-level` method. Use the `go-peer` library. With this method, you can significantly change the work and specifics of the network, including the ability to eliminate traffic anonymization, leaving only E2E encryption. This approach should be chosen only if compatibility with the Hidden Lake anonymous network specification is not required. [Examples](https://github.com/number571/go-peer/tree/master/pkg/anonymity/examples).
