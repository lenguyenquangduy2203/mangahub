package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mangahub/pkg/models/dtos"
	"mangahub/pkg/utils"
	"mangahub/pkg/utils/colors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

func GetMangaById(mangaId string) (*dtos.MangaDetail, error) {
	resp, err := http.Get(utils.BaseURL + "/manga/" + mangaId)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	var result dtos.MangaDetail
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func SearchManga(title, author, genre, status string, limit, page int) (*dtos.PaginatedResponse[dtos.MangaListItem], error) {
	query := "?limit=" + strconv.Itoa(limit) + "&page=" + strconv.Itoa(page)

	if title != "" {
		query += "&title=" + url.QueryEscape(title)
	}

	if author != "" {
		query += "&author=" + url.QueryEscape(author)
	}

	if genre != "" {
		query += "&genre=" + url.QueryEscape(genre)
	}

	if status != "" {
		query += "&status=" + url.QueryEscape(status)
	}

	resp, err := http.Get(utils.BaseURL + "/manga" + query)

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	var result dtos.PaginatedResponse[dtos.MangaListItem]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func PrintMangaDetail(m *dtos.MangaDetail) {
	fmt.Print("\033[H\033[2J")
	fmt.Println("--------------------------------------------------")
	fmt.Printf("%sTitle:%s    %s\n", colors.ColorBlue, colors.ColorReset, m.Title)
	fmt.Printf("%sID:%s       %s\n", colors.ColorBlue, colors.ColorReset, m.MangaID)
	fmt.Printf("%sAuthor:%s   %s\n", colors.ColorBlue, colors.ColorReset, m.Author)
	fmt.Printf("%sStatus:%s   %s\n", colors.ColorBlue, colors.ColorReset, m.Status)
	fmt.Printf("%sChapters:%s %d\n", colors.ColorBlue, colors.ColorReset, m.TotalChapters)
	fmt.Printf("%sGenres:%s   [%s]\n",
		colors.ColorBlue,
		colors.ColorReset,
		strings.Join(m.Genres, ", "),
	)
	fmt.Println("--------------------------------------------------")
	fmt.Printf("%sDescription:%s\n%s\n", colors.ColorBlue, colors.ColorReset, m.Description)
	fmt.Println("--------------------------------------------------")
}

func PrintMangaList(data *dtos.PaginatedResponse[dtos.MangaListItem]) {
	meta := data.Pagination

	fmt.Print("\033[H\033[2J")

	if meta.Total == 0 {
		fmt.Println("No manga found matching your criteria.")
		return
	}

	fmt.Printf("\nFound %d results (Page %d of %d)\n", meta.Total, meta.Page, meta.TotalPages)
	fmt.Println("-------------------------------------------------------------------------------")

	// Initialize TabWriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)

	// Print Header
	_, err := fmt.Fprint(w, "ID\tTITLE\tSTATUS\tCHAPTERS\n")
	if err != nil {
		return
	}

	// Print Rows
	for _, m := range data.Results {
		title := m.Title
		if len(title) > 35 {
			title = title[:32] + "..."
		}

		_, err := fmt.Fprintf(w, "%s\t%s\t%s\t%d\n",
			m.MangaID,
			title,
			m.Status,
			m.TotalChapters,
		)
		if err != nil {
			return
		}
	}

	// Flush buffer to output
	err = w.Flush()
	if err != nil {
		return
	}
	fmt.Println("-------------------------------------------------------------------------------")

	// UX Hint for next page
	if meta.HasNext {
		fmt.Printf("Next Page: Add flag -page(-pg)=%d\n", meta.Page+1)
	} else {
		fmt.Println("\U000F0601 End of results.")
	}
	fmt.Println()
}
