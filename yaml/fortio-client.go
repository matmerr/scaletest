package fortio

type FortioClientDeployment struct {
	Name         string
	Namespace    string
	Replicas     int
	RequestURL   string
	RequestPort  string
	AppLabel     string
	QPS          string
	NodeSelector string
}

func (f FortioClientDeployment) GetTemplate() string {
	return fortioClientDeploymentTemplate
}

const fortioClientDeploymentTemplate = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  replicas: {{ .Replicas }}
  selector:
    matchLabels:
      role: load
      app: {{ .AppLabel }}
  strategy:
    rollingUpdate:
      maxSurge: 20%
      maxUnavailable: 20%
    type: RollingUpdate
  template:
    metadata:
      labels:
        role: load
        app: {{ .AppLabel }}
    spec:
      containers:
        - name: fortio
          args:
            - load
            - -nocatchup
            - -uniform
            - -sequential-warmup
            - -udp-timeout
            - 1500ms
            - -timeout
            - 5s
            - -c
            - "{{ .QPS }}"
            - -qps
            - "{{ .QPS }}"
            - -t
            - "0"
            - http://{{ .RequestURL }}:{{ .RequestPort }}
          image: acnpublic.azurecr.io/fortio:latest
          imagePullPolicy: IfNotPresent
          resources:
            requests:
              cpu: 10m
              memory: 50M
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
      nodeSelector:
        {{ .NodeSelector }}
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      tolerations:
        - effect: NoExecute
          key: node.kubernetes.io/not-ready
          operator: Exists
          tolerationSeconds: 900
        - effect: NoExecute
          key: node.kubernetes.io/unreachable
          operator: Exists
          tolerationSeconds: 900
        - effect: NoSchedule
          key: network-load
          operator: Equal
          value: "true"

`
