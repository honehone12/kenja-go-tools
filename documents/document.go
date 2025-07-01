package documents

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Document interface {
	IdString() string
	Reader() (io.Reader, error)
}

type Parent struct {
	Id           bson.ObjectID `json:"id" bson:"id"`
	Name         string        `json:"name" bson:"name"`
	NameJapanese string        `json:"name_japanese,omitempty" bson:"name_japanese,omitempty"`
}

type FlatDocument struct {
	Id             bson.ObjectID `json:"-" bson:"_id"`
	IdHex          string        `json:"-" bson:"-"`
	ItemType       int32         `json:"item_type" bson:"item_type"`
	Rating         int32         `json:"rating" bson:"rating"`
	Url            string        `json:"url" bson:"url"`
	Parent         Parent        `json:"parent,omitzero" bson:"parent,omitempty"`
	Name           string        `json:"name" bson:"name"`
	NameEnglish    string        `json:"name_english,omitempty" bson:"name_english,omitempty"`
	NameJapanese   string        `json:"name_japanese,omitempty" bson:"name_japanese,omitempty"`
	Aliases        []string      `json:"aliases,omitempty" bson:"aliases,omitempty"`
	Studios        []string      `json:"studios,omitempty" bson:"studios,omitempty"`
	Staff          string        `json:"staff" bson:"staff"`
	Description    string        `json:"description" bson:"description"`
	ImageVector    bson.Vector   `json:"-" bson:"image_vector"`
	ImageVectorF32 []float32     `json:"image_vector" bson:"-"`
	TextVector     bson.Vector   `json:"-" bson:"text_vector"`
	TextVectorF32  []float32     `json:"text_vector" bson:"-"`
	StaffVector    bson.Vector   `json:"-" bson:"staff_vector"`
	StaffVectorF32 []float32     `json:"staff_vector" bson:"-"`
}

func (f *FlatDocument) Convert() error {
	var ok bool
	f.ImageVectorF32, ok = f.ImageVector.Float32OK()
	if !ok {
		return errors.New("failed to convert image vector type to float32")
	}
	f.TextVectorF32, ok = f.TextVector.Float32OK()
	if !ok {
		return errors.New("failed to convert text vector type to float32")
	}
	f.StaffVectorF32, ok = f.StaffVector.Float32OK()
	if !ok {
		return errors.New("failed to convert staff vector type to float32")
	}

	f.IdHex = f.Id.Hex()
	return nil
}

func (f *FlatDocument) Reader() (io.Reader, error) {
	b, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

func (f *FlatDocument) IdString() string {
	return f.IdHex
}
