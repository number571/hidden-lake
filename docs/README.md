# Quick start

Run HLC (composite of HL services). The following services are started by default: HLS, HLM, HLF.

```bash
$ go install github.com/number571/hidden-lake/cmd/hlc@latest
$ hlc -network=oi4r9NW9Le7fKF9d
```

The list of networks can be found here: [hidden-lake/build/networks.yml](https://github.com/number571/hidden-lake/blob/master/build/networks.yml).

After launching, open the browser and go to `localhost:9591` (HLM) or `localhost:9541` (HLF).

To start communicating with someone on the network, you need to follow the following list of actions:
1. Go to `Settings` (with HLM or with HLF) and click on the `Copy key` button. Your public key will be copied to the clipboard,
2. Send your public key to the person you want to contact, as well as receive from him, already his, public key. That is, you just need to `exchange keys` with each other,
3. Log in to `Friends`, enter any nickname in the `Alias` field (this will be the naming of the contact) and insert the public key in the `Key` field. Next, click the ◀ button to add a friend.

As a result, when a friend does the same list of actions, you can start chatting. In order to check if a friend is online, you can send a `ping` request to HLM.

In order to disable HLM or HLF (for example, because it is not needed), the hidden-lake-messenger or hidden-lake-filesharer line should be deleted (or commented) in `hlc.yml`, respectively. It would also be better to delete / comment a similar line in the `hls.yml` file. After that, you just need to restart HLC.

You can also commit your public key here: [number571/hidden-public-keys](https://github.com/number571/hidden-public-keys) to make it easier for people to contact you.
