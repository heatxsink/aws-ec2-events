aws-ec2-events
==============

An extremely simplistic cli tool to emit amazon ec2 events. This might be useful with this whole aws reboot "gate".

usage
-----

	$ make
	$ cd bin
	#
	## To send alert to yourself when there's an event
	#
	$ ./aws-ec2-events \
			-key=<aws access key id> \
			-secret=<aws secret access key> \
			-alert_email=<email where you will recieve alerts> \
			-imap_username=<imap username> \
			-imap_password=<imap password>
	#
	## OR just take a look at what events your instnaces may have ...
	#
	$ ./aws-ec2-events -key=<aws access key id> -secret=<aws secret access key>