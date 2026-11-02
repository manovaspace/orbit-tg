package mock_test

import (
	"context"
	"testing"

	tg "github.com/manovaspace/orbit-tg"
	"github.com/manovaspace/orbit-tg/mock"
)

func TestMockRecordsPhotoAndDocument(t *testing.T) {
	ctx := context.Background()
	m := mock.NewMessenger()

	photo := tg.InputFile{Data: []byte("jpeg-bytes"), Filename: "photo.jpg"}
	if _, err := m.SendPhoto(ctx, -100123, photo, "photo caption"); err != nil {
		t.Fatalf("SendPhoto: %v", err)
	}
	if len(m.Photos) != 1 {
		t.Fatalf("photos = %d, want 1", len(m.Photos))
	}
	gotPhoto := m.Photos[0]
	if gotPhoto.ChatID != -100123 || gotPhoto.Caption != "photo caption" || gotPhoto.Photo.Filename != "photo.jpg" {
		t.Errorf("photo record mismatch: %+v", gotPhoto)
	}

	document := tg.InputFile{Data: []byte("pdf-bytes"), Filename: "receipt.pdf"}
	if _, err := m.SendDocument(ctx, -100123, document, "document caption"); err != nil {
		t.Fatalf("SendDocument: %v", err)
	}
	if len(m.Documents) != 1 {
		t.Fatalf("documents = %d, want 1", len(m.Documents))
	}
	gotDocument := m.Documents[0]
	if gotDocument.ChatID != -100123 || gotDocument.Caption != "document caption" || gotDocument.Document.Filename != "receipt.pdf" {
		t.Errorf("document record mismatch: %+v", gotDocument)
	}
}
