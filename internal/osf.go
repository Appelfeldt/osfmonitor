package osf

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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
