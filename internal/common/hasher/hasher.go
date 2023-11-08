package hasher

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

type ParamsHash struct {
	memory      uint32
	time        uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func GenerateSession(seed string) string {
	hash := sha256.New()
	hash.Write([]byte(seed))
	return hex.EncodeToString(hash.Sum(nil))
}

func GeneratePasswordHash(pwd string) (hashEncoded string, err error) {
	p := ParamsHash{
		memory:      64 * 1024,
		time:        1,
		parallelism: 4,
		saltLength:  16,
		keyLength:   32,
	}

	salt, err := genRandomBytes(p.saltLength)
	if err != nil {
		return "", err
	}

	hashedPassword := argon2.IDKey([]byte(pwd), salt, p.time, p.memory, p.parallelism, p.keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hashedPassword)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		p.memory,
		p.time,
		p.parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

func VerfiyPassword(pwd, encodedHash string) (bool, error) {
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(pwd), salt, p.time, p.memory, p.parallelism, p.keyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func genRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func decodeHash(encodedHash string) (p *ParamsHash, salt, hash []byte, err error) {
	startingPArams := strings.Split(encodedHash, "$")
	if len(startingPArams) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(startingPArams[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &ParamsHash{}
	_, err = fmt.Sscanf(startingPArams[3], "m=%d,t=%d,p=%d", &p.memory, &p.time, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(startingPArams[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(startingPArams[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}

// testPassword
