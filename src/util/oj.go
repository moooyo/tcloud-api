package util

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type PracticeOJResponse struct {
	OJ        uint
	ProblemID uint
	Title     string
	Url       string
}

func pojDispatch(id uint) (*PracticeOJResponse, error) {
	const pojUrl = "http://poj.org"
	url := pojUrl + "/problem?id=" + strconv.Itoa(int(id))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	errNode := doc.Find("font").Nodes[0].FirstChild.Data
	hasProblem := errNode != "Error Occurred"
	if !hasProblem {
		return nil, fmt.Errorf("%s", "no such problem")
	}
	tNode := doc.Find("div.ptt").Nodes
	title := tNode[0].FirstChild.Data
	return &PracticeOJResponse{
		OJ:        1,
		ProblemID: id,
		Url:       url,
		Title:     title,
	}, nil
}

func hduDispatch(id uint) (*PracticeOJResponse, error) {
	const pojUrl = "http://acm.hdu.edu.cn"
	url := pojUrl + "/showproblem.php?pid=" + strconv.Itoa(int(id))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	errNode := doc.Find("tbody").Find("tr").Find("td").Find("div")
	hasProblem := errNode.Nodes[1].FirstChild.Data != "Invalid Parameter."
	if !hasProblem {
		return nil, fmt.Errorf("%s", "no such problem")
	}
	titleNode := doc.Find("tbody").Find("tr").Find("td").Find("h1")
	title := titleNode.Nodes[0].FirstChild.Data
	return &PracticeOJResponse{
		OJ:        2,
		ProblemID: id,
		Url:       url,
		Title:     title,
	}, nil
}

func GetProblem(oj uint, id uint) (*PracticeOJResponse, error) {
	switch oj {
	case 1:
		return pojDispatch(id)
	case 2:
		return hduDispatch(id)
	default:
		return nil, fmt.Errorf("%s", "no such oj")
	}
}
