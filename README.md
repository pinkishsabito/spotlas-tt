# Backend Developer Test - Solution

This repository contains the solution to the Backend Developer Test for Spotlas. The test consists of two tasks: a query task and an endpoint creation task.

## Task 1 - Query

The SQL queries for Task 1 are included in the `queries.txt` file. The queries perform the following tasks:

1. Return spots which have a domain with a count greater than 1.
2. Change the website field, so it only contains the domain.
3. Count how many spots have the same domain.
4. Return 3 columns: spot name, domain, and the count number for the domain.

## Task 2 - Endpoint

The endpoint creation task is implemented in Golang. The main code is provided in the `main.go` file. It sets up an HTTP server with an endpoint `/spots` that accepts the following parameters:

- `latitude`: Latitude of the center point.
- `longitude`: Longitude of the center point.
- `radius`: Radius in meters.
- `type`: Type of shape, either "circle" or "square".

The endpoint calculates the bounding box based on the center point and radius, generates example spots within the bounding box, filters the spots based on the shape (circle or square), and returns the filtered spots as a JSON response.

The code also includes a `main_test.go` file with unit tests for the `calculateBoundingBox`, `filterSpotsInCircle`, and `calculateDistance` functions.

## How to Use

To run the solution, follow these steps:

1. Clone the repository:

git clone https://github.com/username/repository.git

2. Open a terminal and navigate to the cloned repository.

3. Start the HTTP server: go run main.go

4. Send a GET request to the `/spots` endpoint with the required parameters:

http://localhost:8080/spots?latitude=40.7128&longitude=-74.0060&radius=1000&type=circle

Replace the values for `latitude`, `longitude`, `radius`, and `type` as needed.

5. The server will respond with a JSON array of spots that fall within the specified shape and radius.

## Unit Tests

To run the unit tests, use the following command: go test

The tests cover the `calculateBoundingBox`, `filterSpotsInCircle`, and `calculateDistance` functions.

## Project Structure

The project structure is as follows:

├── main.go # Main code implementing the HTTP server and endpoint

├── main_test.go # Unit tests for the main code

├── queries.txt # SQL queries for Task 1

├── README.md # Project documentation (you are here)

## Contact

If you have any questions or need further assistance, please contact me at [wasil@spotlas.com](mailto:wasil@spotlas.com).

Thank you!
