# GitHub Commit Tracker

GitHub Commit Tracker is an application designed to fetch and store commit data from GitHub repositories. It periodically checks for new commits, handles paginated responses, and stores the data in a PostgreSQL database. It also features functionality to fetch and display the top N commit authors based on commit counts.

## Table of Contents

- [GitHub Commit Tracker](#github-commit-tracker)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Installation](#installation)
    - [Prerequisites](#prerequisites)
    - [Installing GitHub Commit Tracker](#installing-github-commit-tracker)
  - [Configuration](#configuration)
  - [Usage](#usage)
    - [Fetching and Storing Commits](#fetching-and-storing-commits)
    - [Top Commit Author](#top-commit-author)
  - [Error Handling](#error-handling)
  - [Contributing](#contributing)
  - [License](#license)

## Features

- **Fetching**: The application fetches new commit data from GitHub repositories.
- **Periodic Fetching**: The application periodically fetches new commit data from GitHub repositories.
- **Paginated Responses**: The application handles paginated responses from the GitHub API.
- **Retry Mechanism**: The application implements a retry mechanism to handle rate limiting using exponential backoff.
- **Commit Data Storage**: The application stores commit data in a PostgreSQL database.
- **Top Commit Author**: The application fetches and displays the top N commit authors based on commit counts.

## Installation

### Prerequisites

Before installing GitHub Commit Tracker, ensure that you have the following prerequisites installed:

- Go 1.22 or higher
- PostgreSQL 15 or higher
- Redis 7.0 or higher

### Installing GitHub Commit Tracker

To install GitHub Commit Tracker, follow these steps:

Clone the GitHub Commit Tracker repository:

```bash
git clone https://github.com/github-commit-tracker/github-commit-tracker.git
```

Change into the GitHub Commit Tracker directory:

```bash
cd github-commit-tracker
```

Install Dependencies:

```bash
go mod tidy
```

## Configuration

GitHub Commit Tracker requires the following environment variables to be set:

- `DATABASE`: The name of the PostgreSQL database.
- `DB_USER`: The username of the PostgreSQL user.
- `DB_PASSWORD`: The password of the PostgreSQL user.
- `DB_HOST`: The hostname or IP address of the PostgreSQL server.
- `DB_PORT`: The port number of the PostgreSQL server.
- `REDIS_HOST`: The hostname or IP address of the Redis server.
- `REDIS_PORT`: The port number of the Redis server.
- `DEFAULT_REPO`: The default repository to track on application startup.
- `DEFAULT_OWNER`: The default owner of the default repository.
- `DEFAULT_PERIOD`: The default period to fetch commit data on application startup.
- `DEFAULT_PAGE_SIZE`: The default page size to fetch commit data on application startup.

## Usage

To use GitHub Commit Tracker, follow these steps:

1. Start the application:

```bash
go run main.go
```

The application will start fetching commit data from the specified repository in the .env file and storing it in the PostgreSQL database.

### Fetching and Storing Commits

The application periodically fetches new commits from the configured GitHub repository that have been indexed and stores them in the PostgreSQL database. By default, this occurs every hour.

To fetch and store commits for a new repository, you can make a GET request to the `/repositories/get/:owner/:repo` endpoint.

```curl
curl -X GET "http://localhost:8080/repositories/get/github/github-commit-tracker"
```

This will fetch and store the commits for the `github-commit-tracker` repository owned by the `github` user.

### Top Commit Author

You can retrieve the top N commit authors using the GetTopCommitAuthors function.

```curl
curl -X GET "http://localhost:8080/commits/top-authors?n=10"
```

This will retrieve the top 10 commit authors based on the commit counts.

## Error Handling

If an error occurs during the execution of the application, it will return a 500 status code with a JSON response containing an error message. Application errors are stored in the `error` file in the logs directory.

## Contributing

Contributions are welcome! If you have any suggestions or improvements, please open an issue or submit a pull request.

1. Fork the repository.
2. Create a new branch for your changes.
3. Make your changes and commit them with a descriptive commit message.
4. Push your changes to your forked repository.
5. Create a pull request to the main repository.

## License

GitHub Commit Tracker is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.
