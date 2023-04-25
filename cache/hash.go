package cache

const (
	offset32 = 2166136261
	offset64 = 14695981039346656037
	prime32  = 16777619
	prime64  = 1099511628211
)

func hash32[T string | []byte](str T) uint32 {
	var hash uint32 = offset32
	for i, n := 0, len(str); i < n; i++ {
		c := str[i]
		hash *= prime32
		hash ^= uint32(c)
	}
	return hash
}

func hash64[T string | []byte](str T) uint64 {
	var hash uint64 = offset64
	for i, n := 0, len(str); i < n; i++ {
		c := str[i]
		hash *= prime64
		hash ^= uint64(c)
	}
	return hash
}
