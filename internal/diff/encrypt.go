package diff

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
)

// EncryptOptions controls which keys are encrypted in the output.
type EncryptOptions struct {
	// Patterns is a list of key substrings that trigger encryption.
	Patterns []string
	// Passphrase is used to derive the AES-256 key.
	Passphrase string
}

// DefaultEncryptOptions returns sensible defaults.
func DefaultEncryptOptions() EncryptOptions {
	return EncryptOptions{
		Patterns:   []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PRIVATE"},
		Passphrase: "",
	}
}

// deriveKey produces a 32-byte AES key from a passphrase via SHA-256.
func deriveKey(passphrase string) []byte {
	h := sha256.Sum256([]byte(passphrase))
	return h[:]
}

// encryptValue encrypts plaintext using AES-256-GCM and returns a base64 string.
func encryptValue(plaintext, passphrase string) (string, error) {
	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptValue decrypts a base64-encoded AES-256-GCM ciphertext.
func DecryptValue(encoded, passphrase string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}
	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt: %w", err)
	}
	return string(plain), nil
}

// shouldEncrypt returns true if the key matches any pattern.
func shouldEncrypt(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

// EncryptEnv encrypts matching values in an env map and returns a new map.
// Encrypted values are prefixed with "enc:" to signal they are ciphertext.
func EncryptEnv(env map[string]string, opts EncryptOptions) (map[string]string, error) {
	if opts.Passphrase == "" {
		return nil, errors.New("passphrase must not be empty")
	}
	out := make(map[string]string, len(env))
	for k, v := range env {
		if shouldEncrypt(k, opts.Patterns) {
			enc, err := encryptValue(v, opts.Passphrase)
			if err != nil {
				return nil, fmt.Errorf("encrypting %s: %w", k, err)
			}
			out[k] = "enc:" + enc
		} else {
			out[k] = v
		}
	}
	return out, nil
}

// WriteEncryptReport writes a human-readable summary of encrypted keys to w.
func WriteEncryptReport(w io.Writer, original, encrypted map[string]string) {
	var changed []string
	for k, v := range encrypted {
		if strings.HasPrefix(v, "enc:") && original[k] != v {
			changed = append(changed, k)
		}
	}
	sort.Strings(changed)
	if len(changed) == 0 {
		fmt.Fprintln(w, "No keys encrypted.")
		return
	}
	fmt.Fprintf(w, "Encrypted %d key(s):\n", len(changed))
	for _, k := range changed {
		fmt.Fprintf(w, "  %s\n", k)
	}
}
