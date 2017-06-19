package slinga

import (
	log "github.com/Sirupsen/logrus"
	"os"
)

// AptomiOject represents an aptomi entity, which gets stored in aptomi DB
type AptomiOject string

const (
	/*
		The following objects can be added to Aptomi
	*/

	// TypeCluster is k8s cluster or any other cluster
	TypeCluster AptomiOject = "cluster"

	// TypeService is service definitions
	TypeService AptomiOject = "service"

	// TypeContext is how service gets allocated
	TypeContext AptomiOject = "context"

	// TypeRules is global rules of the land
	TypeRules AptomiOject = "rules"

	// TypeDependencies is who requested what
	TypeDependencies AptomiOject = "dependencies"

	/*
		The following objects must be configured to point to external resources
	*/

	// TypeUsers is where users are stored (later in AD and LDAP)
	TypeUsers AptomiOject = "users"

	// TypeSecrets is where secret tokens are stored (later in Hashicorp Vault)
	TypeSecrets AptomiOject = "secrets"

	// TypeCharts is where binary charts/images are stored (later in external repo)
	TypeCharts AptomiOject = "charts"

	/*
		The following objects are generated by aptomi during or after dependency resolution via policy
	*/

	// TypePolicyResolution holds usage data for components/dependencies
	TypePolicyResolution AptomiOject = "resolution/usage"

	// TypeLogs contains debug logs
	TypeLogs AptomiOject = "resolution/logs"

	// TypeGraphics contains images generated by graphviz
	TypeGraphics AptomiOject = "resolution/graphics"
)

// AptomiObjectsCanBeModified contains a map of all objects which can be added to aptomi policy
//var AptomiObjectsCanBeModified = map[string]AptomiOject{
//	"cluster":      Clusters,
//	"service":      Services,
//	"context":      Contexts,
//	"rules":        Rules,
//	"dependencies": Dependencies,
//	"users":        Users,
//	"secrets":      Secrets,
//	"chart":        Charts,
//}

// Return aptomi DB directory
func getAptomiEnvVarAsDir(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		debug.WithFields(log.Fields{
			"var": key,
		}).Fatal("Environment variable is not present. Must point to a directory")
	}
	if stat, err := os.Stat(value); err != nil || !stat.IsDir() {
		debug.WithFields(log.Fields{
			"var":       key,
			"directory": value,
			"error":     err,
		}).Fatal("Directory doesn't exist or error encountered")
	}
	return value
}

// GetAptomiBaseDir returns base directory, i.e. the value of APTOMI_DB environment variable
func GetAptomiBaseDir() string {
	return getAptomiEnvVarAsDir("APTOMI_DB")
}

// GetAptomiObjectFilePatternYaml returns file pattern for aptomi objects (so they can be loaded from those files)
func GetAptomiObjectFilePatternYaml(baseDir string, aptomiObject AptomiOject) string {
	return baseDir + "/**/" + string(aptomiObject) + "*.yaml"
}

// GetAptomiObjectFilePatternTgz returns file pattern for tgz objects (so they can be loaded from those files)
func GetAptomiObjectFilePatternTgz(baseDir string, aptomiObject AptomiOject, chartName string) string {
	return baseDir + "/**/" + chartName + ".tgz"
}

// GetAptomiObjectWriteDir returns directory for aptomi objects (so they can be saved to this directory)
func GetAptomiObjectWriteDir(baseDir string, aptomiObject AptomiOject) string {
	return baseDir + "/" + string(aptomiObject)
}
