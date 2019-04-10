package megatec

import (
	"bufio"
	"errors"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"log"
	"strconv"
	"strings"
)

type ups struct {
	port io.ReadWriteCloser
}

type UPS interface {
	Close()
	Query() (QueryResponse, error)
	Test() error
	TestFor(seconds uint)
	TestUntilLow()
	CancelTest()
	ToggleBeeper()
	Shutdown(restore bool)
	CancelShutdown()
	GetInfo()
	GetRating()
	readUntilCR() (string, error)
}

func NewUPS(device string, baud uint, dataBits uint, stopBits uint, parityMode serial.ParityMode) UPS {
	options := serial.OpenOptions{
		PortName:        device,
		BaudRate:        baud,
		DataBits:        dataBits,
		StopBits:        stopBits,
		MinimumReadSize: 1,
		ParityMode:      parityMode,
	}
	port, err := serial.Open(options)
	if err != nil {
		panic(err)
	}
	return &ups{
		port: port,
	}
}

func (u ups) Query() (QueryResponse, error) {
	b := []byte{0x51, 0x31, 0x0D} // Q1<cr>
	_, err := u.port.Write(b)

	data, err := u.readUntilCR()
	if err != nil {
		return QueryResponse{}, err
	}

	parts := strings.Split(data, " ")
	log.Print(len(parts))
	log.Print(parts)
	// || len(parts[7]) != 8
	if len(parts) != 6 {
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
	response.OutputCurrent = uint(outputCurrent)
	response.InputFrequency, _ = strconv.ParseFloat(parts[4], 0)
	response.BatteryVoltage, _ = strconv.ParseFloat(parts[5], 0)
	//response.Temperature, _ = strconv.ParseFloat(parts[6], 0)
	//response.Status.UtilityFail = parts[7][0] == '1'
	//response.Status.BatteryLow = parts[7][1] == '1'
	//response.Status.ByPassOrBuckActive = parts[7][2] == '1'
	//response.Status.UPSFail = parts[7][3] == '1'
	//response.Status.StandBy = parts[7][4] == '1'
	//response.Status.Online = parts[7][4] == '0'
	//response.Status.TestInProgress = parts[7][5] == '1'
	//response.Status.ShutdownActive = parts[7][6] == '1'
	//response.Status.BeeperOn = parts[7][7] == '1'

	//log.Printf("%q", response)

	// TODO: Populate struct with response
	return response, nil
}

func (u ups) Close() {
	u.port.Close()
}

func (u ups) Test() error {
	b := []byte{0x54, 0x0D} // T<cr>
	_, err := u.port.Write(b)
	if err != nil {
		return err
	}
	return nil
}
func (u ups) TestFor(seconds uint) {

}
func (u ups) TestUntilLow() {

}
func (u ups) CancelTest() {

}
func (u ups) ToggleBeeper() {

}
func (u ups) Shutdown(restore bool) {

}
func (u ups) CancelShutdown() {

}
func (u ups) GetInfo() {

}
func (u ups) GetRating() {

}

func (u ups) readUntilCR() (string, error) {
	scanner := bufio.NewScanner(u.port)
	onCR := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for i := 0; i < len(data); i++ {
			log.Print(i, data[i])
			if data[i] == 0x0D {
				return i + 1, data[:i], nil
			}
		}
		return 0, data, bufio.ErrFinalToken
	}
	scanner.Split(onCR)
	var line strings.Builder
	for scanner.Scan() {
		line.WriteString(scanner.Text())
	}
	log.Print(line)
	return line.String(), nil
}
