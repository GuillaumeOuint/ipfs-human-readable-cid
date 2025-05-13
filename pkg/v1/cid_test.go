package v1

import (
	"testing"

	"github.com/GuillaumeOuint/ipfs-human-readable-cid/internal/wordlist"
)

func TestCID(t *testing.T) {
	// ... (main function remains the same, but ensure wordlist is populated)
	// Example: (You MUST replace the wordlist above with your actual 2048 words)
	// For testing, ensure a few words are present in the placeholder:
	if wordlist.Wordlist[0] == "" || wordlist.Wordlist[1] == "" || wordlist.Wordlist[wordlistBase-1] == "" {
		t.Fatal("Wordlist is not sufficiently populated for basic testing. Please fill it.")
	}

	cidStr := "bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi"
	t.Log("Original CID:", cidStr)

	humanReadable, err := CIDv1ToHumanReadable(cidStr)
	if err != nil {
		t.Fatalf("Error converting CID to human-readable: %v", err)
	}
	t.Log("Human-readable CID:", humanReadable)

	reconstructedCidStr, err := HumanReadableToCIDv1(humanReadable)
	if err != nil {
		t.Fatalf("Error converting human-readable to CID: %v", err)
	}
	t.Log("Reconstructed CID:", reconstructedCidStr)

	if cidStr != reconstructedCidStr {
		t.Fatalf("Mismatch! Original: %s, Reconstructed: %s", cidStr, reconstructedCidStr)
	} else {
		t.Log("Success! CIDs match.")
	}
}
