package main

import (
	"bytes"
	"context"
	"fmt"
	"jedi"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ucbrise/jedi-pairing/lang/go/wkdibe"
)

const TestPatternSize = 20

var TestHierarchy = []byte("testHierarchy")

const quote1 = "Imagination is more important than knowledge. --Albert Einstein"
const quote2 = "Today is your day! / Your mountain is waiting. / So... get on your way! --Theodor Seuss Geisel"

func main() {
	ctx := context.Background()
	_, store := NewTestKeyStore()
	encoder := jedi.NewDefaultPatternEncoder(TestPatternSize - jedi.MaxTimeLength)

	// state := NewTestState()
	// now := time.Now()

	// testMessageTransfer(state, TestHierarchy, "a/b/c", now, quote1)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		uri := "a/b/c"
		start := time.Unix(1565119330, 0)
		end := time.Unix(1565219330, 0)
		delegation, err := jedi.Delegate(ctx, store, encoder, TestHierarchy, uri, start, end, jedi.DecryptPermission|jedi.SignPermission)
		if err != nil {
			fmt.Println(err)
		}
		marshalled := delegation.Marshal()

		c.JSON(200, gin.H{
			"data": marshalled,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
	fmt.Println("DONE!")

}

type TestPublicInfo struct {
	params *wkdibe.Params
}

type TestKeyStore struct {
	params *wkdibe.Params
	master *wkdibe.MasterKey
}

func (tpi *TestPublicInfo) ParamsForHierarchy(ctx context.Context, hierarchy []byte) (*wkdibe.Params, error) {
	return tpi.params, nil
}

func testMessageTransfer(state *jedi.ClientState, hierarchy []byte, uri string, timestamp time.Time, message string) {
	var err error
	ctx := context.Background()

	var encrypted []byte
	if encrypted, err = state.Encrypt(ctx, hierarchy, uri, timestamp, []byte(message)); err != nil {
		fmt.Println(err)
	}

	var decrypted []byte
	if decrypted, err = state.Decrypt(ctx, hierarchy, uri, timestamp, encrypted); err != nil {
		fmt.Println(err)
	}

	if !bytes.Equal(decrypted, []byte(message)) {
		fmt.Println("Original and decrypted messages differ")
	}
}

func NewTestKeyStore() (*TestPublicInfo, *TestKeyStore) {
	tks := new(TestKeyStore)
	tks.params, tks.master = wkdibe.Setup(TestPatternSize, true)
	tpi := new(TestPublicInfo)
	tpi.params = tks.params
	return tpi, tks
}

func (tks *TestKeyStore) KeyForPattern(ctx context.Context, hierarchy []byte, pattern jedi.Pattern) (*wkdibe.Params, *wkdibe.SecretKey, error) {
	empty := make(jedi.Pattern, TestPatternSize)
	return tks.params, wkdibe.KeyGen(tks.params, tks.master, empty.ToAttrs()), nil
}

func NewTestState() *jedi.ClientState {
	info, store := NewTestKeyStore()
	encoder := jedi.NewDefaultPatternEncoder(TestPatternSize - jedi.MaxTimeLength)
	return jedi.NewClientState(info, store, encoder, 1<<20)
}
