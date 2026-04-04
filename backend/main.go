package main

import (
"crypto/ed25519"
"encoding/json"
"fmt"
"log"
"net/http"
"time"

"diasoft-auth/database"
"diasoft-auth/security"
"diasoft-auth/storage"

"github.com/google/uuid"
)

type IssueRequest struct {
RealName   string `json:"real_name"`
DiplomaNum string `json:"diploma_num"`
UnivID     string `json:"univ_id"`
StudentID  string `json:"student_id"`
Year       int    `json:"year"`
}

func main() {
db, err := database.Connect()
if err != nil {
log.Fatal("Ошибка подключения к БД:", err)
}
defer db.Close()

pub, privKey, _ := ed25519.GenerateKey(nil)
rdb := storage.ConnectRedis()

// АДМИНКА: Выпуск диплома
http.HandleFunc("/admin/issue", func(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
return
}

if r.Header.Get("X-Admin-Token") != "secret-diasoft-2026" {
http.Error(w, "Доступ запрещен", http.StatusUnauthorized)
return
}

var req IssueRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, "Ошибка данных: "+err.Error(), http.StatusBadRequest)
return
}

maskedName := security.MaskString(req.RealName)
idHash := security.GenerateIdentityHash(req.DiplomaNum, req.UnivID)

entry := storage.DiplomaEntry{
UnivID:       req.UnivID,
StudentID:    req.StudentID,
StudentName:  maskedName,
DiplomaNum:   security.MaskDiploma(req.DiplomaNum),
IdentityHash: idHash,
Signature:    security.CreateDigitalSignature(idHash, privKey),
IssueYear:    req.Year,
}

if err := storage.SaveDiploma(db, entry); err != nil {
http.Error(w, "Ошибка БД: "+err.Error(), http.StatusInternalServerError)
return
}

token := uuid.New().String()
_ = storage.CreatePublicToken(rdb, token, idHash, 24*time.Hour)

w.Header().Set("Content-Type", "application/json")
fmt.Fprintf(w, `{"status":"created", "verify_url":"http://217.60.237.170:8080/verify/%s"}`, token)
log.Printf("Админ выпустил диплом для: %s", req.RealName)
})

// ПРОВЕРКА: Публичный доступ
http.HandleFunc("/verify/", func(w http.ResponseWriter, r *http.Request) {
token := r.URL.Path[len("/verify/"):]
idHash, err := storage.GetHashByToken(rdb, token)
if err != nil {
http.Error(w, "Ссылка протухла или неверная", http.StatusNotFound)
return
}

exists, _ := storage.FindByIdentityHash(db, idHash)

w.Header().Set("Content-Type", "application/json")
fmt.Fprintf(w, `{"valid": %v, "identity_hash": "%s", "issuer_pub_key": "%x"}`, exists, idHash, pub)
})

fmt.Println("🚀 Diasoft Admin API запущен на http://217.60.237.170:8080")
log.Fatal(http.ListenAndServe(":8080", nil))
}
