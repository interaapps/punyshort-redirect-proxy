# Punyshort Redirection Proxy

```bash
docker run -p 80:80 \
  -e PUNYSHORT_BASE_URL='https://api.punyshort.ga' \
  -e PUNYSHORT_ERROR_URL='https://punyshort.ga/error-page' \
  -e PUNYSHORT_KEY='xxx' \
  -e PUNYSHORT_IP_FORWARDING='true' \
  interaapps/punyshort-redirect-proxy
```

## Environment Variables
- PUNYSHORT_BASE_URL
- PUNYSHORT_KEY
- PUNYSHORT_IP_FORWARDING - Allow x-forwarded-for
- PUNYSHORT_ERROR_URL
- PUNYSHORT_USE_SSL: !! Experimental !! - Generating SSL Certificate with Let's Encrypt automatically