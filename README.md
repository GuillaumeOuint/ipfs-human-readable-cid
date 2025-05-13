# IPFS Human-Readable CID

Convert IPFS Content Identifiers (CIDs) to human-readable word sequences and back, aiming for perfect reversibility for common CIDv1 types.

## Overview

IPFS Content Identifiers (CIDs) are powerful but not inherently human-friendly. This project provides tools to convert CIDv1 strings into more memorable sequences of words and then convert them back to CIDs.

The primary focus is on CIDv1 identifiers, particularly those using the `dag-pb` codec and `sha2-256` hashing algorithm (common for IPFS files and objects).

## How It Works

**Conversion to Human-Readable (`CIDv1ToHumanReadable`):**
1. Decodes the input CIDv1 string.
2. Extracts the full digest from the CID's multihash component.
3. Treats this digest as a large integer.
4. Converts this integer to base-2048 (where 2048 is the size of the wordlist).
5. Maps each "digit" in this new base to a corresponding word from an internal 2048-word list.
6. Joins these words with hyphens to form the human-readable string.

**Conversion from Human-Readable (`HumanReadableToCIDv1`):**
1. Splits the hyphenated word string into individual words.
2. Maps each word back to its numerical index from the wordlist.
3. Reconstructs the large integer (representing the original digest) from these base-2048 indices.
4. Converts this integer back into a 32-byte digest.
5. Creates a new CIDv1 using this digest, explicitly assuming:
    - Content Type (Multicodec): `dag-pb` (0x70)
    - Hash Algorithm (Multihash): `sha2-256` (0x12)

## Features

- **Human-Readable**: Converts CIDv1 digests to memorable word sequences.
- **Reversible**: Aims for perfect reconstruction of CIDs, specifically for `dag-pb` / `sha2-256` types.
- **Word-Based**: Utilizes an internal 2048-word list (based on BIP39 English wordlist).

## Installation

```bash
go get github.com/GuillaumeOuint/ipfs-human-readable-cid
```

## Usage

### As a Library

To use the `v1` package:

```go
package main

import (
    "fmt"
    "log"

    humancid "github.com/GuillaumeOuint/ipfs-human-readable-cid/pkg/v1"
)

func main() {
    cidStr := "bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi" // Example CIDv1 dag-pb sha2-256

    // Convert CID to human-readable format
    humanReadable, err := humancid.CIDv1ToHumanReadable(cidStr)
    if err != nil {
        log.Fatalf("Error converting CID to human-readable: %v", err)
    }
    fmt.Printf("Original CID: %s\nHuman-readable: %s\n", cidStr, humanReadable)

    // Convert back to CID
    recoveredCID, err := humancid.HumanReadableToCIDv1(humanReadable)
    if err != nil {
        log.Fatalf("Error converting human-readable to CID: %v", err)
    }
    fmt.Printf("Recovered CID: %s\n", recoveredCID)

    if cidStr == recoveredCID {
        fmt.Println("Success! CIDs match.")
    } else {
        fmt.Println("Mismatch! CIDs do not match.")
    }
}
```

### Example Output

For a CID like `bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi`:

Human-readable equivalent (will vary based on the CID's digest, typically ~24-27 words for a 32-byte digest with a 2048-word list):
```
// Example: absorb-camera-button-mushroom-palace-avail-build-anymore-relax-rural-avid-video-twelve-jacket-ladder-crazy-able-hybrid-student-tomorrow-subway-although-genuine-chalk
```

## Wordlist

The library uses an internal 2048-word list sourced from `github.com/GuillaumeOuint/ipfs-human-readable-cid/internal/wordlist`. This list should be complete and contain unique words for the conversion to be reliable.

## Limitations

- **CIDv1 Focus**: The functions are named `CIDv1ToHumanReadable` and `HumanReadableToCIDv1`, indicating they are designed for CIDv1.
- **Reconstruction Assumptions**: `HumanReadableToCIDv1` specifically reconstructs the CID assuming it should be a `dag-pb` content type with a `sha2-256` multihash. If the original CID used a different content type or hash algorithm (but was still CIDv1), the human-readable string would be based on its digest, but the reconstructed CID would be forced into the `dag-pb`/`sha2-256` format.
- **Digest Length**: The reconstruction assumes a 32-byte digest, standard for SHA2-256.
- **Word Count**: For a 32-byte digest and a 2048-word list (11 bits per word), the human-readable string will be approximately `ceil(256 / 11) = 24` words long.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.