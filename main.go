package main

import (
	"os"
	"path"
	"time"

	"github.com/tkanos/gonfig"

	bolt "go.etcd.io/bbolt"
)

// Configuration discord bot storage
type Configuration struct {
	MinShardCount     int
	DiscordToken      string
	SendAnnouncements bool
}

var (
	// Bot the bot
	Bot    *Discord
	config = Configuration{}

	latestKeys []string
)

func init() {
	wd, _ := os.Getwd()
	err := gonfig.GetConf(path.Join(wd, "conf.json"), &config)
	if err != nil {
		panic(err)
	}

	db, err := bolt.Open("chans.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	db.Update(func(txn *bolt.Tx) error {
		_, _ = txn.CreateBucketIfNotExists([]byte("ChanBucket"))
		return nil
	})
	defer db.Close()
}

// AddToDB add channel id to database
func AddToDB(guild string, channel string) {
	db, err := bolt.Open("chans.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		AddToDB(guild, channel)
	}
	defer db.Close()

	err = db.Update(func(txn *bolt.Tx) error {
		_, err := txn.CreateBucketIfNotExists([]byte("ChanBucket"))
		if err != nil {
			return err
		}
		b := txn.Bucket([]byte("ChanBucket"))
		err = b.Put([]byte(guild), []byte(channel))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		AddToDB(guild, channel)
	}
}

// RetrieveDB retrieve channel id for a guild
func RetrieveDB(guild string) string {
	db, err := bolt.Open("chans.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return RetrieveDB(guild)
	}
	defer db.Close()

	var val []byte
	err = db.View(func(txn *bolt.Tx) error {
		b := txn.Bucket([]byte("ChanBucket"))
		val = b.Get([]byte(guild))

		return nil
	})
	if err != nil {
		return RetrieveDB(guild)
	}

	return string(val)
}

func main() {
	// START UP NOTIFIER \\
	go StartNotifier()

	// START UP DISCORD BOT \\
	Bot = StartDiscord()

	(<-make(chan struct{}))
}
