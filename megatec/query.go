package megatec

type QueryResponse struct {
	InputVoltage      float64
	InputFaultVoltage float64
	OutputVoltage     float64
	OutputCurrent     uint
	InputFrequency    float64
	BatteryVoltage    float64
	Temperature       float64
	Status            UPSStatus
}

type UPSStatus struct {
	UtilityFail        bool // Bit 7
	BatteryLow         bool // Bit 6
	ByPassOrBuckActive bool // Bit 5
	UPSFail            bool // Bit 4
	StandBy            bool // Bit 3 = 1
	Online             bool // Bit 3 = 0
	TestInProgress     bool // Bit 2
	ShutdownActive     bool // Bit 1
	BeeperOn           bool // Bit 0
}
