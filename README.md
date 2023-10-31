# Real-time Chat Application with Bot Integration

This project is a real-time chat application that allows users to communicate instantaneously and includes a decoupled bot for automated message broadcasting. The application is built using Go for the backend, React for the frontend, and Docker for easy setup and deployment.

## Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- Go 1.19
- Node.js 14.17.0
- npm 6.14.13

### Installation

1. Clone the repo
   ```sh
   git clone git@github.com:esirangelomub/go-chat-application.git
    ```
   
2. Build and Run the Docker image

   Run the Docker container
      ```sh
      docker-compose up --build -d  
      ```

3. Instal backend dependencies
   ```sh
   go mod download
   ```
   or
    ```sh
    go get -d ./...
    ```
   or
    ```
   go mod tidy
    ```

4. Install frontend dependencies
    ```sh
   cd frontend
   npm install
   ```
   or
   ```sh
   cd frontend
   yarn install
   ```
   
5. Run the application
    Backend server
   ```sh
   cd cmd/server
    go run main.go
    ```
   Backend chat
   ```sh
   cd cmd/server
    go run main.go
    ```
    Frontend
6. Run the application
    ```sh
   cd frontend
   npm start
    ```
   or
    ```sh
   cd frontend
   yarn start
    ```
   
## Usage

1. Open your browser and go to http://localhost:3000
2. Register a new account
3. Login with your new account
4. Start chatting with other users
5. To test the bot, send a message with the following format: /stock=stock_code

Note: You'll need to open two different browsers (e.g. Chrome and Firefox)

## Contact

Feel free to reach out if you have any questions or would like to discuss this project further.

- **Email:** [eduardo.sirangelo@gmail.com](mailto:eduardo.sirangelo@gmail.com)
- **LinkedIn:** [Eduardo Sirangelo](https://www.linkedin.com/in/eduardosirangelo/?locale=en_US)

## License

Distributed under the MIT License. See `LICENSE` for more information.


   
   
   

