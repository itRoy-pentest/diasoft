package security

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Секретная соль приложения (хранится в .env, не в коде!)
const appPepper = "super-secret-salt-12345"

// GenerateIdentityHash создает защищенный отпечаток для поиска
func GenerateIdentityHash(diplomaNum string, univID string) string {
	// Смешиваем номер диплома, ID вуза и секретный перец
	raw := fmt.Sprintf("%s|%s|%s", diplomaNum, univID, appPepper)
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}

// Создаем ЭЦП (digital_signature) для записи в таблицу diplomas
func CreateDigitalSignature(data string, privKey ed25519.PrivateKey) string {
	hash := sha256.Sum256([]byte(data))
	sig := ed25519.Sign(privKey, hash[:])
	return hex.EncodeToString(sig)
}
