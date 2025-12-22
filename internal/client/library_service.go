package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mangahub/pkg/models/dtos"
	"mangahub/pkg/utils/colors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

func (c *Client) AddMangaOrUpdateToLibrary(mode, token, mangaId string, currentChapter int) (string, error) {
	mangaData := map[string]interface{}{
		"manga_id":        mangaId,
		"current_chapter": currentChapter,
	}
	jsonData, _ := json.Marshal(mangaData)

	var request *http.Request
	var err error

	endpoint := fmt.Sprintf("%s/users/library", c.Config.ServerURL)

	switch mode {
	case "add":
		request, err = http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	case "update":
		request, err = http.NewRequest("PUT", endpoint, bytes.NewBuffer(jsonData))
	default:
		return "", errors.New("invalid mode")
	}

	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+strings.TrimSpace(token))

	resp, err := c.HttpClient.Do(request)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	var result dtos.UserAddOrUpdateLibrary
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Message, nil
}

func (c *Client) GetReadingList(token, status string, limit, page int) (any, error) {
	// Build Query Parameters
	v := url.Values{}
	if status != "" {
		v.Set("status", status)
		v.Set("limit", strconv.Itoa(limit))
		v.Set("page", strconv.Itoa(page))
	}

	// Construct URL
	reqUrl := fmt.Sprintf("%s/users/library", c.Config.ServerURL)
	if len(v) > 0 {
		reqUrl += "?" + v.Encode()
	}

	request, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+strings.TrimSpace(token))

	// Use c.HttpClient
	resp, err := c.HttpClient.Do(request)
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

	if status != "" {
		var result dtos.PaginatedResponse[dtos.UserLibraryItem]
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}
		return &result, nil
	}

	var result dtos.UserLibrary
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func PrintReadingList(data *dtos.UserLibrary) {
	fmt.Print("\033[H\033[2J")

	reading := data.ReadingList.Reading
	ptr := data.ReadingList.PlanToRead
	completed := data.ReadingList.Completed

	fmt.Printf("%sPlan to read:%s\n", colors.ColorYellow, colors.ColorReset)
	printList(ptr)
	fmt.Printf("%sReading:%s\n", colors.ColorGreen, colors.ColorReset)
	printList(reading)
	fmt.Printf("%sCompleted:%s\n", colors.ColorPurple, colors.ColorReset)
	printList(completed)
}

func printList(items []dtos.UserLibraryItem) {
	for _, i := range items {
		fmt.Printf(
			"[ %s%s %s(%s)%s | ch.%d | %s | %s ]\n",
			colors.ColorRed,
			i.Title,
			colors.ColorBlue,
			i.MangaID,
			colors.ColorReset,
			i.CurrentChapter,
			i.Status,
			i.LastUpdated.Format(time.DateTime),
		)
	}
	fmt.Printf("\n\n")
}

func PrintReadingListByStatus(data *dtos.PaginatedResponse[dtos.UserLibraryItem]) {
	meta := data.Pagination

	fmt.Print("\033[H\033[2J")

	if meta.Total == 0 {
		fmt.Println("No readings found matching your criteria.")
		return
	}

	fmt.Printf("\nFound %d results (Page %d of %d)\n", meta.Total, meta.Page, meta.TotalPages)
	fmt.Println("-------------------------------------------------------------------------------")

	// Initialize TabWriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)

	// Print Header
	_, err := fmt.Fprint(w, "ID\tTITLE\tSTATUS\tCHAPTERS\tLAST UPDATED\n")
	if err != nil {
		return
	}

	// Print Rows
	for _, m := range data.Results {
		title := m.Title
		if len(title) > 35 {
			title = title[:32] + "..."
		}

		_, err := fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
			m.MangaID,
			title,
			m.Status,
			m.CurrentChapter,
			m.LastUpdated.Format(time.DateTime),
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
