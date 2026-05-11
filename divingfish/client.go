package divingfish

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Etag       string
}

var defaultBaseURL = "https://www.diving-fish.com/api/maimaidxprober"

func NewClient(etag string) *Client {
	return &Client{
		BaseURL:    defaultBaseURL,
		HTTPClient: &http.Client{},
		Etag:       etag,
	}
}

// 获取所有歌曲数据
func (c *Client) GetMaiMaiMusicData(ctx context.Context) ([]MaiMaiMusicData, string, error) {
	url := fmt.Sprintf("%s/music_data", c.BaseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	req.Header.Set("If-None-Match", c.Etag)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotModified {
			// 没有修改，返回空列表
			return []MaiMaiMusicData{}, "", nil
		}
		return nil, "", fmt.Errorf("status code: %d", resp.StatusCode)
	}

	c.Etag = resp.Header.Get("etag")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	var musicDataList []MaiMaiMusicData
	if err := json.Unmarshal(body, &musicDataList); err != nil {
		return nil, "", err
	}
	return musicDataList, c.Etag, nil
}
