package transport

import "github.com/yazmeyaa/fragmento"

type ByteCodec struct {
}

func (c *ByteCodec) Encode(id uint32, data []byte) []fragmento.Fragment {
	// TODO: Make Protobuffer encoding.
	return fragmento.FragmentData(id, data)
}

func (c *ByteCodec) Decode(data []fragmento.Fragment) []byte {
	return fragmento.FromFragments(data)
}
