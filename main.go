package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// Function to get environment variables and ensure they are not empty
func getEnvOrFail(varName string) string {
	value := os.Getenv(varName)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", varName)
	}
	return value
}

func main() {
	// Path to the JSON credentials file
	credentialsFile := "/credentials/service-account.json"

	// Check if the credentials file exists
	if _, err := os.Stat(credentialsFile); os.IsNotExist(err) {
		log.Fatalf("Credentials file does not exist: %v", err)
	}

	// Get environment variables with validation
	db := getEnvOrFail("MYSQL_DB")
	host := getEnvOrFail("MYSQL_HOST")
	user := getEnvOrFail("MYSQL_USER")
	password := getEnvOrFail("MYSQL_PASSWORD")
	driveFolderId := getEnvOrFail("GDRIVE_FOLDER_ID")
	driveFilePrefix := getEnvOrFail("GDRIVE_FILE_PREFIX")
	backupFrequency := getEnvOrFail("BACKUP_FRECUENCY")

	// Parse the backup frequency
	frequency, err := time.ParseDuration(backupFrequency)
	if err != nil {
		log.Fatalf("Error parsing BACKUP_FRECUENCY: %v", err)
	}

	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	for {
		<-ticker.C
		backupFile := fmt.Sprintf("%s.%s_%s.sql", driveFilePrefix, db, time.Now().Format("20060102_150405"))
		err := backupDatabase(host, user, password, db, backupFile)
		if err != nil {
			log.Printf("Error making database backup: %v", err)
			continue
		}

		err = uploadFileToDrive(credentialsFile, backupFile, driveFolderId)
		if err != nil {
			log.Printf("Error uploading backup to Google Drive: %v", err)
			continue
		}

		log.Printf("Backup completed successfully: %s", backupFile)
	}
}

func backupDatabase(host, user, password, db, backupFile string) error {
	cmd := exec.Command("mysqldump", "-h"+host, "-u"+user, "-p"+password, db)
	outFile, err := os.Create(backupFile)
	if err != nil {
		return fmt.Errorf("error creating backup file: %v", err)
	}
	defer outFile.Close()

	cmd.Stdout = outFile
	err = cmd.Run()
	if err != nil {
		log.Printf("Error running mysqldump: %v", err)
	}
	return err
}

func uploadFileToDrive(credentialsFile, filePath, folderID string) error {
	// Create the Google Drive service
	ctx := context.Background()
	srv, err := drive.NewService(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return fmt.Errorf("could not create Google Drive service: %v", err)
	}

	// Open the file to upload
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	// Configure the file to be uploaded to Google Drive
	driveFile := &drive.File{
		Name:    filePath,
		Parents: []string{folderID}, // ID of the folder where the file will be uploaded
	}

	// Upload the file to Google Drive
	_, err = srv.Files.Create(driveFile).Media(file).Do()
	if err != nil {
		return fmt.Errorf("error uploading file to Google Drive: %v", err)
	}

	log.Printf("File uploaded successfully: %s", filePath)
	return nil
}
