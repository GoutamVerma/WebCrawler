# Web Crawler Service

This project implements a simple web crawler service that crawls a specified domain and returns a sitemap in a tree-like structure. The server receives requests from clients to crawl a URL and responds with the sitemap.

## Features

- Crawls a specified domain, limited to one domain only (e.g., redhat.com/foo/bar will crawl all pages within redhat.com).
- Does not follow external links (e.g., to Facebook and Twitter).
- Returns a sitemap in a tree-like structure.

## Server API Endpoints
- Get SiteMap
 - Endpoint : `/crawl`
 - Method: `GET` 
 - Parameter :
    - `url` : url of site to crawl, example: `http://google.com`
    - `deep` (Optional) : number of pages or links 
 - Example : `curl -X GET 'http://localhost:1234/crawl?url=http://google.com&deep=50'`
 - Response : 
 ```
    google.com
    -intl
        -en
        -about.html
        -locations
        -policies
            -privacy
            -terms
        -products
        -stories
        -hi_in
        
    .... to be continue
 ```

## How to run Server

### Using Source Code (Recommended)

1. Clone the Github repository:
```
    $ git clone https://github.com/GoutamVerma/webcrawler.git
    $ cd webcrawler
```
2. Run the following command to start the server
```
    $ go run app/main.go
```

### Using the Dockerfile

To use the Dockerfile in this project, follow these steps:

1. Make sure you have the Docker installed on your system.
2. Build the Docker image using the Dockerfile, run the following command in the terminal:
```
    $ docker build -t webcrawler:latest .
```
3. Run the docker image.
```
    $ docker run -p 1234:1234 webcrawler:latest

    or 

    $ docker run -p 1234:1234 goutamverma/webcrawler:latest
```

Now you should be able to access your application running inside the Docker container at `http://localhost:1234/crawl`.

## How to run client

## Prerequiste
 - Python3 
 - requests package

## Steps to run client

1. To run client, navigate the current directory to client folder inside of root directory.
```
    cd client
```
2. Client is consist of single python script that will help you to make a GET request, run the script using following command.
```
    python3 client.py
```
3. It will ask you to enter the website URL  and number of pages to crawl. Enter those details and press `Enter`.


*Note: Make sure to run server before running the client.*


