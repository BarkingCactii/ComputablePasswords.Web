package main

import (
	"hash"
)

type Url struct {
	WebURL string
	Salt   string
	Pepper string
	Hash   hash.Hash
	Result string
}
