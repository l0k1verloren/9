package bola

import (
	"github.com/minio/highwayhash"
)

type HHMAC struct {
	// counter allows us to potentially recover to a specific hashchain position
	// or specify a starting point, the latter should be used to extend the
	// useful duration of usability of a master seed key
	// the key for each client should likewise be generated by a derived hash
	// based on the client's nonce value (generated for each run), so it is two
	// steps removed from the original.
	counter int
	key     []byte // must be 32 bytes long
}

// New returns a HighwayHash Message Authentication Cipher
//
// This MAC function uses the very fast HighwayHash128 function to generate a
// message checksum, and the HighwayHash256 function is used with a hashchain
// based on the original preshared key for the server, to generate a difficult
// to reverse authenticated hash to prevent tampering with the message contents
//
// The general concept is similar to the time or sequence based 2FA's used with
// Google Authenticator and similar.
//
// The lifespan of the keys given no leaking should  suffice to not require
// changing of the master key more than once a month for sufficient anti-DoS
// protection for clusters connected over untrusted connections. The rolling
// over of the master key could be automated using SSH scripting.
//
// The master key is intended to be held by both workers and server, and the
// client generates and sends a nonce to start a session, thus creating a new
// key for every run, providing the entropy that should keep the secret safe and
// allow it to be used for longer periods to prevent effective message spoofing
// for DoS or attack purposes.
//
// The security level is not high because the time during which the data it
// protects need only cover a few seconds opportunity to attack the message
// exchange with false and broken and crafted packets. The initial intended use
// case is for securing the connection between workers and a full node for
// solving blocks in a cryptocurrency, and is not encrypted because the data
// that is transmitted is incomplete (no transactions are relayed, only prevhash
// merkle, and the other way, nonce, version and timestamp). Without the tx data
// an attacker cannot formulate the winning block so our only job is to ensure
// it gets published before the attacker might be able to syphon this data out
// by some other method.
//
// nil return indicates key is incorrect or other
func New(key []byte) (hm *HHMAC) {
	if len(key) == 32 {
		copy(hm.key, key)
	}
	return
}

// Inc should be called to advance the counter after a reply is received to
// roll the secret required to get a correct checksum/signature
func (hm *HHMAC) Inc() *HHMAC {
	hm.counter++
	key := highwayhash.Sum(hm.key, hm.key)
	hm.key = key[:]
	return hm
}

// Sum returns a 16 byte hash based on the next element of the hashchain
func (hm *HHMAC) Sum(data []byte) []byte {
	hm.counter++
	key := highwayhash.Sum(hm.key, hm.key)
	hm.key = key[:]
	mhash := highwayhash.Sum128(data, hm.key)
	return mhash[:]
}
