## Это промежуточный API для мониторинга курса евро

Это прослойка для мониторинга текущего курса евро в Альфа-Банке.

Функция параллельно раз в минуту запрашивает курс в Альфа-Банке и кладет в переменную. 

Используется два метода:

/courses - который просто повторяет полученное от Альфа-Банка

/status - для мониторинга в Kubernetes-кластере этого сервиса

Используется переменная окружения ALFA_LINK для того, чтобы задать URL, с которого надо забирать валюту.

Это обычный линк от Альфа-Банка:

https://alfabank.ru/api/v1/scrooge/currencies/alfa-rates?currencyCode.in=EUR&rateType.eq=makeCash&lastActualForDate.eq=true&clientType.eq=standardCC&date.lte=

Но так же можно добавить USD и некоторые другие валюты

### Порт по умолчанию

Сервис слушает по умолчанию порт 8090

### K8S деплоймент

Пример простого deployment для кубера:

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: alfa
  labels:
    app.kubernetes.io/name: alfa
  namespace: monitoring
spec:
  selector:
    matchLabels:
      app: alfa
  template:
    metadata:
      labels:
        app: alfa
    spec:
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
              - matchExpressions:
                  - key: "location"
                    operator: In
                    values:
                      - "russia"
      containers:
        - name: alfa
          image: aladex/alfacourses:latest
          imagePullPolicy: "IfNotPresent"
          command:
            - "/app/app"
          resources:
            requests:
              cpu: 250m
              memory: 256Mi
          env:
            - name: ALFA_LINK
              value: "https://alfabank.ru/api/v1/scrooge/currencies/alfa-rates?currencyCode.in=EUR&rateType.eq=makeCash&lastActualForDate.eq=true&clientType.eq=standardCC&date.lte="
          ports:
            - name: tcp-alfa
              containerPort: 8090
          readinessProbe:
            failureThreshold: 3
            httpGet:
              path: /status
              port: 8090
              scheme: HTTP
            initialDelaySeconds: 10
            periodSeconds: 30
            successThreshold: 1
            timeoutSeconds: 2
          livenessProbe:
            failureThreshold: 3
            initialDelaySeconds: 30
            periodSeconds: 10
            successThreshold: 1
            tcpSocket:
              port: 8090
            timeoutSeconds: 1

---
apiVersion: v1
kind: Service
metadata:
  name: alfa
  namespace: monitoring
spec:
  ports:
    - port: 8090
      protocol: TCP
      targetPort: tcp-alfa
  selector:
    app: alfa
  sessionAffinity: None
  type: ClusterIP
```