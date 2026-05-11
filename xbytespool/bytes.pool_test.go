package xbytespool

import (
	"math"
	"os"
	"testing"
	"unsafe"
)

func TestIndex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   uint32
		want uint32
	}{
		{name: "one", in: 1, want: 0},
		{name: "two", in: 2, want: 1},
		{name: "three", in: 3, want: 2},
		{name: "four", in: 4, want: 2},
		{name: "five", in: 5, want: 3},
		{name: "eight", in: 8, want: 3},
		{name: "nine", in: 9, want: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := index(tt.in); got != tt.want {
				t.Fatalf("index(%d) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestPoolGet(t *testing.T) {
	t.Run("non-positive size returns nil", func(t *testing.T) {
		var p Pool

		for _, size := range []int{0, -1} {
			if got := p.Get(size); got != nil {
				t.Fatalf("Get(%d) = %#v, want nil", size, got)
			}
		}
	})

	t.Run("large size allocates directly", func(t *testing.T) {
		if math.MaxInt <= math.MaxInt32 {
			t.Skip("当前平台无法构造大于 MaxInt32 的切片长度")
		}

		if os.Getenv("XBYTESPOOL_RUN_LARGE_MEMORY_TESTS") == "" {
			t.Skip("默认跳过超过 2GiB 内存分配场景；如需执行请设置 XBYTESPOOL_RUN_LARGE_MEMORY_TESTS=1")
		}

		const size = math.MaxInt32 + 1

		if testing.AllocsPerRun(1, func() { _ = (&Pool{}).Get(size) }) == 0 {
			t.Fatal("Get(>MaxInt32) 应直接分配新的切片")
		}
	})

	t.Run("allocates sized slice when pool is empty", func(t *testing.T) {
		var p Pool

		got := p.Get(3)
		if got == nil {
			t.Fatal("Get(3) returned nil")
		}
		if len(got) != 3 {
			t.Fatalf("len(Get(3)) = %d, want 3", len(got))
		}
		if cap(got) != 4 {
			t.Fatalf("cap(Get(3)) = %d, want 4", cap(got))
		}
	})

	t.Run("reuses pooled pointer", func(t *testing.T) {
		var p Pool
		backing := make([]byte, 8)
		backing[0] = 7

		p.pools[3].Put(unsafe.Pointer(&backing[0]))

		got := p.Get(5)
		if len(got) != 5 {
			t.Fatalf("len(Get(5)) = %d, want 5", len(got))
		}
		if cap(got) != 8 {
			t.Fatalf("cap(Get(5)) = %d, want 8", cap(got))
		}
		if unsafe.Pointer(&got[0]) != unsafe.Pointer(&backing[0]) {
			t.Fatal("Get(5) did not reuse pooled backing array")
		}
		got[0] = 11
		if backing[0] != 11 {
			t.Fatal("Get(5) returned slice not backed by pooled memory")
		}
	})
}

func TestPoolPut(t *testing.T) {
	t.Run("ignores zero capacity slice", func(t *testing.T) {
		var p Pool

		p.Put(nil)

		for idx := range p.pools {
			if ptr := p.pools[idx].Get(); ptr != nil {
				t.Fatalf("pool[%d] unexpectedly stored %#v", idx, ptr)
			}
		}
	})

	t.Run("stores power-of-two capacity in matching bucket", func(t *testing.T) {
		var p Pool
		buf := make([]byte, 2, 4)

		p.Put(buf)

		ptr, _ := p.pools[2].Get().(unsafe.Pointer)
		if ptr == nil {
			t.Fatal("expected pointer in bucket 2")
		}
		if ptr != unsafe.Pointer(&buf[:1][0]) {
			t.Fatal("stored pointer does not match original backing array")
		}
	})

	t.Run("stores non power-of-two capacity in previous bucket", func(t *testing.T) {
		var p Pool
		buf := make([]byte, 2, 5)

		p.Put(buf)

		ptr, _ := p.pools[2].Get().(unsafe.Pointer)
		if ptr == nil {
			t.Fatal("expected pointer in bucket 2")
		}
		if ptr != unsafe.Pointer(&buf[:1][0]) {
			t.Fatal("stored pointer does not match original backing array")
		}
	})
}

func TestBuiltinPoolWrappers(t *testing.T) {
	builtinPool = Pool{}

	t.Run("Get uses builtin pool", func(t *testing.T) {
		builtinPool = Pool{}

		backing := make([]byte, 4)
		builtinPool.pools[2].Put(unsafe.Pointer(&backing[0]))

		got := Get(3)
		if len(got) != 3 {
			t.Fatalf("len(Get(3)) = %d, want 3", len(got))
		}
		if cap(got) != 4 {
			t.Fatalf("cap(Get(3)) = %d, want 4", cap(got))
		}
		if unsafe.Pointer(&got[0]) != unsafe.Pointer(&backing[0]) {
			t.Fatal("Get did not use builtin pool backing array")
		}
	})

	t.Run("Put uses builtin pool", func(t *testing.T) {
		builtinPool = Pool{}

		buf := make([]byte, 1, 2)

		Put(buf)

		ptr, _ := builtinPool.pools[1].Get().(unsafe.Pointer)
		if ptr == nil {
			t.Fatal("expected pointer in builtin bucket 1")
		}
		if ptr != unsafe.Pointer(&buf[:1][0]) {
			t.Fatal("Put did not store original backing array in builtin pool")
		}
	})
}
