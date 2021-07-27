### Webhook test.

A server that accepts POST and then a k8s manifest so you can deploy it in-cluster with credentials.

Just clone this repoistory, set the appropriate values in the values.yaml and install in your cluster.

```
helm install webhook --namespace webhook --values values.yaml --create-namespace ./chart/
```

### Configuring Anchore

Just set up the webhook endpoint as http://webhook-processor.webhook.svc.cluster.local:9000/v1/webhook

Once you do that you can watch the logs for the pod and it'll print out the payload it receives from Anchore.
