# DB Recordum

This Go program automatically backs up a MySQL database and uploads the backup files to a specified Google Drive folder using a service account. It runs periodically based on the configured frequency and is designed to be easily deployed using Docker.

## Overview

The program performs the following tasks:

-   Connects to a MySQL database and creates a backup using `mysqldump`.
-   Uploads the backup file to Google Drive, in a folder shared with a service account.
-   Runs automatically according to the specified frequency.

## Requirements

-   **Docker** and **docker-compose** installed on your system.
-   A Google Cloud service account with access to the Google Drive API.
-   The `JSON` credentials file for the service account downloaded from Google Cloud.
-   A folder in Google Drive shared with the service account.

## Setup

### 1. Create a Service Account in Google Cloud

1. Go to the [Google Cloud Console](https://console.cloud.google.com/).
2. Create a new project (if you haven't done so already).
3. Enable the **Google Drive API**.
4. Go to **APIs & Services > Credentials** and create a **Service Account**.
5. Download the `JSON` key file for the service account.
6. **Important**: Keep this file secure, as it contains sensitive credentials.

### 2. Configure Google Drive

1. Open Google Drive and create a folder for the backups (or use an existing folder).
2. Share the folder with the service account. Use the email address found in the `JSON` key file, which looks like `service-account-name@project-id.iam.gserviceaccount.com`.
3. Copy the **folder ID** from the URL of the Google Drive folder (the part after `/folders/` in the URL).

### 3. Configure the Project

1. Create a folder named `my_drive_credentials` in the root directory of your project and place the `service-account.json` file inside.
2. Create a `docker-compose.yml` file with the following configuration:

    ```yaml
    version: "3.8"

    services:
        backup-service:
            image: unkcode/db-recordum:tagname
            environment:
                - MYSQL_DB=mydb
                - MYSQL_HOST=myhost
                - MYSQL_USER=myuser
                - MYSQL_PASSWORD=mypassword
                - BACKUP_FRECUENCY=24h
                - GDRIVE_FOLDER_ID=your_folder_id
                - GDRIVE_FILE_PREFIX=my_backup
            volumes:
                - ./my_drive_credentials:/my_drive_credentials
    ```

3. Replace the values in the `docker-compose.yml` file:
    - `MYSQL_DB`: The name of your database.
    - `MYSQL_HOST`: The IP address or hostname of your database.
    - `MYSQL_USER`: The database user.
    - `MYSQL_PASSWORD`: The database user's password.
    - `BACKUP_FRECUENCY`: The backup frequency in hours (e.g., `24` for a daily backup).
    - `GDRIVE_FOLDER_ID`: The ID of the Google Drive folder where the backups will be uploaded.
    - `GDRIVE_FILE_PREFIX`: The prefix used in the backup file names.

### 4. Build and Run the Docker Container

1. Build the Docker image (if it is not already available on Docker Hub):

    ```bash
    docker build -t your-docker-image .
    ```

2. Start the service using Docker Compose:

    ```bash
    docker-compose up -d
    ```

3. The container will run in the background, making backups periodically based on the specified frequency.

## Important Notes

-   **mysqldump**: Ensure `mysqldump` is available and properly configured in the Docker container if you are managing dependencies manually.
-   **Security**: Do not upload the `service-account.json` file to any public repository and handle the credentials securely.
-   **Logs**: The logs for the backup and upload processes are recorded in the container's console. You can check the logs using:
    ```bash
    docker logs backup-service
    ```

---

That's it! With this setup, your program will automatically back up and upload files to Google Drive.
