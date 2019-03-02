1. Make a GET on https://api.meethue.com/oauth2/auth?clientid=$CLIENT_ID&appid=$APP_ID&deviceid=$DEVICE_ID&state=$STATE&response_type=code

2. User logs in and approves access

3. A response with a code is sent to the callback url

4. Make a POST to https://api.meethue.com/oauth2/token?code=$CODE&grant_type=authorization_code

5. 401 Auth failed is expected. Will come with nonce, save it.
WWW-Authenticate: Digest realm="oauth2_client@api.meethue.com", nonce="$NONCE"

6. Proceed with Basic or Digest Auth. Digest is more secure.

7. Generate a Response that can pass the challenge-reponse handshake.
`RESPONSE = MD5(HASH1 + “:” + “NONCE” + “:” + HASH2)`
Where ```HASH1 = MD5(“CLIENTID” + “:” + “REALM” + “:” + “CLIENTSECRET”) => MD5($CLIENT_ID:oauth2_client@api.meethue.com:<clientsecret>), HASH2 = MD5(“VERB” + “:” + “PATH”) => MD5("POST:/oauth2/token")```

8. Make a POST to https:api.meethue.com/oauth2/token?code=$CODE&grant_type=authorization_code
Authorization: Digest username="$CLIENT_ID", realm="oauth2_client@api.meethue.com", nonce="$NONCE", uri="/oauth2/token", response="$RESPONSE"

9. Basic Auth is a lot simpler... But I 
Make a POST on https:api.meethue.com/oauth2/token?code=$CODE&grant_type=authorization_code
Authorization: Basic <base64(clientid:clientsecret)>

10. There's this bit on Refresh tokens, but I'll get to that later ¯\\_(ツ)_/¯

