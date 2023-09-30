package privatebin

import (
	"encoding/json"
	"errors"
	"net/http"
)

func GetPaste(instance, id, key string) (string, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		instance+"/?pasteid="+id,
		nil,
	)
	if err != nil {
		return "", err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-type", "application/json")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return "", errors.New("unsuccessful http status code")
	}

	parsed := Response{}
	if err := json.NewDecoder(response.Body).Decode(&parsed); err != nil {
		return "", err
	}

	rawKey, err := deriveKey(
		key,
		parsed.Adata[0].([]any)[1].(string),
		int(parsed.Adata[0].([]any)[2].(float64)),
	)
	if err != nil {
		return "", err
	}

	plainText, err := decryptContent(
		rawKey,
		parsed.Adata[0].([]any)[0].(string),
		parsed.Adata,
		parsed.Ct,
	)
	if err != nil {
		return "", err
	}

	if parsed.Adata[0].([]any)[7].(string) == "zlib" {
		plainText = zlibDecompress(plainText)
	}

	paste := struct {
		Paste string `json:"paste"`
	}{}
	json.Unmarshal(plainText, &paste)

	return paste.Paste, nil
}
