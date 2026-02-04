// Admin CLI tool for bootstrapping users and IP mappings
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/23St/trackstack/internal/common/db"
)

func main() {
	var (
		dataSource = flag.String("db", "./data/trackstack.db", "Database file path")
		ipAddress  = flag.String("ip", "", "IP address to whitelist (required)")
		desc       = flag.String("desc", "", "Description for this device (optional)")
		userID     = flag.String("user", "", "Existing user ID to map IP to (optional, creates new user if not provided)")
	)

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		_, _ = fmt.Fprintf(os.Stderr, "Bootstrap a new user with IP whitelist access.\n\n")
		_, _ = fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		_, _ = fmt.Fprintf(os.Stderr, "\nExamples:\n")
		_, _ = fmt.Fprintf(os.Stderr, "  # Create new user with IP:\n")
		_, _ = fmt.Fprintf(os.Stderr, "  %s -ip=192.168.1.10 -desc=\"Dad's desktop\"\n\n", os.Args[0])
		_, _ = fmt.Fprintf(os.Stderr, "  # Map IP to existing user:\n")
		_, _ = fmt.Fprintf(os.Stderr, "  %s -ip=192.168.1.11 -user=abc-123 -desc=\"Mom's phone\"\n\n", os.Args[0])
	}

	flag.Parse()

	// Validate required flags
	if *ipAddress == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Error: -ip is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Normalize IP (handle IPv6 format)
	*ipAddress = normalizeIP(*ipAddress)

	// Connect to database
	database, err := db.New(*dataSource)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func(database *db.DB) {
		_ = database.Close()
	}(database)

	ctx := context.Background()

	// Start transaction
	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	var targetUser *db.User

	// Create or use existing user
	if *userID != "" {
		// Use existing user
		targetUser, err = database.GetUserByID(ctx, *userID)
		if err != nil {
			log.Fatalf("Failed to find user %s: %v", *userID, err)
		}
		fmt.Printf("Using existing user: %s\n", targetUser.ID)
	} else {
		// Create new user
		targetUser = db.NewUser()
		if err := database.CreateUser(ctx, targetUser); err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
		fmt.Printf("Created new user: %s\n", targetUser.ID)
	}

	// Check if IP already mapped
	exists, err := database.IPMappingExists(ctx, *ipAddress)
	if err != nil {
		log.Fatalf("Failed to check IP mapping: %v", err)
	}
	if exists {
		log.Fatalf("IP %s is already mapped to a user. Use a different IP or delete the existing mapping first.", *ipAddress)
	}

	// Create IP mapping
	mapping := db.NewIPMapping(*ipAddress, targetUser.ID, *desc)
	if err := database.CreateIPMapping(ctx, mapping); err != nil {
		log.Fatalf("Failed to create IP mapping: %v", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Printf("\n✓ Successfully configured access:\n")
	fmt.Printf("  User ID:     %s\n", targetUser.ID)
	fmt.Printf("  IP Address:  %s\n", *ipAddress)
	if *desc != "" {
		fmt.Printf("  Description: %s\n", *desc)
	}
	fmt.Printf("\nUser can now access the app from this IP address.\n")
	fmt.Printf("Session cookies will be valid indefinitely.\n\n")

	// Print environment variable info
	fmt.Printf("To add more devices to this user, run:\n")
	fmt.Printf("  go run cmd/admin/create_user.go -ip=<NEW_IP> -user=%s -desc=\"Description\"\n\n", targetUser.ID)
}

// normalizeIP handles IPv6 localhost and other formats
func normalizeIP(ip string) string {
	ip = strings.TrimSpace(ip)

	// Handle IPv6 localhost
	if ip == "::1" || ip == "[::1]" {
		return "127.0.0.1"
	}

	// Remove brackets from IPv6 addresses
	ip = strings.TrimPrefix(ip, "[")
	ip = strings.TrimSuffix(ip, "]")

	return ip
}
