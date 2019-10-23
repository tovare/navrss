

Google cloud function to convert NAVs API to RSS


gcloud functions deploy RSSFeed --memory=128 --runtime go111 --trigger-http --env-vars-file .env.yaml

Where .env.yaml contains:
```
ARBEIDAPI: "<insert key>"
```