package note

import (
	"testing"
)

func TestNewPostNoteTool(t *testing.T) {
	tool := NewPostNoteTool()

	expectedName := "post_misskey_note"
	if tool.GetName() != expectedName {
		t.Errorf("Expected tool name to be '%s', but got '%s'", expectedName, tool.GetName())
	}

	expectedDescription := "Post a note to Misskey"
	if tool.GetDescription() != expectedDescription {
		t.Errorf("Expected tool description to be '%s', but got '%s'", expectedDescription, tool.GetDescription())
	}
}
