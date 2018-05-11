# webclient
A stable http requests generator.


# Usage

```bash
./_output/webclient --rps 2 --target=http://music.default:8080/cpuwork.php/?cpu=20 
```


# Deployed in Kubernetes
Following example yaml file will create two pods, sending POST requests to target:
   `http://music.default:8080/cpuwork.php/?cpu=20`
Each Pod will try to send 2 requests per second.

   
   
   

```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: httpclient
  namespace: default
  labels:
    purpose: generate-http-load
spec:
  replicas: 2
  selector:
    matchLabels:
      app: httpclient
  template:
    metadata:
      labels:
        app: httpclient
    spec:
      serviceAccount: default
      containers:
      - name: httpclient
        image: beekman9527/webclient:v1
        imagePullPolicy: IfNotPresent
        args:
        - --v=3
        - --threadNum=6
        - --logtostderr
        - --target=http://music.default:8080/cpuwork.php/?cpu=20
        - --rps=2
```
