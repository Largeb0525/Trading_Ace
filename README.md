# Trading_Ace

Trading_Ace is an application designed to integrate Uniswap V2 swap activities, providing rewards based on user trading volume according to task settings. It supports multiple pools and multiple tasks within the same pool. This project is developed in Go and relies on Infura’s Ethereum node service.

## Requirements

- [Go 1.22](https://golang.org/dl/)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Infura account](https://www.infura.io/zh)

## Local Setup

1. **Clone the repository:**

    ```bash
    git clone https://github.com/Largeb0525/Trading_Ace.git
    cd Trading_Ace
    ```

2. **Configure settings:**

    - Edit `config/config.toml` and add your Infura API key under the `[infura]` section:   

        ```toml
        [infura]
        api_key = "YOUR_INFURA_API_KEY"
        ```

3. **Start services:**

    Use Docker Compose to start the application and database:

    ```bash
    docker-compose up --build
    ```
    The application will run on local port 8080.

## API Examples

The following are examples of available API endpoints:

### 1. **Create Campaign**

- **Endpoint:** `POST /Campaign`
- **Payload Parameters:**
    - `name` (string, required): Name of the campaign.
    - `poolAddress` (string, required): Ethereum address of Uniswap V2 pool.
    - `startAt` (int, required): Unix timestamp for when the campaign should start.
    - `onboardingReward` (float, required): Reward amount for the onboarding task.
    - `onboardingThreshold` (float, required): Minimum swap amount in USDC to qualify for the onboarding reward.
    - `pointPool` (float, required): Total points available for distribution in the share pool task.
    - `schedule` (string, required): Interval for each campaign round, formatted as "5m", "1h", "24h", etc.
    - `round` (int, required): Number of rounds to repeat the campaign task.

- **Example Request (using `curl`):**

    ```bash
    curl --location 'localhost:8080/Campaign' \
    --header 'Content-Type: application/json' \
    --data '{
        "name":"test",
        "poolAddress":"0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc",
        "startAt":1731400000,
        "onboardingReward":100,
        "onboardingThreshold":1000,
        "pointPool":10000,
        "schedule":"168h",
        "round":4
    }'
    ```

### 2. **Get User Task Status**

- **Endpoint:** `GET /user/task/status`
- **Query Parameters:**
    - `userAddress` (required): Ethereum address of the user to retrieve task statuses for
    - `userID` (optional): The ID of the user to retrieve task statuses for.

- **Example Request (using `curl`):**

    ```bash
    curl --location --request GET 'localhost:8080/user/task/status?userAddress=0xa69babef1ca67a37ffaf7a485dfff3382056e78c' \
    --header 'Content-Type: application/json'
    ```

- **Example Response:**

    ```bash
    {
        "campaigns": [
            {
                "campaignId": 1,
                "name": "test",
                "poolAddress": "0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc",
                "startTime": 1731345250,
                "endTime": 1731346450,
                "tasks": [
                    {
                        "taskId": 1,
                        "type": "onboarding",
                        "description": "",
                        "completed": true,
                        "amount": 500,
                        "points": 500,
                        "startTime": 1731345250,
                        "endTime": 1731346450
                    },
                    {
                        "taskId": 2,
                        "type": "share_pool",
                        "description": "Round 1",
                        "completed": true,
                        "amount": 9005.202944,
                        "points": 1678.021288,
                        "startTime": 1731345250,
                        "endTime": 1731345550
                    },
                    {
                        "taskId": 3,
                        "type": "share_pool",
                        "description": "Round 2",
                        "completed": true,
                        "amount": 1021.442396,
                        "points": 4376.313797,
                        "startTime": 1731345550,
                        "endTime": 1731345850
                    }
                ]
            }
        ]
    }
    ```

### 3. **Get User Points History**

This endpoint retrieves the points history for a user based on their `userAddress`. It provides a record of points earned through task completions across campaigns.

- **Endpoint:** `GET /user/points`
- **Query Parameters:**
    - `userAddress` (string, required): Ethereum address of the user to retrieve task statuses for
    - `userID` (optional): The ID of the user to retrieve task statuses for.

- **Example Request (using `curl`):**

    ```bash
    curl --location --request GET 'localhost:8080/user/points?userAddress=0xa69babef1ca67a37ffaf7a485dfff3382056e78c' \
    --header 'Content-Type: application/json' \
    --data '{
        "name":"test",
        "poolAddress":"0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc",
        "startAt":1731248500,
        "onboardingReward":150,
        "onboardingThreshold":300,
        "pointPool":10000,
        "schedule":"15m",
        "round":4
    }'
    ```

- **Example Response:**

    ```bash
    {
        "pointsHistory": [
            {
                "campaignId": 5,
                "campaignName": "test",
                "poolAddress": "0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc",
                "taskId": 21,
                "taskType": "onboarding",
                "description": "",
                "points": 150,
                "timestamp": "2024-11-11T17:18:02+08:00"
            },
            {
                "campaignId": 5,
                "campaignName": "test",
                "poolAddress": "0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc",
                "taskId": 24,
                "taskType": "share_pool",
                "description": "Round 3",
                "points": 7103.584443,
                "timestamp": "2024-11-11T17:24:20+08:00"
            }
        ],
        "total": 7253.584443
    }
    ```

## TODO List

1. **Database Design Improvement**: Replace `SERIAL` ID columns with `UUID` for improved scalability and uniqueness across distributed systems.
2. **Unit Testing**: Implement unit tests using the `gomonkey` library to mock calls to third-party APIs.
3. **Caching with Redis**: Use Redis to cache results for repeated queries and reduce redundant calls.
4. **Expand API**: Develop additional UPDATE & DELETE APIs for managing campaign data and leaderboard based on points of distributed tasks.
5. **Error Handling**: Improve error handling for Infura websocket & API requests, ensuring robustness in case of API failures.
6. **Enhanced Logging**: Add detailed logging information for better traceability and debugging.

