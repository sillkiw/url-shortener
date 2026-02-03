# URL Shortener REST API

## 1) Create short link

**Endpoint:** `POST /url`
**Content-Type:** `application/json`

**Request body:**

```json
{
  "url": "https://example.com/path",
  "alias": "custom-id"
}
```

* `alias` — optional. If omitted, the service generates a random one.

**Behavior:**

* Validates `url`
* Generates an alias (if not provided)
* Saves the mapping to the database

**Response:** `201 Created`

```json
{
  "status": "ok",
  "alias": "custom-id",
}
```

---

## 2) Redirect to original URL

**Endpoint:** `GET /{alias}`

**Behavior:**

* Looks up `{alias}` in the database
* Redirects to the original URL

**Response:**

* `302 Found` with `Location: <original_url>`
* `404 Not Found` if alias does not exist

---

## 3) Delete link (admin)

**Endpoint:** `DELETE /admin/{alias}`
**Auth:** Basic Auth

**Behavior:**

* Deletes `{alias}` from the database

**Response:**

* `204 No Content` on success
* `401 Unauthorized` if credentials are missing/invalid
* `404 Not Found` if alias does not exist
