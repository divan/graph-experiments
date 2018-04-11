package whisperv6

/*
import (
	"math/rand"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/whisper/whisperv6"
)

// const from github.com/ethereum/go-ethereum/whisper/whisperv5/doc.go
const (
	aesKeyLength = 32
)

func generateMessage(ttl int) *whisperv6.Envelope {
	// set all the parameters except p.Dst and p.Padding

	buf := make([]byte, 4)
	rand.Read(buf)
	sz := rand.Intn(400)

	params := &whisperv6.MessageParams{
		PoW:      0.01,
		WorkTime: 1,
		Payload:  make([]byte, sz),
		KeySym:   make([]byte, aesKeyLength),
		Topic:    whisperv6.BytesToTopic(buf),
		TTL:      uint32(ttl),
	}
	rand.Read(params.Payload)
	rand.Read(params.KeySym)

	var err error
	params.Src, err = crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	msg, err := whisperv6.NewSentMessage(params)
	if err != nil {
		panic(err)
	}
	env, err := msg.Wrap(params)
	if err != nil {
		panic(err)
	}

	return env
}
*/
