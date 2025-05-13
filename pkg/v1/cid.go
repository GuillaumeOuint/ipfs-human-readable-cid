package v1

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/GuillaumeOuint/ipfs-human-readable-cid/internal/wordlist"
	"github.com/ipfs/go-cid"
	mhash "github.com/multiformats/go-multihash"
)

const wordlistBase = 2048

// Assumed CID components for shortening
const (
	assumedCIDVersion       = 1
	assumedMulticodecType   = cid.DagProtobuf // 0x70
	assumedHashAlgorithm    = mhash.SHA2_256  // 0x12
	assumedHashDigestLength = 32              // Bytes
)

// CIDv1ToHumanReadable converts a CIDv1 to a human-readable string
func CIDv1ToHumanReadable(cidStr string) (string, error) {
	c, err := cid.Decode(cidStr)
	if err != nil {
		return "", fmt.Errorf("error decoding CID: %w", err)
	}

	// Get the multihash bytes
	mhBytes := c.Hash()

	// Decode the multihash to get its components
	decodedMH, err := mhash.Decode(mhBytes)
	if err != nil {
		return "", fmt.Errorf("error decoding multihash: %w", err)
	}

	// Convert the digest to a big.Int
	digestNum := new(big.Int).SetBytes(decodedMH.Digest)

	// Convert to base-2048
	base := big.NewInt(wordlistBase)
	var words []string
	remainder := new(big.Int)
	tempNum := new(big.Int).Set(digestNum) // Copy to avoid modifying digestNum

	for tempNum.Sign() > 0 {
		tempNum.DivMod(tempNum, base, remainder)
		idx := remainder.Int64()
		//fmt.Printf("Remainder %d -> word '%s'\n", idx, wordlist[idx])
		words = append(words, wordlist.Wordlist[idx])
	}

	// Reverse the words
	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i]
	}

	result := strings.Join(words, "-")
	return result, nil
}

// HumanReadableToCIDv1 converts a human-readable string
// back to a specific CIDv1 (DagProtobuf, SHA2-256, 32-byte digest).
func HumanReadableToCIDv1(humanReadable string) (string, error) {

	wordStrings := strings.Split(humanReadable, "-")

	// Build the wordMap for lookup
	wordMap := make(map[string]int64, wordlistBase)
	for i, w := range wordlist.Wordlist {
		if w != "" { // Skip empty entries
			wordMap[w] = int64(i)
		}
	}

	// Special case for wordlist[0]
	if len(wordStrings) == 1 && wordStrings[0] == wordlist.Wordlist[0] {
		digestBytes := make([]byte, 32) // 32 zeros

		reconstructedMH, err := mhash.Sum(digestBytes, mhash.SHA2_256, 32)
		if err != nil {
			return "", fmt.Errorf("error creating zero multihash: %w", err)
		}

		reconstructedCid := cid.NewCidV1(cid.DagProtobuf, reconstructedMH)
		result := reconstructedCid.String()
		return result, nil
	}

	// Rebuild the number from words
	num := big.NewInt(0)
	base := big.NewInt(wordlistBase)

	for _, wordStr := range wordStrings {
		index, found := wordMap[wordStr]
		if !found {
			return "", fmt.Errorf("word '%s' not found in wordlist", wordStr)
		}
		//fmt.Printf("Word '%s' -> index %d\n", wordStr, index)

		num.Mul(num, base)
		num.Add(num, big.NewInt(index))
	}

	// Convert to bytes (fixed length)
	digestBytes := make([]byte, 32) // 32 bytes = 256 bits
	num.FillBytes(digestBytes)

	// Create the multihash by directly encoding our reconstructed digest
	reconstructedMH, err := mhash.Encode(digestBytes, mhash.SHA2_256)
	if err != nil {
		return "", fmt.Errorf("error encoding multihash: %w", err)
	}

	// Create CID
	reconstructedCid := cid.NewCidV1(cid.DagProtobuf, reconstructedMH)

	result := reconstructedCid.String()

	return result, nil
}
