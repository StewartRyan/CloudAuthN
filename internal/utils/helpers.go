package utils

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
)

// RawDeflateBase64Encode compresses the input string using raw deflate algorithm and
// then encodes the compressed data in base64 format.
//
// The function takes a single parameter:
// - input: A string that needs to be compressed and encoded.
//
// The function returns two values:
// - A string representing the compressed and encoded data.
// - An error if any error occurs during the compression or encoding process.
func RawDeflateBase64Encode(input string) (string, error) {
	// Create a buffer to hold the deflated data
	var deflated bytes.Buffer

	// Create a flate writer with no compression (raw deflate)
	writer, err := flate.NewWriter(&deflated, flate.DefaultCompression)
	if err != nil {
		return "", err
	}

	// Write the input data to the flate writer
	_, err = writer.Write([]byte(input))
	if err != nil {
		return "", err
	}

	// Close the writer to flush any remaining data
	err = writer.Close()
	if err != nil {
		return "", err
	}

	// Encode the deflated data in base64
	encoded := base64.StdEncoding.EncodeToString(deflated.Bytes())

	return encoded, nil
}
