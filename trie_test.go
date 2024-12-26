package search_trie

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestTrie_PutAndTraverse(t *testing.T) {
	tests := []struct {
		name     string
		testData map[string]uint
	}{
		{
			name: "English words",
			testData: map[string]uint{
				"iphone":                5,
				"iphone 16":             3,
				"iphone 16 pro":         7,
				"iphone 16 pro max":     2,
				"iphone 16 pro max 256": 1,
				"macbook":               4,
				"macbook air":           6,
				"macbook pro":           8,
				"ipad":                  10,
			},
		},
		{
			name: "Russian words",
			testData: map[string]uint{
				"телефон":         5,
				"телефон 16":      3,
				"телефон 16 про":  7,
				"телефон 16 макс": 2,
				"телефон 256":     1,
				"ноутбук":         4,
				"планшет":         10,
				"ноутбук про":     8,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trie := NewTrie(5)

			// Add keys and their frequencies to the Trie
			for key, freq := range tt.testData {
				trie.Put(key, 0) // Add the key with an initial value
				for i := 0; i < int(freq); i++ {
					trie.Inc(key) // Increment the frequency of the key
				}
			}

			// Run Traverse and collect the output
			output := make(map[string]uint)
			for item := range trie.Traverse() {
				output[item.Key] = item.Frequency
			}

			// Compare the output with the expected data
			for key, expectedFreq := range tt.testData {
				if freq, exists := output[key]; !exists {
					t.Errorf("Key %q not found in output", key)
				} else if freq != expectedFreq {
					t.Errorf("Key %q has incorrect frequency: got %d, want %d", key, freq, expectedFreq)
				}
			}

			// Ensure the output does not contain unexpected keys
			if len(output) != len(tt.testData) {
				t.Errorf("Output contains unexpected keys: got %d keys, want %d keys", len(output), len(tt.testData))
			}
		})
	}
}

func TestTrie_TopK(t *testing.T) {
	tests := []struct {
		name        string
		testData    map[string]uint
		prefix      string
		expectedRes []nodeInfo
	}{
		{
			name: "English words",
			testData: map[string]uint{
				"ipad":                  35,
				"iphone 16 pro":         28,
				"iphone":                30,
				"iphone 16":             45,
				"iphone 16 pro max":     14,
				"iphone 16 pro max 256": 1,
				"macbook":               4,
				"macbook air":           6,
				"macbook pro":           8,
			},
			prefix: "ip",
			expectedRes: []nodeInfo{
				{Key: "iphone 16", Frequency: 45},
				{Key: "ipad", Frequency: 35},
				{Key: "iphone", Frequency: 30},
				{Key: "iphone 16 pro", Frequency: 28},
				{Key: "iphone 16 pro max", Frequency: 14},
			},
		},
		{
			name: "Russian words",
			testData: map[string]uint{
				"айпад":        35,
				"айфон 16 про": 28,
				"айфон":        30,
				"айфон 16":     45,
				"айфон макс":   14,
				"айфон 256":    1,
				"макбук":       4,
				"макбук air":   6,
				"макбук про":   8,
			},
			prefix: "ай",
			expectedRes: []nodeInfo{
				{Key: "айфон 16", Frequency: 45},
				{Key: "айпад", Frequency: 35},
				{Key: "айфон", Frequency: 30},
				{Key: "айфон 16 про", Frequency: 28},
				{Key: "айфон макс", Frequency: 14},
			},
		},
		{
			name: "Mixed words",
			testData: map[string]uint{
				"ipad":          35,
				"iphone 16 pro": 28,
				"iphone":        30,
				"iphone 16":     45,
				"айфон макс":    14,
				"айфон 256":     1,
				"макбук":        4,
				"макбук air":    6,
				"макбук про":    8,
			},
			prefix: "iph",
			expectedRes: []nodeInfo{
				{Key: "iphone 16", Frequency: 45},
				{Key: "iphone", Frequency: 30},
				{Key: "iphone 16 pro", Frequency: 28},
			},
		},
		{
			name: "Mixed words no matches",
			testData: map[string]uint{
				"ipad":          35,
				"iphone 16 pro": 28,
				"iphone":        30,
				"iphone 16":     45,
				"айфон макс":    14,
				"айфон 256":     1,
				"макбук":        4,
				"макбук air":    6,
				"макбук про":    8,
			},
			prefix:      "sams",
			expectedRes: []nodeInfo{},
		},
		{
			name: "Empty Query",
			testData: map[string]uint{
				"ipad":          35,
				"iphone 16 pro": 28,
				"iphone":        30,
				"iphone 16":     45,
				"айфон макс":    14,
				"айфон 256":     1,
				"макбук":        4,
				"макбук air":    6,
				"макбук про":    8,
			},
			prefix:      "",
			expectedRes: []nodeInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trie := NewTrie(5)

			// Add keys and their frequencies to the Trie
			for key, freq := range tt.testData {
				trie.Put(key, 0) // Add the key with an initial value
				for i := 0; i < int(freq); i++ {
					trie.Inc(key) // Increment the frequency of the key
				}
			}

			// Test the TopK method
			res := trie.TopK(tt.prefix)
			if len(res) != len(tt.expectedRes) {
				t.Errorf("TopK() = %v, want %v", res, tt.expectedRes)
			}

			for i, item := range tt.expectedRes {
				if i > len(res)-1 {
					t.Errorf("TopK() = %v, want %v", res, tt.expectedRes)
					break
				}
				if res[i].Key != item.Key || res[i].Frequency != item.Frequency {
					t.Errorf("TopK() = %v, want %v", res[i], tt.expectedRes[i])
				}
			}
		})
	}
}

func TestTrie_Has(t *testing.T) {
	tests := []struct {
		name        string
		testData    map[string]uint
		key         string
		expectedRes bool
	}{
		{
			name: "Prefix exists",
			testData: map[string]uint{
				"ipad":                  35,
				"iphone 16 pro":         28,
				"iphone":                30,
				"iphone 16":             45,
				"iphone 16 pro max":     14,
				"iphone 16 pro max 256": 1,
				"macbook":               4,
				"macbook air":           6,
				"macbook pro":           8,
			},
			key:         "iphone",
			expectedRes: true,
		},
		{
			name: "Prefix does not exist",
			testData: map[string]uint{
				"ipad":                  35,
				"iphone 16 pro":         28,
				"iphone":                30,
				"iphone 16":             45,
				"iphone 16 pro max":     14,
				"iphone 16 pro max 256": 1,
				"macbook":               4,
				"macbook air":           6,
				"macbook pro":           8,
			},
			key:         "samsung",
			expectedRes: false,
		},
		{
			name: "Empty prefix",
			testData: map[string]uint{
				"ipad": 35,
			},
			key:         "",
			expectedRes: false, // Assuming empty prefix always matches
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trie := NewTrie(5)

			// Add keys and their frequencies to the Trie
			for key, freq := range tt.testData {
				trie.Put(key, 0) // Add the key with an initial value
				for i := 0; i < int(freq); i++ {
					trie.Inc(key) // Increment the frequency of the key
				}
			}

			// Test the Has method
			if has := trie.Has(tt.key); has != tt.expectedRes {
				t.Errorf("Has(%q) = %v, expected %v", tt.key, has, tt.expectedRes)
			}
		})
	}
}

//BenchmarkTrie_GetTopK/English_words-8             301017              4525 ns/op             272 B/op          7 allocs/op
//BenchmarkTrie_GetTopK/Russian_words-8             277951              4034 ns/op             304 B/op          9 allocs/op
//BenchmarkTrie_Put/Small_dataset-8                 133527              9283 ns/op            1201 B/op         36 allocs/op
//BenchmarkTrie_Put/Medium_dataset-8                 89562             12807 ns/op            1291 B/op         37 allocs/op
//BenchmarkTrie_Put/Large_dataset-8                  74458             17335 ns/op            1493 B/op         41 allocs/op
//BenchmarkTrie_Put/Small_topk-8                    109454             11641 ns/op            1470 B/op         42 allocs/op
//BenchmarkTrie_Put/Medium_topk-8                    92978             13822 ns/op            1485 B/op         41 allocs/op
//BenchmarkTrie_Put/Large_topk-8                     73536             17225 ns/op            1494 B/op         41 allocs/op
//BenchmarkTrie_Inc/Small_dataset-8                 140379              9047 ns/op            1225 B/op         36 allocs/op
//BenchmarkTrie_Inc/Medium_dataset-8                101518             12461 ns/op            1260 B/op         36 allocs/op
//BenchmarkTrie_Inc/Large_dataset-8                  81650             16850 ns/op            1276 B/op         36 allocs/op
//BenchmarkTrie_Inc/Small_topk-8                    114806             12618 ns/op            1336 B/op         38 allocs/op
//BenchmarkTrie_Inc/Medium_topk-8                    97866             13421 ns/op            1295 B/op         37 allocs/op
//BenchmarkTrie_Inc/Large_topk-8                     77398             16541 ns/op            1276 B/op         36 allocs/op
//BenchmarkTrie_GetTopKParallel-8                   561360              2122 ns/op             146 B/op          5 allocs/op
//BenchmarkTrie_PutParallel/Small_dataset-8                  91459             12926 ns/op            1203 B/op         36 allocs/op
//BenchmarkTrie_PutParallel/Medium_dataset-8                 66752             17856 ns/op            1305 B/op         37 allocs/op
//BenchmarkTrie_PutParallel/Large_dataset-8                  53194             23275 ns/op            1487 B/op         41 allocs/op
//BenchmarkTrie_PutParallel/Small_topk-8                     78505             15391 ns/op            1463 B/op         42 allocs/op
//BenchmarkTrie_PutParallel/Medium_topk-8                    65421             18760 ns/op            1477 B/op         41 allocs/op
//BenchmarkTrie_PutParallel/Large_topk-8                     54554             23091 ns/op            1484 B/op         41 allocs/op
//BenchmarkTrie_IncParallel/Small_dataset-8                  95036             12247 ns/op            1190 B/op         36 allocs/op
//BenchmarkTrie_IncParallel/Medium_dataset-8                 73290             16931 ns/op            1230 B/op         36 allocs/op
//BenchmarkTrie_IncParallel/Large_dataset-8                  58549             21922 ns/op            1263 B/op         36 allocs/op
//BenchmarkTrie_IncParallel/Small_topk-8                     87006             14669 ns/op            1312 B/op         38 allocs/op
//BenchmarkTrie_IncParallel/Medium_topk-8                    69993             17650 ns/op            1287 B/op         37 allocs/op
//BenchmarkTrie_IncParallel/Large_topk-8                     57313             21858 ns/op            1263 B/op         36 allocs/op
//BenchmarkTrie_GetTopKAndInc/Small_dataset-8               146946              8420 ns/op             919 B/op         27 allocs/op
//BenchmarkTrie_GetTopKAndInc/Medium_dataset-8              111247             12673 ns/op             982 B/op         27 allocs/op
//BenchmarkTrie_GetTopKAndInc/Large_dataset-8                92185             14468 ns/op            1082 B/op         27 allocs/op
//BenchmarkTrie_GetTopKAndInc/Small_topk-8                  143139              9362 ns/op            1015 B/op         29 allocs/op
//BenchmarkTrie_GetTopKAndInc/Medium_topk-8                 110716             11883 ns/op            1022 B/op         28 allocs/op
//BenchmarkTrie_GetTopKAndInc/Large_topk-8                   95582             14251 ns/op            1082 B/op         27 allocs/op
//BenchmarkTrie_GetTopKAndIncParallel/Small_dataset-8               109334             11160 ns/op             915 B/op         27 allocs/op
//BenchmarkTrie_GetTopKAndIncParallel/Medium_dataset-8               76646             16896 ns/op             977 B/op         27 allocs/op
//BenchmarkTrie_GetTopKAndIncParallel/Large_dataset-8                65320             19459 ns/op            1070 B/op         27 allocs/op
//BenchmarkTrie_GetTopKAndIncParallel/Small_topk-8                  101413             12367 ns/op             999 B/op         28 allocs/op
//BenchmarkTrie_GetTopKAndIncParallel/Medium_topk-8                  80700             16219 ns/op            1015 B/op         28 allocs/op
//BenchmarkTrie_GetTopKAndIncParallel/Large_topk-8                   64881             19452 ns/op            1072 B/op         27 allocs/op

// Benchmark for getTopK with different datasets
func BenchmarkTrie_GetTopK(b *testing.B) {
	tests := []struct {
		name     string
		testData map[string]uint
		prefix   string
	}{
		{
			name: "English words",
			testData: map[string]uint{
				"iphone":                5,
				"iphone 16":             3,
				"iphone 16 pro":         7,
				"iphone 16 pro max":     2,
				"iphone 16 pro max 256": 1,
				"macbook":               4,
				"macbook air":           6,
				"macbook pro":           8,
				"ipad":                  10,
			},
			prefix: "ip",
		},
		{
			name: "Russian words",
			testData: map[string]uint{
				"телефон":         5,
				"телефон 16":      3,
				"телефон 16 про":  7,
				"телефон 16 макс": 2,
				"телефон 256":     1,
				"ноутбук":         4,
				"ноутбук air":     6,
				"ноутбук про":     8,
				"планшет":         10,
			},
			prefix: "тел",
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			trie := NewTrie(5)

			// Add keys and their frequencies to the Trie
			for key, freq := range tt.testData {
				trie.Put(key, 0) // Add the key with an initial value
				for i := 0; i < int(freq); i++ {
					trie.Inc(key) // Increment the frequency of the key
				}
			}

			// Reset timer to measure only getTopK
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = trie.TopK(tt.prefix) // Call getTopK for the given prefix
			}
		})
	}
}

func BenchmarkTrie_Put(b *testing.B) {
	tests := []struct {
		name    string
		topK    int
		numKeys int
	}{
		{name: "Small dataset", topK: 5, numKeys: 1000},
		{name: "Medium dataset", topK: 10, numKeys: 10000},
		{name: "Large dataset", topK: 20, numKeys: 100000},

		{name: "Small topk", topK: 5, numKeys: 100000},
		{name: "Medium topk", topK: 10, numKeys: 100000},
		{name: "Large topk", topK: 20, numKeys: 100000},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			trie := NewTrie(tt.topK)

			// Generate random keys
			keys := generateRandomKeys(tt.numKeys)

			b.ResetTimer() // Reset timer to exclude setup time
			for i := 0; i < b.N; i++ {
				key := keys[i%len(keys)]
				trie.Put(key, uint(i%100))
			}
		})
	}
}

func BenchmarkTrie_Inc(b *testing.B) {
	tests := []struct {
		name    string
		topK    int
		numKeys int
	}{
		{name: "Small dataset", topK: 5, numKeys: 1000},
		{name: "Medium dataset", topK: 10, numKeys: 10000},
		{name: "Large dataset", topK: 20, numKeys: 100000},

		{name: "Small topk", topK: 5, numKeys: 100000},
		{name: "Medium topk", topK: 10, numKeys: 100000},
		{name: "Large topk", topK: 20, numKeys: 100000},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			trie := NewTrie(tt.topK)

			// Generate random keys and prepopulate the Trie
			keys := generateRandomKeys(tt.numKeys)
			for _, key := range keys {
				trie.Put(key, 1) // Prepopulate with frequency 1
			}

			b.ResetTimer() // Reset timer to exclude setup time
			for i := 0; i < b.N; i++ {
				key := keys[i%len(keys)]
				trie.Inc(key)
			}
		})
	}
}

// Benchmark for getTopK with parallel execution
func BenchmarkTrie_GetTopKParallel(b *testing.B) {
	trie := NewTrie(5)

	// Add data to the Trie (mix of English and Russian words)
	testData := map[string]uint{
		"iphone":      5,
		"iphone 16":   3,
		"телефон":     5,
		"телефон 16":  3,
		"ipad":        10,
		"планшет":     10,
		"macbook pro": 8,
		"ноутбук про": 8,
	}
	for key, freq := range testData {
		trie.Put(key, 0) // Add the key with an initial value
		for i := 0; i < int(freq); i++ {
			trie.Inc(key) // Increment the frequency of the key
		}
	}

	prefixes := []string{"ip", "тел", "ма", "но"}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			prefix := prefixes[rand.Intn(len(prefixes))]
			_ = trie.TopK(prefix) // Call getTopK for random prefixes
		}
	})
}

func BenchmarkTrie_PutParallel(b *testing.B) {
	tests := []struct {
		name    string
		topK    int
		numKeys int
	}{
		{name: "Small dataset", topK: 5, numKeys: 1000},
		{name: "Medium dataset", topK: 10, numKeys: 10000},
		{name: "Large dataset", topK: 20, numKeys: 100000},

		{name: "Small topk", topK: 5, numKeys: 100000},
		{name: "Medium topk", topK: 10, numKeys: 100000},
		{name: "Large topk", topK: 20, numKeys: 100000},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			trie := NewTrie(tt.topK)

			// Generate random keys
			keys := generateRandomKeys(tt.numKeys)

			b.ResetTimer() // Reset timer to exclude setup time
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					key := keys[rand.Intn(len(keys))]
					trie.Put(key, uint(rand.Intn(100)))
				}
			})
		})
	}
}

func BenchmarkTrie_IncParallel(b *testing.B) {
	tests := []struct {
		name    string
		topK    int
		numKeys int
	}{
		{name: "Small dataset", topK: 5, numKeys: 1000},
		{name: "Medium dataset", topK: 10, numKeys: 10000},
		{name: "Large dataset", topK: 20, numKeys: 100000},

		{name: "Small topk", topK: 5, numKeys: 100000},
		{name: "Medium topk", topK: 10, numKeys: 100000},
		{name: "Large topk", topK: 20, numKeys: 100000},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			trie := NewTrie(tt.topK)

			// Generate random keys and prepopulate the Trie
			keys := generateRandomKeys(tt.numKeys)
			for _, key := range keys {
				trie.Put(key, 1) // Prepopulate with frequency 1
			}

			b.ResetTimer() // Reset timer to exclude setup time
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					key := keys[rand.Intn(len(keys))]
					trie.Inc(key)
				}
			})
		})
	}
}

func BenchmarkTrie_GetTopKAndInc(b *testing.B) {
	// Тестовые данные
	tests := []struct {
		name    string
		topK    int
		numKeys int
	}{
		{name: "Small dataset", topK: 5, numKeys: 1000},
		{name: "Medium dataset", topK: 10, numKeys: 10000},
		{name: "Large dataset", topK: 20, numKeys: 100000},

		{name: "Small topk", topK: 5, numKeys: 100000},
		{name: "Medium topk", topK: 10, numKeys: 100000},
		{name: "Large topk", topK: 20, numKeys: 100000},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			trie := NewTrie(tt.topK)

			// Генерация данных
			keys := generateRandomKeys(tt.numKeys)
			for _, key := range keys {
				trie.Put(key, 1) // Предзаполнение Trie
			}

			b.ResetTimer() // Сброс таймера

			for i := 0; i < b.N; i++ {
				// Чередуем вызовы getTopK и Inc
				if rand.Intn(10) <= 2 {
					_ = trie.TopK(keys[i%len(keys)][:2]) // Используем первые 2 символа как префикс
				} else {
					trie.Inc(keys[i%len(keys)])
				}
			}
		})
	}
}

func BenchmarkTrie_GetTopKAndIncParallel(b *testing.B) {
	// Тестовые данные
	tests := []struct {
		name    string
		topK    int
		numKeys int
	}{
		{name: "Small dataset", topK: 5, numKeys: 1000},
		{name: "Medium dataset", topK: 10, numKeys: 10000},
		{name: "Large dataset", topK: 20, numKeys: 100000},

		{name: "Small topk", topK: 5, numKeys: 100000},
		{name: "Medium topk", topK: 10, numKeys: 100000},
		{name: "Large topk", topK: 20, numKeys: 100000},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			trie := NewTrie(tt.topK)

			// Генерация данных
			keys := generateRandomKeys(tt.numKeys)
			for _, key := range keys {
				trie.Put(key, 1) // Предзаполнение Trie
			}

			b.ResetTimer() // Сброс таймера

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					if rand.Intn(10) <= 2 {
						_ = trie.TopK(keys[rand.Intn(len(keys))][:2]) // Используем первые 2 символа как префикс
					} else {
						trie.Inc(keys[rand.Intn(len(keys))])
					}
				}
			})
		})
	}
}

// Utility function to generate random keys
func generateRandomKeys(numKeys int) []string {
	rand.Seed(time.Now().UnixNano())
	keys := make([]string, numKeys)
	for i := 0; i < numKeys; i++ {
		keys[i] = fmt.Sprintf("key-%d", rand.Intn(100000)) // Generate random keys
	}
	return keys
}
