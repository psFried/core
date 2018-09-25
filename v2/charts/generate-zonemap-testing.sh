#!/bin/bash

# Write a `zonemap` ConfigMap which contains a `node-zone.sh` script.
# When invoked with a Node name, it outputs a random selection of zone-A or zone-B.

cat <<EOF
kind: ConfigMap
apiVersion: v1
metadata:
  name: gazette-zonemap
data:
  node-zone.sh: |
    #!/bin/sh
    # This file was generated by generate-zonemap-testing.sh.
    # Do not edit by hand.
    shuf -n 1 << EOF
    zone-A
    zone-B
    EOF
EOF
