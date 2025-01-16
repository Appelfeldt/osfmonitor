package osfm

import (
	"bytes"
	"encoding/binary"
	"log"
	"math"
	"net"
	"time"

	ui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Settings struct {
	Port uint16
}

type OSFFeatures struct {
	EyeLeft                float32
	EyeRight               float32
	EyebrowSteepnessLeft   float32
	EyebrowUpDownLeft      float32
	EyebrowQuirkLeft       float32
	EyebrowSteepnessRight  float32
	EyebrowUpDownRight     float32
	EyebrowQuirkRight      float32
	MouthCornerUpDownLeft  float32
	MouthCornerInOutLeft   float32
	MouthCornerUpDownRight float32
	MouthCornerInOutRight  float32
	MouthOpen              float32
	MouthWide              float32
}

type OSFPacket struct {
	Time             float64
	Id               int32
	CameraResolution rl.Vector2
	RightEyeOpen     float32
	LeftEyeOpen      float32
	Got3DPoints      bool
	Fit3DError       float32
	RawQuaternion    rl.Quaternion
	RawEuler         rl.Vector3
	Translation      rl.Vector3
	Confidence       [68]float32
	Points           [68]rl.Vector2
	Points3D         [70]rl.Vector3
	Features         OSFFeatures
}

func (p OSFPacket) GetRotation() rl.Vector3 {
	r := p.RawEuler
	r.X = float32(math.Mod(float64(-(r.X + 180)), 360))
	r.Z = float32(math.Mod(float64(r.Z-90), 360))
	return r
}

func Run(settings Settings) {
	rl.SetConfigFlags(rl.FlagWindowTransparent)
	rl.InitWindow(500, 500, "OpenSeeFace Monitor")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	//BG toggle button
	bgDefaultColor := rl.Black
	bgColor := bgDefaultColor
	bgTransparent := false
	bgBtnBounds := rl.Rectangle{X: 24, Y: 8, Width: 110, Height: 20}
	bgBtnText := "Make transparent"

	//Points size slider
	pRadiusLblBounds := rl.Rectangle{X: 24, Y: 8 + bgBtnBounds.Height + bgBtnBounds.Y, Width: 110, Height: 20}
	pRadiusLblText := "Point Size:"
	pRadiusSldrBounds := rl.Rectangle{X: 24, Y: 0 + pRadiusLblBounds.Height + pRadiusLblBounds.Y, Width: 110, Height: 20}
	var pRadius float32 = 1.0

	//Waiting for data message
	noInputMsg := "Waiting for connection..."
	noInputMsgLen := len(noInputMsg)
	noInputMsgSize := rl.MeasureTextEx(rl.GetFontDefault(), noInputMsg, 20, 2)

	//Start listening for osf data
	c := make(chan OSFPacket)
	go listen(settings, c)
	var packet OSFPacket
	var lastPacket time.Time
	var connected = false

	for !rl.WindowShouldClose() {

		//Toggle background button
		if ui.Button(bgBtnBounds, bgBtnText) {
			if bgTransparent {
				bgColor = bgDefaultColor
				bgBtnText = "Make transparent"
			} else {
				bgColor = rl.Blank
				bgBtnText = "Make opaque"
			}
			bgTransparent = !bgTransparent
		}

		//Circle size slider
		ui.Label(pRadiusLblBounds, pRadiusLblText)
		pRadius = ui.Slider(pRadiusSldrBounds, "0", "10", pRadius, 0, 10)

		//Check for packets
		select {
		case packet = <-c:
			lastPacket = time.Now()
		default:
		}
		//Consider the connection active if a packet was received the last 2 seconds
		connected = !lastPacket.Add(time.Second * 2).Before(time.Now())

		rl.BeginDrawing()
		rl.ClearBackground(bgColor)

		if connected {
			rl.SetWindowSize(int(packet.CameraResolution.X), int(packet.CameraResolution.Y))
			for _, p := range packet.Points {
				rl.DrawCircle(int32(p.X), int32(p.Y), pRadius, rl.Red)
			}
		} else {
			rl.SetWindowSize(500, 500)
			rl.DrawText(noInputMsg[:noInputMsgLen-3+int(int64(rl.GetTime())%4)],
				int32(rl.GetScreenWidth()/2-int(noInputMsgSize.X/2)),
				int32(rl.GetScreenHeight()/2-int(noInputMsgSize.Y/2)),
				20,
				rl.GetColor(0xAA0000FF),
			)
		}

		rl.EndDrawing()
	}
}

func listen(settings Settings, c chan OSFPacket) {

	//Start listening for OSF traffic
	addr := &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: int(settings.Port),
		Zone: "",
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	//Write data to struct, send struct to channel
	buffer := make([]byte, 8+4+2*4+2*4+1+4+3*4+3*4+4*4+4*68+4*2*68+4*3*70+4*14)
	for !rl.WindowShouldClose() {
		_, err := conn.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		buf := bytes.NewBuffer(buffer)

		datagram := OSFPacket{}
		err = binary.Read(buf, binary.LittleEndian, &datagram)
		if err != nil {
			log.Fatal(err)
		}
		c <- datagram
	}
}
