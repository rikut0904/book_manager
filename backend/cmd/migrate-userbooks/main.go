package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"book_manager/backend/internal/db"
	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/idgen"
	"book_manager/backend/internal/repository"
	"book_manager/backend/internal/repository/gormrepo"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")

	userID := flag.String("user-id", "", "internal user ID to migrate (omit to list all users)")
	cleanup := flag.Bool("cleanup", false, "delete all data from all tables (users are preserved)")
	cleanupAll := flag.Bool("cleanup-all", false, "delete all data from all tables including users")
	flag.Parse()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	dbConn, err := db.Open(databaseURL)
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}

	if *cleanupAll {
		cleanupAllTables(dbConn)
		return
	}

	if *cleanup {
		cleanupDataTables(dbConn)
		return
	}

	userRepo := gormrepo.NewUserRepository(dbConn)
	bookRepo := gormrepo.NewBookRepository(dbConn)
	userBookRepo := gormrepo.NewUserBookRepository(dbConn)

	if *userID == "" {
		listUsers(userRepo)
		return
	}

	migrateUserBooks(*userID, bookRepo, userBookRepo)
}

func listUsers(userRepo repository.UserRepository) {
	users := userRepo.List()
	if len(users) == 0 {
		fmt.Println("No users found.")
		return
	}
	fmt.Println("Users:")
	for _, u := range users {
		fmt.Printf("  ID: %s  UserID: %s  Email: %s  DisplayName: %s\n", u.ID, u.UserID, u.Email, u.DisplayName)
	}
	fmt.Println("\nRun with --user-id=<ID> to migrate user_books for a specific user.")
}

func migrateUserBooks(userID string, bookRepo repository.BookRepository, userBookRepo repository.UserBookRepository) {
	allBooks := bookRepo.List()
	userBooks := userBookRepo.ListByUser(userID)

	ownedBookIDs := make(map[string]struct{}, len(userBooks))
	for _, ub := range userBooks {
		ownedBookIDs[ub.BookID] = struct{}{}
	}

	var created, skipped int
	for _, book := range allBooks {
		if _, exists := ownedBookIDs[book.ID]; exists {
			continue
		}

		ub := domain.UserBook{
			ID:     idgen.NewUserBook(),
			UserID: userID,
			BookID: book.ID,
		}
		if err := userBookRepo.Create(ub); err != nil {
			if errors.Is(err, repository.ErrUserBookExists) {
				skipped++
				continue
			}
			log.Printf("error creating user_book for book %s: %v", book.ID, err)
			skipped++
			continue
		}
		created++
	}

	fmt.Printf("Migration complete for user %s\n", userID)
	fmt.Printf("  Total books: %d\n", len(allBooks))
	fmt.Printf("  Already owned: %d\n", len(userBooks))
	fmt.Printf("  Created: %d\n", created)
	fmt.Printf("  Skipped: %d\n", skipped)
}

// cleanupAllTables deletes all data from all tables including users.
func cleanupAllTables(dbConn *gorm.DB) {
	tables := []string{
		"recommendations",
		"favorites",
		"next_to_buy_manuals",
		"user_books",
		"books",
		"series",
		"audit_logs",
		"admin_invitations",
		"admin_users",
		"profile_settings",
		"open_ai_keys",
		"isbn_caches",
		"users",
	}
	deleteFromTables(dbConn, tables)
}

// cleanupDataTables deletes all data except users and profile settings.
func cleanupDataTables(dbConn *gorm.DB) {
	tables := []string{
		"recommendations",
		"favorites",
		"next_to_buy_manuals",
		"user_books",
		"books",
		"series",
		"audit_logs",
		"isbn_caches",
	}
	deleteFromTables(dbConn, tables)
}

func deleteFromTables(dbConn *gorm.DB, tables []string) {
	for _, table := range tables {
		result := dbConn.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if result.Error != nil {
			log.Printf("  %s: error: %v", table, result.Error)
		} else {
			fmt.Printf("  %s: deleted %d rows\n", table, result.RowsAffected)
		}
	}
	fmt.Println("Cleanup complete.")
}
