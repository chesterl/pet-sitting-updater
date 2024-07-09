package main

import (
	"fmt"
  "log"
	"bytes"
	"encoding/json"
  // http requests
	"net/http"
  // handle compressed response
  "compress/gzip"
  "io"
  "os"
  //   load envr variables
  "github.com/joho/godotenv"

)

type Credentials struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
	ClientType        string `json:"clientType"`
}

type Response struct {
	IDToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"`
}

func main() {
  err := godotenv.Load()

  if err != nil {
    log.Fatal("Error loading .env file")
  }

  psKey := os.Getenv("PAWSHAKE_KEY")
  psUsername := os.Getenv("PAWSHAKE_USERNAME")
  psPassword := os.Getenv("PAWSHAKE_PASSWORD")
  //TODO add error handling if values are blank

  apiURL := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", psKey)

	creds := Credentials{
		Email:             psUsername,
		Password:          psPassword,
		ReturnSecureToken: true,
		ClientType:        "CLIENT_TYPE_WEB",
	}

	jsonData, err := json.Marshal(creds)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

  //TODO check what is ncessary
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Accept", "*/*")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	// req.Header.Set("Accept-Language", "en-GB,en;q=0.9,en-US;q=0.8,cs;q=0.7")
	// req.Header.Set("Dnt", "1")
	// req.Header.Set("Origin", "https://www.pawshake.com.au")
	// req.Header.Set("Priority", "u=1, i")
	// req.Header.Set("Sec-Ch-Ua", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	// req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	// req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	// req.Header.Set("Sec-Fetch-Dest", "empty")
	// req.Header.Set("Sec-Fetch-Mode", "cors")
	// req.Header.Set("Sec-Fetch-Site", "cross-site")
	// req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	// req.Header.Set("X-Client-Version", "Chrome/JsCore/9.23.0/FirebaseCore-web")
	// req.Header.Set("X-Firebase-Gmpid", "1:100067502341:web:3e30142aa87eaa9c9e2f6c")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

  var reader io.ReadCloser

  switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println("Error creating gzip reader:", err)
			return
		}
		defer reader.Close()
  default:
    reader = resp.Body
  }

	var response Response
	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	fmt.Println("ID Token:", response.IDToken)
	fmt.Println("Email:", response.Email)
	fmt.Println("Refresh Token:", response.RefreshToken)
	fmt.Println("Expires In:", response.ExpiresIn)
  fmt.Println("Local ID:", response.LocalID)
}
