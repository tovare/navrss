

Google cloud function to convert NAVs API to RSS (https://github.com/navikt/pam-public-feed)

gcloud functions deploy RSSFeed --memory=128 --runtime go111 --trigger-http --env-vars-file .env.yaml

Where .env.yaml contains:
```
ARBEIDAPI: "<insert key>"
```

The live version can be found:
https://us-central1-nav-stillinger.cloudfunctions.net/RSSFeed 

