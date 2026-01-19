package slack

import (
	"testing"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantChan  string
		wantMsgTS string
		wantThTS  string
		wantErr   bool
	}{
		{
			name:      "simple message URL",
			url:       "https://example.slack.com/archives/C12345678/p1716192523567890",
			wantChan:  "C12345678",
			wantMsgTS: "1716192523.567890",
			wantThTS:  "",
			wantErr:   false,
		},
		{
			name:      "thread message URL",
			url:       "https://example.slack.com/archives/C12345678/p1716192523567890?thread_ts=1716192500.123456",
			wantChan:  "C12345678",
			wantMsgTS: "1716192523.567890",
			wantThTS:  "1716192500.123456",
			wantErr:   false,
		},
		{
			name:      "thread with additional params",
			url:       "https://example.slack.com/archives/C12345678/p1716192523567890?thread_ts=1716192500.123456&cid=C12345678",
			wantChan:  "C12345678",
			wantMsgTS: "1716192523.567890",
			wantThTS:  "1716192500.123456",
			wantErr:   false,
		},
		{
			name:    "invalid URL",
			url:     "https://example.com/not-slack",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.ChannelID != tt.wantChan {
				t.Errorf("ParseURL() ChannelID = %v, want %v", got.ChannelID, tt.wantChan)
			}
			if got.MessageTS != tt.wantMsgTS {
				t.Errorf("ParseURL() MessageTS = %v, want %v", got.MessageTS, tt.wantMsgTS)
			}
			if got.ThreadTS != tt.wantThTS {
				t.Errorf("ParseURL() ThreadTS = %v, want %v", got.ThreadTS, tt.wantThTS)
			}
		})
	}
}

func TestNormalizeTimestamp(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"1716192523567890", "1716192523.567890"},
		{"1716192523.567890", "1716192523.567890"},
		{"123456", "123456"},
	}

	for _, tt := range tests {
		got := normalizeTimestamp(tt.input)
		if got != tt.want {
			t.Errorf("normalizeTimestamp(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
