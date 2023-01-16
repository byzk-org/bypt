package vos

type SystemType string
type PlatformType string

const (
	LINUX   SystemType = "linux"
	WINDOWS SystemType = "windows"
	MAC     SystemType = "darwin"
)

const (
	P386     PlatformType = "386"
	AMD64    PlatformType = "amd64"
	ARM      PlatformType = "arm"
	ARM64    PlatformType = "arm64"
	MIPS     PlatformType = "mips"
	MIPS64   PlatformType = "mips64"
	MIPS64le PlatformType = "mips64le"
	MIPSle   PlatformType = "mipsle"
	PPC64    PlatformType = "ppc"
	PPC64le  PlatformType = "ppc64le"
	RISCV64  PlatformType = "riscv64"
	S390X    PlatformType = "s390x"
)

var (
	SystemPlatformMap = map[SystemType][]PlatformType{
		LINUX:   {P386, AMD64, ARM, ARM64, MIPS, MIPS64, MIPS64le, MIPSle, PPC64, PPC64le, RISCV64, S390X},
		MAC:     {P386, AMD64, ARM, ARM64},
		WINDOWS: {P386, AMD64, ARM},
	}
	AllSystem = []SystemType{LINUX, MAC, WINDOWS}
)