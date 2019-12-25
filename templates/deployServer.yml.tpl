apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: geoipdb-server-claim
spec:
  storageClassName: geoipdb-server-localstorage
  accessModes:
    - ReadOnlyMany
  resources:
    requests:
      storage: 1Gi
---
apiVersion: v1
kind: Service
metadata:
  name: geoipserver
spec:
  selector:
    app: geoipserver
  ports:
    - name: geoip-api
      protocol: TCP
      port: LOCALHOST_PORT
      targetPort: 8080
  type: LoadBalancer
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: geoipserver
spec:
  replicas: NUM_SERVER_REPLICA
  selector:
    matchLabels:
      app: geoipserver
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: geoipserver
    spec:
      restartPolicy: Always
      containers:
        - name: geoipserver
          image: geoipserver:latest
          imagePullPolicy: Never
          command: ["/entrypoint.sh"]
          args:
            - --verbosity 6
            - /usr/local/share/GeoIP/GeoLite2-City.mmdb
          ports:
            - name: geoip-api
              containerPort: 8080
          volumeMounts:
            - name: geoipdb
              mountPath: /usr/local/share/GeoIP
              readOnly: true
            - name: entrypoint
              mountPath: /entrypoint.sh
              readOnly: true
            - name: loop
              mountPath: /loop.sh
              readOnly: true
      volumes:
        - name: geoipdb
          persistentVolumeClaim:
            claimName: geoipdb-server-claim
            readOnly: true
        - name: entrypoint
          hostPath:
            path: "SERVER_ENTRYPOINT_SCRIPT"
            type: File
        - name: loop
          hostPath:
            path: "SERVER_LOOP_SCRIPT"
            type: File
