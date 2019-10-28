

# RSS 2.0 feed from NAV job ads public API

A google cloud function to convert NAVs API to RSS 2.0 (https://github.com/navikt/pam-public-feed)

Intended to be suitable for tools like Feedly, IFTTT, Flow.


## Link to the feed 

https://tovare.com/jobb/rss

(Example live feed)

## Supported filters

(see original API for docs)

* size
* municipal
* county

https://tovare.com/jobb/rss?size=100&municipal=BERGEN

## Configure and deploy

See the quickstart documentation https://cloud.google.com/functions/

### With development token:
gcloud functions deploy RSSFeed --memory=128 --runtime go111 --trigger-http

### With production token from NAV:
gcloud functions deploy RSSFeed --memory=128 --runtime go111 --trigger-http --env-vars-file .env.yaml

Create a filee for enviromeent configuration .env.yaml contains:
```
ARBEIDAPI: "<insert key>"
```

# Harvesting to firebase

harvestproc.go contains a pub/sub functions which stores into firebase.
The model defined in feed.go is mapped to both firebase and json results.