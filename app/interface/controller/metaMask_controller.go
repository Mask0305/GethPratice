package controller

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/util/log"
)

var (
	ErrUserNotExists  = errors.New("user does not exist")
	ErrUserExists     = errors.New("user already exists")
	ErrInvalidAddress = errors.New("invalid address")
	ErrInvalidNonce   = errors.New("invalid nonce")
	ErrMissingSig     = errors.New("signature is missing")
	ErrAuthError      = errors.New("authentication error")
)

var (
	hexRegex   *regexp.Regexp = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
	nonceRegex *regexp.Regexp = regexp.MustCompile(`^[0-9]+$`)
)

type User struct {
	Address string
	Nonce   string
}
type MemStorage struct {
	lock  sync.RWMutex
	users map[string]User
}

func NewMemStorage() *MemStorage {
	ans := MemStorage{
		users: make(map[string]User),
	}
	return &ans
}

func (m *MemStorage) Get(address string) (User, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	u, exists := m.users[address]
	if !exists {
		return u, ErrUserNotExists
	}
	return u, nil
}
func (m *MemStorage) Update(user User) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.users[user.Address] = user
	return nil
}

type metaMaskcontroller struct {
}

type MetaMaskController interface {
	Router(app *gin.Engine)
}

func NewMetaMaskController() MetaMaskController {
	return &metaMaskcontroller{}
}

func (c *metaMaskcontroller) Router(app *gin.Engine) {

	// init
	storage := NewMemStorage()
	n, _ := GetNonce()
	storage.users["0x28cf3988e6b9a0a121dcbb9a530e9050f10a076f"] = User{
		Address: "0x28cf3988e6b9a0a121dcbb9a530e9050f10a076f",
		Nonce:   n,
	}

	Group := app.Group("/metamask")
	{
		Group.POST("/signin", c.Signin(storage))
		Group.GET("/:address/nonce", c.GetUserNonce(storage))
		Group.OPTIONS("/:address/nonce", c.GetUserNonce(storage))

	}
}

type SigninPayload struct {
	Address string `json:"address"`
	Nonce   string `json:"nonce"`
	Sig     string `json:"sig"`
}

func (s SigninPayload) Validate() error {

	if !hexRegex.MatchString(s.Address) {
		return ErrInvalidAddress
	}
	if !nonceRegex.MatchString(s.Nonce) {
		return ErrInvalidNonce
	}
	if len(s.Sig) == 0 {
		return ErrMissingSig
	}
	return nil
}
func Authenticate(storage *MemStorage, address string, nonce string, sigHex string) (User, error) {
	user, err := storage.Get(address)
	if err != nil {
		return user, err
	}
	if user.Nonce != nonce {
		return user, ErrAuthError
	}

	sig := hexutil.MustDecode(sigHex)
	// https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L516
	// check here why I am subtracting 27 from the last byte
	sig[crypto.RecoveryIDOffset] -= 27
	msg := accounts.TextHash([]byte(nonce))
	recovered, err := crypto.SigToPub(msg, sig)
	if err != nil {
		return user, err
	}
	recoveredAddr := crypto.PubkeyToAddress(*recovered)

	if user.Address != strings.ToLower(recoveredAddr.Hex()) {
		return user, ErrAuthError
	}

	// update the nonce here so that the signature cannot be resused
	nonce, err = GetNonce()
	if err != nil {
		return user, err
	}
	user.Nonce = nonce
	storage.Update(user)

	return user, nil
}

var (
	max  *big.Int
	once sync.Once
)

func GetNonce() (string, error) {
	once.Do(func() {
		max = new(big.Int)
		max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))
	})
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return n.Text(10), nil
}

func (c *metaMaskcontroller) Signin(storage *MemStorage) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		p := new(SigninPayload)
		if err := ctx.BindJSON(p); err != nil {
			ctx.Error(err)
			return
		}

		if err := p.Validate(); err != nil {
			ctx.Error(err)
			return
		}

		address := strings.ToLower(p.Address)
		user, err := Authenticate(storage, address, p.Nonce, p.Sig)
		if err != nil {
			log.Error(err)
			return
		}
		spew.Dump(user)

		ctx.JSON(200, user)
	}

}

func (c *metaMaskcontroller) GetUserNonce(storage *MemStorage) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		fmt.Println(11)
		address := ctx.Param("address")
		if !hexRegex.MatchString(address) {
			ctx.Error(errors.New("地址正則錯誤"))
			return
		}
		user, err := storage.Get(strings.ToLower(address))
		if err != nil {
			log.Error(err)
			return
		}

		ctx.JSON(200, map[string]interface{}{
			"Nonce": user.Nonce,
		})
	}
}
