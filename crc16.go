package vietqr

import (
	"hash"
	"math/bits"
)

const (
	paddingCrc = 4
)

type CrcParams struct {
	Poly   uint16
	Init   uint16
	RefIn  bool
	RefOut bool
	XorOut uint16
	Check  uint16
	Name   string
}

// predefined CRC-16 algorithms: http://reveng.sourceforge.net/crc-catalogue/16.htm
var (
	CRC16_ARC         = CrcParams{0x8005, 0x0000, true, true, 0x0000, 0xBB3D, "CRC-16/ARC"}
	CRC16_AUG_CCITT   = CrcParams{0x1021, 0x1D0F, false, false, 0x0000, 0xE5CC, "CRC-16/AUG-CCITT"}
	CRC16_BUYPASS     = CrcParams{0x8005, 0x0000, false, false, 0x0000, 0xFEE8, "CRC-16/BUYPASS"}
	CRC16_CCITT_FALSE = CrcParams{0x1021, 0xFFFF, false, false, 0x0000, 0x29B1, "CRC-16/CCITT-FALSE"}
	CRC16_CDMA2000    = CrcParams{0xC867, 0xFFFF, false, false, 0x0000, 0x4C06, "CRC-16/CDMA2000"}
	CRC16_DDS_110     = CrcParams{0x8005, 0x800D, false, false, 0x0000, 0x9ECF, "CRC-16/DDS-110"}
	CRC16_DECT_R      = CrcParams{0x0589, 0x0000, false, false, 0x0001, 0x007E, "CRC-16/DECT-R"}
	CRC16_DECT_X      = CrcParams{0x0589, 0x0000, false, false, 0x0000, 0x007F, "CRC-16/DECT-X"}
	CRC16_DNP         = CrcParams{0x3D65, 0x0000, true, true, 0xFFFF, 0xEA82, "CRC-16/DNP"}
	CRC16_EN_13757    = CrcParams{0x3D65, 0x0000, false, false, 0xFFFF, 0xC2B7, "CRC-16/EN-13757"}
	CRC16_GENIBUS     = CrcParams{0x1021, 0xFFFF, false, false, 0xFFFF, 0xD64E, "CRC-16/GENIBUS"}
	CRC16_MAXIM       = CrcParams{0x8005, 0x0000, true, true, 0xFFFF, 0x44C2, "CRC-16/MAXIM"}
	CRC16_MCRF4XX     = CrcParams{0x1021, 0xFFFF, true, true, 0x0000, 0x6F91, "CRC-16/MCRF4XX"}
	CRC16_RIELLO      = CrcParams{0x1021, 0xB2AA, true, true, 0x0000, 0x63D0, "CRC-16/RIELLO"}
	CRC16_T10_DIF     = CrcParams{0x8BB7, 0x0000, false, false, 0x0000, 0xD0DB, "CRC-16/T10-DIF"}
	CRC16_TELEDISK    = CrcParams{0xA097, 0x0000, false, false, 0x0000, 0x0FB3, "CRC-16/TELEDISK"}
	CRC16_TMS37157    = CrcParams{0x1021, 0x89EC, true, true, 0x0000, 0x26B1, "CRC-16/TMS37157"}
	CRC16_USB         = CrcParams{0x8005, 0xFFFF, true, true, 0xFFFF, 0xB4C8, "CRC-16/USB"}
	CRC16_CRC_A       = CrcParams{0x1021, 0xC6C6, true, true, 0x0000, 0xBF05, "CRC-16/CRC-A"}
	CRC16_KERMIT      = CrcParams{0x1021, 0x0000, true, true, 0x0000, 0x2189, "CRC-16/KERMIT"}
	CRC16_MODBUS      = CrcParams{0x8005, 0xFFFF, true, true, 0x0000, 0x4B37, "CRC-16/MODBUS"}
	CRC16_X_25        = CrcParams{0x1021, 0xFFFF, true, true, 0xFFFF, 0x906E, "CRC-16/X-25"}
	CRC16_XMODEM      = CrcParams{0x1021, 0x0000, false, false, 0x0000, 0x31C3, "CRC-16/XMODEM"}
)

// compatible with hash.Hash
type Hash16 interface {
	hash.Hash
	Sum16() uint16
}

var _ Hash16 = &digest{}

type digest struct {
	algo  CrcParams
	table [256]uint16
	sum   uint16
}

func NewCrc16(algo CrcParams) Hash16 {
	crcTable := [256]uint16{}

	for n := 0; n < 256; n++ {
		crc := uint16(n) << 8
		for i := 0; i < 8; i++ {
			bit := (crc & 0x8000) != 0
			crc <<= 1
			if bit {
				crc ^= algo.Poly
			}
		}

		crcTable[n] = crc
	}

	return &digest{
		algo:  algo,
		table: crcTable,
		sum:   algo.Init,
	}
}

func (h *digest) Write(data []byte) (int, error) {
	for _, d := range data {
		if h.algo.RefIn {
			d = bits.Reverse8(d)
		}
		h.sum = h.sum<<8 ^ h.table[byte(h.sum>>8)^d]
	}

	return len(data), nil
}

func (h digest) Sum(b []byte) []byte {
	s := h.Sum16()
	return append(b, byte(s>>8), byte(s))
}

func (h *digest) Reset() {
	h.sum = h.algo.Init
}

func (h digest) Size() int {
	return 2
}

func (h digest) BlockSize() int {
	return 1
}

func (h digest) Sum16() uint16 {
	if h.algo.RefOut {
		return bits.Reverse16(h.sum) ^ h.algo.XorOut
	}
	return h.sum ^ h.algo.XorOut
}
