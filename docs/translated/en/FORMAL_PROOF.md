# Formal proof

> A more detailed description of the anonymous Hidden Lake network can be found in the work - [hidden_lake_anonymous_network.pdf](hidden_lake_anonymous_network.pdf)

To assess anonymity in __Hidden Lake__ (QB networks), the apparatus of probability theory and entropy can be used. The main goal is to prove that for an outside observer, including a global one, the uncertainty about the source or recipient of the message is maximum.

## Threat model

Let the network be represented as a set of nodes: $S = \{s_1, s_2, ..., s_N\}$ , $|S| = N$ .

- Global Passive Adversary: The enemy can view all network traffic (all edges of the network graph),
- Nodes: All nodes are considered potentially compromised, except for the sender $Sender$ and recipient $Recipient$.

The threat model is limited exclusively to the network part of the interaction and does not affect the security of the work environment when starting/operating the node.

## Working principle

1. Let there be a __queue__ $Q$ storing ciphertexts for sending: $C = (c_1, c_2, ..., c_x)$ ,
2. If a __true message__$M$ appears, it is encrypted with the recipient's public key $c_i =E_{pk}(M)$ and placed in the queue $Q\leftarrow c_i$ ,
3. If the queue is empty, a __false message__ $R$ is generated, which is encrypted with $c_j=E_{pk_{R}}(R)$, where $pk_{R}\leftarrow CSPRNG$, and placed in the queue of $Q\leftarrow c_j$ ,
4. When the __ time period__ $\Delta t$ arrives, one ciphertext $c\leftarrow Q$ is taken from the queue and __sent to the entire set of nodes $S$ ,
5. If the ciphertext $c$ has been received from the network, then __ an attempt is made to decrypt __ $D_{sk}(c)$ with the private key $sk$. If the ciphertext has been successfully decrypted into plaintext $M$, it means that it has been successfully received. Otherwise, the ciphertext is ignored and considered a node for generated noise.

## Probabilistic assessment

In the system described above, each node transmits packets to all neighbors. Assume that all packets are of fixed length $L$ and are encrypted with a scheme having the property __IND-CPA__ (or higher). If node $A$ sends a message and node $B$ receives it, then for the global observer:

- __Sender probability__: Since all nodes generate packets with the same frequency $\Delta t$ (constant noise), the probability that a particular node $s_i$ is the sender of the message $M$ is: $P(Sender = s_i) =\frac{1}{N}$ ,

- __Recipient probability__: Since the packet is encrypted asymmetrically (using the recipient's key) and sent to everyone, it is impossible to determine who exactly was able to decrypt it. The probability that node $s_j$ is the true recipient of message $M$ is: $P(Recipient =s_j) = \frac{1}{N}$ .

## The entropy metric

The anonymity level of the system is calculated using the Shannon entropy $H(X)$. If the probability distribution between the participants is uniform (which was proved earlier), then the entropy is considered maximum: $H(X) =\sum_{i=1}^{N} p_{i} log_{2}(p_{i}) = log_{2}(N)= H_{max}$ . At maximum entropy, __the degree of anonymity __$d$ becomes equal to $\frac{H(X)}{H_{max}} = 1$, which indicates the presence of theoretical anonymity. In Tor or I2P networks, this indicator is reduced to $d < 1$ due to the possibility of statistical analysis of delays and routes.

Anonymity is considered __perfect__ if the a posteriori probability that the message was sent by $s_i$ is equal to the a priori probability: $P(Sender = s_i | Traffic) = P(Sender = s_i) = \frac{1}{N}$ . In the system described above, this is achieved by sequencing generation, which is why the traffic vector $T = (C_1, C_2, ..., C_x)$, at any given time, has an entropy independent of the initiator of communication: $H(Sender=s_i|T) = H(Sender=s_j|T)$ , $\forall s_i,s_j \in S \cup \varnothing$ . Therefore, the amount of mutual information between the sender and the observed traffic is zero: $I(Sender;T) = H(Sender) - H(Sender|T) = 0$.

Since $H(T)$ depends only on the system parameters $(L,\Delta t,Q)$ and is not a function of the presence of real messages $M$, the mutual information $I(M;T)$ between the fact of sending and the observed traffic is: $I(M;T) = H(M) - H(M|T) = 0$ . This means that knowing the traffic of $T$ does not reduce the uncertainty about the message of $M$, which indicates that the very fact of sending is hidden.
