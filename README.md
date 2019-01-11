# collections
[![Build Status](https://www.travis-ci.org/FelixSeptem/collections.svg?branch=master)](https://www.travis-ci.org/FelixSeptem/collections)
[![Coverage Status](https://coveralls.io/repos/github/FelixSeptem/collections/badge.svg?branch=master)](https://coveralls.io/github/FelixSeptem/collections?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/FelixSeptem/collections)](https://goreportcard.com/report/github.com/FelixSeptem/collections)

some useful datatypes inspired by [collections](https://docs.python.org/2/library/collections.html) and [boltons](https://github.com/mahmoud/boltons)
## Install
```shell
go get -u github.com/FelixSeptem/collections
```

### Cache
- LRU [![GoDoc](http://godoc.org/github.com/FelixSeptem/collections/lru?status.svg)](http://godoc.org/github.com/FelixSeptem/collections/lru)
implement a thread safe `Least Recently Used` [ref](https://en.wikipedia.org/wiki/Cache_replacement_policies#Least_recently_used_(LRU)) [Code](https://github.com/FelixSeptem/collections/tree/master/lru)
- LFU [![GoDoc](http://godoc.org/github.com/FelixSeptem/collections/lfu?status.svg)](http://godoc.org/github.com/FelixSeptem/collections/lfu)
implement a thread safe `Least Frequently Used` [ref](https://en.wikipedia.org/wiki/Cache_replacement_policies#Least-frequently_used_(LFU)) [Code](https://github.com/FelixSeptem/collections/tree/master/lfu)
- ARC [![GoDoc](http://godoc.org/github.com/FelixSeptem/collections/arc?status.svg)](http://godoc.org/github.com/FelixSeptem/collections/arc)
implement a thread safe `Adaptive Replacement Cache` [ref](https://en.wikipedia.org/wiki/Adaptive_replacement_cache) Paper:[[1]](https://www.usenix.org/legacy/events/fast03/tech/full_papers/megiddo/megiddo.pdf)[[2]](https://arxiv.org/pdf/1503.07624.pdf) [Code]()

### Others
