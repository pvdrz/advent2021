package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

type PayloadKind uint8

const (
	sumKind PayloadKind = iota
	prodKind
	minKind
	maxKind
	litKind
	gtKind
	ltKind
	eqKind
)

type Packet struct {
	version      uint8
	payload_kind PayloadKind
	payload      interface{}
}

func auxVersionSum(packet *Packet, acc uint) uint {
	acc += uint(packet.version)
	if packet.payload_kind != litKind {
		packets := packet.payload.([]Packet)
		for _, sub_packet := range packets {
			acc = auxVersionSum(&sub_packet, acc)
		}
	}
	return acc
}

func (packet *Packet) versionSum() uint {
	return auxVersionSum(packet, 0)
}

func (packet *Packet) eval() (uint64, error) {
	switch packet.payload_kind {
	case sumKind:
		sum := uint64(0)
		packets := packet.payload.([]Packet)
		for _, sub_packet := range packets {
			val, err := sub_packet.eval()
			if err != nil {
				return sum, err
			}
			sum += val
		}
		return sum, nil
	case prodKind:
		prod := uint64(1)
		packets := packet.payload.([]Packet)
		for _, sub_packet := range packets {
			val, err := sub_packet.eval()
			if err != nil {
				return prod, err
			}
			prod *= val
		}
		return prod, nil
	case minKind:
		min := uint64(math.MaxUint64)
		packets := packet.payload.([]Packet)
		if len(packets) == 0 {
			return min, fmt.Errorf("no operands for min")
		}
		for _, sub_packet := range packets {
			val, err := sub_packet.eval()
			if err != nil {
				return min, err
			}
			if min > val {
				min = val
			}
		}
		return min, nil
	case maxKind:
		max := uint64(0)
		packets := packet.payload.([]Packet)
		if len(packets) == 0 {
			return max, fmt.Errorf("no operands for max")
		}
		for _, sub_packet := range packets {
			val, err := sub_packet.eval()
			if err != nil {
				return max, err
			}
			if max < val {
				max = val
			}
		}
		return max, nil
	case litKind:
		return packet.payload.(uint64), nil
	case gtKind:
		packets := packet.payload.([]Packet)
		if len(packets) != 2 {
			return 0, fmt.Errorf("invalid number of operands for gt")
		}
		val1, err := packets[0].eval()
		if err != nil {
			return val1, err
		}
		val2, err := packets[1].eval()
		if err != nil {
			return val2, err
		}

		if val1 > val2 {
			return 1, nil
		} else {
			return 0, nil
		}
	case ltKind:
		packets := packet.payload.([]Packet)
		if len(packets) != 2 {
			return 0, fmt.Errorf("invalid number of operands for lt")
		}
		val1, err := packets[0].eval()
		if err != nil {
			return val1, err
		}
		val2, err := packets[1].eval()
		if err != nil {
			return val2, err
		}

		if val1 < val2 {
			return 1, nil
		} else {
			return 0, nil
		}
	case eqKind:
		packets := packet.payload.([]Packet)
		if len(packets) != 2 {
			return 0, fmt.Errorf("invalid number of operands for eq")
		}
		val1, err := packets[0].eval()
		if err != nil {
			return val1, err
		}
		val2, err := packets[1].eval()
		if err != nil {
			return val2, err
		}

		if val1 == val2 {
			return 1, nil
		} else {
			return 0, nil
		}
	default:
		panic("unreachable")
	}
}

type Bits struct {
	buf       []uint8
	idx       uint8
	read_bits uint
}

func parseBits(input string) (Bits, error) {
	bytes := []uint8{}
	for i := 0; i < len(input); i += 2 {
		hex := input[i : i+2]
		b, err := strconv.ParseUint(hex, 16, 8)
		if err != nil {
			return (Bits{}), err
		}
		bytes = append(bytes, uint8(b))
	}

	return Bits{buf: bytes, idx: 7, read_bits: 0}, nil
}

func (bits *Bits) readBits() uint {
	return bits.read_bits
}

func (bits *Bits) next() (bool, error) {
	if len(bits.buf) == 0 {
		return false, fmt.Errorf("no more bytes")
	}

	bit := bits.buf[0] >> bits.idx

	if bits.idx == 0 {
		bits.buf = bits.buf[1:]
		bits.idx = 7
	} else {
		bits.idx -= 1
	}

	bits.read_bits += 1

	return (bit & 1) == 1, nil
}

func (bits *Bits) getUint(size uint8) (uint, error) {
	value := uint(0)
	for i := size; i > 0; i -= 1 {
		bit_is_one, err := bits.next()
		if err != nil {
			return value, err
		}

        if bit_is_one {
			value += 1 << (i - 1)
		}
	}

	return value, nil
}

func (bits *Bits) parsePacket() (Packet, error) {
	version, err := bits.getUint(3)
	if err != nil {
		return (Packet{}), err
	}

	type_id, err := bits.getUint(3)
	if err != nil {
		return (Packet{}), err
	}

	payload_kind := PayloadKind(type_id)
	if payload_kind > eqKind {
		return (Packet{}), fmt.Errorf("invalid payload kind")
	}

	var payload interface{}
	if payload_kind == litKind {
		literal := uint64(0)

		for {
			keep_going, err := bits.next()
			if err != nil {
				return (Packet{}), err
			}

			bits, err := bits.getUint(4)
			if err != nil {
				return (Packet{}), err
			}

			literal = literal << 4
			literal += uint64(bits)

			if !keep_going {
				break
			}
		}

		payload = literal
	} else {
		length_type_id, err := bits.next()
		if err != nil {
			return (Packet{}), err
		}

		if length_type_id {
			sub_packets_len, err := bits.getUint(11)
			if err != nil {
				return (Packet{}), err
			}

			sub_packets := make([]Packet, sub_packets_len)
			for i := 0; i < int(sub_packets_len); i += 1 {
				sub_packet, err := bits.parsePacket()
				if err != nil {
					return (Packet{}), err
				}
				sub_packets[i] = sub_packet
			}

			payload = sub_packets
		} else {
			bits_len, err := bits.getUint(15)
			if err != nil {
				return (Packet{}), err
			}

			i := uint(0)

			sub_packets := []Packet{}
			curr_bits := bits.readBits()
			for i < bits_len {
				sub_packet, err := bits.parsePacket()
				if err != nil {
					return (Packet{}), err
				}

				prev_bits := curr_bits
				curr_bits = bits.readBits()
				i += curr_bits - prev_bits

				sub_packets = append(sub_packets, sub_packet)
			}
			payload = sub_packets
		}
	}
	return Packet{
		version:      uint8(version),
		payload_kind: payload_kind,
		payload:      payload,
	}, nil
}

func main() {
	bytes, err := ioutil.ReadFile("./input")
	if err != nil {
		panic(err)
	}

	input := strings.Trim(string(bytes), "\n")

	bits, err := parseBits(input)
	if err != nil {
		panic(err)
	}

	packet, err := bits.parsePacket()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Part 1: %d\n", packet.versionSum())

	result, err := packet.eval()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Part 2: %d\n", result)
}
