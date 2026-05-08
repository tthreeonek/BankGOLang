package model

type Card struct {
	ID           int    `json:"id"`
	AccountID    int    `json:"account_id"`
	EncryptedPAN string `json:"-"`          // зашифрованный номер
	PANHMAC      string `json:"-"`          // HMAC для проверки целостности
	EncryptedExp string `json:"-"`          // зашифрованный срок (MM/YY)
	CVVHash      string `json:"-"`          // bcrypt хэш CVV
	MaskedPAN    string `json:"masked_pan"` // последние 4 цифры
}
