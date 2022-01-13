package packageManager

type PackageManager interface {
	IsPackageDependencyRequirementFile(path string) bool
	GetPackageName() string
}
