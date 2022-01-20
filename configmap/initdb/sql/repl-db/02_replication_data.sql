USE qservReplica;

-----------------------------------------------------------
-- Preload configuration parameters for testing purposes --
-----------------------------------------------------------

-- qserv-worker-1.qserv-worker.default.svc.cluster.local {{.WorkerDn}} {{.WorkerReplicas}}
-- {{- range $val := Iterate .WorkerReplicas}}
-- INSERT INTO `config_worker` VALUES ({{$.WorkerDn}}-{{$val}});
-- {{- end}}

{{- range $val := Iterate .WorkerReplicas}}
{{$workerId := print $.WorkerDn "-" $val}}
{{$workerFqdn := print $workerId "." $.WorkerDn}}
INSERT INTO `config_worker` VALUES ('{{$workerId}}', 1, 0, '{{$workerFqdn}}', NULL,
                                    '{{$workerFqdn}}', NULL, NULL,
                                    '{{$workerFqdn}}', NULL, NULL,
                                    '{{$workerFqdn}}', NULL, NULL,
                                    '{{$workerFqdn}}', NULL, NULL);
{{- end}}
