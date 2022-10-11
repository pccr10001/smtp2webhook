SMTP to Webhook Forwarder
--
Receive Email letters and forward to `n8n.io` webhook service.

### Configuration
* `isTest: true` -> Send to `webhook-test` instead of `webhook` endpoint
* Forward `hello@example.com` to `https://n8n.example.com/webhook-test/0000000-0000-0000-0000-0000000000`
```yaml
- alias:
    - hello
  isTest: true
  host: example.com
  webhook:
    id: 00000000-0000-0000-0000-0000000000
    host: https://n8n.example.com
```

### DNS settings
* Create `A` or `CNAME` record and point to SMTP server IP
* Launch server and allow forwarding port `25` to `server:2525`
  * Default listening on `2525`
* Test with `telnet server 25` to check server is running and working
* Create `MX` record as Email hostname with priority `10` and point to record name at step 1.
* Send Email to `ALIAS@HOSTNAME`