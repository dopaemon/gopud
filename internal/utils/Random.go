package utils

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"os"
	"time"
)

func gatherEntropy() []byte {
	seed := make([]byte, 32)
	_, _ = rand.Read(seed)

	h := sha256.New()
	hn, _ := os.Hostname()
	binary.Write(h, binary.LittleEndian, time.Now().UnixNano())
	h.Write([]byte(hn))
	pid := int32(os.Getpid())
	binary.Write(h, binary.LittleEndian, pid)

	sum := h.Sum(nil)
	for i := 0; i < 32; i++ {
		seed[i] ^= sum[i%len(sum)]
	}
	return seed
}

func aesCTRStream(key, iv []byte, outLen int) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		out := make([]byte, outLen)
		_, _ = rand.Read(out)
		return out
	}
	stream := make([]byte, outLen)
	ctr := make([]byte, aes.BlockSize)
	copy(ctr, iv)
	for i := 0; i < outLen; i += aes.BlockSize {
		block.Encrypt(ctr, ctr)
		chunk := aes.BlockSize
		if i+chunk > outLen {
			chunk = outLen - i
		}
		for j := 0; j < chunk; j++ {
			stream[i+j] = ctr[j] ^ byte((i+j)^int(ctr[j]))
		}
		for k := 0; k < len(ctr); k++ {
			ctr[k]++
			if ctr[k] != 0 {
				break
			}
		}
	}
	return stream
}

func mixAndShuffle(stream []byte) []byte {
	n := 32
	arr := make([]byte, n)
	for i := 0; i < n; i++ {
		arr[i] = byte(i)
	}
	for i := n - 1; i > 0; i-- {
		idx := (i * 4) % len(stream)
		v := binary.LittleEndian.Uint32(stream[idx : idx+4])
		j := int(v % uint32(i+1))
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func complicatedDigits(length int) string {
	seed := gatherEntropy()
	key := seed[:16]
	iv := seed[16:32]

	stream := aesCTRStream(key, iv, 128)
	extraHash := sha256.Sum256(stream)
	shuffle := mixAndShuffle(stream)

	digits := make([]byte, length)
	for i := 0; i < length; i++ {
		b1 := stream[i%len(stream)]
		b2 := stream[(i+37)%len(stream)]
		b3 := extraHash[i%len(extraHash)]
		b4 := shuffle[i%len(shuffle)]
		x := uint32(b1) ^ (uint32(b2)<<3 | uint32(b2)>>5)
		x += uint32(b3)*uint32(b4) + uint32((i+17)*(i+31))
		t := uint32(time.Now().UnixNano()>>8) ^ uint32(seed[i%len(seed)])
		x ^= t
		d := byte(x % 10)
		digits[i] = '0' + d
	}

	final := make([]byte, length)
	copy(final, digits)
	hash2 := sha256.Sum256(final)
	for i := 0; i < length/2; i++ {
		a := int(hash2[i]) % length
		b := int(hash2[i+16]) % length
		final[a], final[b] = final[b], final[a]
	}

	return string(final)
}

func Random(n int) string {
	return complicatedDigits(n)
}
