# FAQ for Amur

## How to add prometheus auth to my prometheus config?
In your scrape config, add the following:
```yaml
scrape_configs:
  - job_name: amur
    basic_auth:
      username: {YOUR USERNAME OR DEFAULT USERNAME  ADMIN}
      password: {YOUR PASSWORD}
    static_configs:
      - targets:
          - "localhost:3000/metrics" // This could be also "amur:3000/metrics" if you're running a docker container
```