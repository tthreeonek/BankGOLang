package service

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/beevik/etree"
	"github.com/sirupsen/logrus"
)

type CBRService struct {
	url string
}

func NewCBRService() *CBRService {
	return &CBRService{url: "https://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx"}
}

func (s *CBRService) GetKeyRate() (float64, error) {
	soapBody := s.buildRequest()
	req, _ := http.NewRequest("POST", s.url, bytes.NewBuffer([]byte(soapBody)))
	req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")
	req.Header.Set("SOAPAction", "http://web.cbr.ru/KeyRate")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("CBR request failed: ", err)
		return 0, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	rate, err := parseRate(body)
	if err != nil {
		return 0, err
	}
	// маржа банка +5%
	return rate + 5, nil
}

func (s *CBRService) buildRequest() string {
	from := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	to := time.Now().Format("2006-01-02")
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
    <soap12:Envelope xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
      <soap12:Body>
        <KeyRate xmlns="http://web.cbr.ru/">
          <fromDate>%s</fromDate>
          <ToDate>%s</ToDate>
        </KeyRate>
      </soap12:Body>
    </soap12:Envelope>`, from, to)
}

func parseRate(xmlData []byte) (float64, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlData); err != nil {
		return 0, err
	}
	elements := doc.FindElements("//KeyRate/KR")
	if len(elements) == 0 {
		return 0, fmt.Errorf("rate not found")
	}
	rateElem := elements[0].FindElement("./Rate")
	if rateElem == nil {
		return 0, fmt.Errorf("Rate missing")
	}
	var rate float64
	fmt.Sscanf(rateElem.Text(), "%f", &rate)
	return rate, nil
}
