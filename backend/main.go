package main

import (
	"crypto/ed25519"
	"fmt"
	"log"

	"diasoft-auth/database"
	"diasoft-auth/security"
	"diasoft-auth/storage"
)

func main() {
	// 1. Подключаемся к БД
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	// 2. Генерируем ключи (в реальности берем из БД или Vault)
	pub, privKey, _ := ed25519.GenerateKey(nil)

	// 3. Данные студента
	realName := "Сергей Игоревич"
	realNum := "106104-777"
	uID := "550e8400-e29b-41d4-a716-446655440000" // Тестовый UUID

	// 4. Маскирование и хеширование
	maskedName := security.MaskString(realName)
	idHash := security.GenerateIdentityHash(realNum, uID)

	// 5. Формируем запись для вставки
	entry := storage.DiplomaEntry{
		UnivID:       uID,
		StudentID:    "00000000-0000-0000-0000-000000000000", // Заглушка
		StudentName:  maskedName,
		DiplomaNum:   security.MaskDiploma(realNum),
		IdentityHash: idHash,
		Signature:    security.CreateDigitalSignature(idHash, privKey),
		IssueYear:    2026,
	}

	// 6. Сохранение
	err = storage.SaveDiploma(db, entry)
	if err != nil {
		log.Printf("Ошибка сохранения в БД: %v", err)
	} else {
		fmt.Println("✅ Запись успешно создана!")
		fmt.Printf("Публичный ключ ВУЗа: %x\n", pub)
		fmt.Printf("Поиск по хешу %s...\n", idHash[:10])
		
		// 7. Проверка поиска
		exists, _ := storage.FindByIdentityHash(db, idHash)
		fmt.Printf("Диплом найден в базе: %v\n", exists)
	}
}
