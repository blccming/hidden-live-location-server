package db

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	glide "github.com/valkey-io/valkey-glide/go/v2"
	"github.com/valkey-io/valkey-glide/go/v2/config"
)

// sessionKey returns the Valkey key used to store session data for the token.
func sessionKey(token string) string {
	return "session:" + token
}

// locationKey returns the Valkey key used to store location data for the token.
func locationKey(token string) string {
	return "session:" + token + ":loc"
}

// TESTING ONLY
func Test(c *glide.Client) {
	token := "AECD16"
	var exists bool
	exists, _ = SessionExists(c, token)
	fmt.Println(exists)
	AddSession(c, token, 3600, 30)
	exists, _ = SessionExists(c, token)
	fmt.Println(exists)
	loc, _ := GetLocation(c, token)
	fmt.Println(loc)
	SetLocation(c, token, 123, 456)
	loc, _ = GetLocation(c, token)
	fmt.Println(loc)
	RemoveSession(c, token)
}

// Connect initializes a Valkey client, verifies connectivity with PING,
// applies basic configuration, and returns the connected client.
func Connect() (*glide.Client, error) {
	// TODO: make these configurable via environment variable
	host := "localhost"
	port := 6379
	password := "YOUR_PASSWORD_HERE"

	config := config.NewClientConfiguration().
		WithAddress(&config.NodeAddress{Host: host, Port: port}).
		WithCredentials(config.NewServerCredentials("", password)).
		WithRequestTimeout(5 * time.Second)

	client, err := glide.NewClient(config)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Failed to initialize database client.")
		return nil, err
	}

	if _, err = client.Ping(context.Background()); err != nil {
		log.Error().Stack().Err(err).Msg("Failed to ping database.")
		return nil, err
	}
	log.Info().Msgf("Connected to database! Server responded to ping.")

	if err = configureDatabase(client); err != nil {
		log.Error().Stack().Err(err).Msg("Failed to configure database.")
	}
	log.Info().Msg("Database configured successfully.")

	return client, nil
}

// configureDatabase sets basic Valkey configuration options such as maxmemory
// and disables persistence.
func configureDatabase(c *glide.Client) error {
	maxmem := "4GB" // TODO: make this configurable via environment variable
	_, err := c.ConfigSet(context.Background(), map[string]string{
		"maxmemory":  maxmem,
		"save":       "",
		"appendonly": "no",
	})
	log.Debug().Msgf("Database configured with maxmemory %s, persistence disabled.", maxmem)
	return err
}

// AddSession creates a new session for the token, storing the location TTL
// and setting the session key to expire after sessionTimeout seconds.
func AddSession(c *glide.Client, token string, sessionTimeout int, locationTTL int) error {
	// Store the location TTL in this session dataset, using "session:token" as the key.
	// We keep session data and location data in separate datasets so we can use
	// Valkey's expiration feature independently for each of them.
	// This lets the session timeout and the location TTL expire on their own schedules.

	if sessionExists, _ := SessionExists(c, token); sessionExists {
		log.Warn().Msgf("Skipping session creation: Session already exists for token %s.", token)
		return fmt.Errorf("session already exists")
	}

	if sessionTimeout < 1 {
		return fmt.Errorf("session timeout must be greater than 0")
	}
	if locationTTL < 1 {
		return fmt.Errorf("location TTL must be greater than 0")
	}

	data := strconv.Itoa(locationTTL)
	if _, err := c.Set(context.Background(), sessionKey(token), data); err != nil {
		log.Panic().Stack().Err(err).Msg("Failed to add session.")
		return err
	}

	if didSetExp, err := c.Expire(context.Background(), sessionKey(token), time.Duration(sessionTimeout)*time.Second); err != nil {
		log.Panic().Stack().Err(err).Msg("Failed to set session expiration.")
		return err
	} else if !didSetExp { // just doing the else to be still in scope of didSetExp
		log.Panic().Msg("Failed to set session expiration.")
		return fmt.Errorf("failed to set session expiration")
	}

	log.Debug().Msgf("Session created for token %s with timeout %ds and location TTL %ds.", token, sessionTimeout, locationTTL)
	return nil
}

// SessionExists returns true if a session for the given token exists in Valkey.
func SessionExists(c *glide.Client, token string) (bool, error) {
	exists, err := c.Exists(context.Background(), []string{sessionKey(token)})
	if err != nil {
		return false, err
	}
	if exists > 1 {
		log.Fatal().Msgf("Unintended behavior: Multiple sessions found for token %s.", token)
		return true, fmt.Errorf("multiple sessions found for token %s", token)
	}

	if exists == 0 {
		log.Debug().Msgf("Session does not exist for token %s.", token)
		return false, nil
	}
	log.Debug().Msgf("Session exists for token %s.", token)
	return true, nil
}

// getLocationTTL retrieves and parses the locationTTL for the given session token from Valkey.
func getLocationTTL(c *glide.Client, token string) (int, error) {
	locationTTLStr, err := c.Get(context.Background(), sessionKey(token))
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Failed to get locationTTL for token %s.", token)
		return -1, err
	}

	locationTTL, err := strconv.Atoi(locationTTLStr.Value())
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Failed to parse locationTTL for token %s.", token)
		return -1, err
	}

	log.Debug().Msgf("Location TTL for token %s is %ds.", token, locationTTL)
	return locationTTL, nil
}

// SetLocation stores the given longitude and latitude for the session token in Valkey
// and sets the key’s expiration based on the current location TTL stored for the session in Valkey.
func SetLocation(c *glide.Client, token string, longitude float64, latitude float64) error {
	// First, get locationTTL
	locationTTL, err := getLocationTTL(c, token)
	if err != nil {
		return err
	}

	// Assume locationTTL is valid (condition checked when creating session),
	// now proceed with setting location in Valkey
	data := map[string]string{
		"longitude":   strconv.FormatFloat(longitude, 'f', -1, 64),
		"latitude":    strconv.FormatFloat(latitude, 'f', -1, 64),
		"lastChanged": time.Now().Format(time.RFC3339),
	}

	if _, err := c.HSet(context.Background(), locationKey(token), data); err != nil {
		log.Error().Stack().Err(err).Msgf("Failed to set location for token %s.", token)
		return err
	}
	if _, err := c.Expire(context.Background(), locationKey(token), time.Duration(locationTTL)*time.Second); err != nil {
		log.Error().Stack().Err(err).Msgf("Failed to set location expiration for token %s.", token)
		return err
	}

	log.Debug().Msgf("Location set for token %s.", token)
	return nil
}

// getLocation retrieves the location for a given token from the database.
func GetLocation(c *glide.Client, token string) (map[string]string, error) {
	result, err := c.HGetAll(context.Background(), locationKey(token))

	if err != nil {
		log.Error().Stack().Err(err).Msgf("Failed to get location for token %s.", token)
		return nil, err
	}
	if len(result) == 0 {
		log.Error().Msgf("Location not found for token %s.", token)
		return nil, fmt.Errorf("location not found for token %s", token)
	}
	if result["longitude"] == "" || result["latitude"] == "" || result["lastChanged"] == "" {
		log.Error().Msgf("Location data incomplete for token %s.", token)
		return nil, fmt.Errorf("location data incomplete for token %s", token)
	}

	log.Debug().Msgf("Location retrieved for token %s.", token)
	return result, nil
}

func RemoveSession(c *glide.Client, token string) error {
	keys := []string{sessionKey(token), locationKey(token)}

	if _, err := c.Del(context.Background(), keys); err != nil {
		log.Error().Stack().Err(err).Msgf("Failed to remove session for token %s.", token)
		return err
	}

	log.Debug().Msgf("Session and location removed for token %s.", token)
	return nil
}
