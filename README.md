# pooled-bitset

*Go library to manipulate a pool of bitsets*

[![Master Build Status](https://secure.travis-ci.org/mcuelenaere/pooled-bitset.png?branch=master)](https://travis-ci.org/mcuelenaere/pooled-bitset?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mcuelenaere/pooled-bitset)](https://goreportcard.com/report/github.com/mcuelenaere/pooled-bitset)
[![GoDoc](https://godoc.org/github.com/mcuelenaere/pooled-bitset?status.svg)](http://godoc.org/github.com/mcuelenaere/pooled-bitset)

Based on [willf/bitset](https://github.com/willf/bitset).

## Description

A library used for manipulating large'ish, dense bitsets.
It is optimized for efficient use of memory (by reusing memory from an object pool) and for performance (by implementing
the actual bit operations in hand-written assembly, using processor extensions like SSE2/AVX).

Currently it provides the following operations:
  * setting/clearing/flipping/testing individual bits
  * bitset operations:
    * AND
    * OR
    * XOR
    * NOT

## Installation

```bash
go get github.com/mcuelenaere/pooled-bitset
```