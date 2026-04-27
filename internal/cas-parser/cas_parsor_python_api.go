package casparser

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	casparsermodel "gin-investment-tracker/internal/cas-parser/model"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"
)

type CasParserPythonApi struct{}

func NewCasParserPythonApi() *CasParserPythonApi {
	return &CasParserPythonApi{}
}

func (p *CasParserPythonApi) ProcessCasFile(ctx context.Context, file *multipart.FileHeader, filePassword string, userID int64) (*casparsermodel.CASStatement, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		slog.Error("Got error while creating file in form", "error", err.Error())
		return nil, err
	}

	openedFile, err := file.Open()
	if err != nil {
		slog.Error("Got error while opening file", "error", err.Error())
		return nil, err
	}
	io.Copy(part, openedFile)

	writer.WriteField("password", filePassword)
	writer.Close()

	req, err := http.NewRequest("POST", "http://localhost:8000/parse-cas", &body)
	if err != nil {
		slog.Error("Got error while creating request", "error", err.Error())
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Got error while seding request", "error", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	var statement casparsermodel.CASStatement
	if err := json.Unmarshal(respBody, &statement); err != nil {
		return nil, fmt.Errorf("cas-parser: unmarshal: %w", err)
	}
	return &statement, nil
}
