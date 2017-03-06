// Package apm provides an interface to fetch and scrape data out of Amazon Prime Movies.
package apm

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	baseURL = "https://www.amazon.de/gp/product/"
)

// A Movie holds the data of an Amazon Prime movie and methods to parse it.
type Movie struct {
	Title       string   `json:"title"`
	ReleaseYear int      `json:"release_year"`
	Actors      []string `json:"actors"`
	Poster      string   `json:"poster"`
	SimilarIDs  []string `json:"similar_ids"`
}

// IsEmpty checks if Movie is empty by it's Title field.
// An empty Title should be enough to consider the Movie empty.
func (m Movie) IsEmpty() bool {
	return m.Title == ""
}

// findTitle scrapes node and sets the Movie Title.
func (m *Movie) findTitle(node *html.Node) {
	title, ok := scrape.Find(node, scrape.ById("aiv-content-title"))
	if ok {
		// title text also includes the year
		// thus the first node text is used
		m.Title = scrape.TextJoin(title, firstNodeText)
	}
}

// findReleaseYear scrapes node and sets the Movie ReleaseYear.
func (m *Movie) findReleaseYear(node *html.Node) {
	y, ok := scrape.Find(node, scrape.ByClass("release-year"))
	if ok {
		if yn, err := strconv.Atoi(scrape.Text(y)); err == nil {
			m.ReleaseYear = yn
		}
	}
}

// findActors scrapes node and sets the Movie Actors.
func (m *Movie) findActors(node *html.Node) {
	ac, ok := scrape.Find(node, scrape.ByClass("dv-meta-info"))
	if ok {
		// the actors are in the first <dd> inside .dv-meta-info
		// scrape.Find returns the first matching node
		a, ok := scrape.Find(ac, scrape.ByTag(atom.Dd))
		if ok {
			m.Actors = strings.Split(scrape.Text(a), ", ")
		}
	}
}

// findPoster scrapes node and sets the Movie Poster.
func (m *Movie) findPoster(node *html.Node) {
	pc, ok := scrape.Find(node, scrape.ById("dv-dp-left-content"))
	if ok {
		// the poster is the first <img> inside #dv-dp-left-content
		// scrape.Find returns the first matching node
		a, ok := scrape.Find(pc, scrape.ByTag(atom.Img))
		if ok {
			m.Poster = scrape.Attr(a, "src")
		}
	}
}

// findSimilarIDs scrapes node and sets the Movie SimilarIDs.
func (m *Movie) findSimilarIDs(node *html.Node) {
	ulNode, ok := scrape.Find(node, scrape.ByTag(atom.Ul))
	if ok {
		similarsNodes := scrape.FindAll(ulNode, scrape.ByTag(atom.A))
		for _, s := range similarsNodes {
			// scrape.Attr returns "" when the attribute is empty or not present
			if url := scrape.Attr(s, "href"); url != "" {
				m.SimilarIDs = append(m.SimilarIDs, parseID(url))
			}
		}
	}
}

// FetchMovie fetches an Amazon Prime movie page for the amazonID,
// scrapes it and returns the movie data.
//
// This is the page structure and where the movie information is placed:
//
// 	#aiv-main-content
// 		#dv-dp-title-content
// 			h1#aiv-content-title ---  movie title
// 				span.release-year  ---  movie year
// 		#dv-dp-left-content
// 			img                  ---  movie poster
// 			img
// 		.dv-info
// 			dl.dv-meta-info
// 				dt
// 				dd                 --- movie actors
// 				dt
// 				dd
// 	#dv-sims
// 		ul
// 			li
// 				a                  --- related movie
// 			li
// 				a                  --- related movie
// 			...
//
// Example movie: https://www.amazon.de/gp/product/B00K19SD8Q.
func FetchMovie(amazonID string) (Movie, error) {
	resp, err := http.Get(baseURL + amazonID)
	if err != nil {
		return Movie{}, fmt.Errorf("Error trying to fetch amazon_id: %v", amazonID)
	}

	root, err := html.Parse(resp.Body)
	if err != nil {
		return Movie{}, fmt.Errorf("Error trying to fetch amazon_id: %v", amazonID)
	}

	var m Movie

	dataNode, ok := scrape.Find(root, scrape.ById("aiv-main-content"))
	if ok {
		m.findTitle(dataNode)

		m.findReleaseYear(dataNode)

		m.findActors(dataNode)

		m.findPoster(dataNode)
	}

	similarsNode, ok := scrape.Find(root, scrape.ById("dv-sims"))
	if ok {
		m.findSimilarIDs(similarsNode)
	}

	return m, nil
}

// firstNodeText is used to return the text of the first node.
func firstNodeText(nodesText []string) string {
	if len(nodesText) > 0 {
		return strings.TrimSpace(nodesText[0])
	}
	return ""
}

// parseID returns the amazon_id from a Amazon Prime product URL.
func parseID(url string) string {
	if strs := strings.Split(url, "product/"); len(strs) >= 2 {
		return strings.Split(strs[1], "/")[0]
	}
	return ""
}
