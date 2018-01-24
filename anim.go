


/*

anim.go by George Loo 15.12.2017



jj
*/


package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
	"path/filepath"
	"image/color"
	"image"

)

const (
	version = "0.1"
	//datafolder = "data"
	//screenwidth = 800
	//screenheight = 400
)


type AnimationType struct {
	sequence *ebiten.Image 
	x, y float64 
	width, height int  // of one frame
	scale float64
	looping bool
	run bool
	numFrames int 
	startAt int  // which frame
	speed int   // 60 is one second 
	numberOfPlays int 
	count int 
	currF int 

}

var (
  //mousedownState bool
  //animseq *ebiten.Image
  anim1 AnimationType
  anim2 AnimationType
  anim3 AnimationType

)

func (a *AnimationType) init(name string, w int, h int, numF int) {
	a.sequence = readimg(name)
	a.count = 1
	a.run = false
	a.width = w
	a.height = h 
	a.startAt = 0
	a.numFrames = numF

}

func (a *AnimationType) animate(screen *ebiten.Image, x float64, y float64) {
	var x1, y1, x2, y2 int

	if !a.run {
		return
	}
	a.count++
	x1 = a.currF * a.width
	y1 = 0
	x2 = x1 + a.width
	y2 = a.height
	r := image.Rect(x1,y1,x2,y2)
	//fmt.Printf("%d %d %d %d \n",x1,y1,x2,y2)
	//if a.count == 1 {
	draw(screen, a.sequence, x, y, r)
		//}
	if a.count < a.speed {
		return
	}
	//fmt.Print(a.count," anim \n")
	a.count = 1

	//fmt.Print(a.count," anim \n")
	a.currF += 1
	if a.currF > a.numFrames - 1 {
		a.currF = 0
		if a.looping == false {
			a.run = false
		}
	}

}

func (a *AnimationType) plate() {

}

//////////////////////////////////////////////////

// func mouseLeftdown() {
// 	fmt.Print("mousedown \n")
// }

// func mouseLeftup() {
// 	fmt.Print("mouseup \n")
// }


func update0(screen *ebiten.Image) error {

	if ebiten.IsRunningSlowly() {
		return nil
		//fmt.Print("running slowly! \n")
	}

	screen.Fill(color.NRGBA{255, 255, 0, 0xff})  // yellow

  	if mousedownState {
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			mousedownState = false
			mouseLeftup()
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !mousedownState {
			mousedownState = true
			mouseLeftdown()
		}
	}

	//draw(screen, animseq, 100, 150)

	anim1.animate(screen, 200, 150)
	anim2.animate(screen, 200, 250)
	anim3.animate(screen, 300, 50)
	
  	return nil

}

func draw(screen *ebiten.Image, image2draw *ebiten.Image, x float64, y float64, r image.Rectangle) {
	//w, h := image.Size()
	//fmt.Printf("w %d h %d \n",w,h)
	//var r image.Rectangle 
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Reset()
	//opts.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	//opts.GeoM.Rotate(float64(l.pointedDir % 360) * 2 * math.Pi / 360)
	//opts.GeoM.Scale( 1.0, 1.0 )
	opts.GeoM.Scale( 1.0, 1.0 )
	opts.GeoM.Translate(x, y)

	//r := image.Rect(0, 0, 50, 50)
	opts.SourceRect = &r

	screen.DrawImage(image2draw, opts)

}


func readimg0(fn string) *ebiten.Image {
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

func initprog() {

	//animseq = readimg("sixframes.png")
	anim1.init("sixframes.png", 50,50, 6)
	anim1.looping = true
	anim1.speed = 15
	anim1.run = true
	anim1.currF = 0

	anim2.init("sixframes.png", 50,50, 6)
	anim2.looping = true
	anim2.speed = 60
	anim2.run = true
	anim2.currF = 2

	anim3.init("explosion01.png", 100,100, 15)
	anim3.looping = true
	anim3.speed = 15
	anim3.run = true
	anim3.currF = 0

}


func anim() {

	initprog()

    scale := 1.0
    // Initialize Ebiten, and loop the update() function
    if err := ebiten.Run(update, screenwidth, screenheight, scale, "Animation test 0.0 by George Loo"); err != nil {
      panic(err)
    }
    fmt.Printf("Program ended -----------------\n")

}
