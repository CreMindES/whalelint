package parser

import Utils "github.com/cremindes/whalelint/utils"

func IsDebPackageManager(bin string) bool {
	return Utils.EqualsEither(bin, []string{"apt-get", "apt", "snap"})
}

func IsPythonPackageManager(bin string) bool {
	return Utils.EqualsEither(bin, []string{"pip", "conda"})
}

func IsRpmPackageManager(bin string) bool {
	return bin == "yum"
}

func IsNpmPackageManager(bin string) bool {
	return bin == "npm"
}

func IsApkPackageManager(bin string) bool {
	return bin == "npm"
}

func IsRubyPackageManager(bin string) bool {
	return bin == "gem"
}

func IsZyppPackageManager(bin string) bool {
	return bin == "zypper"
}

func IsDnfPackageManager(bin string) bool {
	return bin == "dnf"
}

func IsArchPackageManager(bin string) bool {
	// source: https://wiki.archlinux.org/index.php/AUR_helpers
	archPackageManagerSlice := []string{"pacman", "paru", "yay", "aura", "pacaur", "pakku", "pikaur", "trizen"}

	return Utils.EqualsEither(bin, archPackageManagerSlice)
}

const (
	installStr = "install"
)

func HasPackageUpdateCommand(packageManager string, bashCommand BashCommand) bool {
	switch packageManager {
	case "apt":
		return bashCommand.Bin() == "apt" && bashCommand.subCommand == "update"
	case "apt-get":
		return bashCommand.Bin() == "apt-get" && bashCommand.subCommand == "update"
	case "apk":
		hasUpdate := bashCommand.SubCommand() == "update"
		hasAddWithUpdate := bashCommand.SubCommand() == "add" &&
			Utils.SliceContains(bashCommand.OptionKeyList(), "--update")

		return bashCommand.Bin() == "apk" && (hasUpdate || hasAddWithUpdate)
	case "pip":
		return bashCommand.Bin() == "pip" &&
			Utils.SliceContains(bashCommand.OptionKeyList(), []string{"-U", "--update"})
	case "yum":
		return bashCommand.Bin() == "yum" && bashCommand.subCommand == "update"
	case "zypper":
		return bashCommand.Bin() == "zypper" && bashCommand.subCommand == "refresh"
	case "dnf":
		return bashCommand.Bin() == "dnf" && bashCommand.subCommand == "update"
	}

	return false
}

func IsDebPackageInstall(bashCommand BashCommand) bool {
	return IsDebPackageManager(bashCommand.Bin()) && bashCommand.SubCommand() == installStr
}

func IsPythonPackageInstall(bashCommand BashCommand) bool {
	return IsPythonPackageManager(bashCommand.Bin()) && bashCommand.SubCommand() == installStr
}

func IsRpmPackageInstall(bashCommand BashCommand) bool {
	return IsRpmPackageManager(bashCommand.Bin()) && bashCommand.SubCommand() == installStr
}

func IsNpmPackageInstall(bashCommand BashCommand) bool {
	return IsNpmPackageManager(bashCommand.Bin()) &&
		(bashCommand.SubCommand() == "install" || bashCommand.SubCommand() == "i")
}

func IsApkPackageInstall(bashCommand BashCommand) bool {
	return IsApkPackageManager(bashCommand.Bin()) && bashCommand.SubCommand() == "add"
}

func IsRubyPackageInstall(bashCommand BashCommand) bool {
	return IsRubyPackageManager(bashCommand.Bin()) && bashCommand.SubCommand() == installStr
}

func IsSusePackageInstall(bashCommand BashCommand) bool {
	return IsZyppPackageManager(bashCommand.Bin()) &&
		(bashCommand.SubCommand() == "install" || bashCommand.SubCommand() == "in")
}

func IsFedoraPackageInstall(bashCommand BashCommand) bool {
	return IsDnfPackageManager(bashCommand.Bin()) && bashCommand.SubCommand() == installStr
}
