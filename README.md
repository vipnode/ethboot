# ethboot

Bootnode for Ethereum.

Protip: Don't run with `--nodiscover` when trying to use `--bootnodes` for
discovery.

**Status**: Not usable. Went through a bunch of code with a machete and left a bunch of logging to try and reverse engineer the client behaviour. Might be the basis for something in the future.

## Appendix

### References

- https://medium.com/coinmonks/data-structure-in-ethereum-episode-1-recursive-length-prefix-rlp-encoding-decoding-d1016832f919
- https://github.com/ethereum/devp2p/blob/master/rlpx.md
- https://godoc.org/github.com/ethereum/go-ethereum/p2p/discv5
- https://github.com/ethereum/go-ethereum/tree/ccc0debb63124ee99906c2cfff6125de30e8c62f/cmd/bootnode

## License

MIT
