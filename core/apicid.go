package core

import (
	"encoding/json"

	cid "gx/ipfs/QmPSQnBKM9g7BaUcZCvswUJVscQ1ipjmwxN5PXCjkp9EQ7/go-cid"
	mbase "gx/ipfs/QmekxXDhCxCJRNuzmHreuaT3BsuJcsjcXWNrtV9C8DRHtd/go-multibase"
	//path "gx/ipfs/QmX7uSbkNz76yNwBhuwYwRbhihLnJqM73VTCjS3UMJud9A/go-path"
)

// CidEncoder is a type used to encode or recode Cid as the user
// specifies
type CidEncoder interface {
	Encode(c cid.Cid) string
	Recode(v string) (string, error)
}

// BasicCidEncoder is a basic CidEncoder that will encode Cid's using
// a specifed base, optionally upgrading a Cid if is Version 0
type BasicCidEncoder struct {
	Base    mbase.Encoder
	Upgrade bool
}

var DefaultCidEncoder = BasicCidEncoder{
	Base:    mbase.MustNewEncoder(mbase.Base58BTC),
	Upgrade: false,
}

// CidJSONBase is the base to use when Encoding into JSON.
//var CidJSONBase mbase.Encoder = mbase.MustNewEncoder(mbase.Base58BTC)
var CidJSONBase mbase.Encoder = mbase.MustNewEncoder(mbase.Base32)

// APICid is a type to respesnt CID in the API
type APICid struct {
	str string // always in CidJSONBase
}

// FromCid created an APICid from a Cid
func FromCid(c cid.Cid) APICid {
	return APICid{c.Encode(CidJSONBase)}
}

// Cid converts an APICid to a CID
func (c APICid) Cid() (cid.Cid, error) {
	return cid.Decode(c.str)
}

func (c APICid) String() string {
	return c.Encode(DefaultCidEncoder)
}

func (c APICid) Encode(enc CidEncoder) string {
	if c.str == "" {
		return ""
	}
	str, err := enc.Recode(c.str)
	if err != nil {
		return c.str
	}
	return str
}

func (c *APICid) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &c.str)
}

func (c APICid) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.str)
}

func (enc BasicCidEncoder) Encode(c cid.Cid) string {
	if c.Version() == 0 {
		c = cid.NewCidV1(c.Type(), c.Hash())
	}
	return c.Encode(enc.Base)
}

func (enc BasicCidEncoder) Recode(v string) (string, error) {
	skip, err := enc.NoopRecode(v)
	if skip || err != nil {
		return v, err
	}

	c, err := cid.Decode(v)
	if err != nil {
		return v, err
	}

	return enc.Encode(c), nil
}

func (enc BasicCidEncoder) NoopRecode(v string) (bool, error) {
	if len(v) < 2 {
		return false, cid.ErrCidTooShort
	}
	ver := cidVer(v)
	skip := ver == 0 && !enc.Upgrade || ver == 1 && v[0] == byte(enc.Base.Encoding())
	return skip, nil
}

func cidVer(v string) int {
	if len(v) == 46 && v[:2] == "Qm" {
		return 0
	} else {
		return 1
	}
}

// func (enc *CidEncoder) Scan(cids ...string) {
// 	if enc.Override == nil {
// 		enc.Override = map[cid.Cid]string{}
// 	}
// 	for _, p := range cids {
// 		//segs := path.FromString(p).Segments()
// 		//v := segs[0]
// 		//if v == "ipfs" && len(segs) > 0 {
// 		//	v = segs[1]
// 		//}
// 		v := p
// 		skip, err := enc.noopRecode(v)
// 		if skip || err != nil {
// 			continue
// 		}
// 		c, err := cid.Decode(v)
// 		if err != nil {
// 			continue
// 		}
// 		enc.Override[c] = v
// 	}
// }
