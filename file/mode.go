package file

const (
	// FileModeNewDirectory is the default directory mode
	FileModeNewDirectory = 0o750

	// FileModeDefault is the default file mode
	FileModeNewFile = 0o640

	// FileModeMinValid is the minimum valid file mode
	FileModeMinValid = 0o000

	// FileModeMaxValid is the maximum valid file mode
	FileModeMaxValid = 0o777
)
