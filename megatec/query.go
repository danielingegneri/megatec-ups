package megatec

type QueryResponse struct {
	InputVoltage      float64   `json:"input_voltage"`
	InputFaultVoltage float64   `json:"input_fault_voltage"`
	OutputVoltage     float64   `json:"output_voltage"`
	OutputCurrent     uint64    `json:"output_current"`
	InputFrequency    float64   `json:"input_frequency"`
	BatteryVoltage    float64   `json:"battery_voltage"`
	Temperature       float64   `json:"temperature"`
	Status            UPSStatus `json:"status"`
}

type UPSStatus struct {
	UtilityFail        bool `json:"utility_fail"`          // Bit 7
	BatteryLow         bool `json:"battery_low"`           // Bit 6
	ByPassOrBuckActive bool `json:"bypass_or_buck_active"` // Bit 5
	UPSFail            bool `json:"ups_fail"`              // Bit 4
	StandBy            bool `json:"standby"`               // Bit 3 == 1
	Online             bool `json:"online"`                // Bit 3 == 0
	TestInProgress     bool `json:"testing_in_progress"`   // Bit 2
	ShutdownActive     bool `json:"shutdown_active"`       // Bit 1
	BeeperOn           bool `json:"beeper_on"`             // Bit 0
}
