
This is a static site hosted as a cloudflare worker.


# Technical information

wrangler.toml file is required to publish with the wrangler publish command.

    name = "site"
    type = "webpack"
    account_id = "xxx"
    workers_dev = false
    route = "tovare.com/jobb/"
    zone_id = "xxx"

    [site]
    bucket = "./public"
    entry-point = "workers-site"
