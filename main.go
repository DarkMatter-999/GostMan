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
    window.Resize(fyne.NewSize(800, 800))

	httpMethods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
    methodEntry := widget.NewSelect(httpMethods, func(selected string) {
	})
    // methodEntry.SetPlaceHolder("Method (GET, POST, etc.)")

    urlEntry := widget.NewEntry()
    urlEntry.SetPlaceHolder("URL")

    headersEntry := widget.NewMultiLineEntry()
    headersEntry.SetPlaceHolder("Headers (JSON format)")

    bodyEntry := widget.NewMultiLineEntry()
    bodyEntry.SetPlaceHolder("Body (JSON format)")

	responseEntry := widget.NewMultiLineEntry()
    // responseEntry.Disable()

    status := widget.NewLabel("")

    sendButton := widget.NewButton("Send", func() {
        req := Request{
            Method:  methodEntry.Selected,
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
			responseEntry.SetText(fmt.Sprintf("%s", res.Body))
            status.SetText(fmt.Sprintf("Status: %d", res.StatusCode))
		}
	})

	sendButton.Importance = widget.HighImportance

	topBar := container.NewBorder(nil, nil, methodEntry, sendButton, urlEntry,)

    dataTabs := container.NewAppTabs(
        container.NewTabItem("Body", bodyEntry),
        container.NewTabItem("Headers", headersEntry),
		container.NewTabItem("Auth", widget.NewLabel("Authorization Section")),
    )

    form := container.NewVBox(
        topBar,
        dataTabs,
    )

    tabs := container.NewAppTabs(
		container.NewTabItem("Response", container.NewBorder(
            status, 
            nil, nil, nil,
            container.NewVScroll(responseEntry),
            )),
	)

    tabs.SetTabLocation(container.TabLocationTop)

	content := container.NewBorder(form, nil, nil, nil, tabs)

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

