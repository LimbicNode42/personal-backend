// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

type DeletePost struct {
	ID          string            `json:"id"`
	Published   *bool             `json:"published,omitempty"`
	Title       string            `json:"title"`
	Text        *string           `json:"text,omitempty"`
	Attachments []*graphql.Upload `json:"attachments,omitempty"`
	Tags        []*Tags           `json:"tags,omitempty"`
}

type EditPost struct {
	ID                   string            `json:"id"`
	Published            bool              `json:"published"`
	Title                string            `json:"title"`
	Text                 string            `json:"text"`
	UnchangedAttachments []*string         `json:"unchangedAttachments,omitempty"`
	NewAttachments       []*graphql.Upload `json:"newAttachments,omitempty"`
	DeletedAttachments   []*string         `json:"deletedAttachments,omitempty"`
	Tags                 []*Tags           `json:"tags,omitempty"`
}

type Mutation struct {
}

type NewPost struct {
	Published   bool              `json:"published"`
	Title       string            `json:"title"`
	Text        string            `json:"text"`
	Attachments []*graphql.Upload `json:"attachments,omitempty"`
	Tags        []*Tags           `json:"tags,omitempty"`
}

type Post struct {
	ID          string    `json:"id"`
	Published   bool      `json:"published"`
	Title       string    `json:"title"`
	Text        string    `json:"text"`
	Tags        []*Tags   `json:"tags,omitempty"`
	Attachments []*string `json:"attachments,omitempty"`
}

type Query struct {
}

type Tags string

const (
	TagsCoding             Tags = "Coding"
	TagsSystemArchitecture Tags = "System_Architecture"
	TagsBook               Tags = "Book"
)

var AllTags = []Tags{
	TagsCoding,
	TagsSystemArchitecture,
	TagsBook,
}

func (e Tags) IsValid() bool {
	switch e {
	case TagsCoding, TagsSystemArchitecture, TagsBook:
		return true
	}
	return false
}

func (e Tags) String() string {
	return string(e)
}

func (e *Tags) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Tags(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Tags", str)
	}
	return nil
}

func (e Tags) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
