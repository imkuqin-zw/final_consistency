#!/bin/bash

./etcdctl put zw.com/final_consistency/msg_api/app '{"name":"zw.com.final_consistency.msg_api","version":"v1","address":"", "msg_version":1}'

./etcdctl put zw.com/final_consistency/msg_api/zap '{"level":"debug","development":true}'

./etcdctl put zw.com/final_consistency/msg_api/etcd '{"addrs":["127.0.0.1:2379"]}'