# A blockchain implementation in golang

## Blockchain Module

## Network Module

## Client Module

## Wallet Module

The wallet module is used to create the identity of the participant using the Elliptic Curve Digital Signature Algorithm(ECDSA). The details of the ECDSA can be learned from this [website](https://www.certicom.com/content/certicom/en/10-introduction.html). 

As show in the following picture, the main function of the wallet is to generate the identity of the participant. The identity is a string which is called *address*. The address is hash of the public key. The private key and public key pair is generated using the Elliptic Curve Digital Signature Algorithm(ECDSA). The details of the ECDSA can be learned from this [website](https://www.certicom.com/content/certicom/en/10-introduction.html). The public is derived from the the private key and can be distributed to other entities. The private key is used to sign the transactions or blocks. 

 <img src="./imgs/wallet_1.png" alt="wallet_1" style="zoom:50%;" />

The details of the address generation is shown as following:

<img src="./imgs/wallet_2.png" alt="wallet_2"  />

