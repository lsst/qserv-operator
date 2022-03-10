USE qservReplica;

-----------------------------------------------------------
-- Preload configuration parameters for testing purposes --
-----------------------------------------------------------

-- WorkerDN: {{.WorkerDN}}, WorkerReplicas: {{.WorkerReplicas}}

{{- range $val := Iterate .WorkerReplicas}}
{{$workerId := print $.WorkerDN "-" $val}}
{{$workerFQDN := print $workerId "." $.WorkerDN}}
INSERT INTO `config_worker` VALUES ('{{$workerId}}', 1, 0, '{{$workerFQDN}}', NULL,
                                    '{{$workerFQDN}}', NULL, NULL,
                                    '{{$workerFQDN}}', NULL, NULL,
                                    '{{$workerFQDN}}', NULL, NULL,
                                    '{{$workerFQDN}}', NULL, NULL);
{{- end}}
