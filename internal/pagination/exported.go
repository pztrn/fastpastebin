package pagination

import (
	"strconv"
	"strings"

	"go.dev.pztrn.name/fastpastebin/assets"
)

// CreateHTML creates pagination HTML based on passed parameters.
func CreateHTML(currentPage int, pages int, linksBase string) string {
	// Load templates.
	paginationHTMLRaw, err := assets.Data.ReadFile("pagination.html")
	if err != nil {
		return "Missing pagination.html"
	}

	paginationLinkRaw, err1 := assets.Data.ReadFile("pagination_link.html")
	if err1 != nil {
		return "Missing pagination_link.html"
	}

	paginationLinkCurrentRaw, err2 := assets.Data.ReadFile("pagination_link_current.html")
	if err2 != nil {
		return "Missing pagination_link_current.html"
	}

	paginationEllipsisRaw, err3 := assets.Data.ReadFile("pagination_ellipsis.html")
	if err3 != nil {
		return "Missing pagination_ellipsis.html"
	}

	// First page should always be visible.
	var paginationString string
	if currentPage == 1 {
		paginationString = strings.Replace(string(paginationLinkCurrentRaw), "{pageNum}", strconv.Itoa(currentPage), -1)
	} else {
		paginationString = strings.Replace(string(paginationLinkRaw), "{pageNum}", "1", -1)
		paginationString = strings.Replace(paginationString, "{paginationLink}", linksBase+"1", -1)
	}

	var (
		ellipsisStartAdded = false
		ellipsisEndAdded   = false
		// nolint:varnamelen
		i = 2
	)

	for i <= pages {
		// ToDo: fix it!
		// nolint:nestif
		if pages > 5 {
			if currentPage-3 < i && currentPage+3 > i || i == pages {
				paginationItemRaw := string(paginationLinkRaw)
				if i == currentPage {
					paginationItemRaw = string(paginationLinkCurrentRaw)
				}

				paginationItem := strings.Replace(paginationItemRaw, "{pageNum}", strconv.Itoa(i), -1)
				paginationItem = strings.Replace(paginationItem, "{paginationLink}", linksBase+strconv.Itoa(i), 1)
				paginationString += paginationItem
			} else {
				if currentPage-3 < i && !ellipsisStartAdded {
					paginationString += string(paginationEllipsisRaw)
					ellipsisStartAdded = true
				} else if currentPage+3 > i && !ellipsisEndAdded {
					paginationString += string(paginationEllipsisRaw)
					ellipsisEndAdded = true
				}
			}
		} else {
			paginationItemRaw := string(paginationLinkRaw)
			if i == currentPage {
				paginationItemRaw = string(paginationLinkCurrentRaw)
			}
			paginationItem := strings.Replace(paginationItemRaw, "{pageNum}", strconv.Itoa(i), -1)
			paginationItem = strings.Replace(paginationItem, "{paginationLink}", linksBase+strconv.Itoa(i), 1)
			paginationString += paginationItem
		}

		i++
	}

	pagination := strings.Replace(string(paginationHTMLRaw), "{paginationLinks}", paginationString, 1)
	if currentPage+1 <= pages {
		pagination = strings.Replace(pagination, "{nextPageLink}", linksBase+strconv.Itoa(currentPage+1), 1)
	} else {
		pagination = strings.Replace(pagination, "{nextPageLink}", linksBase+strconv.Itoa(pages), 1)
	}

	if currentPage-1 > 1 {
		pagination = strings.Replace(pagination, "{previousPageLink}", linksBase+strconv.Itoa(currentPage-1), 1)
	} else {
		pagination = strings.Replace(pagination, "{previousPageLink}", linksBase, 1)
	}

	return pagination
}
