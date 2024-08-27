# api-okta-login
Login with Okta

# commit 1

# commit 2

# commit 3

# commit 4

```go
export MONGODB_URI="mongodb://localhost:27017"
export MONGODB_DB="priceAlerts"
export OKTA_CLIENT_ID="0oagllbqkndAe2ArM5d7"
export OKTA_CLIENT_SECRET="QG10g2WJuee00WBYxeaa0Bpetbxu3mDr89u1eMvuOjy7apqVHHSGQytEXMOQCOAN"
export OKTA_ISSUER="https://dev-96634040.okta.com/oauth2/default"
export JWT_SECRET_KEY="oYUMa8rT8EHUfOzZ0U0Ul5FzMBgM0DO4"
export JWT_EXPIRED_AFTER=846000
export OKTA_NONCE="h7JmSXePP0c25ch1ieNZqfn8ABwZpEGZ"
```

# Build with docker file

```bash
docker build -t chibk/api-okta-login .
```

Run with docker file

```bash

docker run -d -p 8080:1323 --name api-okta-login \
-e MONGODB_URI="mongodb://localhost:27017" \
-e MONGODB_DB="priceAlerts" \
-e OKTA_CLIENT_ID="0oagllbqkndAe2ArM5d7" \
-e OKTA_CLIENT_SECRET="QG10g2WJuee00WBYxeaa0Bpetbxu3mDr89u1eMvuOjy7apqVHHSGQytEXMOQCOAN" \
-e OKTA_ISSUER="https://dev-96634040.okta.com/oauth2/default" \
-e JWT_SECRET_KEY="oYUMa8rT8EHUfOzZ0U0Ul5FzMBgM0DO4" \
-e JWT_EXPIRED_AFTER=846000 \
-e OKTA_NONCE="h7JmSXePP0c25ch1ieNZqfn8ABwZpEGZ" \
chibk/api-okta-login

```
