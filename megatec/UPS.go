package megatec

import (
	"bufio"
	"errors"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"strconv"
	"strings"
)

type ups struct {
	port   io.ReadWriteCloser
	reader bufio.Reader
}

type UPS interface {
	Close()
	Query() (QueryResponse, error)
	Test() error
	TestFor(minutes uint64) error
	TestUntilLow() error
	CancelTest() error
	ToggleBeeper() error
	Shutdown(shutdownDelayMin float64) error
	ShutdownRestore(shutdownDelayMin float64, restoreDelayMin uint) error
	CancelShutdown() error
	GetInfo() error
	GetRating() error
	readUntilCR() (string, error)
	write(line string) error
}

func NewUPS(device string, baud uint, dataBits uint, stopBits uint, parityMode serial.ParityMode) UPS {
	options := serial.OpenOptions{
		PortName:              device,
		BaudRate:              baud,
		DataBits:              dataBits,
		StopBits:              stopBits,
		InterCharacterTimeout: 5000,
		ParityMode:            parityMode,
	}
	port, err := serial.Open(options)
	reader := bufio.NewReader(port)
	if err != nil {
		panic(err)
	}
	return &ups{
		port:   port,
		reader: *reader,
	}
}

func (u ups) Query() (QueryResponse, error) {
	err := u.write("Q1")
	if err != nil {
		return QueryResponse{}, err
	}
	data, err := u.readUntilCR()
	if err != nil {
		return QueryResponse{}, err
	}

	parts := strings.Split(data, " ")
	if len(parts) != 8 || len(parts[7]) != 8 {
		return QueryResponse{}, errors.New("invalid response")
	}

	var response QueryResponse
	inputVoltage := parts[0]
	if parts[0][0] == '(' {
		inputVoltage = parts[0][1:] // For some reason it responds with a bracket ( at the start?
	}
	response.InputVoltage, _ = strconv.ParseFloat(inputVoltage, 0)
	response.InputFaultVoltage, _ = strconv.ParseFloat(parts[1], 0)
	response.OutputVoltage, _ = strconv.ParseFloat(parts[2], 0)
	outputCurrent, _ := strconv.ParseUint(parts[3], 10, 0)
	response.OutputCurrent = uint64(outputCurrent)
	response.InputFrequency, _ = strconv.ParseFloat(parts[4], 0)
	response.BatteryVoltage, _ = strconv.ParseFloat(parts[5], 0)
	response.Temperature, _ = strconv.ParseFloat(parts[6], 0)
	response.Status.UtilityFail = parts[7][0] == '1'
	response.Status.BatteryLow = parts[7][1] == '1'
	response.Status.ByPassOrBuckActive = parts[7][2] == '1'
	response.Status.UPSFail = parts[7][3] == '1'
	response.Status.StandBy = parts[7][4] == '1'
	response.Status.Online = parts[7][4] == '0'
	response.Status.TestInProgress = parts[7][5] == '1'
	response.Status.ShutdownActive = parts[7][6] == '1'
	response.Status.BeeperOn = parts[7][7] == '1'

	return response, nil
}

func (u ups) Close() {
	u.port.Close()
}

func (u ups) Test() error {
	return u.write("T")
}
func (u ups) TestFor(minutes uint64) error {
	// TODO: TEST
	return u.write("T" + strconv.FormatUint(minutes, 10))
}
func (u ups) TestUntilLow() error {
	// TODO: TEST
	return u.write("TL")
}
func (u ups) CancelTest() error {
	// TODO: TEST
	return u.write("CT")
}
func (u ups) ToggleBeeper() error {
	// TODO: TEST
	return u.write("Q")
}
func (u ups) Shutdown(shutdownDelayMin float64) error {
	// TODO: TEST
	return u.write("S" + strconv.FormatFloat(shutdownDelayMin, 'f', 1, 32))
}
func (u ups) ShutdownRestore(shutdownDelayMin float64, restoreDelayMin uint) error {
	// TODO: TEST
	if shutdownDelayMin > 10 {
		return errors.New("shutdownDelayMin limited to 10")
	}
	if restoreDelayMin > 9999 {
		return errors.New("restoreDelayMin limited to 9999")
	}
	// TODO: Check if shutdownDelayMin needs to be formatted exactly: "is a number ranging from .2, .3, …, 01, 02, …, up to 10"
	// TODO: Check if restoreDelayMin needs to be formatted exactly: "is a number ranging from 0001 to 9999"
	return u.write("S" + strconv.FormatFloat(shutdownDelayMin, 'f', 1, 32) + "C" + strconv.Itoa(int(restoreDelayMin)))
}
func (u ups) CancelShutdown() error {
	// TODO: TEST
	return u.write("C")
}
func (u ups) GetInfo() error {
	// TODO: Implement info
	return u.write("I")
}
func (u ups) GetRating() error {
	// TODO: Implement rating stuff
	return u.write("F")
}

func (u ups) readUntilCR() (string, error) {
	line, err := u.reader.ReadString('\r')
	if err != nil {
		return "", err
	}
	return strings.Trim(line, "\r"), nil
}

func (u ups) write(line string) error {
	b := []byte(line + "\r")
	if _, err := u.port.Write(b); err != nil {
		return err
	}
	return nil
}
