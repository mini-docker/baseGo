package cache

import (
	"testing"
	"time"
)

func TestBGet(t *testing.T) {

}

// 会分配接近 500 MB 的内存.
// 83271  ___TestBSet_ 0.0       00:00.05 6     0     15     26M+   0B     0B     19483 83270 sleeping *0[1]            0.00000 0.00000    501  7130+      46       12
func TestBSet(t *testing.T) {
	InitCache(time.Second * 1)
	var data = map[string][]byte{
		"k1": []byte(`123`),
		"k2": []byte(`456`),
	}

	for k, v := range data {
		BSet(k, v)
	}

	b1, err := BGet("k1")
	if err != nil {
		t.Error(`缓存不存在.`)
	}
	t.Log(b1, err)
	b, err := BGet("k2")
	if err != nil {
		t.Error(`缓存不存在.`)
	}
	t.Log(b, err)

	time.Sleep(7 * time.Second)

	b1, err = BGet("k1")
	if err == nil {
		t.Error(`缓存存在.`)
	}
	t.Log(b1, err)
	b, err = BGet("k2")
	if err == nil {
		t.Error(`缓存存在.`)
	}
	t.Log(b, err)

}
