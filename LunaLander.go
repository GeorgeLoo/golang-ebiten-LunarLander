

/*

George Loo

Luna Lander 26.8.2017

*/



//jj
package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	// "image"
	// "image/color"
	"log"
	"os"
	"math"
	"path/filepath"
)

const (
	screenwidth = 600
	screenheight = 400
	datafolder = "lunarLanderData"
	sampleRate   = 44100
	clockwise = 200
	anticlockwise = 202
	notRotating = 204
	kLunarLanderHeight = 50
	kMoonOrbitalSpeed = 1600
	kCommandModule = 801
	kLunarModule = 802
)

type landerData struct {
	pointedDir int 
	flyingDir int 
	speed int 
	height int 
	fuel int
	thrust float64 
	x, y float64 
	w, h int
	image *ebiten.Image
	flameimg *ebiten.Image
	rocketEngine bool
	vertSpeed int
	horSpeed int 
	rotateSpeed int 
	shiprotdir int
	key string 
	keyleft bool // prevent this key from repeating
	keyright bool // prevent this key from repeating
	docked bool
	shipname string
	dummy int
	retrox float64
	retroy float64
}

type landingZone struct {
	lunarSurfaceimage *ebiten.Image 
	x, y float64 

}

type soundData struct {
	mute bool
	//audioContext    *audio.Context
	audioPlayer     *audio.Player
	soundArr       []audio.Player

}

var (
	surface landingZone
	ship landerData
	commMod landerData
	count int
	sound soundData
	soundjab int
	soundboom int 
	soundmainrocket int 
	canChangeFullscreen bool
	engineloop *audio.Player
	audioContext    *audio.Context
	shipFocus int 

	keyStates    = map[ebiten.Key]int{
		ebiten.KeyUp:    0,
		ebiten.KeyDown:  0,
		ebiten.KeyLeft:  0,
		ebiten.KeyRight: 0,
		ebiten.KeyA:     0,
		ebiten.KeyS:     0,
		ebiten.KeyW:     0,
		ebiten.KeyD:     0,
	}
)

func loopsoundinit() {
	var err error
	fmt.Print("hello init\n")
	audioContext, err = audio.NewContext(sampleRate)
	if err != nil {
		log.Fatal(err)
	}
}

func loadloop(fn string) *audio.Player {

	wavF, err := ebitenutil.OpenFile(filepath.Join(datafolder, fn))
	if err != nil {
		log.Fatal(err)
	}

	wavS, err := wav.Decode(audioContext, wavF)
	if err != nil {
		log.Fatal(err)
	}

	s := audio.NewInfiniteLoop(wavS, wavS.Size())

	player, err := audio.NewPlayer(audioContext, s)
	if err != nil {
		log.Fatal(err)
	}
	return player
}

func (s *soundData) load(fn string) int {
	var err error
	
	//var audioPlayer     *audio.Player
	
	f, err := os.Open(filepath.Join(datafolder, fn))
	if err != nil {
		log.Fatal(err)
	}

	d, err := wav.Decode(audioContext, f)
	if err != nil {
		log.Fatal(err)
	}

	s.audioPlayer, err = audio.NewPlayer(audioContext, d)
	if err != nil {
		log.Fatal(err)
	}
	s.soundArr = append(s.soundArr, *s.audioPlayer)
	i := len(s.soundArr) - 1
	return i // index to the sound

}

func (s *soundData) play(idx int) error {
	//var err error

	if s.mute {
		return nil
	}
	ap := s.soundArr[idx]
	if !ap.IsPlaying() {
		//fmt.Print("sound or not?\n")
		ap.Rewind()
		err := ap.Play()
		if err != nil {
			panic(err)
		}
	}

	if err := audioContext.Update(); err != nil {
		fmt.Print(" !!!!!!!!!!!! SOUND ERROR \n")
		return err
	}
	return nil 
}

func (s *soundData) nosound(m bool) {
	s.mute = m 
}


/*
only need one audio context, or so I think...
*/
func (s *soundData) init() {
	const sampleRate  = 44100
	//var err error
	s.mute = false
	/*s.audioContext, err = audio.NewContext(sampleRate)
	if err != nil {
		log.Fatal(err)
	}*/

}

func readimg(fn string) *ebiten.Image {
	var err error
	var fname string
	fname = filepath.Join(datafolder, fn)
	img, _, err := ebitenutil.NewImageFromFile(
		fname,
		ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	return img

}

func (z *landingZone) init(w,h int ) {
	z.lunarSurfaceimage = readimg("Lunar-Landscape.png")
	z.x = 0
	z.y = 300
}

func (z *landingZone) draw(screen *ebiten.Image ) {

	w, h := z.lunarSurfaceimage.Size()
	//fmt.Print("w h ", w,h,"\n")
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Reset()
	xs := (float64(screenwidth) / float64(w)) // making the image fit the screen
//	ys := float64(screenheight) / float64(h)
//	xs :=  float64(w) / float64(screenwidth) 
//	ys :=  float64(h) / float64(screenheight) / 2
	ys := 1.0
	z.y = float64(screenheight) - float64(h)
	//fmt.Print("xs ys ", xs,", ",ys,"\n")
	opts.GeoM.Scale( xs, ys )
	opts.GeoM.Translate(z.x, z.y)

	screen.DrawImage(z.lunarSurfaceimage, opts)

}

func (l *landerData) init(shipFILEname string,
						  flamename string,
						  pointedDir int,
						  shipname string,
						  x float64,
						  y float64  ) {

	// var err error
	// var fname string
	// fname = filepath.Join(datafolder, "lander.png")
	// img, _, err := ebitenutil.NewImageFromFile(
	// 	fname,
	// 	ebiten.FilterNearest)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//l.image = readimg("lander.png")
	//l.flameimg = readimg("flame.png")
	l.image = readimg(shipFILEname)
	l.flameimg = readimg(flamename)
	l.shipname = shipname

	l.pointedDir = pointedDir // "north" is 0
	l.x = x 
	l.y = y
	l.rocketEngine = false
	l.rotateSpeed = 0
	l.shiprotdir = notRotating
	l.horSpeed = kMoonOrbitalSpeed
	l.vertSpeed = 0
	l.height = 200
	l.thrust = 0
	l.flyingDir = 90
	l.key = "none"
	l.fuel = 500000
	l.docked = true 
	l.retrox = 0
	l.retroy = 0

}

func (l *landerData) draw(screen *ebiten.Image) {

	

	w, h := l.image.Size()
	//fmt.Printf("w %d h %d \n",w,h)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Reset()
	opts.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	opts.GeoM.Rotate(float64(l.pointedDir % 360) * 2 * math.Pi / 360)
	opts.GeoM.Scale( 1.0, 1.0 )
	opts.GeoM.Translate(l.x, l.y)

	screen.DrawImage(l.image, opts)
	if l.rocketEngine {
		screen.DrawImage(l.flameimg, opts)
	}

	// msg := fmt.Sprintf("Fuel %d HSpeed %d VSpeed %d\n dir %d",
	// 	l.fuel, l.horSpeed, l.vertSpeed, l.pointedDir )
	// ebitenutil.DebugPrint(screen, msg)

}

func (l *landerData) drawStatus(screen *ebiten.Image) {
	var dockstr string 

	if l.docked {
		dockstr = "DOCKED"
	} else {
		dockstr = "UNDOCKED"
	}
	msg := fmt.Sprintf("[%s] [%s] Fuel %d HSpeed %d VSpeed %d\n dir %d %f %f",
		l.shipname, dockstr,  l.fuel, l.horSpeed, l.vertSpeed, l.pointedDir, l.x, l.y )

	ebitenutil.DebugPrint(screen, msg)
}

func approx(n,l,h int) bool {

	// if math.Abs(float64(n)-float64(m)) < 10 {
	// 	return true
	// }
	if n > l && n < h {
		return true
	}

	// if n > 355 && m == 0 {
	// 	return true //handle 0 case 
	// }
	return false 
}

func (l *landerData) physics(screen *ebiten.Image) {

	const (
		gravityVal = 0.1
	)

	l.x += float64(l.retrox)
	l.y += float64(l.retroy)

	if l.horSpeed > 0 && l.flyingDir == 90 {
		l.x += 1
	}

	if l.fuel < 1 {
		l.rocketEngine = false
	}

	if approx(l.pointedDir, 260,280) && l.rocketEngine {
		l.horSpeed -= 1
		l.fuel -= int(10 * l.thrust * 10)
		if l.horSpeed < 1 {
			l.horSpeed = 0
		}
	}

	if l.horSpeed < kMoonOrbitalSpeed {
		l.y += gravityVal  // going down by gravity
		l.vertSpeed += 1
		if l.vertSpeed > 10 {
			l.vertSpeed = 10
		}
	}

	w, h := screen.Size()
	if l.x > float64(w) {		// right edge of the screen
		l.x = -0.0
	}

	if l.y < 25 {
		l.y = 25 // top of screen
	}

	if l.x < 0 {
		l.x = float64(w)  // loop left to right of screen
	}

	if approx(l.pointedDir, 350, 360) || approx(l.pointedDir, -1, 10) {

		if l.rocketEngine {

			l.y -= 0.0 + float64(l.thrust)
			l.vertSpeed -= 1
			l.fuel -= int(10 * l.thrust * 10)
			if l.vertSpeed < 1 {
				l.vertSpeed = 0
			}
		}
	}

	if l.rocketEngine {  // climbing back to orbit
		if approx(l.pointedDir, 25,75)  {
			l.x += float64(l.thrust)
			l.y -= float64(l.thrust)
			l.fuel -= int(10 * l.thrust * 10)
			l.horSpeed += 1
		}
	}

	if approx(l.pointedDir, 70,110) && l.rocketEngine {
		l.x += float64(l.thrust)
		l.fuel -= int(10 * l.thrust * 10)
		l.horSpeed += 1
	}

	if l.rocketEngine {  // slowing down for landing
		if approx(l.pointedDir, 270, 350)  {
			//l.x += float64(l.thrust)
			l.y -= gravityVal  // go up
			l.fuel -= int(10 * l.thrust * 10)
			l.horSpeed -= 1
		}
	}

	if l.rotateSpeed > 0 {
		l.pointedDir += l.rotateSpeed
	 	if l.pointedDir > 359 {
	 		l.pointedDir = 0
	 	}

	} else if l.rotateSpeed < 0 {
	 	l.pointedDir += l.rotateSpeed
	 	if l.pointedDir < 0 {
	 		l.pointedDir = 359
		}
	}

	if l.y > float64(h) - kLunarLanderHeight { //

		l.y = float64(h) - kLunarLanderHeight //

		// check if too fast then crashed

		l.horSpeed = 0
		l.vertSpeed = 0 
	}
}

func (l *landerData) control() {
	for key := range keyStates {
		if !ebiten.IsKeyPressed(key) {
			keyStates[key] = 0
			continue
		}
		keyStates[key]++
	}


	if ebiten.IsKeyPressed(ebiten.KeySlash) && l.key != "KeySlash" {
		l.key = "KeySlash"
	} 	

	if keyStates[ebiten.KeyA] == 1 { //ebiten.IsKeyPressed(ebiten.KeyA) {
		sound.play(soundboom)
		l.retrox -= 0.1
		//l.x -= 1
	}
	if keyStates[ebiten.KeyW] == 1 { //ebiten.IsKeyPressed(ebiten.KeyW) {
		sound.play(soundboom)
		l.retroy -= 0.1
		//l.y -= 1
	}

	if keyStates[ebiten.KeyS] == 1 { // ebiten.IsKeyPressed(ebiten.KeyS) {
		sound.play(soundboom)
		l.retroy += 0.1
		//l.y += 1
	}

	if keyStates[ebiten.KeyD] == 1 { // ebiten.IsKeyPressed(ebiten.KeyD) {
		sound.play(soundboom)
		l.retrox += 0.1
		//l.x += 1
	}

	if keyStates[ebiten.KeyLeft] == 1 { //ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if l.keyleft  {
			//l.key = "keyleft"
			l.keyleft = false
			sound.play(soundboom)
			l.rotateSpeed += 1
			l.shiprotdir = anticlockwise
		}
	} else {
		l.keyleft = true
		//l.key = "none"
	}

	if keyStates[ebiten.KeyRight] == 1 { // ebiten.IsKeyPressed(ebiten.KeyRight)  {
		if l.keyright {
			l.rotateSpeed -= 1
			//l.key = "keyright"		
			l.keyright = false
		 	sound.play(soundboom)
			l.shiprotdir = clockwise
	 	}
	} else {
		l.keyright = true
	}

	if l.fuel < 1 {
		engineloop.Pause()
	}

	if keyStates[ebiten.KeyUp] == 1 {  // reduce thrust
		l.thrust -= 0.1
		if l.thrust < 0 {
			l.thrust = 0
			l.rocketEngine = false
			engineloop.Pause()
		}
	}

	if keyStates[ebiten.KeyDown] == 1 {  // increase thrust
		l.thrust += 0.1
		if l.thrust > 5 {
			l.thrust = 5
		}
	 	l.rocketEngine = true
		if l.rocketEngine {
			//sound.play(soundmainrocket)
			engineloop.Play()
		}
	 	
	}

	// if l.rotateSpeed > 0 {
	// 	l.pointedDir += l.rotateSpeed
	//  	if l.pointedDir > 359 {
	//  		l.pointedDir = 0
	//  	}

	// } else if l.rotateSpeed < 0 {
	//  	l.pointedDir += l.rotateSpeed
	//  	if l.pointedDir < 0 {
	//  		l.pointedDir = 359
	// 	}
	// }


	togglFullscreen()
}

func togglFullscreen() {
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		if canChangeFullscreen {
			ebiten.SetFullscreen(!ebiten.IsFullscreen())
			canChangeFullscreen = false
		}
	} else {
		canChangeFullscreen = true
	}
}


func update(screen *ebiten.Image) error {

	if ebiten.IsRunningSlowly() {
		return nil
		//fmt.Print("running slowly! \n")
	}

	if ebiten.IsKeyPressed(ebiten.Key1) {
		shipFocus = kCommandModule
		
	}

	if ebiten.IsKeyPressed(ebiten.Key2) {
		shipFocus = kLunarModule
	}

	if ebiten.IsKeyPressed(ebiten.KeyU) {
		ship.docked = false
		commMod.docked = false
		if ship.x < commMod.x {
			ship.x -= 3
		} else {
			ship.x += 3
		}

	}

	if ship.x < commMod.x {

		if (int(commMod.x) - int(ship.x) == 71) && (int(commMod.y) == int(ship.y)) {
			ship.docked = true
			commMod.docked = true
			ship.pointedDir = 90
			commMod.pointedDir = 270
			ship.retrox = 0
			ship.retroy = 0
		}
	} else {
		if (int(commMod.x) - int(ship.x) == -71) && (int(commMod.y) == int(ship.y)) {
			ship.docked = true
			commMod.docked = true
			ship.pointedDir = 270
			commMod.pointedDir = 90
			ship.retrox = 0
			ship.retroy = 0
		}
	}


	//fmt.Print(" command x y ",commMod.x,", ",commMod.y,"\n")
	//fmt.Print(" luna x y ", ship.x,", ",ship.y,"\n")

	surface.draw(screen)

	count++
	ship.draw(screen)
	commMod.draw(screen)

	if shipFocus == kLunarModule {
		ship.control()
		ship.drawStatus(screen)
	} else if shipFocus == kCommandModule {
		commMod.control()
		commMod.drawStatus(screen)
	}

	ship.physics(screen)
	commMod.physics(screen)

	if ship.docked {

		if commMod.rocketEngine {
			ship.x += float64(commMod.thrust)
			ship.horSpeed += 1
		}
	}

	return nil 
}



func main() {
	loopsoundinit()  // has sound context init code in it, only one instance per program
	sound.init()
	soundjab = sound.load("jab.wav")
	soundboom = sound.load("boom.wav")
	//soundmainrocket = sound.load("explosion.wav")
	engineloop = loadloop("explosion.wav")
	ebiten.SetRunnableInBackground(true)
	ebiten.SetFullscreen(false)
	count = 0
	surface.init(screenwidth, screenheight)
	
	ship.init("lander.png",
				"flame.png",
				90, // pointed direction
				"LUNAR MODULE",
				10,
				40)

	ship.dummy = 1234
	
	commMod.init("command.png",
				"commandflame.png",
				270, //
				"COMMAND MODULE",
				80, //x
				40)  //y
	
	shipFocus = kLunarModule

	
	scale := 2.0
	// Initialize Ebiten, and loop the update() function
	if err := ebiten.Run(update, screenwidth, screenheight, scale, "Luna Lander 0.0 by George Loo"); err != nil {
		panic(err)
	}
	fmt.Printf("Lunar Lander Program ended -----------------\n")

}

