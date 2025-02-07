# receipt-processor-challenge

Receipt Processor
This Go application provides two endpoints that process receipt data in memory and compute point totals based on several rules.

Overview
All of the logic for this service resides in main.go. It uses the following key components:

Data Models

Receipt struct with fields retailer, purchaseDate, purchaseTime, total, and a slice of Item.
Item struct with fields shortDescription and price.
In-Memory Storage

A global map[string]int called pointsStore that maps a UUID (id) to a points total.
Endpoints

POST /receipts/process
Accepts a JSON receipt in the request body.
Calculates the receipt’s points using the calculatePoints function.
Generates a UUID, stores the points in the in-memory map, and returns the UUID in JSON.
GET /receipts/{id}/points
Looks up the points in the in-memory map by id.
Returns the points in JSON (e.g., {"points": 35}).
Points Calculation
The function calculatePoints applies the following rules to compute a final point value for each receipt:

Alphanumeric Characters: +1 point for each [A-Za-z0-9] in the retailer name.
Round Dollar Total: +50 points if the total has no cents (e.g., 9.00).
Multiple of 0.25: +25 points if the total is divisible by 0.25.
Item Pairs: +5 points for every two items in the receipt.
Item Description Multiple of 3: For each item with a trimmed description length that’s a multiple of 3,
add ceil(item price * 0.2) points.
LLM-Specific Rule: +5 points if the total is over 10.00.
Odd Purchase Day: +6 points if the day of the month is odd.
Afternoon Purchase: +10 points if the purchase time is between 2:00pm and 4:00pm (14:00–15:59).
Server Setup

Uses Gorilla Mux for routing.
Starts an HTTP server on port 8080.
1. Running the Application
Install Dependencies (in the same folder as main.go):
go get github.com/gorilla/mux
go get github.com/google/uuid
go mod tidy

2. Run the Application:
go run main.go

3. Send Requests:
POST http://localhost:8080/receipts/process with a JSON body:
{
  "retailer": "Target",
  "purchaseDate": "2022-01-02",
  "purchaseTime": "13:13",
  "total": "1.25",
  "items": [
    {
      "shortDescription": "Pepsi - 12-oz",
      "price": "1.25"
    }
  ]
}

GET http://localhost:8080/receipts/{id}/points to retrieve points for the given id.
