package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
	Videos      []*PreviewVideo `json:"videos,omitempty"`
}

//PreviewVideo represents summary properties for a web page
type PreviewVideo struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
}

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	/*TODO: add code and additional functions to do the following:
	- Add an HTTP header to the response with the name
	 `Access-Control-Allow-Origin` and a value of `*`. This will
	  allow cross-origin AJAX requests to your server.
	- Get the `url` query string parameter value from the request.
	  If not supplied, respond with an http.StatusBadRequest error.
	- Call fetchHTML() to fetch the requested URL. See comments in that
	  function for more details.
	- Call extractSummary() to extract the page summary meta-data,
	  as directed in the assignment. See comments in that function
	  for more details
	- Close the response HTML stream so that you don't leak resources.
	- Finally, respond with a JSON-encoded version of the PageSummary
	  struct. That way the client can easily parse the JSON back into
	  an object. Remember to tell the client that the response content
	  type is JSON.
	Helpful Links:
	https://golang.org/pkg/net/http/#Request.FormValue
	https://golang.org/pkg/net/http/#Error
	https://golang.org/pkg/encoding/json/#NewEncoder
	*/
	w.Header().Set("Access-Control-Allow-Origin", "*")
	query, ok := r.URL.Query()["url"]

	if !ok || len(query) < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	readerCloser, err := fetchHTML(query[0])

	if err != nil {
		http.Error(w, "Can't Fetch", 503)
		return
	} else {
		pageExtract, err := extractSummary(query[0], readerCloser)
		if err != nil {
			http.Error(w, "Something wrong to extract Summary", 500)
		} else {
			w.Header().Set("Content-Type", "application/json")
			encoder := json.NewEncoder(os.Stdout)
			encoder.Encode(pageExtract)

			json, _ := json.Marshal(pageExtract)
			w.Write([]byte(string(json)))
		}
		defer readerCloser.Close()
		return
	}
}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	/*TODO: Do an HTTP GET for the page URL. If the response status
	code is >= 400, return a nil stream and an error. If the response
	content type does not indicate that the content is a web page, return
	a nil stream and an error. Otherwise return the response body and
	no (nil) error.
	To test your implementation of this function, run the TestFetchHTML
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestFetchHTML
	Helpful Links:
	https://golang.org/pkg/net/http/#Get
	*/

	_, err := url.ParseRequestURI(pageURL)

	if err != nil {
		return nil, errors.New("Not a valid url")
	}

	content, err := http.Get(pageURL)

	if err != nil {
		return nil, errors.New("Can't Fetch")
	}

	contentType := content.Header.Get("Content-type")
	if !strings.HasPrefix(contentType, "text/html") {
		return nil, errors.New("Not an HTML Page")
	}

	if content.StatusCode >= 400 {
		return nil, errors.New("Status Code Great Than or Equal to 400")
	}

	return ioutil.NopCloser(content.Body), nil
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	/*TODO: tokenize the `htmlStream` and extract the page summary meta-data
	according to the assignment description.
	To test your implementation of this function, run the TestExtractSummary
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestExtractSummary
	Helpful Links:
	https://drstearns.github.io/tutorials/tokenizing/
	http://ogp.me/
	https://developers.facebook.com/docs/reference/opengraph/
	https://golang.org/pkg/net/url/#URL.ResolveReference
	*/

	tokenizer := html.NewTokenizer(htmlStream)

	result := PageSummary{}
	images := []*PreviewImage{}
	videos := []*PreviewVideo{}

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}

			return nil, errors.New("error tokenizing HTML")
		}

		token := tokenizer.Token()

		if token.Data == "meta" {

			propName, propVal := getMetaAttr(&token)

			if propName != "" && propVal != "" {
				if propName == "og:type" {
					result.Type = propVal
				} else if propName == "og:url" {
					result.URL = propVal
				} else if propName == "og:title" {
					result.Title = propVal
				} else if propName == "og:site_name" {
					result.SiteName = propVal
				} else if propName == "og:description" {
					result.Description = propVal
				} else if propName == "description" && result.Description == "" {
					result.Description = propVal
				} else if propName == "author" {
					result.Author = propVal
				} else if propName == "keywords" {
					regex := regexp.MustCompile(`, *`)
					result.Keywords = regex.Split(propVal, -1)
				} else if propName == "og:image" {
					image := PreviewImage{}
					image.URL = getURL(pageURL, propVal)
					images = append(images, &image)
				} else if propName == "og:image:secure_url" {
					images[len(images)-1].SecureURL = propVal
				} else if propName == "og:image:type" {
					images[len(images)-1].Type = propVal
				} else if propName == "og:image:width" {
					width, _ := strconv.Atoi(propVal)
					images[len(images)-1].Width = width
				} else if propName == "og:image:height" {
					height, _ := strconv.Atoi(propVal)
					images[len(images)-1].Height = height
				} else if propName == "og:image:alt" {
					images[len(images)-1].Alt = propVal
				} else if propName == "og:video" {
					video := PreviewVideo{}
					video.URL = getURL(pageURL, propVal)
					videos = append(videos, &video)
				} else if propName == "og:video:secure_url" {
					videos[len(videos)-1].SecureURL = propVal
				} else if propName == "og:video:type" {
					videos[len(videos)-1].Type = propVal
				} else if propName == "og:video:width" {
					width, _ := strconv.Atoi(propVal)
					videos[len(videos)-1].Width = width
				} else if propName == "og:video:height" {
					height, _ := strconv.Atoi(propVal)
					videos[len(videos)-1].Height = height
				}
			}

		}

		if token.Data == "title" {
			tokenType = tokenizer.Next()
			if tokenType == html.TextToken {
				title := tokenizer.Token().Data
				if result.Title == "" {
					result.Title = title
				}
			}
		}

		if token.Data == "link" {
			processLink(&result, &token, pageURL)
		}
	}

	if len(images) > 0 {
		result.Images = images
	}

	if len(videos) > 0 {
		result.Videos = videos
	}

	return &result, nil
}

func getMetaAttr(token *html.Token) (prop string, val string) {
	allAttrs := (*token).Attr
	propVal := ""
	propName := ""
	for i := 0; i < len(allAttrs); i++ {
		if allAttrs[i].Key == "content" {
			propVal = allAttrs[i].Val
		}

		if allAttrs[i].Key == "name" || allAttrs[i].Key == "property" {
			propName = allAttrs[i].Val
		}
	}
	return propName, propVal
}

func processLink(summary *PageSummary, token *html.Token, pageURL string) {
	rel := ""
	url := ""
	linkType := ""
	height := 0
	width := 0

	attrs := (*token).Attr
	for i := 0; i < len(attrs); i++ {
		attr := attrs[i]
		if attr.Key == "rel" {
			rel = attr.Val
		} else if attr.Key == "href" {
			url = attr.Val
		} else if attr.Key == "type" {
			linkType = attr.Val
		} else if attr.Key == "sizes" {
			size := attr.Val
			if size != "any" {
				sizeSplit := strings.Split(attr.Val, "x")
				height, _ = strconv.Atoi(sizeSplit[0])
				width, _ = strconv.Atoi(sizeSplit[1])
			}
		}
	}

	if rel == "icon" {
		icon := PreviewImage{}
		icon.URL = getURL(pageURL, url)

		if linkType != "" {
			icon.Type = linkType
		}

		if height != 0 {
			icon.Height = height
		}

		if width != 0 {
			icon.Width = width
		}

		(*summary).Icon = &icon
	}
}

// getURL will parse and get the correct photo url
func getURL(pageURL string, testURL string) string {

	isURL, _ := regexp.MatchString("^http.*://", testURL)

	res := testURL

	if !isURL {
		url, _ := url.Parse(pageURL)
		pathSplit := strings.Split(url.Path, "/")
		resPath := "/"

		for i := 0; i < len(pathSplit); i++ {
			cur := pathSplit[i]
			if cur != "" && !strings.Contains(cur, ".html") {
				resPath = resPath + cur + "/"
			}
		}

		if testURL[0] == '/' {
			resPath = resPath[:len(resPath)-1]
		}

		res = url.Scheme + "://" + url.Host + resPath + testURL
	}

	return res
}
