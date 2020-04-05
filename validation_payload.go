package main

const orgJSON = `{
	"metadata": {
	   "guid": "7fbbd854-8851-452c-82fb-0ff12fde5e0f",
	   "url": "/v2/organizations/7fbbd854-8851-452c-82fb-0ff12fde5e0f",
	   "created_at": "2019-11-10T00:44:46Z",
	   "updated_at": "2019-11-10T00:44:46Z"
	},
	"entity": {
	   "name": "test",
	   "billing_enabled": false,
	   "quota_definition_guid": "af19e86c-05e6-444d-8fee-952866b782bc",
	   "status": "active",
	   "default_isolation_segment_guid": null,
	   "quota_definition_url": "/v2/quota_definitions/af19e86c-05e6-444d-8fee-952866b782bc",
	   "spaces_url": "/v2/organizations/7fbbd854-8851-452c-82fb-0ff12fde5e0f/spaces",
	   "domains_url": "/v2/organizations/7fbbd854-8851-452c-82fb-0ff12fde5e0f/domains",
	   "private_domains_url": "/v2/organizations/7fbbd854-8851-452c-82fb-0ff12fde5e0f/private_domains",
	   "users_url": "/v2/organizations/7fbbd854-8851-452c-82fb-0ff12fde5e0f/users",
	   "managers_url": "/v2/organizations/7fbbd854-8851-452c-82fb-0ff12fde5e0f/managers",
	   "billing_managers_url": "/v2/organizations/7fbbd854-8851-452c-82fb-0ff12fde5e0f/billing_managers",
	   "auditors_url": "/v2/organizations/7fbbd854-8851-452c-82fb-0ff12fde5e0f/auditors",
	   "app_events_url": "/v2/organizations/7fbbd854-8851-452c-82fb-0ff12fde5e0f/app_events",
	   "space_quota_definitions_url": "/v2/organizations/7fbbd854-8851-452c-82fb-0ff12fde5e0f/space_quota_definitions"
	}
 }`

const spaceJSON = `{
	"metadata": {
	   "guid": "be3b78dc-fdb2-4a99-90c8-ff52839d3214",
	   "url": "/v2/spaces/be3b78dc-fdb2-4a99-90c8-ff52839d3214",
	   "created_at": "2019-11-16T17:48:37Z",
	   "updated_at": "2019-11-16T17:48:37Z"
	},
	"entity": {
	   "name": "test",
	   "organization_guid": "7fbbd854-8851-452c-82fb-0ff12fde5e0f",
	   "space_quota_definition_guid": null,
	   "isolation_segment_guid": null,
	   "allow_ssh": true,
	   "organization_url": "/v2/organizations/7fbbd854-8851-452c-82fb-0ff12fde5e0f",
	   "developers_url": "/v2/spaces/be3b78dc-fdb2-4a99-90c8-ff52839d3214/developers",
	   "managers_url": "/v2/spaces/be3b78dc-fdb2-4a99-90c8-ff52839d3214/managers",
	   "auditors_url": "/v2/spaces/be3b78dc-fdb2-4a99-90c8-ff52839d3214/auditors",
	   "apps_url": "/v2/spaces/be3b78dc-fdb2-4a99-90c8-ff52839d3214/apps",
	   "routes_url": "/v2/spaces/be3b78dc-fdb2-4a99-90c8-ff52839d3214/routes",
	   "domains_url": "/v2/spaces/be3b78dc-fdb2-4a99-90c8-ff52839d3214/domains",
	   "service_instances_url": "/v2/spaces/be3b78dc-fdb2-4a99-90c8-ff52839d3214/service_instances",
	   "app_events_url": "/v2/spaces/be3b78dc-fdb2-4a99-90c8-ff52839d3214/app_events",
	   "events_url": "/v2/spaces/be3b78dc-fdb2-4a99-90c8-ff52839d3214/events",
	   "security_groups_url": "/v2/spaces/be3b78dc-fdb2-4a99-90c8-ff52839d3214/security_groups",
	   "staging_security_groups_url": "/v2/spaces/be3b78dc-fdb2-4a99-90c8-ff52839d3214/staging_security_groups"
	}
 }`

const applicationJSON = `{
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
