# ezadmis-deny-ns-deletion

This admission webhook will deny namespace deletion, unless the namespace has an annotation:

```
ezadmis.yankeguo.github.io/deletion-allowed: "true"
```

## Installation

**Assuming we are installing to namespace `autoops`**

1. complete [RBAC Initialization for ezadmis-install](../ezadmis-install)
2. deploy YAML resources

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: install-ezadmis-deny-ns-deletion-cfg
  namespace: autoops
data:
  config.json: |
    {
      "name": "ezadmis-deny-ns-deletion",
      "namespace": "autoops",
      "mutating": false,
      "admissionRules": [
        {
          "apiGroups": [""],
          "apiVersions": ["*"],
          "resources": ["namespaces"],
          "operations": ["DELETE"]
        }
      ],
      "sideEffect": "None",
      "image": "yankeguo/ezadmis-deny-ns-deletion"
    }
---
# Job
apiVersion: batch/v1
kind: Job
metadata:
  name: install-ezadmis-deny-ns-deletion
  namespace: autoops
spec:
  template:
    spec:
      serviceAccount: ezadmis-install
      containers:
        - name: install-ezadmis-deny-ns-deletion
          image: yankeguo/ezadmis-install
          args:
            - /ezadmis-install
            - -conf
            - /config.json
          volumeMounts:
            - name: vol-cfg
              mountPath: /config.json
              subPath: config.json
      volumes:
        - name: vol-cfg
          configMap:
            name: install-ezadmis-deny-ns-deletion-cfg
      restartPolicy: OnFailure
```

## Credits

GUO YANKE, MIT License
