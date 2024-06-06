# Script debugging
Experience 2 quick ways to debug bitcoin scripts

Using example: 
create p2wsh script that allows unlocking with the following conditions
- alice's signature with message "VN wins"
- bob's signature with message "TL wins"

then spend that tx on alice


## I. Use btcd lib
- build
```
go build main.go
```
- run
```
./main
```

Result:
```
First tx:  02000000019baf9334d86ce0892c6e24f73993c927e1a9b6e164ed0d335d52dc839b8af4af000000006b483045022100847ece5c342d2cfaeba910fc7a7a2dd1d6c8818fa3d69e6090f58f65de229c5702204ae03d60b89a5bebef776c30bf465a7354d2304f4d4bfce0ba3097eff4e1a80f03210278eb65c765fd52c21e8e8e791537b58aa57385a47d72ac4a5c6757fe7bd0fd8d000000000180ce341d00000000220020ca794c879de127406762c63cb973a26b511763d3016f223d230f74f11c76770400000000
Second tx:  0200000000010137538e5a883535b633127f641e7f7dde3b0cf4616418202e975471c3bc4c906c0000000000000000000100389c1c000000001976a914cc9a79488f73e5b89bde90489df6b4594483c3b488ac04473044022065d2ad829f57919a0ac6972a2646092cbba4498035425ca32337d2ba682a1043022043d85cd3de21a79d59f8529e5fc014dd81a3d37abbce80c0d417e125f0fbd99a03210278eb65c765fd52c21e8e8e791537b58aa57385a47d72ac4a5c6757fe7bd0fd8d07564e2077696e737ba8762039a4acd8ca965d4996f0ab768ca77a5d9c003f8ac5039c1b60e90d537ada45fa87637576a914cc9a79488f73e5b89bde90489df6b4594483c3b4886720538982ea73e9d58e3180797389e63d4def9e7d962d0fdb2093f0260fe106a5ce8876a914c2c69741747f9dc330b5769af0febade430022608868ac00000000
verify success
```

## II. Use btcdeb tool
download: https://github.com/bitcoin-core/btcdeb/tree/master

guide: https://github.com/bitcoin-core/btcdeb/blob/master/doc/btcdeb.md

- Test with 2 above tx 
```
btcdeb --tx=0200000000010137538e5a883535b633127f641e7f7dde3b0cf4616418202e975471c3bc4c906c0000000000000000000100389c1c000000001976a914cc9a79488f73e5b89bde90489df6b4594483c3b488ac04473044022065d2ad829f57919a0ac6972a2646092cbba4498035425ca32337d2ba682a1043022043d85cd3de21a79d59f8529e5fc014dd81a3d37abbce80c0d417e125f0fbd99a03210278eb65c765fd52c21e8e8e791537b58aa57385a47d72ac4a5c6757fe7bd0fd8d07564e2077696e737ba8762039a4acd8ca965d4996f0ab768ca77a5d9c003f8ac5039c1b60e90d537ada45fa87637576a914cc9a79488f73e5b89bde90489df6b4594483c3b4886720538982ea73e9d58e3180797389e63d4def9e7d962d0fdb2093f0260fe106a5ce8876a914c2c69741747f9dc330b5769af0febade430022608868ac00000000 --txin=02000000019baf9334d86ce0892c6e24f73993c927e1a9b6e164ed0d335d52dc839b8af4af000000006b483045022100847ece5c342d2cfaeba910fc7a7a2dd1d6c8818fa3d69e6090f58f65de229c5702204ae03d60b89a5bebef776c30bf465a7354d2304f4d4bfce0ba3097eff4e1a80f03210278eb65c765fd52c21e8e8e791537b58aa57385a47d72ac4a5c6757fe7bd0fd8d000000000180ce341d00000000220020ca794c879de127406762c63cb973a26b511763d3016f223d230f74f11c76770400000000
```

Result:
```
LOG: signing segwit taproot
notice: btcdeb has gotten quieter; use --verbose if necessary (this message is temporary)
input tx index = 0; tx input vout = 0; value = 490000000
got witness stack of size 4
34 bytes (v0=P2WSH, v1=taproot/tapscript)
valid script
- generating prevout hash from 1 ins
[+] COutPoint(6c904cbcc3, 0)
19 op script loaded. type `help` for usage information
script                                                           |
       stack
-----------------------------------------------------------------+-------------------------------------------------------------------
OP_SHA256                                                        |                                                     564e2077696e73
OP_DUP                                                           | 0278eb65c765fd52c21e8e8e791537b58aa57385a47d72ac4a5c6757fe7bd0fd8d
39a4acd8ca965d4996f0ab768ca77a5d9c003f8ac5039c1b60e90d537ada45fa | 3044022065d2ad829f57919a0ac6972a2646092cbba4498035425ca32337d2b...
OP_EQUAL                                                         |
OP_IF                                                            |
OP_DROP                                                          |
OP_DUP                                                           |
OP_HASH160                                                       |
cc9a79488f73e5b89bde90489df6b4594483c3b4                         |
OP_EQUALVERIFY                                                   |
OP_ELSE                                                          |
538982ea73e9d58e3180797389e63d4def9e7d962d0fdb2093f0260fe106a5ce |
OP_EQUALVERIFY                                                   |
OP_DUP                                                           |
OP_HASH160                                                       |
c2c69741747f9dc330b5769af0febade43002260                         |
OP_EQUALVERIFY                                                   |
OP_ENDIF                                                         |
OP_CHECKSIG                                                      |
#0000 OP_SHA256
btcdeb>
```