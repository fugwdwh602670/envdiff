package diff

import (
	"bytes"
	"strings"
	"testing"
)

func makeEncryptEnv() map[string]string {
	return map[string]string{
		"APP_NAME":      "myapp",
		"DB_PASSWORD":   "s3cr3t",
		"API_TOKEN":     "tok_abc123",
		"LOG_LEVEL":     "info",
		"PRIVATE_KEY":   "-----BEGIN RSA-----",
	}
}

func TestEncryptEnv_EncryptsMatchingKeys(t *testing.T) {
	env := makeEncryptEnv()
	opts := DefaultEncryptOptions()
	opts.Passphrase = "testpass"

	out, err := EncryptEnv(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, k := range []string{"DB_PASSWORD", "API_TOKEN", "PRIVATE_KEY"} {
		if !strings.HasPrefix(out[k], "enc:") {
			t.Errorf("expected %s to be encrypted, got %q", k, out[k])
		}
	}

	for _, k := range []string{"APP_NAME", "LOG_LEVEL"} {
		if strings.HasPrefix(out[k], "enc:") {
			t.Errorf("expected %s to be plain, got %q", k, out[k])
		}
		if out[k] != env[k] {
			t.Errorf("expected %s=%q, got %q", k, env[k], out[k])
		}
	}
}

func TestEncryptEnv_EmptyPassphrase(t *testing.T) {
	env := makeEncryptEnv()
	opts := DefaultEncryptOptions()
	opts.Passphrase = ""

	_, err := EncryptEnv(env, opts)
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

func TestDecryptValue_RoundTrip(t *testing.T) {
	plain := "super-secret-value"
	passphrase := "mypassphrase"

	enc, err := encryptValue(plain, passphrase)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	decrypted, err := DecryptValue(enc, passphrase)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if decrypted != plain {
		t.Errorf("expected %q, got %q", plain, decrypted)
	}
}

func TestDecryptValue_WrongPassphrase(t *testing.T) {
	enc, _ := encryptValue("secret", "correctpass")
	_, err := DecryptValue(enc, "wrongpass")
	if err == nil {
		t.Fatal("expected error with wrong passphrase")
	}
}

func TestWriteEncryptReport_WithChanges(t *testing.T) {
	original := map[string]string{"DB_PASSWORD": "plain", "APP_NAME": "app"}
	encrypted := map[string]string{"DB_PASSWORD": "enc:abc123==", "APP_NAME": "app"}

	var buf bytes.Buffer
	WriteEncryptReport(&buf, original, encrypted)
	out := buf.String()

	if !strings.Contains(out, "DB_PASSWORD") {
		t.Errorf("expected DB_PASSWORD in report, got: %s", out)
	}
	if !strings.Contains(out, "1 key") {
		t.Errorf("expected count in report, got: %s", out)
	}
}

func TestWriteEncryptReport_NoChanges(t *testing.T) {
	original := map[string]string{"APP_NAME": "app"}
	encrypted := map[string]string{"APP_NAME": "app"}

	var buf bytes.Buffer
	WriteEncryptReport(&buf, original, encrypted)
	out := buf.String()

	if !strings.Contains(out, "No keys encrypted") {
		t.Errorf("expected no-op message, got: %s", out)
	}
}
