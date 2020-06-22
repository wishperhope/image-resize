## Rest Api image resize

Simple http rest api that resize image based on query params, using [lilliput](https://github.com/discordapp/lilliput) to resize image

### Setup

```bash
git clone https://github.com/wishperhope/rest-api-image-resizer
cd rest-api-image-resizer
cp .env.example .env
docker-compose up -d
```

### Example Storing Image

```bash
curl --request POST \
  --url http://localhost:8080/new \
  --header 'authorization: secure_password' \
  --header 'content-type: multipart/form-data; boundary=---011000010111000001101001' \
  --form img='yourimage'
```

Will response something like :

```json
{
  "success": true,
  "urlPath": /image/uploads/648550049.jpg
}
```

### Example get image

```bash
curl --request GET \
  --url 'http://localhost:8080/image/uploads/648550049.jpg?h=1000&w=1200'
```

Make sure to not call api from client side application without refactoring the auth with necessary security.

### Replication

You can replicate the micro services if necessary, make sure storage volume is replicated across node.

```json
docker-compose up --scale image-resize=5
```
