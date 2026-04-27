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
	"os"
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
	defer openedFile.Close()
	if _, err := io.Copy(part, openedFile); err != nil {
		slog.Error("Got error while copying file contents", "error", err.Error())
		return nil, err
	}
	if err := writer.WriteField("password", filePassword); err != nil {
		slog.Error("Got error while writing password field", "error", err.Error())
		return nil, err
	}
	if err := writer.Close(); err != nil {
		slog.Error("Got error while closing multipart writer", "error", err.Error())
		return nil, err
	}

	parserUrl := os.Getenv("CAS_PARSER_API")
	if parserUrl == "" {
		parserUrl = "http://localhost:8000/parse-cas"
	}

	req, err := http.NewRequest("POST", parserUrl, &body)
	if err != nil {
		slog.Error("Got error while creating request", "error", err.Error())
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Got error while seding request", "error", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cas-parser: read response body: %w", err)
	}

	var statement casparsermodel.CASStatement
	if err := json.Unmarshal(respBody, &statement); err != nil {
		return nil, fmt.Errorf("cas-parser: unmarshal: %w", err)
	}
	return &statement, nil
}
