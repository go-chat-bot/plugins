package dedup

import (
	"github.com/go-chat-bot/bot"
	"hash/fnv"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	dedupConfigEnv   = "DEDUP_TIMEOUT"
	defaultDedupTime = "5"
)

var (
	pastMessageMap     map[uint32]time.Time
	pastMessageMapLock = sync.RWMutex{}
	dedupConfig        time.Duration
	fnvHash            = fnv.New32a()
)

func messageHash(msg, target string) uint32 {
	fnvHash.Write([]byte(msg))
	fnvHash.Write([]byte(target))
	defer fnvHash.Reset()
	return fnvHash.Sum32()
}

func recordMessage(msgHash uint32, target string) {
	until := time.Now().UTC().Add(dedupConfig)
	pastMessageMapLock.Lock()
	pastMessageMap[msgHash] = until
	pastMessageMapLock.Unlock()
	go func() {
		time.Sleep(dedupConfig)
		pastMessageMapLock.Lock()
		delete(pastMessageMap, msgHash)
		pastMessageMapLock.Unlock()
	}()
}

func dedupFilter(cmd *bot.FilterCmd) (string, error) {
	msgHash := messageHash(cmd.Message, cmd.Target)
	pastMessageMapLock.RLock()
	_, found := pastMessageMap[msgHash]
	pastMessageMapLock.RUnlock()
	if !found {
		// No past message like this, record and send
		recordMessage(msgHash, cmd.Target)
		return cmd.Message, nil
	}

	// Past message found, filter out!
	log.Printf("Deduplicating message in %s\n", cmd.Target)
	return "", nil
}

func init() {
	pastMessageMap = make(map[uint32]time.Time)
	dedupVar := os.Getenv(dedupConfigEnv)
	if dedupVar == "" {
		dedupVar = defaultDedupTime
	}
	min, err := strconv.Atoi(dedupVar)
	if err != nil {
		log.Printf("Failed to load dedup configuration. Falling back to	default")
		min, _ = strconv.Atoi(defaultDedupTime)
	}
	dedupConfig = time.Duration(min) * time.Minute

	bot.RegisterFilterCommand(
		"dedup",
		dedupFilter)
}
