# Resource Monitor
## TL;DR
If you don't want to read a long README file, the below section is enough to get you started.
```shell
$ # 1) Deploy.
$ git clone git@drt-it-github-prod-1.eng.nutanix.com:rohith-subramanyam/experimental.git
$ scp experimental/resource_mon/resource_mon nutanix@<node_ip_or_name>:~/  # Copy resource_mon over to the node you want to diagnose.
$ ssh nutanix@<node_ip_or_node>

nutanix@<node>:~$ # 2) Install.
nutanix@<node>:~$ ./resource_mon [--cluster] install  # Install to ~nutanix/bin.
nutanix@<node>:~$ rm ./resource_mon  # Delete the installer.
nutanix@<node>:~$ resource_mon --help
nutanix@<node>:~$ # 3) Run.
nutanix@<node>:~$ resource_mon [--cluster] start|status|stop|restart  # Control the service.
nutanix@<node>:~$ # or 
nutanix@<node>:~$ sudo /home/nutanix/bin/resource_mon [--cluster] start|status|stop|restart  # Run as sudo to gets stats of processes owned by all users including root.
nutanix@<node>:~$ # 4) Output.
nutanix@<node>:~$ vi data/logs/resource_mon.csv.out  # Check output in data/logs/resource_mon.csv.out.

nutanix@<node>:~$ # 5) Uninstall.
nutanix@<node>:~$ resource_mon [--cluster] uninstall
```
If you skip the `--cluster` option, the default behavior is to perform the operation only on the node.<br/><br/>
Default interval between readings is `120` seconds. If you want to adjust it pass the command-line option `--interval=m` seconds.<br/><br/>
P.S.: It might consume high CPU for a few seconds every interval seconds if you run it as sudo.

## What is it?
Resource Monitor is a tool to monitor the memory, CPU and other resources on a Nutanix CVM at a system and process level.

### System
Gets the following system-level stats:
1. **total memory:** total physical memory
2. **available memory:** the memory that can be given instantly to processes without the system going into swap
3. **free memory:** memory not being used at all (zeroed) that is readily available
4. **cpu percent:** current system-wide CPU utilization as a percentage for each CPU
5. **num_fds:** total number of file descriptors currently opened by all the processes in the system that I have access to
6. **num_threads:** total number of threads currently used by all the processes in system that I have access to

### Process
Gets the following process-level stats of all the running processes in
the system/cl that it has access to:
1. **ip:** IP address of the node in which the process is run
2. **uid:** name of the user that owns the process
3. **pid:** process ID of the process
4. **ppid:** process ID of the parent process
5. **name:** name of the process (deciphers nutanix service name from its command-line)
6. **pss:** aka `Proportional Set Size`, is the amount of memory shared with other processes, accounted in a way that the amount is divided evenly between the processes that share it. I.e., if a process has 10 MBs all to itself and 10 MBs shared with another process its PSS will be 15 MBs.
7. **uss:** aka `Unique Set Size`, this is the memory which is unique to a process and which would be freed if the process was terminated right now
8. **rss:** aka `Resident Set Size`, this is the non-swapped physical memory a process has used. It matches top's RES column.
9. **vms:** aka `Virtual Memory Size`, this is the total amount of virtual memory used by the process. It matches top's VIRT column.
10. **swap:** amount of memory that has been swapped out to disk
11. **num_fds:** number of file descriptors currently opened by the process (non cumulative)
12. **num_threads:** number of threads currently used by the process (non cumulative)
13. **num_children:** number of children of the process (non-recursive)
14. **cpu_percent:** process CPU utilization as a percentage which can also be > 100.0 in case of a process running multiple threads on different CPUs. Can be disabled with --noprocess_cpu option as it adds a delay of 0.1 seconds per process.
15. **leader:** this is set to TRUE if the process is a Nutanix service and the service leader. Can be disabled with --noleadership option.
16. **timestamp:** the epoch at which the above stats were collected

## Install
Like everything at Nutanix, it is simple and 1-click.
```shell
~ $ git clone git@drt-it-github-prod-1.eng.nutanix.com:rohith-subramanyam/experimental.git
~ $ ./experimental/resource_mon/resource_mon [--cluster] install  # Copies the file to ~nutanix/bin.
~ $ rm -rf experimental  # Delete the installer.
~ $ resource_mon --help
```

## Run
It can run in 2 modes:
### Background/Daemon
Runs as a daemon in the background, gets stats once every "interval" seconds until it is stopped and appends it to `<output_dir>/resource_mon.csv.out`.
```shell
resource_mon [--cluster] [--interval=N] [--[no]leadership] [--niceness=S] [--output_dir=/home/nutanix] start | restart
```
Check the status of the daemon and stop the daemon as shown below:
```shell
resource_mon [--cluster] status | stop
```
### Foreground process
Runs in the foreground, gets stats once every "interval" seconds "count" number of times and writes stats to a new file `output_dir/resource_mon.IP_YYYYMMDD_HHMMSS.csv.out` for each iteration of count.
```shell
resource_mon [--cluster] [--count=M] [--interval=N] [--[no]leadership] [--niceness=S] [--output_dir=/home/nutanix]
```
Examples:-
1) Run it as a foreground process and get one reading:
```shell
$ resource_mon
```
2) Run it as a foreground process and get n readings:
```shell
$ resource_mon [--count=n]
```

## Output
The output is in CSV format and is written to `output_dir/resource_mon.csv.out` when running in background mode.<br/>
When running in foreground mode, output is written to `output_dir/resource_mon.IP_YYYYMMDD_HHMMSS.csv.out`.<br/>
`<output_dir>` by default is `/home/nutanix/data/logs`.
### Scavenger
The output files are rotated by scavenger by default without needing any change in scavenger.
### Analysis
ELK (ElasticSearch Logstash Kibana) stack can be used to visualize `resource_mon` output. The output CSV files can be transformed using [this logstash config](logstash.conf), indexed in Elastic Search and visualized in Kibana.

## Uninstall
Removes the installed file.
```shell
resource_mon [--cluster] uninstall
```

## Usage
```shell
resource_mon [flags] [install | start | status | restart | stop | uninstall | version]

positional arguments:
install:   install the executable file to /home/nutanix/bin
start:     start the daemon
status:    return the pids of all the processes of the daemon
           currently running
restart:   restart the daemom
stop:      stop the daemon
uninstall: remove the file from /home/nutanix/bin
version:   print version of the program and exit

flags:
```shell
resource_mon:
  -c,--[no]cluster: Operation is run on the all the CVMs in the cluster.
    (default: 'false')
  -n,--count: Number of times to collect system and process stats. Only applicable when running in foreground process mode.
    (default: '1')
    (an integer)
  -?,--[no]help: show this help
  --[no]helpshort: show usage only for this module
  --[no]helpxml: like --help, but generates XML output
  -i,--interval: time in seconds between which the system and process stats are collected. Note that this interval is not guaranteed if each iteration takes more than interval seconds.
    (default: '120')
    (an integer)
  -l,--[no]leadership: Adds a boolean column 'leadership' which checks if the nutanix process running on the node is a leader or not.
    (default: 'true')
  -s,--niceness: Set niceness which affects process scheduling. It ranges from -20 (most favorable scheduling) to 19 (least favorable).
    (default: '10')
    (an integer)
  -o,--output_dir: Path to the directory where the output files are written
    (default: '/home/nutanix/data/logs')
  -u,--[no]process_cpu: Adds a float column 'cpu_percent' which is the process CPU utilization as a percentage which can also be > 100.0 in case of a process running multiple threads on different CPUs. Adds a delay of 0.1 seconds per process.
    (default: 'true')

util.base.log:
  --[no]debug: If True, enable DEBUG log messages.
    (default: 'false')
  --debug_trace_level: Debug trace level applies only if debug is enabled.
    (default: '0')
    (an integer)
  --[no]log_thread_id: If True, log last 8 digits thread id to log line.
    (default: 'false')
  --logfile: If specified, logfile to write to.
  --[no]logtostderr: If True, log to stderr instead of a log file.
    (default: 'false')
  --[no]use_sys_exit_on_fatal: If True, use sys.exit() rather than os._exit() to end the process on FATAL errors.
    (default: 'false')

gflags:
  --flagfile: Insert flag definitions from the given file into the command line.
    (default: '')
  --undefok: comma-separated list of flag names that it is okay to specify on the command line even if the program does not define a flag with that name. IMPORTANT: flags in this list that have arguments MUST use the --flag=value format.
    (default: '')
```
