package misskey

import (
	misskey "github.com/yitsushi/go-misskey"
	misskey_sdk "github.com/yitsushi/go-misskey/services/notes"
)

// Client is an interface for a Misskey client.
type Client interface {
	Notes() NotesService
}

// NotesService is an interface for the notes service.
type NotesService interface {
	Create(misskey_sdk.CreateRequest) (*misskey_sdk.CreateResponse, error)
}

// misskeyGoClientWrapper wraps the go-misskey client.
type misskeyGoClientWrapper struct {
	*misskey.Client
}

// NewMisskeyGoClientWrapper creates a new misskeyGoClientWrapper.
func NewMisskeyGoClientWrapper(client *misskey.Client) Client {
	return &misskeyGoClientWrapper{Client: client}
}

// Notes returns a wrapper for the notes service.
func (c *misskeyGoClientWrapper) Notes() NotesService {
	return &notesServiceWrapper{NotesService: c.Client.Notes()} // Added ()
}

// notesServiceWrapper wraps the go-misskey notes service.
type notesServiceWrapper struct {
	NotesService *misskey_sdk.Service
}

// Create calls the underlying SDK's Create method.
func (s *notesServiceWrapper) Create(req misskey_sdk.CreateRequest) (*misskey_sdk.CreateResponse, error) {
	res, err := s.NotesService.Create(req)
	if err != nil {
		return nil, err
	}
	return &res, nil // Return address of res
}
