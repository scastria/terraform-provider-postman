# Postman Provider
The Postman provider is used to interact with the Postman API.  The provider
needs to be configured with the proper credentials before it can be used.

This provider does NOT cover 100% of the Postman API.  If there is something missing
that you would like to be added, please submit an Issue in corresponding GitHub repo.

**IMPORTANT:** The Postman API does not support simultaneous writes to the same resource.  Therefore, you **MUST** use `-parallelism=1` for any apply or destroy commands.
## Example Usage
```hcl
terraform {
  required_providers {
    postman = {
      source  = "scastria/postman"
      version = "~> 0.1.0"
    }
  }
}

# Configure the Postman Provider
provider "postman" {
  api_key = "XXXX"
  num_retries = 3
  retry_delay = 30
}
```
## Argument Reference
* `api_key` - **(Required, String)** Your API key obtained via Postman UI. Can be specified via env variable `POSTMAN_API_KEY`.
* `num_retries` - **(Optional, Integer)** Number of retries for each Postman API call in case of 400-Bad Request, 429-Too Many Requests, or any 5XX status code. Can be specified via env variable `POSTMAN_NUM_RETRIES`. Default: 3.
* `retry_delay` - **(Optional, Integer)** How long to wait (in seconds) in between retries. Can be specified via env variable `POSTMAN_RETRY_DELAY`. Default: 30.
