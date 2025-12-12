package main

import (
	"flag"
	"fmt"
	"log"
	"mangahub/internal/client"
	"mangahub/pkg/models/dtos"
	"mangahub/pkg/pagination"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var aliasMap = map[string]string{
	"-h":   "-help",
	"-m":   "-mode",
	"-u":   "-user",
	"-p":   "-pass",
	"-r":   "-room",
	"-t":   "-title",
	"-a":   "-author",
	"-g":   "-genre",
	"-id":  "-mangaId",
	"-s":   "-status",
	"-l":   "-lib",
	"-c":   "-chapter",
	"-pg":  "-page",
	"-lim": "-limit",
}

func normalizeArgs() {
	for i, arg := range os.Args {
		if rl, ok := aliasMap[arg]; ok {
			os.Args[i] = rl
		}
		for alias, rl := range aliasMap {
			if strings.HasPrefix(arg, alias+"=") {
				os.Args[i] = rl + "=" + strings.TrimPrefix(arg, alias+"=")
			}
		}
	}
}

func main() {
	normalizeArgs()

	// -----------------------------
	// Help
	// -----------------------------
	help := flag.Bool("help", false, "Show help (alias: -h)")

	// -----------------------------
	// App Mode
	// -----------------------------
	mode := flag.String("mode", "search", "Mode: register | lib | chat | search (alias: -m)")

	// -----------------------------
	// Auth Flags
	// -----------------------------
	username := flag.String("user", "", "Username (alias: -u)")
	password := flag.String("pass", "", "Password (alias: -p)")

	// -----------------------------
	// Chat Flags
	// -----------------------------
	room := flag.String("room", "general", "Chat room (alias: -r)")

	// -----------------------------
	// Manga Search Flags
	// -----------------------------
	searchTitle := flag.String("title", "", "Search manga by title (alias: -t)")
	searchAuthor := flag.String("author", "", "Search manga by author (alias: -a)")
	searchGenre := flag.String("genre", "", "Search manga by genre (alias: -g)")

	// -----------------------------
	// Shared Manga/Library Flags
	// -----------------------------
	mangaId := flag.String("mangaId", "", "Manga ID (alias: -id)")
	searchStatus := flag.String("status", "", "Status filter (alias: -s)")

	// -----------------------------
	// Library Actions
	// -----------------------------
	libAction := flag.String("lib", "get", "Library action: add | update | get (alias: -l)")
	currentChapter := flag.Int("chapter", 0, "Current chapter (alias: -c)")

	// -----------------------------
	// Pagination
	// -----------------------------
	page := flag.Int("page", 1, "Page number (alias: -pg)")
	limit := flag.Int("limit", pagination.DEFAULT_LIMIT,
		fmt.Sprintf("Max results per page (max %d) (alias: -lim)", pagination.MAX_LIMIT))

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	// Switch Modes
	switch *mode {
	case "register":
		if *username == "" || *password == "" {
			log.Fatal("Error: -user and -pass are required")
		}

		_, err := client.Register(*username, *password)
		if err != nil {
			log.Fatalf("Register failed: %v", err)
		}
		fmt.Println("Register success")

	case "search":
		if *mangaId == "" {
			fmt.Println("--- Manga Search Mode ---")

			if *searchTitle == "" && *searchAuthor == "" && *searchGenre == "" && *searchStatus == "" {
				fmt.Println("Warning: No search filters provided. Listing recent manga...")
			}

			fmt.Println("--- Searching ---", *searchTitle)
			mangas, err := client.SearchManga(*searchTitle, *searchAuthor, *searchGenre, *searchStatus, *limit, *page)

			if err != nil {
				log.Fatalf("Search failed: %v", err)
			}

			client.PrintMangaList(mangas)

		} else {
			fmt.Printf("--- Fetching Details for Manga ID: %s ---\n", *mangaId)

			manga, err := client.GetMangaById(*mangaId)
			if err != nil {
				log.Fatalf("Fetch Details failed: %v", err)
			}

			client.PrintMangaDetail(manga)
		}

	case "lib":
		token := getAuth(username, password)

		if *libAction == "add" || *libAction == "update" {
			caser := cases.Title(language.English)
			fmt.Printf("--- %s manga in user library ---\n", caser.String(*libAction))

			msg, err := client.AddMangaOrUpdateToLibrary(*libAction, token, *mangaId, *currentChapter)

			if err != nil {
				log.Fatalf("%s failed: %v", caser.String(*libAction), err)
			}

			fmt.Println(msg)

		} else if *libAction == "get" {
			fmt.Println("--- Get reading lists in user library ---")
			list, err := client.GetReadingList(token, *searchStatus, *limit, *page)

			if err != nil {
				log.Fatalf("Get Reading lists failed: %v", err)
			}

			if data, ok := list.(*dtos.PaginatedResponse[dtos.UserLibraryItem]); ok {
				client.PrintReadingListByStatus(data)
			}

			if data, ok := list.(*dtos.UserLibrary); ok {
				client.PrintReadingList(data)
			}

		} else {
			log.Fatal("Error: Please use correct mode: 'add', 'update', 'get'")
		}

	case "chat":
		token := getAuth(username, password)
		client.StartChat(token, *room, *username)

	default:
		log.Fatalf("Unknown mode: %s.", *mode)
	}
}

func getAuth(username, password *string) string {
	// Validation
	if *username == "" || *password == "" {
		log.Fatal("Error: -user and -pass are required")
	}

	fmt.Println("\uF09C Authenticating...")
	token, err := client.Login(*username, *password)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	fmt.Println("Login successful!")

	return token
}
