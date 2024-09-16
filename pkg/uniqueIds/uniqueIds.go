// pkg/uniqueids/uniqueids.go

package uniqueids

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func init() {

	// Get Redis address from environment variable
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379" // Default to Docker service name
	}

	// Initialize Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	// Test Redis connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Start the routine to reset counts every minute
	go resetCountEveryMinute()
}

// AddID adds an ID to the Redis set for the current minute
func AddID(id int) error {
	// Use current minute as key
	key := time.Now().Format("200601021504") // Format: YYYYMMDDHHMM

	// Add ID to Redis set
	err := rdb.SAdd(ctx, key, id).Err()
	if err != nil {
		log.Printf("Error adding ID to Redis: %v", err)
		return err
	}

	// Set expiration to clean up old keys
	err = rdb.Expire(ctx, key, 2*time.Minute).Err()
	if err != nil {
		log.Printf("Error setting expiration on Redis key: %v", err)
		return err
	}

	return nil
}

// GetCurrentCount returns the count of unique IDs for the current minute
func GetCurrentCount() int {
	key := time.Now().Format("200601021504")

	count, err := rdb.SCard(ctx, key).Result()
	if err != nil {
		log.Printf("Error getting unique IDs count from Redis: %v", err)
		return 0
	}

	return int(count)
}

func resetCountEveryMinute() {
	for {
		// Sleep until the start of the next minute
		sleepUntilNextMinute()

		// Get key for the previous minute
		key := getPreviousMinuteKey()

		// Get the count of unique IDs from the previous minute
		count, err := rdb.SCard(ctx, key).Result()
		if err != nil {
			log.Printf("Error getting unique IDs count from Redis: %v", err)
			continue
		}

		// Log the count
		log.Printf("Unique IDs received in the last minute: %d", count)

		// Send the count to a distributed streaming service
		sendCountToStreamingService(int(count))

		// Delete the key to clean up
		err = rdb.Del(ctx, key).Err()
		if err != nil {
			log.Printf("Error deleting Redis key: %v", err)
		}
	}
}

// sleepUntilNextMinute sleeps until the start of the next minute
func sleepUntilNextMinute() {
	now := time.Now()
	nextMinute := now.Truncate(time.Minute).Add(time.Minute)
	time.Sleep(time.Until(nextMinute))
}

// getPreviousMinuteKey returns the Redis key for the previous minute
func getPreviousMinuteKey() string {
	previousMinute := time.Now().Add(-time.Minute)
	return previousMinute.Format("200601021504")
}

// sendCountToStreamingService sends the count to a streaming service
func sendCountToStreamingService(count int) {
	// Kafka configuration
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	// Get Kafka brokers from environment variable
	brokersEnv := os.Getenv("KAFKA_BROKERS")
	var brokers []string
	if brokersEnv == "" {
		brokers = []string{"kafka:9092"} // Default to Docker service name
	} else {
		brokers = []string{brokersEnv}
	}

	// Create a new sync producer
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Printf("Error creating Kafka producer: %v", err)
		return
	}
	defer producer.Close()

	// Prepare the message
	msg := &sarama.ProducerMessage{
		Topic: "unique_id_counts", // Kafka topic
		Value: sarama.StringEncoder(strconv.Itoa(count)),
	}

	// Send the message
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("Error sending message to Kafka: %v", err)
		return
	}

	log.Printf("Message sent to Kafka partition %d at offset %d", partition, offset)
}
