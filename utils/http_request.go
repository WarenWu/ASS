package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	URL "net/url"
	"os"
	"strconv"
	"strings"
)

const (
	GET  = "GET"
	POST = "POST"
)

func UnescapeUnicode(raw []byte) (json string, err error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return "", err
	}
	return str, nil
}

func Request(method string, url string, params map[string]string, file []*os.File, header map[string]string) (body []byte, err error) {
	if method == "" || url == "" {
		return nil, nil
	}
	reqURL := url
	reqBody := &bytes.Buffer{}

	fncDo := func(method string, url string, body io.Reader) ([]byte, error) {
		client := http.Client{}
		req, err := http.NewRequest(method, url, body)
		if err != nil {
			return nil, err
		}
		for k, v := range header {
			req.Header.Add(k, v)
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		return ioutil.ReadAll(resp.Body)
	}

	if method == GET {
		{
			parseURL, err := URL.Parse(url)
			if err != nil {
				return nil, err
			}

			if params != nil {
				urlParams := &URL.Values{}
				for k, v := range params {
					urlParams.Set(k, v)
				}
				parseURL.RawQuery = urlParams.Encode()
			}

			reqURL = parseURL.String()
		}

		return fncDo(http.MethodGet, reqURL, reqBody)

	} else if method == POST {

		{
			writer := multipart.NewWriter(reqBody)
			if params != nil {
				for key, val := range params {
					err := writer.WriteField(key, val)
					if err != nil {
						return nil, err
					}
				}
			}
			if file != nil {
				for _, val := range file {
					fileWrite, err := writer.CreateFormFile("file", val.Name())
					_, err = io.Copy(fileWrite, val)
					if err != nil {
						return nil, err
					}
				}

			}
			return fncDo(http.MethodPost, reqURL, reqBody)
		}

	} else {
		return nil, nil
	}

	return nil, nil
}
