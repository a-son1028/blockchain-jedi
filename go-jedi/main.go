package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"jedi"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/ucbrise/jedi-pairing/lang/go/wkdibe"
)

const TestPatternSize = 20

var TestHierarchy = []byte("testHierarchy")

const quote1 = "Imagination is more important than knowledge. --Albert Einstein"
const quote2 = "Today is your day! / Your mountain is waiting. / So... get on your way! --Theodor Seuss Geisel"

type DecryptRequest struct {
	URI              string `string:"uri" binding:"required"`
	ENCRYPTEDMESSAGE string `string:"encryptedMessage" binding:"required"`
	KEY              string `string:"key" binding:"required"`
}

type EncryptRequest struct {
	URI     string `string:"uri" binding:"required"`
	MESSAGE string `string:"message" binding:"required"`
}

type MeasureUsage struct {
	memory uint64
	cpuPercentage float64
}

func measureMemoryUsage()(MeasureUsage) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	fmt.Printf("Number of CPUs: %d\n", runtime.NumCPU())
	fmt.Printf("Alloc = %v MiB\n", memStats.Alloc/1024/1024)
	fmt.Printf("TotalAlloc = %v MiB\n", memStats.TotalAlloc/1024/1024)
	fmt.Printf("Sys = %v MiB\n", memStats.Sys/1024/1024)
	fmt.Printf("NumGC = %v\n\n", memStats.NumGC)
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		fmt.Println(err)
	}

	return MeasureUsage{ memory: memStats.Alloc, cpuPercentage: percentages[0] };
}

func main() {
	ctx := context.Background()
	_, store := NewTestKeyStore()
	encoder := jedi.NewDefaultPatternEncoder(TestPatternSize - jedi.MaxTimeLength)

	state := NewTestState()
	now := time.Now()

	// testMessageTransfer(state, TestHierarchy, "a/b/c", now, quote1)

	r := gin.Default()
	r.GET("/jedi-private-key", func(c *gin.Context) {
		uri := "a/b/c"

		start := time.Unix(1565119330, 0)
		end := time.Unix(1565219330, 0)

		parent := c.DefaultQuery("parent", "")
		if parent != "" {
			// delegation := new(jedi.Delegation)

			// if !delegation.Unmarshal([]byte(parent)) {
			// 	fmt.Println("Could not unmarshal delegation")
			// }

			// fmt.Println(1)
			// fmt.Println(delegation)

			// c.JSON(200, gin.H{
			// 	"data": "Test",
			// })
			uri := "a/b/c/d"
			delegation, err := jedi.Delegate(ctx, store, encoder, TestHierarchy, uri, start, end, jedi.DecryptPermission|jedi.SignPermission)
			if err != nil {
				fmt.Println(err)
			}

			marshalled := delegation.Marshal()

			c.JSON(200, gin.H{
				"data": marshalled,
			})
			return
		}

		startTime := time.Now()
		delegation, err := jedi.Delegate(ctx, store, encoder, TestHierarchy, uri, start, end, jedi.DecryptPermission|jedi.SignPermission)
		if err != nil {
			fmt.Println(err)
		}
		endTime := time.Now()

		marshalled := delegation.Marshal()

		c.JSON(200, gin.H{
			"time": endTime.Sub(startTime).Microseconds(),
			"data": marshalled,
		})
	})

	r.POST("/encrypt", func(c *gin.Context) {
		measureUsageBefore := measureMemoryUsage()
		fmt.Println(measureUsageBefore)
		var err error

		var encryptRequest EncryptRequest
		c.BindJSON(&encryptRequest)

		message := encryptRequest.MESSAGE
		uri := encryptRequest.URI

		startTime := time.Now()
		var encrypted []byte
		if encrypted, err = state.Encrypt(ctx, TestHierarchy, uri, now, []byte(message)); err != nil {
			fmt.Println(err)
		}
		endTime := time.Now()
		measureUsage := measureMemoryUsage()

		var decrypted []byte
		if decrypted, err = state.Decrypt(ctx, TestHierarchy, uri, now, encrypted); err != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		if !bytes.Equal(decrypted, []byte(message)) {
			fmt.Println("Original and decrypted messages differ")
		}

		

		c.JSON(200, gin.H{
			"time": endTime.Sub(startTime).Microseconds(),
			"memoryUsage": measureUsage.memory,
			"cpuPercentage": measureUsage.cpuPercentage,
			"data": base64.StdEncoding.EncodeToString(encrypted),
		})
	})

	r.POST("/decrypt", func(c *gin.Context) {
		measureMemoryUsage()
		var err error
		var decryptRequest DecryptRequest
		c.BindJSON(&decryptRequest)

		encrypted, err := base64.StdEncoding.DecodeString(decryptRequest.ENCRYPTEDMESSAGE)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		uri := decryptRequest.URI

		startTime := time.Now()
		var decrypted []byte
		if decrypted, err = state.Decrypt(ctx, TestHierarchy, uri, now, encrypted); err != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		endTime := time.Now()
		measureUsage := measureMemoryUsage()

		if !bytes.Equal(decrypted, []byte("test message")) {
			fmt.Println("Original and decrypted messages differ")
		}

		str := string(decrypted)

		c.JSON(200, gin.H{
			"time": endTime.Sub(startTime).Microseconds(),
			"memoryUsage": measureUsage.memory,
			"cpuPercentage": measureUsage.cpuPercentage,
			"data": str,
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
