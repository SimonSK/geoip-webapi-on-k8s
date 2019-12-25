apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: geoipdb-updater-claim
spec:
  storageClassName: geoipdb-updater-localstorage
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
---
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: geoipupdate
spec:
  # MaxMind releases new database files at some time on every Tuesday.
  # Let's schedule an update every 3 hours on Tuesdays (in UTC).
  schedule: "0 */3 * * 2"
  # Newer job is preferred
  concurrencyPolicy: Replace
  failedJobsHistoryLimit: 3
  successfulJobsHistoryLimit: 5
  jobTemplate:
    spec:
      activeDeadlineSeconds: 600
      template:
        spec:
          restartPolicy: OnFailure
          containers:
            - name: geoipupdate
              image: geoipupdate:latest
              imagePullPolicy: Never
              command: ["geoipupdate"]
              args:
                - -v
                - -d
                - /usr/local/share/GeoIP
                - -f
                - /usr/local/etc/GeoIP.conf
              volumeMounts:
                - name: geoipdb
                  mountPath: /usr/local/share/GeoIP
                  readOnly: false
                - name: geoip-config
                  mountPath: /usr/local/etc/GeoIP.conf
                  readOnly: true
          volumes:
            - name: geoipdb
              persistentVolumeClaim:
                claimName: geoipdb-updater-claim
                readOnly: false
            - name: geoip-config
              hostPath:
                path: "UPDATER_CONF"
                type: File
