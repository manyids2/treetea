clean:
	rm -r /tmp/task

unset-fish:
	# For fish shell
	echo 'set -gx TASKDATA'
	echo 'set -gx TASKRC'

unset-bash:
	# For fish shell
	echo 'unset gx TASKDATA'
	echo 'unset gx TASKRC'

setup-fish:
	rm -r /tmp/task || true
	mkdir /tmp/task
	echo 'set -gx TASKDATA /tmp/task'
	echo 'set -gx TASKRC /tmp/task/taskrc'

setup-bash:
	rm -r /tmp/task || true
	mkdir /tmp/task
	echo 'export TASKDATA=/tmp/task'
	echo 'export TASKRC=/tmp/task/taskrc'
