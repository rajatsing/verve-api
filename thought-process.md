Source: google/stackoverflow

Design Considerations:
1. The use of environment variables for configuration allows flexibility and portability across different environments.
2. Seprating logics in different filenames and packages for easy ubderstanding
3. using docker compose to run the application in a containerized environment

`main.go`

1. Initializes the HTTP server and sets up the routing for the API endpoint.
3. Registers the /api/verve/accept endpoint to the AcceptHandler method in the handlers package.
   
`accept.go` in `handlers` package

1. Extracts and validates the mandatory id query parameter, returning "failed" if missing or invalid.
2. Checks for an optional endpoint parameter to determine if a POST request should be made.
3. calls `uniqueids.AddID(id)` after validating the id.
   


`external.go` in `external` package
1. Creates a CountData struct to define the JSON payload for POST requests.
2. Uses json.Marshal to serialize the count data, handling any potential errors.
3. makes a POST request with the JSON payload to the provided endpoint, including appropriate headers.
4. Sets a timeout for the HTTP client to avoid hanging requests.
5. Logs the HTTP status code 

`uniqueids.go` in `uniqueids` package
1.  Chooses Redis to store unique IDs for deduplication across multiple instances, ensuring high performance.
2.  Uses the current minute formatted as YYYYMMDDHHMM as the Redis key to group IDs per minute.
3.  Sets a 2-minute expiration on Redis keys to automatically clean up old data and manage memory usage.
4.  Implements `resetCountEveryMinute` as a goroutine to handle counting and resetting IDs every minute.
5.   Uses the Sarama library to send the unique ID counts to a Kafka topic, handling any connection issues.
