package mediawiki

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
)

type Parser struct {
	Title    string   `json:"title,omitempty"`
	PageId   int64    `json:"pageid,omitempty"`
	Images   []string `json:"images,omitempty"`
	WikiText string   `json:"wikitext,omitempty"`
}

type WikiResponse struct {
	Parser `json:"parse"`
}

const (
	apiURL           = "https://wiki.wizard101central.com/wiki/api.php"
	apiImageRedirect = "https://wiki.wizard101central.com/wiki/Special:FilePath/"
)

// WikiService provides methods for interacting with the Wizard101 Central Wiki.
type WikiService struct {
	HttpClient *HttpClient
}

// NewWikiService creates a new instance of WikiService.
func NewWikiService() *WikiService {
	return &WikiService{
		HttpClient: NewHttpClient(),
	}
}

func (s *WikiService) WikiText(pageID string) (WikiResponse, error) {
	url := fmt.Sprintf("%s?action=parse&page=%s&prop=wikitext|images&formatversion=2&format=json", apiURL, pageID)
	http, err := s.HttpClient.Get(url)
	if err != nil {
		return WikiResponse{}, err
	}
	defer http.Body.Close()
	api, err := s.wikiText(bodyReader(http.Body))
	if err != nil {
		return WikiResponse{}, err
	}
	return *api, nil
}

func bodyReader(body io.ReadCloser) []byte {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return []byte{}
	}
	return bodyBytes
}

// wikiText parses the response body into a WikiResponse struct.
func (s *WikiService) wikiText(body []byte) (w *WikiResponse, err error) {
	err = json.Unmarshal(body, &w)
	if err != nil {
		return &WikiResponse{}, err
	}
	return w, nil
}

// Json converts the infobox in the WikiText content to a JSON string.
func (s *WikiService) Json(pageID string) (bytes []byte, err error) {
	var wiki WikiResponse
	wiki, err = s.WikiText(pageID)
	if err != nil {
		return nil, err
	}
	header := FindHeader(wiki.Parser.WikiText)
	infobox := ReplaceInfoboxHeader(wiki.Parser.WikiText, header)
	data := extractKeyValuePairs(infobox)
	data["name"] = s.AppendName(data, pageID)(err)
	data["image"] = s.GetImageURLFunc(data, pageID)(err)
	return s.JsonMarshal(data)(err), nil
}

func (s *WikiService) AppendName(data map[string]string, pageID string) func(error) string {
	return func(err error) string {
		if err != nil {
			log.Printf("error getting image URL: %s", err)
		}
		re := regexp.MustCompile(`:\s*([^,]+)`)
		matches := re.FindStringSubmatch(pageID)
		data["name"] = matches[1]
		return matches[1]
	}
}

func (s *WikiService) GetImageURLFunc(data map[string]string, pageID string) func(error) string {
	return func(err error) string {
		if err != nil {
			log.Printf("error getting image URL: %s", err)
		}
		imageURL, err := s.GetImageURL(pageID)
		if err != nil {
			log.Printf("error getting image URL: %s", err)
		}
		data["image"] = imageURL
		return imageURL
	}
}

func (s *WikiService) JsonMarshal(data map[string]string) func(error) []byte {
	return func(err error) []byte {
		if err != nil {
			log.Printf("error marshaling JSON: %s", err)
		}
		bytes, err := json.Marshal(data)
		if err != nil {
			log.Printf("error marshaling JSON: %s", err)
		}
		return bytes
	}
}

// GetImageURL returns the URL of the specified image.
func (s *WikiService) GetImageURL(pageID string) (string, error) {
	var PageIdParser string
	switch {
	case strings.Contains(pageID, "Pet:"):
		PageIdParser = handlePet(pageID)

	case strings.Contains(pageID, "TreasureCard:"):
		PageIdParser = handleTreasureCard(pageID)
	default:
		PageIdParser = pageID
	}
	return s.getHeader(apiImageRedirect, PageIdParser), nil
}

func extractKeyValuePairs(data string) map[string]string {
	result := make(map[string]string)
	re := regexp.MustCompile(`\| (\w+)\s*=\s*([^|]+)`)
	indexes := re.FindAllStringSubmatchIndex(data, -1)
	for _, index := range indexes {
		key := data[index[2]:index[3]]
		value := data[index[4]:index[5]]
		result[key] = strings.TrimSpace(value)
	}
	return result
}

func (s *WikiService) getHeader(redirect, pageID string) string {
	res, err := s.HttpClient.Client.Head(redirect + pageID)
	if err != nil {
		log.Printf("error getting image URL: %s", err)
	}
	return res.Request.URL.String()
}

func handlePet(pageID string) string {
	re := regexp.MustCompile(`(.*):(.*)`)
	resd := re.ReplaceAllString(pageID, "($1)_$2.png")
	return resd
}

func handleTreasureCard(pageID string) string {
	re := regexp.MustCompile(`(Treasure)(Card):(.*)`)
	return re.ReplaceAllString(pageID, "(Treasure_Card)_$3.png")
}

// ReplaceInfoboxHeader removes the infobox header and footer from the WikiText content.
func ReplaceInfoboxHeader(data, template string) string {
	data = strings.TrimPrefix(data, fmt.Sprintf("{{%s\n", template))
	data = strings.TrimSuffix(data, "}}")
	data = strings.TrimSpace(data)
	return data
}

// FindHeader returns the infobox header from the WikiText content.
func FindHeader(data string) string {
	header := regexp.MustCompile(`{{(.+?)\n`)
	headerMatches := header.FindStringSubmatch(data)
	if len(headerMatches) != 2 {
		panic("invalid infobox")
	}
	return headerMatches[1]
}
