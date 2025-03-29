# SwiftCodeApp

## Overview
SwiftCodeApp is a RESTful API service that parses, stores, and retrieves SWIFT codes from a spreadsheet. The application processes SWIFT data, stores it in a database, and provides endpoints for retrieving SWIFT codes by country and individual lookup.

## Technologies Used
- **Go** (Preferred language for implementation)
- **PostgreSQL** (Database for storing SWIFT codes)
- **Docker & Docker Compose** (Containerized setup)
- **Google API Client** (For interacting with Google Sheets)

## Prerequisites
Before running the application, ensure you have the following installed:

- [Go](https://go.dev/dl/) (Version 1.24.1 recommended)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Git](https://git-scm.com/downloads)
- Google API credentials (explained below)

## Setup Instructions

### 1. Clone the Repository
```sh
git clone https://github.com/mswiatek12/SwiftCodeApp.git
cd SwiftCodeApp
```

### 2. Obtain Google API Credentials
The application interacts with Google Sheets. You need to download Google API credentials:

1. Visit [Google Cloud Console](https://console.cloud.google.com/).
2. Create or select a project.
3. Enable "Google Sheets API" and "Google Drive API".
4. Create a new **service account** and generate JSON key.
5. Download the JSON file and place it in the project root as `credentials.json`.

### 3. Run the Application
To start the service using Docker:
```sh
docker-compose up --build
```
This will:
- Build the Go application.
- Set up the PostgreSQL database.
- Expose the API at `http://localhost:8080`.

### 4. Running Tests
To run tests inside the container:
```sh
docker exec -it swiftcode_app /bin/sh
cd /app
go mod init SwiftCodeApp
go mod tidy
go test -v ./tests/...
```

## API Endpoints

### Retrieve details of a single SWIFT code
**GET** `/v1/swift-codes/{swift-code}`

#### Response (Headquarter SWIFT code):
```json
{
  "address":        "string",
  "bankName":       "string",
  "countryISO2":    "string",
  "countryName":    "string",
  "isHeadquarter":  "bool",
  "swiftCode":      "string",
  "branches": [
    { 
        "address":          "string",
        "bankName":         "string",
        "countryISO2":      "string",
        "isHeadquarter":    "bool",
        "swiftCode":        "string"
    }
  ]
}
```
#### Response (Branch SWIFT code):
```json
{
  "address":        "string",
  "bankName":       "string",
  "countryISO2":    "string",
  "countryName":    "string",
  "isHeadquarter":  "bool",
  "swiftCode":      "string"
}
```

### Retrieve all SWIFT codes for a country
**GET** `/v1/swift-codes/country/{countryISO2code}`
```json
{
    "countryISO2": "string",
    "countryName": "string",
    "swiftCodes": [
        {
            "address": "string",
    		 "bankName": "string",
    		 "countryISO2": "string",
    		 "isHeadquarter": "bool",
    		 "swiftCode": "string"
        },
        {
            "address": "string",
    		 "bankName": "string",
    		 "countryISO2": "string",
    		 "isHeadquarter": "bool",
    		 "swiftCode": "string"
        }, . . .
    ]
}
```

### Add a new SWIFT code
**POST** `/v1/swift-codes`
#### Request Body:
```json
{
  "address":        "string",
  "bankName":       "string",
  "countryISO2":    "string",
  "countryName":    "string",
  "isHeadquarter":  "bool",
  "swiftCode":      "string"
}
```
### Delete a SWIFT code
**DELETE** `/v1/swift-codes/{swift-code}`