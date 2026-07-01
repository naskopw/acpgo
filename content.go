package acp

// ContentBlock is a piece of content in a message (text, image, audio, etc.).
type ContentBlock struct {
	Type        string            `json:"type"`
	Text        string            `json:"text,omitempty"`
	Data        string            `json:"data,omitempty"`
	MimeType    string            `json:"mimeType,omitempty"`
	URI         string            `json:"uri,omitempty"`
	Name        string            `json:"name,omitempty"`
	Title       string            `json:"title,omitempty"`
	Description string            `json:"description,omitempty"`
	Size        int64             `json:"size,omitempty"`
	Annotations *Annotations      `json:"annotations,omitempty"`
	Resource    *EmbeddedResource `json:"resource,omitempty"`
	Meta        map[string]any    `json:"_meta,omitempty"`
}

// EmbeddedResource is a resource embedded directly in a message.
type EmbeddedResource struct {
	MimeType string         `json:"mimeType,omitempty"`
	Text     string         `json:"text,omitempty"`
	Blob     string         `json:"blob,omitempty"`
	URI      string         `json:"uri"`
	Meta     map[string]any `json:"_meta,omitempty"`
}

// Annotations provide optional metadata for content blocks.
type Annotations struct {
	Audience     []string      `json:"audience,omitempty"`
	LastModified string        `json:"lastModified,omitempty"`
	Priority     float64       `json:"priority,omitempty"`
	Meta         map[string]any `json:"_meta,omitempty"`
}

// ContentChunk is a streamed chunk of content during a prompt turn.
type ContentChunk struct {
	Content   *ContentBlock  `json:"content"`
	MessageID string         `json:"messageId,omitempty"`
	Meta      map[string]any `json:"_meta,omitempty"`
}

