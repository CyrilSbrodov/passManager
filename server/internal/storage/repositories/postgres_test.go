package repositories

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CyrilSbrodov/passManager.git/server/cmd/config"
	"github.com/CyrilSbrodov/passManager.git/server/internal/models"
)

var (
	CFG config.Config
)

func TestMain(m *testing.M) {
	cfg := config.ConfigInit()
	CFG = *cfg
	CFG.DatabaseDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	os.Exit(m.Run())
}

func TestStore_Register(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users", "cards", "binary_table", "text_table", "passwords")

	tests := []struct {
		name        string
		user        models.User
		expectedUID string
	}{
		{
			name: "ok",
			user: models.User{
				UID:      "",
				Login:    "Login",
				Password: "Password",
			},
			expectedUID: "1",
		},
		{
			name: "false",
			user: models.User{
				UID:      "",
				Login:    "Login",
				Password: "Password",
			},
			expectedUID: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uid, _ := s.Register(&tt.user)
			assert.Equal(t, tt.expectedUID, uid)
		})
	}
}

func TestStore_Login(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users", "cards", "binary_table", "text_table", "passwords")
	uid, err := s.Register(&models.User{
		Login:    "test",
		Password: "testPass",
	})
	assert.NoError(t, err)
	tests := []struct {
		name        string
		user        models.User
		expectedUID string
	}{
		{
			name: "ok",
			user: models.User{
				UID:      "",
				Login:    "test",
				Password: "testPass",
			},
			expectedUID: uid,
		},
		{
			name: "wrong login",
			user: models.User{
				UID:      "",
				Login:    "Login",
				Password: "Password",
			},
			expectedUID: "",
		},
		{
			name: "wrong pass",
			user: models.User{
				UID:      "",
				Login:    "test",
				Password: "Password",
			},
			expectedUID: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, _ := s.Login(&tt.user)
			assert.Equal(t, tt.expectedUID, id)
		})
	}
}

func TestStore_CollectBinary(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users", "cards", "binary_table", "text_table", "passwords")
	uid, err := s.Register(&models.User{
		Login:    "test",
		Password: "testPass",
	})
	status, err := s.CollectBinary(&models.CryptoBinaryData{
		Data: nil,
	}, uid)
	assert.NoError(t, err)
	assert.Equal(t, 200, status)
}

func TestStore_GetBinary(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users", "cards", "binary_table", "text_table", "passwords")
	uid, err := s.Register(&models.User{
		Login:    "test",
		Password: "testPass",
	})
	status, err := s.CollectBinary(&models.CryptoBinaryData{
		Data: []byte("dada"),
	}, uid)
	statusCode, b, err := s.GetBinary(uid)
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, 200, status)
	assert.NotNil(t, b)
}

func TestStore_DeleteBinary(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users", "cards", "binary_table", "text_table", "passwords")
	uid, err := s.Register(&models.User{
		Login:    "test",
		Password: "testPass",
	})
	status, err := s.CollectBinary(&models.CryptoBinaryData{
		Data: []byte("dada"),
	}, uid)
	statusCode, b, err := s.GetBinary(uid)
	sc, errDel := s.DeleteBinary(&b[0], uid)
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, 200, status)
	assert.NotNil(t, b)
	assert.NoError(t, errDel)
	assert.Equal(t, 200, sc)
}

func TestStore_UpdateBinary(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users", "cards", "binary_table", "text_table", "passwords")
	uid, err := s.Register(&models.User{
		Login:    "test",
		Password: "testPass",
	})
	status, err := s.CollectBinary(&models.CryptoBinaryData{
		Data: []byte("dada"),
	}, uid)
	statusCode, b, err := s.GetBinary(uid)
	statusUpd, errUpd := s.UpdateBinary(&models.CryptoBinaryData{
		UID:  b[0].UID,
		Data: []byte("dada"),
	}, uid)
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, 200, status)
	assert.NotNil(t, b)
	assert.NoError(t, errUpd)
	assert.Equal(t, 200, statusUpd)
}

func TestStore_CollectText(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users", "cards", "binary_table", "text_table", "passwords")
	uid, err := s.Register(&models.User{
		Login:    "test",
		Password: "testPass",
	})
	status, err := s.CollectText(&models.CryptoTextData{
		Text: []byte("text"),
	}, uid)
	assert.NoError(t, err)
	assert.Equal(t, 200, status)
}

func TestStore_GetText(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users", "cards", "binary_table", "text_table", "passwords")
	uid, err := s.Register(&models.User{
		Login:    "test",
		Password: "testPass",
	})
	status, err := s.CollectText(&models.CryptoTextData{
		Text: []byte("text"),
	}, uid)
	statusGet, text, errGet := s.GetText(uid)
	assert.NoError(t, err)
	assert.NoError(t, errGet)
	assert.Equal(t, 200, status)
	assert.Equal(t, 200, statusGet)
	assert.NotNil(t, text)
}

func TestStore_DeleteText(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users", "cards", "binary_table", "text_table", "passwords")
	uid, err := s.Register(&models.User{
		Login:    "test",
		Password: "testPass",
	})
	status, err := s.CollectText(&models.CryptoTextData{
		Text: []byte("text"),
	}, uid)
	statusGet, text, errGet := s.GetText(uid)
	statusDel, errDel := s.DeleteText(&text[0], uid)
	assert.NoError(t, err)
	assert.NoError(t, errGet)
	assert.NoError(t, errDel)
	assert.Equal(t, 200, status)
	assert.Equal(t, 200, statusGet)
	assert.Equal(t, 200, statusDel)
	assert.NotNil(t, text)
}

func TestStore_UpdateText(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users", "cards", "binary_table", "text_table", "passwords")
	uid, err := s.Register(&models.User{
		Login:    "test",
		Password: "testPass",
	})
	status, err := s.CollectText(&models.CryptoTextData{
		Text: []byte("text"),
	}, uid)
	statusGet, text, errGet := s.GetText(uid)
	statusDel, errDel := s.UpdateText(&models.CryptoTextData{
		UID:  text[0].UID,
		Text: []byte("text"),
	}, uid)
	assert.NoError(t, err)
	assert.NoError(t, errGet)
	assert.NoError(t, errDel)
	assert.Equal(t, 200, status)
	assert.Equal(t, 200, statusGet)
	assert.Equal(t, 200, statusDel)
	assert.NotNil(t, text)
}

func TestStore_CollectPassword(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users", "cards", "binary_table", "text_table", "passwords")
	uid, err := s.Register(&models.User{
		Login:    "test",
		Password: "testPass",
	})
	status, err := s.CollectPassword(&models.CryptoPassword{
		Login: []byte("login"),
		Pass:  []byte("pass"),
	}, uid)
	assert.NoError(t, err)
	assert.Equal(t, 200, status)
}

func TestStore_GetPassword(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users", "cards", "binary_table", "text_table", "passwords")
	uid, err := s.Register(&models.User{
		Login:    "test",
		Password: "testPass",
	})
	status, err := s.CollectPassword(&models.CryptoPassword{
		Login: []byte("login"),
		Pass:  []byte("pass"),
	}, uid)
	statusGet, p, errGet := s.GetPassword(uid)
	assert.NoError(t, err)
	assert.NoError(t, errGet)
	assert.Equal(t, 200, status)
	assert.Equal(t, 200, statusGet)
	assert.NotNil(t, p)
}

func TestStore_DeletePassword(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users", "cards", "binary_table", "text_table", "passwords")
	uid, err := s.Register(&models.User{
		Login:    "test",
		Password: "testPass",
	})
	status, err := s.CollectPassword(&models.CryptoPassword{
		Login: []byte("login"),
		Pass:  []byte("pass"),
	}, uid)
	statusGet, p, errGet := s.GetPassword(uid)
	statusDel, errDel := s.DeletePassword(&p[0], uid)
	assert.NoError(t, err)
	assert.NoError(t, errGet)
	assert.NoError(t, errDel)
	assert.Equal(t, 200, status)
	assert.Equal(t, 200, statusGet)
	assert.Equal(t, 200, statusDel)
	assert.NotNil(t, p)
}

func TestStore_UpdatePassword(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users", "cards", "binary_table", "text_table", "passwords")
	uid, err := s.Register(&models.User{
		Login:    "test",
		Password: "testPass",
	})
	status, err := s.CollectPassword(&models.CryptoPassword{
		Login: []byte("login"),
		Pass:  []byte("pass"),
	}, uid)
	statusGet, p, errGet := s.GetPassword(uid)
	statusDel, errDel := s.UpdatePassword(&models.CryptoPassword{
		UID:   p[0].UID,
		Login: []byte("loginsda"),
		Pass:  []byte("pass"),
	}, uid)
	assert.NoError(t, err)
	assert.NoError(t, errGet)
	assert.NoError(t, errDel)
	assert.Equal(t, 200, status)
	assert.Equal(t, 200, statusGet)
	assert.Equal(t, 200, statusDel)
	assert.NotNil(t, p)
}

func TestStore_CollectCard(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users", "cards", "binary_table", "text_table", "passwords")
	uid, err := s.Register(&models.User{
		Login:    "test",
		Password: "testPass",
	})
	status, err := s.CollectCard(&models.CryptoCard{
		Number: nil,
		Name:   nil,
		CVC:    nil,
	}, uid)
	assert.NoError(t, err)
	assert.Equal(t, 200, status)
}

func TestStore_GetCards(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users", "cards", "binary_table", "text_table", "passwords")
	uid, err := s.Register(&models.User{
		Login:    "test",
		Password: "testPass",
	})
	status, err := s.CollectCard(&models.CryptoCard{
		Number: nil,
		Name:   nil,
		CVC:    nil,
	}, uid)
	statusGet, c, errGet := s.GetCards(uid)
	assert.NoError(t, err)
	assert.NoError(t, errGet)
	assert.Equal(t, 200, status)
	assert.Equal(t, 200, statusGet)
	assert.NotNil(t, c)
}

func TestStore_DeleteCard(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users", "cards", "binary_table", "text_table", "passwords")
	uid, err := s.Register(&models.User{
		Login:    "test",
		Password: "testPass",
	})
	status, err := s.CollectCard(&models.CryptoCard{
		Number: nil,
		Name:   nil,
		CVC:    nil,
	}, uid)
	statusGet, c, errGet := s.GetCards(uid)
	statusDel, errDel := s.DeleteCard(&c[0], uid)
	assert.NoError(t, err)
	assert.NoError(t, errGet)
	assert.NoError(t, errDel)
	assert.Equal(t, 200, status)
	assert.Equal(t, 200, statusGet)
	assert.Equal(t, 200, statusDel)
	assert.NotNil(t, c)
}

func TestStore_UpdateCard(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users", "cards", "binary_table", "text_table", "passwords")
	uid, err := s.Register(&models.User{
		Login:    "test",
		Password: "testPass",
	})
	status, err := s.CollectCard(&models.CryptoCard{
		Number: nil,
		Name:   nil,
		CVC:    nil,
	}, uid)
	statusGet, c, errGet := s.GetCards(uid)
	statusDel, errDel := s.UpdateCard(&models.CryptoCard{
		UID:    c[0].UID,
		Number: []byte("fdsfsdf"),
		Name:   []byte("234234"),
		CVC:    []byte("24234"),
	}, uid)
	assert.NoError(t, err)
	assert.NoError(t, errGet)
	assert.NoError(t, errDel)
	assert.Equal(t, 200, status)
	assert.Equal(t, 200, statusGet)
	assert.Equal(t, 200, statusDel)
	assert.NotNil(t, c)
}
