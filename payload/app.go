package payload

const ApplicationJSON = `{
	"metadata": {
	   "guid": "a2cffb66-3524-406d-bc22-169028b1e7fc",
	   "url": "/v2/apps/a2cffb66-3524-406d-bc22-169028b1e7fc",
	   "created_at": "2020-03-28T17:54:35Z",
	   "updated_at": "2020-03-28T18:00:07Z"
	},
	"entity": {
	   "name": "tester",
	   "production": false,
	   "space_guid": "be3b78dc-fdb2-4a99-90c8-ff52839d3214",
	   "stack_guid": "9725be32-44a2-4591-9c0e-1e6cfbba0f00",
	   "buildpack": "java_buildpack",
	   "detected_buildpack": "",
	   "detected_buildpack_guid": "efb2bdf6-8ba2-41fd-a342-2c613bdcc933",
	   "environment_json": {
		  "JBP_CONFIG_OPEN_JDK_JRE": "{ jre: { version: 11.+ } }"
	   },
	   "memory": 768,
	   "instances": 1,
	   "disk_quota": 1024,
	   "state": "STARTED",
	   "version": "8d18da0b-3251-42eb-90b6-3923e5722921",
	   "command": null,
	   "console": false,
	   "debug": null,
	   "staging_task_id": "e1ae3957-8e4b-45ed-aff6-0c10ad7ab6f4",
	   "package_state": "STAGED",
	   "health_check_type": "port",
	   "health_check_timeout": null,
	   "health_check_http_endpoint": "",
	   "staging_failed_reason": null,
	   "staging_failed_description": null,
	   "diego": true,
	   "docker_image": null,
	   "docker_credentials": {
		  "username": null,
		  "password": null
	   },
	   "package_updated_at": "2020-03-28T18:00:00Z",
	   "detected_start_command": "JAVA_OPTS=\"-agentpath:$PWD/.java-buildpack/open_jdk_jre/bin/jvmkill-1.16.0_RELEASE=printHeapHistogram=1 -Djava.io.tmpdir=$TMPDIR -XX:ActiveProcessorCount=$(nproc) -Djava.ext.dirs= -Djava.security.properties=$PWD/.java-buildpack/java_security/java.security $JAVA_OPTS\" && CALCULATED_MEMORY=$($PWD/.java-buildpack/open_jdk_jre/bin/java-buildpack-memory-calculator-3.13.0_RELEASE -totMemory=$MEMORY_LIMIT -loadedClasses=19798 -poolType=metaspace -stackThreads=250 -vmOptions=\"$JAVA_OPTS\") && echo JVM Memory Configuration: $CALCULATED_MEMORY && JAVA_OPTS=\"$JAVA_OPTS $CALCULATED_MEMORY\" && MALLOC_ARENA_MAX=2 SERVER_PORT=$PORT eval exec $PWD/.java-buildpack/open_jdk_jre/bin/java $JAVA_OPTS -cp $PWD/.:$PWD/.java-buildpack/container_security_provider/container_security_provider-1.16.0_RELEASE.jar org.springframework.boot.loader.JarLauncher",
	   "enable_ssh": true,
	   "ports": [
		  8080
	   ],
	   "space_url": "/v2/spaces/be3b78dc-fdb2-4a99-90c8-ff52839d3214",
	   "stack_url": "/v2/stacks/9725be32-44a2-4591-9c0e-1e6cfbba0f00",
	   "routes_url": "/v2/apps/a2cffb66-3524-406d-bc22-169028b1e7fc/routes",
	   "events_url": "/v2/apps/a2cffb66-3524-406d-bc22-169028b1e7fc/events",
	   "service_bindings_url": "/v2/apps/a2cffb66-3524-406d-bc22-169028b1e7fc/service_bindings",
	   "route_mappings_url": "/v2/apps/a2cffb66-3524-406d-bc22-169028b1e7fc/route_mappings"
	}
 }`
