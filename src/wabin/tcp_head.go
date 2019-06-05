package wabin

const MESSAGE_LEN = 4

type PackHead struct {
	Len uint32
	Cmd uint32
	Uid uint64
	Sid uint64
}

func ByteToHead(l []byte) *PackHead {
	return &PackHead{
		//Len: byteToUint32(l[:4]),
		Cmd: byteToUint32(l[:4]),
		Uid: byteToUint64(l[4:12]),
		Sid: byteToUint64(l[12:20]),
	}
}

func HeadToByte(ph *PackHead, b []byte) {
	Uint32ToByte(ph.Len, b[0:4])
	Uint32ToByte(ph.Cmd, b[4:8])
	Uint64ToByte(ph.Uid, b[8:16])
	Uint64ToByte(ph.Sid, b[16:24])
}
