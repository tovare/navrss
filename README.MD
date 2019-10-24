

# RSS 2.0 feed from NAV job ads public API

A google cloud function to convert NAVs API to RSS 2.0 (https://github.com/navikt/pam-public-feed)

Intended to be suitable for tools like Feedly, IFTTT, Flow.


## Link to the feed 
https://us-central1-nav-stillinger.cloudfunctions.net/RSSFeed

(Example live feed)

## Supported filters

(see original API for docs)

* size
* municipal
* county

https://us-central1-nav-stillinger.cloudfunctions.net/RSSFeed?size=100&municipal=BERGEN

## Configure and deploy

See the quickstart documentation https://cloud.google.com/functions/

If you don´t want to configure your PC with development tools, you can copy and paste the contents of the feed.go source file.

### With development token:
gcloud functions deploy RSSFeed --memory=128 --runtime go111 --trigger-http

### With production token from NAV:
gcloud functions deploy RSSFeed --memory=128 --runtime go111 --trigger-http --env-vars-file .env.yaml

Create a filee for enviromeent configuration .env.yaml contains:
```
ARBEIDAPI: "<insert key>"
```


