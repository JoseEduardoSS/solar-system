package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"
)

// dummyImage é utilizada como textura para desenhar triângulos.
var dummyImage *ebiten.Image

func init() {
	dummyImage = ebiten.NewImage(1, 1)
	dummyImage.Fill(color.White)
	rand.Seed(time.Now().UnixNano())
}

// lerpColor interpola linearmente entre duas cores.
func lerpColor(c1, c2 color.RGBA, t float64) color.RGBA {
	return color.RGBA{
		R: uint8(float64(c1.R) + t*(float64(c2.R)-float64(c1.R))),
		G: uint8(float64(c1.G) + t*(float64(c2.G)-float64(c1.G))),
		B: uint8(float64(c1.B) + t*(float64(c2.B)-float64(c1.B))),
		A: uint8(float64(c1.A) + t*(float64(c2.A)-float64(c1.A))),
	}
}

// drawFilledCircle desenha um círculo preenchido aproximando-o por um fan de triângulos.
func drawFilledCircle(screen *ebiten.Image, cx, cy, radius float64, clr color.RGBA) {
	const segments = 30
	n := segments

	vertices := make([]ebiten.Vertex, n+2)
	rf, gf, bf, af := float32(clr.R)/255, float32(clr.G)/255, float32(clr.B)/255, float32(clr.A)/255

	// Centro do círculo.
	vertices[0] = ebiten.Vertex{
		DstX:   float32(cx),
		DstY:   float32(cy),
		ColorR: rf, ColorG: gf, ColorB: bf, ColorA: af,
		SrcX: 0, SrcY: 0,
	}

	// Vértices do contorno.
	for i := 0; i <= n; i++ {
		theta := 2 * math.Pi * float64(i) / float64(n)
		x := cx + radius*math.Cos(theta)
		y := cy + radius*math.Sin(theta)
		vertices[i+1] = ebiten.Vertex{
			DstX:   float32(x),
			DstY:   float32(y),
			ColorR: rf, ColorG: gf, ColorB: bf, ColorA: af,
			SrcX: 0, SrcY: 0,
		}
	}

	indices := make([]uint16, n*3)
	for i := 0; i < n; i++ {
		indices[i*3] = 0
		indices[i*3+1] = uint16(i + 1)
		indices[i*3+2] = uint16(i + 2)
	}
	screen.DrawTriangles(vertices, indices, dummyImage, nil)
}

// drawThickLine desenha uma linha grossa entre dois pontos.
func drawThickLine(screen *ebiten.Image, x1, y1, x2, y2, thickness float64, clr color.RGBA) {
	dx := x2 - x1
	dy := y2 - y1
	length := math.Hypot(dx, dy)
	if length == 0 {
		return
	}
	nx := -dy / length
	ny := dx / length
	half := thickness / 2

	x1a, y1a := x1+nx*half, y1+ny*half
	x1b, y1b := x1-nx*half, y1-ny*half
	x2a, y2a := x2+nx*half, y2+ny*half
	x2b, y2b := x2-nx*half, y2-ny*half

	rf, gf, bf, af := float32(clr.R)/255, float32(clr.G)/255, float32(clr.B)/255, float32(clr.A)/255
	vertices := []ebiten.Vertex{
		{DstX: float32(x1a), DstY: float32(y1a), SrcX: 0, SrcY: 0, ColorR: rf, ColorG: gf, ColorB: bf, ColorA: af},
		{DstX: float32(x2a), DstY: float32(y2a), SrcX: 0, SrcY: 0, ColorR: rf, ColorG: gf, ColorB: bf, ColorA: af},
		{DstX: float32(x2b), DstY: float32(y2b), SrcX: 0, SrcY: 0, ColorR: rf, ColorG: gf, ColorB: bf, ColorA: af},
		{DstX: float32(x1b), DstY: float32(y1b), SrcX: 0, SrcY: 0, ColorR: rf, ColorG: gf, ColorB: bf, ColorA: af},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, dummyImage, nil)
}

// drawGlowingLine desenha uma linha com três camadas (espessuras e cores diferentes).
func drawGlowingLine(screen *ebiten.Image, x1, y1, x2, y2 float64,
	glowColor, midColor, coreColor color.RGBA) {
	drawThickLine(screen, x1, y1, x2, y2, 4, glowColor)
	drawThickLine(screen, x1, y1, x2, y2, 2, midColor)
	drawThickLine(screen, x1, y1, x2, y2, 1, coreColor)
}

// drawCircleOutline desenha o contorno de um círculo, aproximado por segmentos.
func drawCircleOutline(screen *ebiten.Image, cx, cy, radius, thickness float64, clr color.RGBA) {
	const segments = 60
	points := make([][2]float64, segments)
	for i := 0; i < segments; i++ {
		theta := 2 * math.Pi * float64(i) / float64(segments)
		x := cx + radius*math.Cos(theta)
		y := cy + radius*math.Sin(theta)
		points[i] = [2]float64{x, y}
	}
	for i := 0; i < segments; i++ {
		next := (i + 1) % segments
		drawThickLine(screen, points[i][0], points[i][1], points[next][0], points[next][1], thickness, clr)
	}
}

// drawSunGradient desenha o sol com gradiente radial e pulsação.
func drawSunGradient(screen *ebiten.Image, cx, cy, baseRadius, t float64) {
	steps := 30
	// O sol pulsa levemente (variação de ±10% no raio).
	pulse := 1 + 0.1*math.Sin(t*2)
	sunRadius := baseRadius * pulse

	yellow := color.RGBA{255, 255, 0, 255}
	midColor := color.RGBA{255, 140, 0, 200}
	outerColor := color.RGBA{255, 140, 0, 0}

	for i := 0; i <= steps; i++ {
		f := float64(i) / float64(steps)
		// Raio varia do solRadius (centro) até 4 vezes maior (contorno).
		r := sunRadius * (1 + 3*(1-f))
		var clr color.RGBA
		if f < 0.3 {
			clr = lerpColor(outerColor, midColor, f/0.3)
		} else {
			clr = lerpColor(midColor, yellow, (f-0.3)/0.7)
		}
		drawFilledCircle(screen, cx, cy, r, clr)
	}
}

// drawPlanetGradient desenha um planeta com gradiente do contorno (cor externa) até o centro (cor interna).
func drawPlanetGradient(screen *ebiten.Image, cx, cy, radius float64, innerColor, outerColor color.RGBA) {
	steps := 20
	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps-1)
		r := radius * (1 - t)
		clr := lerpColor(outerColor, innerColor, t)
		drawFilledCircle(screen, cx, cy, r, clr)
	}
}

// drawSaturnRings desenha os anéis de Saturno com inclinação.
func drawSaturnRings(screen *ebiten.Image, cx, cy, planetRadius float64) {
	segments := 60
	tilt := 20 * math.Pi / 180.0 // 20 graus de inclinação
	outerX := planetRadius * 2.0
	outerY := planetRadius * 1.0
	innerX := planetRadius * 1.6
	innerY := planetRadius * 0.8

	ringColor := color.RGBA{210, 180, 140, 180}
	outerPoints := make([][2]float64, segments)
	innerPoints := make([][2]float64, segments)
	for i := 0; i < segments; i++ {
		theta := 2 * math.Pi * float64(i) / float64(segments)
		ox := outerX * math.Cos(theta)
		oy := outerY * math.Sin(theta)
		rx := ox*math.Cos(tilt) - oy*math.Sin(tilt)
		ry := ox*math.Sin(tilt) + oy*math.Cos(tilt)
		outerPoints[i] = [2]float64{cx + rx, cy + ry}

		ix := innerX * math.Cos(theta)
		iy := innerY * math.Sin(theta)
		rx2 := ix*math.Cos(tilt) - iy*math.Sin(tilt)
		ry2 := ix*math.Sin(tilt) + iy*math.Cos(tilt)
		innerPoints[i] = [2]float64{cx + rx2, cy + ry2}
	}

	vertices := []ebiten.Vertex{}
	indices := []uint16{}
	for i := 0; i < segments; i++ {
		next := (i + 1) % segments
		base := uint16(len(vertices))
		vertices = append(vertices,
			ebiten.Vertex{
				DstX: float32(outerPoints[i][0]), DstY: float32(outerPoints[i][1]),
				SrcX: 0, SrcY: 0,
				ColorR: float32(ringColor.R) / 255, ColorG: float32(ringColor.G) / 255,
				ColorB: float32(ringColor.B) / 255, ColorA: float32(ringColor.A) / 255,
			},
			ebiten.Vertex{
				DstX: float32(outerPoints[next][0]), DstY: float32(outerPoints[next][1]),
				SrcX: 0, SrcY: 0,
				ColorR: float32(ringColor.R) / 255, ColorG: float32(ringColor.G) / 255,
				ColorB: float32(ringColor.B) / 255, ColorA: float32(ringColor.A) / 255,
			},
			ebiten.Vertex{
				DstX: float32(innerPoints[next][0]), DstY: float32(innerPoints[next][1]),
				SrcX: 0, SrcY: 0,
				ColorR: float32(ringColor.R) / 255, ColorG: float32(ringColor.G) / 255,
				ColorB: float32(ringColor.B) / 255, ColorA: float32(ringColor.A) / 255,
			},
			ebiten.Vertex{
				DstX: float32(innerPoints[i][0]), DstY: float32(innerPoints[i][1]),
				SrcX: 0, SrcY: 0,
				ColorR: float32(ringColor.R) / 255, ColorG: float32(ringColor.G) / 255,
				ColorB: float32(ringColor.B) / 255, ColorA: float32(ringColor.A) / 255,
			},
		)
		indices = append(indices, base, base+1, base+2, base, base+2, base+3)
	}
	screen.DrawTriangles(vertices, indices, dummyImage, nil)
}

// -------------------------
// ESTRUTURAS ADICIONAIS
// -------------------------

// Star representa uma estrela com brilho oscilante.
type Star struct {
	X, Y           float64
	Phase, Speed   float64
	BaseBrightness uint8
}

// Moon representa uma lua orbitando um planeta.
type Moon struct {
	OrbitRadius float64
	Angle       float64
	OrbitSpeed  float64
	Radius      float64
	InnerColor  color.RGBA
	OuterColor  color.RGBA
}

// Update atualiza a posição da lua.
func (m *Moon) Update() {
	m.Angle += m.OrbitSpeed
}

// Asteroid representa uma partícula do cinturão de asteroides.
type Asteroid struct {
	OrbitRadius float64
	Angle       float64
	OrbitSpeed  float64
	Radius      float64
}

// Update atualiza a posição do asteroide.
func (a *Asteroid) Update() {
	a.Angle += a.OrbitSpeed
}

// Comet representa um cometa com cauda dinâmica.
type Comet struct {
	X, Y          float64
	Angle         float64
	Speed         float64
	TailPoints    [][2]float64
	TailMaxLength int
}

// Update atualiza a posição do cometa e gerencia a cauda.
func (c *Comet) Update() {
	// Atualiza posição
	c.X += c.Speed * math.Cos(c.Angle)
	c.Y += c.Speed * math.Sin(c.Angle)

	// Insere a posição atual na cauda (no início)
	c.TailPoints = append([][2]float64{{c.X, c.Y}}, c.TailPoints...)

	// Limita o comprimento da cauda
	if len(c.TailPoints) > c.TailMaxLength {
		c.TailPoints = c.TailPoints[:c.TailMaxLength]
	}
}

// -------------------------
// Planet – inclui possível lua(s)
// -------------------------
type Planet struct {
	Name        string
	OrbitRadius float64
	Radius      float64
	Angle       float64
	OrbitSpeed  float64
	X, Y        float64
	InnerColor  color.RGBA
	OuterColor  color.RGBA
	Draggable   bool
	IsDragged   bool
	Moons       []*Moon
}

// Update atualiza o ângulo da órbita do planeta.
func (p *Planet) Update() {
	if !p.IsDragged {
		p.Angle += p.OrbitSpeed
	}
	// Atualiza as luas
	for _, m := range p.Moons {
		m.Update()
	}
}

// -------------------------
// Simulation – estado geral da simulação
// -------------------------
type Simulation struct {
	sunX, sunY               float64
	sunRadius                float64
	planets                  []*Planet
	draggedPlanet            *Planet
	dragOffsetX, dragOffsetY float64
	stars                    []Star
	asteroids                []Asteroid
	comet                    *Comet
	time                     float64
}

// NewSimulation inicializa a simulação com planetas, estrelas, cinturão de asteroides e cometa.
func NewSimulation() *Simulation {
	sim := &Simulation{
		sunRadius: 40,
		planets:   make([]*Planet, 0),
		time:      0,
	}

	// --- Planetas ---
	// Mercúrio
	sim.planets = append(sim.planets, &Planet{
		Name:        "Mercúrio",
		OrbitRadius: 80,
		Radius:      6,
		Angle:       0,
		OrbitSpeed:  0.04,
		InnerColor:  color.RGBA{169, 169, 169, 255},
		OuterColor:  color.RGBA{105, 105, 105, 255},
		Draggable:   false,
	})
	// Vênus
	sim.planets = append(sim.planets, &Planet{
		Name:        "Vênus",
		OrbitRadius: 120,
		Radius:      8,
		Angle:       0,
		OrbitSpeed:  0.03,
		InnerColor:  color.RGBA{255, 215, 0, 255},
		OuterColor:  color.RGBA{218, 165, 32, 255},
		Draggable:   false,
	})
	// Terra (arrastável) – com uma lua
	terra := &Planet{
		Name:        "Terra",
		OrbitRadius: 160,
		Radius:      10,
		Angle:       0,
		OrbitSpeed:  0.02,
		InnerColor:  color.RGBA{100, 149, 237, 255},
		OuterColor:  color.RGBA{25, 25, 112, 255},
		Draggable:   true,
	}
	terra.Moons = []*Moon{
		{
			OrbitRadius: 20,
			Angle:       0,
			OrbitSpeed:  0.05,
			Radius:      3,
			InnerColor:  color.RGBA{240, 240, 240, 255},
			OuterColor:  color.RGBA{160, 160, 160, 255},
		},
	}
	sim.planets = append(sim.planets, terra)
	// Marte
	sim.planets = append(sim.planets, &Planet{
		Name:        "Marte",
		OrbitRadius: 200,
		Radius:      7,
		Angle:       0,
		OrbitSpeed:  0.015,
		InnerColor:  color.RGBA{205, 92, 92, 255},
		OuterColor:  color.RGBA{139, 69, 19, 255},
		Draggable:   false,
	})
	// Júpiter
	sim.planets = append(sim.planets, &Planet{
		Name:        "Júpiter",
		OrbitRadius: 250,
		Radius:      14,
		Angle:       0,
		OrbitSpeed:  0.01,
		InnerColor:  color.RGBA{222, 184, 135, 255},
		OuterColor:  color.RGBA{160, 82, 45, 255},
		Draggable:   false,
	})
	// Saturno (com anéis)
	sim.planets = append(sim.planets, &Planet{
		Name:        "Saturn",
		OrbitRadius: 300,
		Radius:      12,
		Angle:       0,
		OrbitSpeed:  0.008,
		InnerColor:  color.RGBA{222, 203, 164, 255},
		OuterColor:  color.RGBA{210, 180, 140, 255},
		Draggable:   false,
	})
	// Urano
	sim.planets = append(sim.planets, &Planet{
		Name:        "Urano",
		OrbitRadius: 350,
		Radius:      10,
		Angle:       0,
		OrbitSpeed:  0.006,
		InnerColor:  color.RGBA{175, 238, 238, 255},
		OuterColor:  color.RGBA{72, 209, 204, 255},
		Draggable:   false,
	})
	// Netuno
	sim.planets = append(sim.planets, &Planet{
		Name:        "Netuno",
		OrbitRadius: 400,
		Radius:      10,
		Angle:       0,
		OrbitSpeed:  0.005,
		InnerColor:  color.RGBA{65, 105, 225, 255},
		OuterColor:  color.RGBA{25, 25, 112, 255},
		Draggable:   false,
	})
	// Plutão
	sim.planets = append(sim.planets, &Planet{
		Name:        "Plutão",
		OrbitRadius: 450,
		Radius:      4,
		Angle:       0,
		OrbitSpeed:  0.004,
		InnerColor:  color.RGBA{205, 133, 63, 255},
		OuterColor:  color.RGBA{139, 69, 19, 255},
		Draggable:   false,
	})

	// --- Campo de Estrelas ---
	w, h := ebiten.WindowSize()
	starCount := 200
	sim.stars = make([]Star, starCount)
	for i := 0; i < starCount; i++ {
		sim.stars[i] = Star{
			X:              float64(rand.Intn(w)),
			Y:              float64(rand.Intn(h)),
			Phase:          rand.Float64() * 2 * math.Pi,
			Speed:          0.005 + rand.Float64()*0.005,
			BaseBrightness: uint8(100 + rand.Intn(155)),
		}
	}

	// --- Cinturão de Asteroides ---
	asteroidCount := 150
	sim.asteroids = make([]Asteroid, asteroidCount)
	for i := 0; i < asteroidCount; i++ {
		// Distribuídos entre 210 e 240 (entre Marte e Júpiter)
		orbitRadius := 210 + rand.Float64()*30
		sim.asteroids[i] = Asteroid{
			OrbitRadius: orbitRadius,
			Angle:       rand.Float64() * 2 * math.Pi,
			OrbitSpeed:  0.008 + rand.Float64()*0.004,
			Radius:      1 + rand.Float64()*1.5,
		}
	}

	// --- Cometa ---
	sim.comet = &Comet{
		// Começa fora da tela, na parte superior esquerda
		X:             -50,
		Y:             -50,
		Angle:         math.Pi / 4, // 45° em direção à direita/inferior
		Speed:         4.0,
		TailPoints:    make([][2]float64, 0),
		TailMaxLength: 20,
	}

	return sim
}

// Update é chamado a cada frame.
func (sim *Simulation) Update() error {
	w, h := ebiten.WindowSize()
	sim.sunX = float64(w) / 2
	sim.sunY = float64(h) / 2
	sim.time += 0.016 // Aproximadamente 60 fps

	// Atualiza as estrelas (twinkling)
	for i := range sim.stars {
		sim.stars[i].Phase += sim.stars[i].Speed
	}

	// Atualiza os asteroides
	for i := range sim.asteroids {
		sim.asteroids[i].Update()
	}

	// Atualiza o cometa
	sim.comet.Update()
	// Se o cometa sair da tela, reinicia sua posição
	if sim.comet.X > float64(w)+50 || sim.comet.Y > float64(h)+50 {
		sim.comet.X = -50
		sim.comet.Y = -50
		sim.comet.TailPoints = sim.comet.TailPoints[:0]
	}

	// Processa entrada do mouse para planetas arrastáveis
	mx, my := ebiten.CursorPosition()
	mouseX := float64(mx)
	mouseY := float64(my)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for _, p := range sim.planets {
			if p.Draggable {
				dx := mouseX - p.X
				dy := mouseY - p.Y
				if math.Hypot(dx, dy) <= p.Radius {
					sim.draggedPlanet = p
					p.IsDragged = true
					sim.dragOffsetX = p.X - mouseX
					sim.dragOffsetY = p.Y - mouseY
					break
				}
			}
		}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if sim.draggedPlanet != nil {
			sim.draggedPlanet.IsDragged = false
			sim.draggedPlanet = nil
		}
	}
	if sim.draggedPlanet != nil {
		sim.draggedPlanet.X = mouseX + sim.dragOffsetX
		sim.draggedPlanet.Y = mouseY + sim.dragOffsetY
		dx := sim.draggedPlanet.X - sim.sunX
		dy := sim.draggedPlanet.Y - sim.sunY
		sim.draggedPlanet.OrbitRadius = math.Hypot(dx, dy)
		sim.draggedPlanet.Angle = math.Atan2(dy, dx)
	}

	// Atualiza planetas e suas luas
	for _, p := range sim.planets {
		if !p.IsDragged {
			p.Update()
			p.X = sim.sunX + p.OrbitRadius*math.Cos(p.Angle)
			p.Y = sim.sunY + p.OrbitRadius*math.Sin(p.Angle)
		}
	}

	return nil
}

// Draw é chamado a cada frame para renderizar a cena.
func (sim *Simulation) Draw(screen *ebiten.Image) {
	_, _ = screen.Size()
	// Fundo espacial
	screen.Fill(color.RGBA{10, 10, 30, 255})

	// Desenha as estrelas com brilho oscilante
	for _, star := range sim.stars {
		brightness := 128 + 127*math.Sin(star.Phase)
		if brightness < 0 {
			brightness = 0
		} else if brightness > 255 {
			brightness = 255
		}
		starColor := color.RGBA{255, 255, 255, uint8(brightness)}
		drawFilledCircle(screen, star.X, star.Y, 1, starColor)
	}

	// Desenha o cinturão de asteroides
	for _, a := range sim.asteroids {
		ax := sim.sunX + a.OrbitRadius*math.Cos(a.Angle)
		ay := sim.sunY + a.OrbitRadius*math.Sin(a.Angle)
		asteroidColor := color.RGBA{169, 169, 169, 200}
		drawFilledCircle(screen, ax, ay, a.Radius, asteroidColor)
	}

	// Desenha o sol com pulsação
	drawSunGradient(screen, sim.sunX, sim.sunY, sim.sunRadius, sim.time)

	// Desenha as órbitas dos planetas
	orbitColor := color.RGBA{200, 200, 200, 50}
	for _, p := range sim.planets {
		drawCircleOutline(screen, sim.sunX, sim.sunY, p.OrbitRadius, 1, orbitColor)
	}

	// Desenha os planetas e, se houver, suas luas
	for _, p := range sim.planets {
		// "Halo" do planeta
		glowColor := color.RGBA{0, 0, 0, 100}
		drawFilledCircle(screen, p.X, p.Y, float64(p.Radius)*1.4, glowColor)
		// Planeta com gradiente
		drawPlanetGradient(screen, p.X, p.Y, float64(p.Radius), p.InnerColor, p.OuterColor)
		// Se for Saturno, desenha os anéis
		if p.Name == "Saturn" {
			drawSaturnRings(screen, p.X, p.Y, float64(p.Radius))
		}
		// Desenha as luas, se houver
		for _, m := range p.Moons {
			// A posição da lua é relativa ao planeta
			mx := p.X + m.OrbitRadius*math.Cos(m.Angle)
			my := p.Y + m.OrbitRadius*math.Sin(m.Angle)
			drawPlanetGradient(screen, mx, my, m.Radius, m.InnerColor, m.OuterColor)
		}
	}

	// Desenha os "raios" de luz que partem do sol (efeito de raios intermitentes)
	for angleDeg := 0; angleDeg < 360; angleDeg++ {
		theta := float64(angleDeg) * math.Pi / 180.0
		dx := math.Cos(theta)
		dy := math.Sin(theta)
		ox := sim.sunX
		oy := sim.sunY
		bestT := math.MaxFloat64
		var hitPlanet *Planet
		var hitX, hitY float64
		for _, p := range sim.planets {
			cx := p.X
			cy := p.Y
			r := float64(p.Radius)
			ocx := ox - cx
			ocy := oy - cy
			b := 2 * (dx*ocx + dy*ocy)
			c := ocx*ocx + ocy*ocy - r*r
			disc := b*b - 4*c
			if disc < 0 {
				continue
			}
			sqrtDisc := math.Sqrt(disc)
			t1 := (-b - sqrtDisc) / 2
			t2 := (-b + sqrtDisc) / 2
			var t float64
			if t1 > 0 {
				t = t1
			} else {
				t = t2
			}
			if t > 0 && t < bestT {
				bestT = t
				hitPlanet = p
				hitX = ox + dx*t
				hitY = oy + dy*t
			}
		}
		var endX, endY float64
		if hitPlanet == nil {
			endX = ox + dx*1000
			endY = oy + dy*1000
		} else {
			endX = hitX
			endY = hitY
		}
		drawGlowingLine(screen, ox, oy, endX, endY,
			color.RGBA{255, 255, 200, 60},
			color.RGBA{255, 255, 170, 120},
			color.RGBA{255, 255, 150, 200})
	}

	// Desenha o cometa e sua cauda
	// Desenha a cauda (linha conectando pontos, com opacidade decrescente)
	tailLen := len(sim.comet.TailPoints)
	for i := 0; i < tailLen-1; i++ {
		alpha := uint8(200 * (1 - float64(i)/float64(tailLen)))
		c1 := color.RGBA{255, 255, 255, alpha}
		c2 := color.RGBA{255, 255, 255, alpha / 2}
		x1 := sim.comet.TailPoints[i][0]
		y1 := sim.comet.TailPoints[i][1]
		x2 := sim.comet.TailPoints[i+1][0]
		y2 := sim.comet.TailPoints[i+1][1]
		drawGlowingLine(screen, x1, y1, x2, y2, c1, c1, c2)
	}
	// Desenha o núcleo do cometa
	drawFilledCircle(screen, sim.comet.X, sim.comet.Y, 4, color.RGBA{255, 255, 255, 255})
}

// Layout define o tamanho da janela.
func (sim *Simulation) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowTitle("Simulação Avançada do Sistema Solar")
	ebiten.SetFullscreen(true)
	sim := NewSimulation()
	if err := ebiten.RunGame(sim); err != nil {
		log.Fatal(err)
	}
}
