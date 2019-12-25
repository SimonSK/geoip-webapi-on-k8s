kind: StorageClass
apiVersion: storage.k8s.io/v1beta1
metadata:
  name: geoipdb-updater-localstorage
provisioner: docker.io/hostpath
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: geoipdb-updater-pv
  labels:
    type: local
spec:
  storageClassName: geoipdb-updater-localstorage
  capacity:
    storage: 1Gi
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  hostPath:
    path: "DATABASE_DIR"
---
kind: StorageClass
apiVersion: storage.k8s.io/v1beta1
metadata:
  name: geoipdb-server-localstorage
provisioner: docker.io/hostpath
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: geoipdb-server-pv
  labels:
    type: local
spec:
  storageClassName: geoipdb-server-localstorage
  capacity:
    storage: 1Gi
  accessModes:
    - ReadOnlyMany
  persistentVolumeReclaimPolicy: Retain
  hostPath:
    path: "DATABASE_DIR"
