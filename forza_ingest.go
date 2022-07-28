package forza5

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

const maxBufferSize = 1024

type FH5Telemetry interface {
	printDetails()
}

// packet 324 bytes?
type FH5Telemetry struct {
	IsRaceOn                             int32   `json:"isRaceOn"` // 0 or 1
	TimeStampMS                          uint32  `json:"timeStampMS"`
	EngineMaxRPM                         float32 `json:"engineMaxRPM"`
	EngineIdleRPM                        float32 `json:"engineIdleRPM"`
	CurrentEngineRPM                     float32 `json:"currentEngineRPM"`
	AccelerationX                        float32 `json:"accelerationX"`
	AccelerationY                        float32 `json:"accelerationY"`
	AccelerationZ                        float32 `json:"accelerationZ"`
	VelocityX                            float32 `json:"velocityX"`
	VelocityY                            float32 `json:"velocityY"`
	VelocityZ                            float32 `json:"velocityZ"`
	AngularVelocityX                     float32 `json:"angularVelocityX"`
	AngularVelocityY                     float32 `json:"angularVelocityY"`
	AngularVelocityZ                     float32 `json:"angularVelocityZ"`
	Yaw                                  float32 `json:"yaw"`
	Pitch                                float32 `json:"pitch"`
	Roll                                 float32 `json:"roll"`
	NormalizedSuspensionTravelFrontLeft  float32 `json:"normalizedSuspensionTravelFrontLeft"`
	NormalizedSuspensionTravelFrontRight float32 `json:"normalizedSuspensionTravelFrontRight"`
	NormalizedSuspensionTravelRearLeft   float32 `json:"normalizedSuspensionTravelRearLeft"`
	NormalizedSuspensionTravelRearRight  float32 `json:"normalizedSuspensionTravelRearRight"`
	TireSlipRatioFrontLeft               float32 `json:"tireSlipRatioFrontLeft"`
	TireSlipRatioFrontRight              float32 `json:"tireSlipRatioFrontRight"`
	TireSlipRatioRearLeft                float32 `json:"tireSlipRatioRearLeft"`
	TireSlipRatioRearRight               float32 `json:"tireSlipRatioRearRight"`
	WheelRotationSpeedFrontLeft          float32 `json:"wheelRotationSpeedFrontLeft"`
	WheelRotationSpeedFrontRight         float32 `json:"wheelRotationSpeedFrontRight"`
	WheelRotationSpeedRearLeft           float32 `json:"wheelRotationSpeedRearLeft"`
	WheelRotationSpeedRearRight          float32 `json:"wheelRotationSpeedRearRight"`
	WheelOnRumbleStripFrontLeft          int32   `json:"wheelOnRumbleStripFrontLeft"`
	WheelOnRumbleStripFrontRight         int32   `json:"wheelOnRumbleStripFrontRight"`
	WheelOnRumbleStripRearLeft           int32   `json:"wheelOnRumbleStripRearLeft"`
	WheelOnRumbleStripRearRight          int32   `json:"wheelOnRumbleStripRearRight"`
	WheelInPuddleDepthFrontLeft          float32 `json:"wheelInPuddleDepthFrontLeft"`
	WheelInPuddleDepthFrontRight         float32 `json:"wheelInPuddleDepthFrontRight"`
	WheelInPuddleDepthRearLeft           float32 `json:"wheelInPuddleDepthRearLeft"`
	WheelInPuddleDepthRearRight          float32 `json:"wheelInPuddleDepthRearRight"`
	SurfaceRumbleFrontLeft               float32 `json:"surfaceRumbleFrontLeft"`
	SurfaceRumbleFrontRight              float32 `json:"surfaceRumbleFrontRight"`
	SurfaceRumbleRearLeft                float32 `json:"surfaceRumbleRearLeft"`
	SurfaceRumbleRearRight               float32 `json:"surfaceRumbleRearRight"`
	TireSlipAngleFrontLeft               float32 `json:"tireSlipAngleFrontLeft"`
	TireSlipAngleFrontRight              float32 `json:"tireSlipAngleFrontRight"`
	TireSlipAngleRearLeft                float32 `json:"tireSlipAngleRearLeft"`
	TireSlipAngleRearRight               float32 `json:"tireSlipAngleRearRight"`
	TireCombinedSlipFrontLeft            float32 `json:"tireCombinedSlipFrontLeft"`
	TireCombinedSlipFrontRight           float32 `json:"tireCombinedSlipFrontRight"`
	TireCombinedSlipRearLeft             float32 `json:"tireCombinedSlipRearLeft"`
	TireCombinedSlipRearRight            float32 `json:"tireCombinedSlipRearRight"`
	SuspensionTravelMetersFrontLeft      float32 `json:"suspensionTravelMetersFrontLeft"`
	SuspensionTravelMetersFrontRight     float32 `json:"suspensionTravelMetersFrontRight"`
	SuspensionTravelMetersRearLeft       float32 `json:"suspensionTravelMetersRearLeft"`
	SuspensionTravelMetersRearRight      float32 `json:"suspensionTravelMetersRearRight"`
	CarClass                             int32   `json:"carClass"`
	CarPerformanceIndex                  int32   `json:"carPerformanceIndex"`
	DrivetrainType                       int32   `json:"drivetrainType"`
	NumCylinders                         int32   `json:"numCylinders"`
	CarType                              int32   `json:"carType"`
	Unknown1                             byte    `json:"unknown1"`
	Unknown2                             byte    `json:"unknown2"`
	Unknown3                             byte    `json:"unknown3"`
	Unknown4                             byte    `json:"unknown4"`
	Unknown5                             byte    `json:"unknown5"`
	Unknown6                             byte    `json:"unknown6"`
	Unknown7                             byte    `json:"unknown7"`
	Unknown8                             byte    `json:"unknown8"`
	CarOrdinal                           int32   `json:"carOrdinal"`
	PositionX                            float32 `json:"positionX"`
	PositionY                            float32 `json:"positionY"`
	PositionZ                            float32 `json:"positionZ"`
	Speed                                float32 `json:"speed"`
	Power                                float32 `json:"power"`
	Torque                               float32 `json:"torque"`
	TireTempFrontLeft                    float32 `json:"tireTempFrontLeft"`
	TireTempFrontRight                   float32 `json:"tireTempFrontRight"`
	TireTempRearLeft                     float32 `json:"tireTempRearLeft"`
	TireTempRearRight                    float32 `json:"tireTempRearRight"`
	Boost                                float32 `json:"boost"`
	Fuel                                 float32 `json:"fuel"`
	DistanceTraveled                     float32 `json:"distanceTraveled"`
	BestLap                              float32 `json:"bestLap"`
	LastLap                              float32 `json:"lastLap"`
	CurrentLap                           float32 `json:"currentLap"`
	CurrentRaceTime                      float32 `json:"currentRaceTime"`
	LapNumber                            uint16  `json:"lapNumber"`
	RacePosition                         uint8   `json:"racePosition"`
	Throttle                             uint8   `json:"throttle"`  // normalize255to1
	Brake                                uint8   `json:"brake"`     // normalize255to1
	Clutch                               uint8   `json:"clutch"`    // normalize255to1
	HandBrake                            uint8   `json:"handBrake"` // normalize255to1
	Gear                                 uint8   `json:"gear"`
	Steer                                uint8   `json:"steer"` // normalize255to1
	NormalizedDrivingLine                uint8   `json:"normalizedDrivingLine"`
	NormalizedAIBrakeDifference          uint8   `json:"normalizedAIBrakeDifference"`
}

func (telemetry *FH5Telemetry) ReadBuffer(b []byte, len int) {
	telemetry.IsRaceOn = binary.LittleEndian.Uint32(b[0:4])
	telemetry.TimeStampMS = binary.LittleEndian.Uint32(b[4:8])
}

func (value T) ReadBufferFromOffset(b []byte, offset int, size int) {
	binary.Read(bytes.NewReader(b[0:4]), binary.LittleEndian, &value)
}

func Server(ctx context.Context, address string) (err error) {

	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		return
	}

	defer conn.Close()

	doneChan := make(chan error, 1)
	buffer := make([]byte, maxBufferSize)

	go func() {
		count := 0
		for {
			n, addr, err := conn.ReadFrom(buffer)
			if err != nil {
				doneChan <- err
				return
			}

			fmt.Printf("Packet received (Server) %d: bytes=%d from=%s buffer=%x\n", count, n, addr.String(), buffer[:n])

			var telemetry FH5Telemetry
			telemetry.ReadBuffer(buffer, n)
			fmt.Printf("Telem is race on: %d\n", telemetry.IsRaceOn)
			//err = json.Unmarshal(buffer[:n], &telemetry)
			//if err != nil {
			//	fmt.Println(err)
			//	doneChan <- err
			//	return
			//}

			//fmt.Printf("Buffer (Server): ", buffer)

			deadline := time.Now().Add(15 * time.Second)
			err = conn.SetWriteDeadline(deadline)
			if err != nil {
				doneChan <- err
				return
			}

			//n, err = conn.WriteTo(buffer[:n], addr)
			//if err != nil {
			//	doneChan <- err
			//	return
			//}
			//
			//fmt.Printf("Packet written (Server): bytes=%d to=%s\n", n, addr.String())

			count++
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Cancelled")
		err = ctx.Err()
	case err = <-doneChan:
	}

	return nil
}

func Client(ctx context.Context, address string, reader io.Reader) (err error) {
	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return
	}

	defer conn.Close()

	doneChan := make(chan error, 1)

	go func() {
		n, err := io.Copy(conn, reader)
		if err != nil {
			doneChan <- err
			return
		}

		fmt.Printf("Packet Written (Client): bytes=%d\n", n)

		buffer := make([]byte, maxBufferSize)
		deadline := time.Now().Add(15 * time.Second)
		err = conn.SetReadDeadline(deadline)
		if err != nil {
			doneChan <- err
			return
		}

		nRead, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			doneChan <- err
			return
		}

		fmt.Printf("Packet Received (Client): bytes=%d from=%s\n", nRead, addr.String())
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Cancelled")
		err = ctx.Err()
	case err = <-doneChan:
	}

	return nil
}
