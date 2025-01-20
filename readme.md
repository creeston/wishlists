<div align="center">
    <img src="views/static/icons/logo.png" height="64px" />
    <h1 align="center">Wish! - A Wishlist Creation App</h1>
</div>

[![Azure Demo](https://img.shields.io/badge/demo-online-008000.svg)](https://wishlists-app.azurewebsites.net)

## Overview

**Wish!** is a simple yet fully functional wishlist creation app, built as a learning project to explore the Go programming language and experiment with lightweight front-end frameworks such as HTMX and Alpine.js. The app is deployed on Azure and available for demo purposes.

### Features

- **Create and Share Wishlists**: Easily create wishlists and share them with others.
- **Reserve Items**: Visitors can mark items as "reserved," ensuring no one else can reserve them.
- **Immutable Reserved Items**: Reserved items cannot be deleted or modified by the owner, ensuring fairness.

### Technologies Used

- **Backend**: [Go](https://go.dev/)
- **Frontend**: [HTMX](https://htmx.org/), [Alpine.js](https://alpinejs.dev/)
- **Database**: [SQLite](https://www.sqlite.org/)
- **CSS Framework**: [Pico CSS](https://picocss.com/)

---

## Development

### Pre-requisites

Before running the app locally, make sure you have the following installed:

- **Go** (version 1.23 or higher)
- **Air** (for live reloading)
- **Docker** (optional, for containerized development)
- **dlv** (optional, for debugging)

---

### Running the App

#### Using Docker

To run the app with Docker, execute:

```bash
docker compose up -d
```

#### Running Locally in Terminal

Use the air tool for live reloading:

```bash
air
```

By default, the app will run on port 8000 and will be accessible at http://localhost:8000.


### Extracting Texts for Translation

To extract texts for translation into supported languages:

```bash
gotext -srclang=en-GB update -out="catalog.go" -lang="en-GB,pl-PL,ru-RU,be-BY" creeston/lists/internal/handlers
```

### Building a Docker Image

To build the Docker image:

```bash
docker build -t wishlist-app .
```

### Deploying to Azure

This project is deployed on Azure App Service for demo purposes. To deploy it yourself:

1. Run the Terraform scripts located in the infrastructure directory to create the required Azure resources.
2. The application will be automatically deployed to Azure App Service from DockerHub.

#### Database Configuration in Azure
The SQLite3 database is stored in an Azure File Share, mounted to the Azure Web App. However, SQLite's default configuration can cause "Database locked" errors. To resolve this:

1. Download the database file from Azure File Share.

2. Open the database locally and execute:

```sql
PRAGMA journal_mode = WAL;
```

3. Upload the updated database file back to Azure File Share.


> Note: Using SQLite3 in production is not recommended, but it is used here for simplicity.

