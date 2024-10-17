package controllers

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	uID := uuid.New().String()
	uploadPath := filepath.Join("uploads", file.Filename)
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}

	key, err := GenerateKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Key generation failed"})
		return
	}

	cmd := exec.Command("/bin/bash", "encrypt.sh", uploadPath, uID, uID, base64.StdEncoding.EncodeToString(key))
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("Encryption failed: %s\n", string(output))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Encryption failed"})
		return
	}

	err = UpdateContentInLCP(uID, key, uploadPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Encryption failed"})
		return
	}

	ext := filepath.Ext(file.Filename)
	mimeType := mime.TypeByExtension(ext)
	license, err := CreateLicense(uID, "user@example.com", key, mimeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "License creation failed"})
		return
	}

	fileData, err := os.Create(fmt.Sprintf("/root/server/backend/uploads/%s.lcpl", uID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "License creation failed"})
		return
	}
	defer fileData.Close()

	_, err = fileData.Write(license)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "License creation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "File encrypted and license created successfully",
		"license":       string(license),
		"download_link": fmt.Sprintf("/uploads/%s.lcpl", uID),
	})
}

func CreateLicense(contentID, userEmail string, contentKey []byte, mime string) ([]byte, error) {
	now := time.Now()

	hashBytes := sha256.Sum256([]byte(userEmail))
	hash := base64.StdEncoding.EncodeToString(hashBytes[:])
	licensePayload := map[string]interface{}{
		"provider": "https://lsp.dwbank.org/lcpserver",
		"id":       contentID,
		"encryption": map[string]interface{}{
			"user_key": map[string]interface{}{
				"text_hint": "email",
				"algorithm": "http://www.w3.org/2001/04/xmlenc#sha256",
				"value":     hash,
			},
		},
		"links": []map[string]interface{}{
			{
				"rel":  "publication",
				"href": fmt.Sprintf("https://lsp.dwbank.org/uploads/%s", contentID),
				"type": mime,
			},
		},
		"user": map[string]interface{}{
			"id":        userEmail,
			"email":     userEmail,
			"name":      userEmail,
			"encrypted": []string{"email"},
		},
		"rights": map[string]interface{}{
			"copy":  1000,
			"print": 10,
			"start": now.UTC().Format(time.RFC3339),
			"end":   now.AddDate(0, 1, 0).UTC().Format(time.RFC3339),
		},
	}

	jsonData, _ := json.Marshal(licensePayload)

	url := fmt.Sprintf("http://127.0.0.1:8989/contents/%s/licenses", contentID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	req.SetBasicAuth("admin", "admin")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GenerateUserKey(passphrase string) []byte {
	hash := sha256.Sum256([]byte(passphrase))
	return hash[:]
}

func UpdateContentInLCP(contentID string, key []byte, filePath string) error {
	url := fmt.Sprintf("http://127.0.0.1:8989/contents/%s", contentID)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to stat file: %v", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return fmt.Errorf("failed to calculate hash: %v", err)
	}
	hash := hasher.Sum(nil)

	ext := filepath.Ext(fileInfo.Name())
	mimeType := mime.TypeByExtension(ext)
	data := map[string]interface{}{
		"content-id":                    contentID,
		"storage-mode":                  0,
		"protected-content-type":        mimeType,
		"protected-content-length":      fileInfo.Size(),
		"protected-content-sha256":      fmt.Sprintf("%x", hash),
		"protected-content-location":    fmt.Sprintf("https://lsp.dwbank.org/uploads/%s", contentID),
		"protected-content-disposition": contentID,
		"content-encryption-key":        base64.StdEncoding.EncodeToString(key),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	req.SetBasicAuth("admin", "admin")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update content: %s", string(body))
	}

	return nil
}

func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
