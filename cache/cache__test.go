package cache

import "testing"

func TestWrite(t *testing.T) {
	t.Log("With some content")
	_, err := Write("test_hssh", "Lorem ipsum")

	if err != nil {
		t.Errorf("Should not return any error")
	}

	t.Log("With empty content")
	_, err = Write("test_hssh", "")

	if err != nil {
		t.Errorf("Should not return any error")
	}
}

func TestRead(t *testing.T) {
	t.Log("With some content")
	fileName := "test_hssh"
	Write(fileName, "Lorem ipsum")
	content, err := Read(fileName)

	if err != nil {
		t.Errorf("Should not return any error")
	}

	if content != "Lorem ipsum" {
		t.Errorf("Should return the content provided instead of %s", content)
	}
}

func TestPreventMultiple(t *testing.T) {
	t.Log("try to write multiple files with same pattern name")
	fileName := "test_hssh"
	Write(fileName, "Handshakes")
	Write(fileName, "Hugs")
	content, err := Read(fileName)

	if err != nil {
		t.Errorf("Should not return any error")
	}

	if content != "Hugs" {
		t.Errorf("Should return the latest content provided instead of %s", content)
	}
}
