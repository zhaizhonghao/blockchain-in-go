# A blockchain implementation in golang

## Blockchain Module

## Network Module

## Client Module

## Wallet Module

The wallet module is used to create the identity of the participant using the Elliptic Curve Digital Signature Algorithm(ECDSA). The details of the ECDSA can be learned from this [website](https://www.certicom.com/content/certicom/en/10-introduction.html). 

The identity(wallet) mainly consists of the private key and public key. The public is derived from the the private key and can be distributed to other entities. The private key is used to sign the transactions or blocks. ![wallet_1](./imgs/wallet_1.png)

