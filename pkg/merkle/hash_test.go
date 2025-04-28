package merkle

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestHash(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "empty input",
			input:    []byte{},
			expected: "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
		},
		{
			name:     "simple string",
			input:    []byte("hello world"),
			expected: "47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad",
		},
		{
			name:     "binary data",
			input:    []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			expected: "7d87c5ea75f7378bb701e404c50639161af3eff66293e9f375b5f17eb50476f4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Hash(tt.input)
			resultHex := hex.EncodeToString(result)
			if resultHex != tt.expected {
				t.Errorf("Hash(%v) = %s, want %s", tt.input, resultHex, tt.expected)
			}
		})
	}

	// Test consistency.
	t.Run("consistency", func(t *testing.T) {
		input := []byte("test consistency")
		result1 := Hash(input)
		result2 := Hash(input)
		if !bytes.Equal(result1, result2) {
			t.Errorf("Hash is not consistent for the same input")
		}
	})

	// Test different inputs produce different outputs.
	t.Run("different inputs", func(t *testing.T) {
		input1 := []byte("input1")
		input2 := []byte("input2")
		result1 := Hash(input1)
		result2 := Hash(input2)
		if bytes.Equal(result1, result2) {
			t.Errorf("Different inputs produced the same hash")
		}
	})
}

func TestHashNode(t *testing.T) {
	tests := []struct {
		name     string
		left     []byte
		right    []byte
		expected string
	}{
		{
			name:     "empty left and right",
			left:     []byte{},
			right:    []byte{},
			expected: "8c83ac90634bd25e7214fbf0e3d45b5e8633bb010cb64411c122b53131c7c431",
		},
		{
			name:     "non-empty left, empty right",
			left:     []byte{0x01},
			right:    []byte{},
			expected: "83f78f037fd069ea3b5c402c21b4d1e122d206ba5c847ba67ca62ebe75bbf51f",
		},
		{
			name:     "empty left, non-empty right",
			left:     []byte{},
			right:    []byte{0x01},
			expected: "83f78f037fd069ea3b5c402c21b4d1e122d206ba5c847ba67ca62ebe75bbf51f",
		},
		{
			name:     "non-empty left and right",
			left:     []byte{0x01, 0x02},
			right:    []byte{0x03, 0x04},
			expected: "3db27c66a80724227ca3d27d202f29487959ec0912b0db2b9dbbf01953a7ab49",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HashNode(tt.left, tt.right)
			resultHex := hex.EncodeToString(result)
			if resultHex != tt.expected {
				t.Errorf(
					"HashNode(%v, %v) = %s, want %s",
					tt.left,
					tt.right,
					resultHex,
					tt.expected,
				)
			}
		})
	}

	// Test that node prefix is used
	t.Run("node prefix is used", func(t *testing.T) {
		data := []byte{0x01, 0x02}
		nodeHash := HashNode(data, []byte{})
		directHash := Hash(append([]byte(NODE_PREFIX), append(data, []byte{}...)...))
		if !bytes.Equal(nodeHash, directHash) {
			t.Errorf("Node hash doesn't match manual construction with prefix")
		}
	})

	// Test that HashNode is different from HashLeaf with same data
	t.Run("node hash differs from leaf hash", func(t *testing.T) {
		data1 := []byte{0x01}
		data2 := []byte{0x02}
		nodeHash := HashNode(data1, data2)
		combinedData := append(data1, data2...)
		leafHash := HashLeaf(combinedData)
		if bytes.Equal(nodeHash, leafHash) {
			t.Errorf("Node hash equals leaf hash for same data")
		}
	})
}

func TestHashLeaf(t *testing.T) {
	tests := []struct {
		name     string
		leaf     []byte
		expected string
	}{
		{
			name:     "empty leaf",
			leaf:     []byte{},
			expected: "3ef0000fc8752f5372eb9bcff2d75ad56ac4dc0824bb0dffcf7e454001558bf7",
		},
		{
			name:     "simple leaf",
			leaf:     []byte("test leaf"),
			expected: "a072d672610b40ba2ad9429f65421dddb525911c302c5933737ce7e62cc1da26",
		},
		{
			name:     "binary data",
			leaf:     []byte{0x01, 0x02, 0x03},
			expected: "7c7b7f152b1492aadc0df2d4ab2bf9cc171d5359082fa840b18ef9304a3886e2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HashLeaf(tt.leaf)
			resultHex := hex.EncodeToString(result)
			if resultHex != tt.expected {
				t.Errorf("HashLeaf(%v) = %s, want %s", tt.leaf, resultHex, tt.expected)
			}
		})
	}

	// Test that leaf prefix is used.
	t.Run("HashLeaf uses leaf prefix", func(t *testing.T) {
		data := []byte{0x01, 0x02}
		leafHash := HashLeaf(data)
		directHash := Hash(append([]byte(LEAF_PREFIX), data...))
		if !bytes.Equal(leafHash, directHash) {
			t.Errorf("Leaf hash doesn't match manual construction with prefix")
		}
	})
}

func TestHashEmptyLeaf(t *testing.T) {
	// Test that HashEmptyLeaf returns the same result as HashLeaf with empty data.
	t.Run("equals HashLeaf with empty data", func(t *testing.T) {
		emptyLeafHash := HashEmptyLeaf()
		manualEmptyLeafHash := HashLeaf([]byte{})
		if !bytes.Equal(emptyLeafHash, manualEmptyLeafHash) {
			t.Errorf("HashEmptyLeaf() doesn't match HashLeaf([]byte{})")
		}
	})

	// HashEmpty leaf known value test.
	t.Run("known value", func(t *testing.T) {
		expected := "3ef0000fc8752f5372eb9bcff2d75ad56ac4dc0824bb0dffcf7e454001558bf7"
		emptyLeafHash := HashEmptyLeaf()
		resultHex := hex.EncodeToString(emptyLeafHash)
		if resultHex != expected {
			t.Errorf("HashEmptyLeaf() = %s, want %s", resultHex, expected)
		}
	})
}

func TestHashDomainSeparation(t *testing.T) {
	// Test that leaf and node hashes are different even with the same data.
	t.Run("leaf and node hash separation", func(t *testing.T) {
		data := []byte{0x01, 0x02}
		leafHash := HashLeaf(data)
		nodeHash := HashNode(data, []byte{})
		if bytes.Equal(leafHash, nodeHash) {
			t.Errorf("Leaf and node hashes are not properly domain-separated")
		}
	})

	// Test that different prefixes result in different hashes.
	t.Run("prefix separation", func(t *testing.T) {
		data := []byte{0x01, 0x02}
		hash1 := Hash(append([]byte("prefix1|"), data...))
		hash2 := Hash(append([]byte("prefix2|"), data...))
		if bytes.Equal(hash1, hash2) {
			t.Errorf("Different prefixes should result in different hashes")
		}
	})
}
