
# Redditracker :Real-Time Top Post Up Votes & Top Posting Users (Coding Challenge)

## Overview
This application is designed to monitor a chosen subreddit.

### Installation


1. Clone the repository:
   ```bash
   git clone git@github.com:marcelluseasley/redditracker.git
   ```
2. To build the application, run the following command:
   ```bash
   make build
   ```
   This will compile source code and create a binary named `redditracker` in the `cmd/redditracker` directory.



4. To remove the binary, run the following command:
   ```bash
   make clean
   ```

### Configuration

1. Configure the API Token in a `.env` file (you will need to create it, since its in `.gitignore`):
   - Navigate to `cmd/redditracker/.env`.
   ```bash
    REDDIT_CLIENT_ID=''
    REDDIT_CLIENT_SECRET=''
    USER_AGENT=''
   ```

### Running the Application
 To run the application, use the following command:
   ```bash
   make run PORT=8080 SUBREDDIT=TaylorSwift
   ```
   This will run the `redditracker` binary with the specified port and subreddit. You can replace `8080` and `TaylorSwift` with the port number and subreddit you want to track.

## Usage
Once the application is running, it will begin monitoring the specified subreddit. The statistics will be reported by navigating to (fill in your port):
```bash
    http://localhost:[PORT]/
```
You can also make a curl request to get a JSON response of results:
```bash
   curl --location 'http://localhost:[PORT]/data'
```


## License
This project is licensed under the [MIT License](LICENSE).

## Acknowledgments
- [Reddit API](https://www.reddit.com/dev/api) for providing the data used by this application.
``````
