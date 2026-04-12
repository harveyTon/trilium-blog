package etapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type Note struct {
	NoteID       string      `json:"noteId"`
	Title        string      `json:"title"`
	DateModified string      `json:"dateModified"`
	Type         string      `json:"type"`
	Mime         string      `json:"mime"`
	Attributes   []Attribute `json:"attributes"`
}

type Attribute struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type NotesResponse struct {
	Results []Note `json:"results"`
}

func (c *Client) GetNotes(search string) ([]Note, error) {
	url := fmt.Sprintf("%s/etapi/notes?search=%s&orderBy=utcDateModified", c.baseURL, search)
	var resp NotesResponse
	if err := c.doRequest(url, &resp); err != nil {
		return nil, err
	}
	return resp.Results, nil
}

func (c *Client) GetNote(noteID string) (*Note, error) {
	url := fmt.Sprintf("%s/etapi/notes/%s", c.baseURL, noteID)
	var note Note
	if err := c.doRequest(url, &note); err != nil {
		return nil, err
	}
	return &note, nil
}

func (c *Client) GetNoteContent(noteID string) (string, error) {
	url := fmt.Sprintf("%s/etapi/notes/%s/content", c.baseURL, noteID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get note content: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

type Attachment struct {
	OwnerID string `json:"ownerId"`
	Mime    string `json:"mime"`
}

func (c *Client) GetAttachment(attachmentID string) (*Attachment, error) {
	url := fmt.Sprintf("%s/etapi/attachments/%s", c.baseURL, attachmentID)
	var att Attachment
	if err := c.doRequest(url, &att); err != nil {
		return nil, err
	}
	return &att, nil
}

func (c *Client) GetAttachmentContent(attachmentID string) ([]byte, string, error) {
	att, err := c.GetAttachment(attachmentID)
	if err != nil {
		return nil, "", err
	}

	url := fmt.Sprintf("%s/etapi/attachments/%s/content", c.baseURL, attachmentID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to get attachment content: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	return body, att.Mime, nil
}

func (c *Client) doRequest(url string, target interface{}) error {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
