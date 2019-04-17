package crypto

import (
	"../models"
)

func DecryptEmail(key []byte, email *models.Email) {
    if email.Title != "" { email.Title = Decrypt(key, email.Title) }
    if email.Body  != "" { email.Body  = Decrypt(key, email.Body)  }
    if email.Date  != "" { email.Date  = Decrypt(key, email.Date)  }
}
