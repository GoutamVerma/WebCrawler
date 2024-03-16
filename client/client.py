import requests

while True:
    url = input("Enter the URL of the website to crawl: ")
    deep = input("Enter the depth of URLs (optional): ")

    endpoint = f"http://localhost:1234/crawl?url={url}"
    if deep:
        endpoint += f"&deep={deep}"

    print("Generating response...")
  
    try:
        response = requests.get(endpoint)
        if response.status_code == 200:
            print("Response received:")
            print(response.text)
        else:
            print(f"Failed to get response. Status code: {response.status_code}")
    except requests.exceptions.RequestException as e:
        print(f"An error occurred: {e}")

    retry = input("Do you want to retry? (y/n): ")
    if retry.lower() != 'y':
        break
