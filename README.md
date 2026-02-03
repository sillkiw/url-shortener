# URL shortener REST API
1. Create Short Link
    Endpoint: `POST /url`
    Payload: 
    ```json
    {"url": "https://example.com/path", "alias": "custom-id"} (alias is optional)
    ```
    Behavior: Validates the URL, generates a random alias if not provided, and saves to DB.
    Returns: 201 Created with the short link data.

2. Redirect to Original
    Endpoint: `GET /{alias}`
    Behavior: Look up the original URL by alias and perform a 302 Found redirect.
    Returns: 404 Not Found if the alias does not exist.

3. Delete Link (Admin)
    Endpoint: `DELETE /admin/{alias}`
    Authentication: Basic Auth required.
    Behavior: Removes the entry from the database.
    Returns: 204 No Content or 401 Unauthorized.