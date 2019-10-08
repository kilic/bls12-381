package blssig

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"
)

type testVectorDataAggregatePubkeys struct {
	Input  []string `yaml:"input"`
	Output string   `yaml:"output"`
}

func testFiles(path string) [][]byte {
	var testVectors [][]byte
	if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.Name() == "data.yaml" {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			testVectors = append(testVectors, data)
		}
		return nil
	}); err != nil {
		panic(err)
	}
	return testVectors
}

func TestAggregatePubkeys(t *testing.T) {
	testVectors := testFiles("tests/aggregate_pubkeys/small")
	for i, data := range testVectors {
		var v testVectorDataAggregatePubkeys
		if err := yaml.Unmarshal(data, &v); err != nil {
			panic(err)
		}
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			aggregated := new(PublicKey)
			expectedOutputString := toBytes(48, v.Output)
			for i := 0; i < len(v.Input); i++ {
				pubkey, err := NewPublicKeyFromCompresssed(toBytes(48, v.Input[i]))
				if err != nil {
					t.Fatal(err)
				}
				AggregatePublicKey(aggregated, pubkey)
			}
			aggregatedBytes := PublicKeyToCompressed(aggregated)
			if !bytes.Equal(expectedOutputString, aggregatedBytes) {
				t.Fatalf("\nwant: %x\nhave: %x\n", expectedOutputString, aggregatedBytes)
			}
		})
	}
}

type testVectorDataAggregateSignatures struct {
	Input  []string `yaml:"input"`
	Output string   `yaml:"output"`
}

func TestAggregateSignatures(t *testing.T) {
	testVectors := testFiles("tests/aggregate_sigs/small")
	for i, data := range testVectors {
		var v testVectorDataAggregateSignatures
		if err := yaml.Unmarshal(data, &v); err != nil {
			panic(err)
		}
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			aggregated := new(Signature)
			expectedOutputString := toBytes(96, v.Output)
			for i := 0; i < len(v.Input); i++ {
				signature, err := NewSignatureFromCompresssed(toBytes(96, v.Input[i]))
				if err != nil {
					t.Fatal(err)
				}
				AggregateSignature(aggregated, signature)
			}
			aggregatedBytes := SignatureToCompressed(aggregated)
			if !bytes.Equal(expectedOutputString, aggregatedBytes) {
				t.Fatalf("\nwant: %x\nhave: %x\n", expectedOutputString, aggregatedBytes)
			}
		})
	}
}

type testVectorDataPrivToPub struct {
	Input  string `yaml:"input"`
	Output string `yaml:"output"`
}

func TestPrivToPub(t *testing.T) {
	testVectors := testFiles("tests/priv_to_pub/small")
	for i, data := range testVectors {
		var v testVectorDataPrivToPub
		if err := yaml.Unmarshal(data, &v); err != nil {
			panic(err)
		}
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			expectedOutputString := toBytes(48, v.Output)
			secretKey, err := SecretKeyFromBytes(toBytes(32, v.Input))
			if err != nil {
				panic(err)
			}
			pub := PublicKeyFromSecretKey(secretKey)
			pubBytes := PublicKeyToCompressed(pub)
			if !bytes.Equal(expectedOutputString, pubBytes) {
				t.Fatalf("\nwant: %x\nhave: %x\n", expectedOutputString, pubBytes)
			}
		})
	}
}

type testVectorDataHashCompressed struct {
	Input struct {
		Message string `yaml:"message"`
		Domain  string `yaml:"domain"`
	} `yaml:"input"`
	Output [2]string `yaml:"output"`
}

func TestHashCompressed(t *testing.T) {
	testVectors := testFiles("tests/msg_hash_compressed/small")
	for i, data := range testVectors {
		var v testVectorDataHashCompressed
		if err := yaml.Unmarshal(data, &v); err != nil {
			panic(err)
		}
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			expectedOutputString := toBytes(48, v.Output[0], v.Output[1])
			msg, domain := [32]byte{}, [8]byte{}
			copy(msg[:], toBytes(32, v.Input.Message))
			copy(domain[:], toBytes(8, v.Input.Domain))
			msgHash := HashWithDomain(msg, domain)
			msgHashCompressed := g2ToCompressed(msgHash)
			if !bytes.Equal(expectedOutputString, msgHashCompressed) {
				t.Fatalf("\nwant: %x\nhave: %x\n", expectedOutputString, msgHashCompressed)
			}
		})
	}
}

type testVectorSignMessage struct {
	Input struct {
		Secret  string `yaml:"privkey"`
		Message string `yaml:"message"`
		Domain  string `yaml:"domain"`
	} `yaml:"input"`
	Output string `yaml:"output"`
}

func TestSignMessage(t *testing.T) {
	testVectors := testFiles("tests/sign_msg/small")
	for i, data := range testVectors {
		var v testVectorSignMessage
		if err := yaml.Unmarshal(data, &v); err != nil {
			panic(err)
		}
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			expectedOutputString := toBytes(96, v.Output)
			msg, domain := [32]byte{}, [8]byte{}
			copy(msg[:], toBytes(32, v.Input.Message))
			copy(domain[:], toBytes(8, v.Input.Domain))
			secretKey, err := SecretKeyFromBytes(toBytes(32, v.Input.Secret))
			if err != nil {
				t.Fatal(err)
			}
			signature := Sign(msg, domain, secretKey)
			signatureBytes := SignatureToCompressed(signature)
			if !bytes.Equal(expectedOutputString, signatureBytes) {
				t.Fatalf("\nwant: %x\nhave: %x\n", expectedOutputString, signatureBytes)
			}
		})
	}
}

func TestVerifySignature(t *testing.T) {
	msg := [32]byte{1}
	domain := [8]byte{1}
	secretKey, err := RandSecretKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	pubKey := PublicKeyFromSecretKey(secretKey)
	signature := Sign(msg, domain, secretKey)
	if !Verify(msg, domain, signature, pubKey) {
		t.Fatal("verification fails")
	}
}

func TestVerifySignatureAggregatedCommon(t *testing.T) {
	msg := [32]byte{1}
	domain := [8]byte{1}
	signerSize := 128
	pubkeys := make([]*PublicKey, signerSize)
	signatures := make([]*Signature, signerSize)
	for i := 0; i < signerSize; i++ {
		secretKey, err := RandSecretKey(rand.Reader)
		if err != nil {
			panic(err)
		}
		pubkeys[i] = PublicKeyFromSecretKey(secretKey)
		signatures[i] = Sign(msg, domain, secretKey)
	}
	signature := AggregateSignatures(signatures)
	if !VerifyAggregateCommon(msg, domain, pubkeys, signature) {
		t.Fatal("verification fails")
	}
}

func TestVerifySignatureAggregated(t *testing.T) {
	msgs := make([][32]byte, 10)
	domain := [8]byte{1}
	signerSize := 10
	pubkeys := make([]*PublicKey, signerSize)
	signatures := make([]*Signature, signerSize)
	for i := 0; i < signerSize; i++ {
		msgs[i][0] = byte(i)
		secretKey, err := RandSecretKey(rand.Reader)
		if err != nil {
			panic(err)
		}
		pubkeys[i] = PublicKeyFromSecretKey(secretKey)
		signatures[i] = Sign(msgs[i], domain, secretKey)
	}
	signature := AggregateSignatures(signatures)
	if !VerifyAggregate(msgs, domain, pubkeys, signature) {
		t.Fatal("verification fails")
	}
}

func BenchmarkVerifySignatureAggregatedCommon(t *testing.B) {
	msg := [32]byte{1}
	domain := [8]byte{1}
	signerSize := 128
	pubkeys := make([]*PublicKey, signerSize)
	signatures := make([]*Signature, signerSize)
	for i := 0; i < signerSize; i++ {
		secretKey, err := RandSecretKey(rand.Reader)
		if err != nil {
			panic(err)
		}
		pubkeys[i] = PublicKeyFromSecretKey(secretKey)
		signatures[i] = Sign(msg, domain, secretKey)
	}
	signature := AggregateSignatures(signatures)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		VerifyAggregateCommon(msg, domain, pubkeys, signature)
	}
}

func toBytes(size int, hexStrs ...string) []byte {
	var out []byte
	if size > 0 {
		out = make([]byte, size*len(hexStrs))
	}
	for i := 0; i < len(hexStrs); i++ {
		hexStr := hexStrs[i]
		if hexStr[:2] == "0x" {
			hexStr = hexStr[2:]
		}
		if len(hexStr)%2 == 1 {
			hexStr = "0" + hexStr
		}
		bytes, err := hex.DecodeString(hexStr)
		if err != nil {
			panic(err)
		}
		if size <= 0 {
			out = append(out, bytes...)
		} else {
			if len(bytes) > size {
				panic(fmt.Sprintf("bad input string\ninput: %x\nsize: %d\nlenght: %d\n", bytes, size, len(bytes)))
			}
			offset := i*size + (size - len(bytes))
			copy(out[offset:], bytes)
		}
	}
	return out
}
