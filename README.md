High speed bls12-381 implementation in Go.

#### Pairing Instance

A Group instance or a pairing engine instance _is not_ suitable for concurrent processing since an instance has its own preallocated memory for temporary variables. A new instance must be created for each thread.

#### Base Field

x86 optimized base field is generated with [kilic/fp](https://github.com/kilic/fp) and for native go is generated with [goff](https://github.com/ConsenSys/goff). Generated codes are slightly edited in both for further requirements.

#### Scalar Field

Standart big.Int module is currently used for scalar field implementation. x86 optimized faster field implementation is planned to be added.

#### Serialization

Point serialization is in line with [zkcrypto library](https://github.com/zkcrypto/pairing/tree/master/src/bls12_381#serialization).

#### Benchmarks

on _2.7 GHz i5_

```
BenchmarkPairing  1145435 ns/op
```
