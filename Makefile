colima-start:
	colima start -c 2 -m 2 --disk 30 --root-disk 40 --ssh-port 2222 --network-address --network-mode bridged --save-config

colima-delete:
	yes | colima delete

colima-delete-data:
	yes | colima delete --data
