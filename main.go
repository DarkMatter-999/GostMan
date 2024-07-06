package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Request struct {
    Method  string            `json:"method"`
    URL     string            `json:"url"`
    Headers map[string]string `json:"headers"`
    Body    string            `json:"body"`
}

type Response struct {
    StatusCode int               `json:"statusCode"`
    Headers    map[string]string `json:"headers"`
    Body       string            `json:"body"`
}

func main() {
    app := app.New()

    window := app.NewWindow("GostMan REST API Client")
    window.Resize(fyne.NewSize(600, 800))

    methodEntry := widget.NewEntry()
    methodEntry.SetPlaceHolder("Method (GET, POST, etc.)")

    urlEntry := widget.NewEntry()
    urlEntry.SetPlaceHolder("URL")

    headersEntry := widget.NewMultiLineEntry()
    headersEntry.SetPlaceHolder("Headers (JSON format)")

    bodyEntry := widget.NewMultiLineEntry()
    bodyEntry.SetPlaceHolder("Body (JSON format)")

	responseEntry := widget.NewMultiLineEntry()
    responseEntry.Disable()

    sendButton := widget.NewButton("Send Request", func() {
        req := Request{
            Method:  methodEntry.Text,
            URL:     urlEntry.Text,
            Headers: make(map[string]string),
            Body:    bodyEntry.Text,
        }

        // Parse headers
        if headersEntry.Text != "" {
            if err := json.Unmarshal([]byte(headersEntry.Text), &req.Headers); err != nil {
                fmt.Println("Error parsing headers:", err)
                return
            }
        }

		res, err := sendRequest(req)
		if err != nil {
			responseEntry.SetText(fmt.Sprintf("Error: %s", err))
		} else {
			responseEntry.SetText(fmt.Sprintf("Status: %d\n\n%s", res.StatusCode, res.Body))    
		}
	})

    form := container.NewVBox(
        widget.NewForm(
            widget.NewFormItem("Method", methodEntry),
            widget.NewFormItem("URL", urlEntry),
            widget.NewFormItem("Headers", headersEntry),
            widget.NewFormItem("Body", bodyEntry),
        ),
        sendButton,
		widget.NewLabel("Response"),
    )

	content := container.NewBorder(form, nil, nil, nil, container.NewVScroll(responseEntry))

    window.SetContent(content)

    window.ShowAndRun()
}

func sendRequest(req Request) (Response, error) {
    // Create an HTTP client
    client := &http.Client{
		Timeout: 15 * time.Second,
	}

    httpRequest, err := http.NewRequest(req.Method, req.URL, bytes.NewBufferString(req.Body))
    if err != nil {
        fmt.Println("Error creating HTTP request:", err)
        return Response{}, err
    }

    for key, value := range req.Headers {
        httpRequest.Header.Set(key, value)
    }

    // Send the HTTP request
    httpResponse, err := client.Do(httpRequest)
    if err != nil {
        fmt.Println("Error sending HTTP request:", err)
        return Response{}, err
    }
    defer httpResponse.Body.Close()

    // Read the response body
    var responseBody bytes.Buffer
    _, err = responseBody.ReadFrom(httpResponse.Body)
    if err != nil {
        fmt.Println("Error reading response body:", err)
        return Response{}, err
    }

    // Create a Response struct
	response := Response{
        StatusCode: httpResponse.StatusCode,
        Headers:    make(map[string]string),
        Body:       responseBody.String(),
    }

    for key, value := range httpResponse.Header {
        response.Headers[key] = value[0]
    }

    // fmt.Println("Response:")
    // fmt.Printf("Status Code: %d\n", response.StatusCode)
    // fmt.Println("Headers:", response.Headers)
    // fmt.Println("Body:", response.Body)

	return response, nil
}
