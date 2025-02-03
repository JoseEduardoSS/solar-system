package main

import (
	"math"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Definição de uma cor Silver (já que rl.Silver não está definida)
var Silver = rl.NewColor(192, 192, 192, 255)

// Definindo constantes para os modos de câmera (para os modos "normais")
const (
	CameraFree    = 1 // Modo livre (free)
	CameraOrbital = 2 // Modo orbital
)

// ─────────────────────────────────────────────
// Estruturas da simulação
type Star struct {
	Position       rl.Vector3
	Phase          float64
	Speed          float64
	BaseBrightness int
}

type Moon struct {
	OrbitRadius float64
	Angle       float64
	OrbitSpeed  float64
	Radius      float32
	Color       rl.Color
}

type Planet struct {
	Name        string
	OrbitRadius float64
	Radius      float32
	Angle       float64
	OrbitSpeed  float64
	Color       rl.Color
	Moons       []*Moon
}

type Asteroid struct {
	OrbitRadius float64
	Angle       float64
	OrbitSpeed  float64
	Radius      float32
}

type Comet struct {
	Position      rl.Vector3
	Angle         float64
	Speed         float64
	TailPoints    []rl.Vector3
	TailMaxLength int
}

type Simulation struct {
	SunRadius float32
	Planets   []*Planet
	Stars     []Star
	Asteroids []Asteroid // Inclui cinturão principal e o Kuiper Belt
	Comet     Comet
	Time      float64

	// Campos para o efeito de explosão (impacto)
	ExplosionActive   bool
	ExplosionTime     float64
	ExplosionDuration float64
	ExplosionPosition rl.Vector3
}

// NewSimulation cria e inicializa os corpos celestes
func NewSimulation() *Simulation {
	sim := &Simulation{
		SunRadius:         40,
		Planets:           make([]*Planet, 0),
		Time:              0,
		ExplosionActive:   false,
		ExplosionDuration: 1.0, // duração da explosão em segundos
	}

	// Planetas – cada um com seus parâmetros
	sim.Planets = append(sim.Planets, &Planet{
		Name:        "Mercurio",
		OrbitRadius: 80,
		Radius:      6,
		Angle:       0,
		OrbitSpeed:  0.04,
		Color:       rl.Gray,
	})
	sim.Planets = append(sim.Planets, &Planet{
		Name:        "Venus",
		OrbitRadius: 120,
		Radius:      8,
		Angle:       0,
		OrbitSpeed:  0.03,
		Color:       rl.Orange,
	})
	terra := &Planet{
		Name:        "Terra",
		OrbitRadius: 160,
		Radius:      10,
		Angle:       0,
		OrbitSpeed:  0.02,
		Color:       rl.Blue,
	}
	terra.Moons = []*Moon{
		{
			OrbitRadius: 20,
			Angle:       0,
			OrbitSpeed:  0.05,
			Radius:      3,
			Color:       rl.LightGray,
		},
	}
	sim.Planets = append(sim.Planets, terra)
	sim.Planets = append(sim.Planets, &Planet{
		Name:        "Marte",
		OrbitRadius: 200,
		Radius:      7,
		Angle:       0,
		OrbitSpeed:  0.015,
		Color:       rl.Red,
	})
	// Jupiter com múltiplas luas
	jupiter := &Planet{
		Name:        "Jupiter",
		OrbitRadius: 250,
		Radius:      14,
		Angle:       0,
		OrbitSpeed:  0.01,
		Color:       rl.Brown,
	}
	jupiter.Moons = []*Moon{
		{
			OrbitRadius: 20,
			Angle:       0,
			OrbitSpeed:  0.06,
			Radius:      3,
			Color:       rl.LightGray,
		},
		{
			OrbitRadius: 30,
			Angle:       1,
			OrbitSpeed:  0.04,
			Radius:      2,
			Color:       Silver,
		},
		{
			OrbitRadius: 40,
			Angle:       2,
			OrbitSpeed:  0.035,
			Radius:      2,
			Color:       rl.LightGray,
		},
	}
	sim.Planets = append(sim.Planets, jupiter)
	sim.Planets = append(sim.Planets, &Planet{
		Name:        "Saturn",
		OrbitRadius: 300,
		Radius:      12,
		Angle:       0,
		OrbitSpeed:  0.008,
		Color:       rl.Gold,
	})
	sim.Planets = append(sim.Planets, &Planet{
		Name:        "Urano",
		OrbitRadius: 350,
		Radius:      10,
		Angle:       0,
		OrbitSpeed:  0.006,
		Color:       rl.NewColor(173, 216, 230, 255),
	})
	sim.Planets = append(sim.Planets, &Planet{
		Name:        "Netuno",
		OrbitRadius: 400,
		Radius:      10,
		Angle:       0,
		OrbitSpeed:  0.005,
		Color:       rl.DarkBlue,
	})
	sim.Planets = append(sim.Planets, &Planet{
		Name:        "Plutao",
		OrbitRadius: 450,
		Radius:      4,
		Angle:       0,
		OrbitSpeed:  0.004,
		Color:       rl.Brown,
	})

	// Estrelas distribuídas numa casca esférica distante
	starCount := 200
	sim.Stars = make([]Star, starCount)
	for i := 0; i < starCount; i++ {
		r := 600 + rand.Float64()*200
		theta := rand.Float64() * 2 * math.Pi
		phi := rand.Float64() * math.Pi
		x := r * math.Sin(phi) * math.Cos(theta)
		y := r * math.Cos(phi)
		z := r * math.Sin(phi) * math.Sin(theta)
		sim.Stars[i] = Star{
			Position:       rl.NewVector3(float32(x), float32(y), float32(z)),
			Phase:          rand.Float64() * 2 * math.Pi,
			Speed:          0.005 + rand.Float64()*0.005,
			BaseBrightness: 100 + rand.Intn(155),
		}
	}

	// Cinturão de Asteroides (principal e Kuiper Belt)
	asteroidCount := 150
	sim.Asteroids = make([]Asteroid, 0, asteroidCount+50)
	for i := 0; i < asteroidCount; i++ {
		orbitRadius := 210 + rand.Float64()*30
		sim.Asteroids = append(sim.Asteroids, Asteroid{
			OrbitRadius: orbitRadius,
			Angle:       rand.Float64() * 2 * math.Pi,
			OrbitSpeed:  0.008 + rand.Float64()*0.004,
			Radius:      1 + float32(rand.Float64()*1.5),
		})
	}
	// Kuiper Belt
	kuiperCount := 50
	for i := 0; i < kuiperCount; i++ {
		orbitRadius := 500 + rand.Float64()*100
		sim.Asteroids = append(sim.Asteroids, Asteroid{
			OrbitRadius: orbitRadius,
			Angle:       rand.Float64() * 2 * math.Pi,
			OrbitSpeed:  0.003 + rand.Float64()*0.002,
			Radius:      0.5 + float32(rand.Float64()*1.0),
		})
	}

	// Inicializa o cometa com posição e direção aleatórias (na periferia)
	resetComet(sim)

	rand.Seed(time.Now().UnixNano())
	return sim
}

// resetComet define uma nova posição e direção para o cometa.
// O cometa é posicionado aleatoriamente em um anel na região externa (raio entre 600 e 1000)
// e sua direção é calculada para apontar aproximadamente para o centro (0,0,0).
func resetComet(sim *Simulation) {
	rMin, rMax := 600.0, 1000.0
	radius := rMin + rand.Float64()*(rMax-rMin)
	angle := rand.Float64() * 2 * math.Pi
	sim.Comet.Position = rl.NewVector3(float32(radius*math.Cos(angle)), 0, float32(radius*math.Sin(angle)))
	// Direção para o centro
	sim.Comet.Angle = math.Atan2(-float64(sim.Comet.Position.Z), -float64(sim.Comet.Position.X))
	sim.Comet.Speed = 4.0
	sim.Comet.TailPoints = make([]rl.Vector3, 0)
}

// distance retorna a distância Euclidiana entre dois pontos 3D.
func distance(v1, v2 rl.Vector3) float32 {
	dx := v1.X - v2.X
	dy := v1.Y - v2.Y
	dz := v1.Z - v2.Z
	return float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}

// CheckCollisions verifica se o cometa colide com algum objeto (planeta ou asteroide)
// (exceto o Sol). Se houver colisão, ativa um efeito de explosão e reinicia o cometa.
func (sim *Simulation) CheckCollisions() {
	cometRadius := float32(4)
	// Colisão com planetas
	for _, p := range sim.Planets {
		planetPos := rl.NewVector3(
			float32(p.OrbitRadius*math.Cos(p.Angle)),
			0,
			float32(p.OrbitRadius*math.Sin(p.Angle)),
		)
		if distance(sim.Comet.Position, planetPos) < (cometRadius + p.Radius) {
			sim.ExplosionActive = true
			sim.ExplosionTime = 0
			sim.ExplosionPosition = sim.Comet.Position
			resetComet(sim)
			return
		}
	}
	// Colisão com asteroides
	for i := range sim.Asteroids {
		a := &sim.Asteroids[i]
		asteroidPos := rl.NewVector3(
			float32(a.OrbitRadius*math.Cos(a.Angle)),
			0,
			float32(a.OrbitRadius*math.Sin(a.Angle)),
		)
		if distance(sim.Comet.Position, asteroidPos) < (cometRadius + a.Radius) {
			sim.ExplosionActive = true
			sim.ExplosionTime = 0
			sim.ExplosionPosition = sim.Comet.Position
			resetComet(sim)
			return
		}
	}
}

func (sim *Simulation) Update() {
	dt := 1.0 / 60.0
	sim.Time += dt

	// Atualiza as fases das estrelas (cintilação)
	for i := range sim.Stars {
		sim.Stars[i].Phase += sim.Stars[i].Speed
	}

	// Atualiza os ângulos dos asteroides
	for i := range sim.Asteroids {
		sim.Asteroids[i].Angle += sim.Asteroids[i].OrbitSpeed
	}

	// Atualiza a posição e o rastro do cometa
	sim.Comet.Position.X += float32(sim.Comet.Speed * math.Cos(sim.Comet.Angle))
	sim.Comet.Position.Z += float32(sim.Comet.Speed * math.Sin(sim.Comet.Angle))
	// Atualiza o rastro: adiciona a posição atual no início da lista
	sim.Comet.TailPoints = append([]rl.Vector3{sim.Comet.Position}, sim.Comet.TailPoints...)
	if len(sim.Comet.TailPoints) > sim.Comet.TailMaxLength {
		sim.Comet.TailPoints = sim.Comet.TailPoints[:sim.Comet.TailMaxLength]
	}
	// Se o cometa sair da região, reinicia
	if sim.Comet.Position.X > 800 || sim.Comet.Position.Z > 800 {
		resetComet(sim)
	}

	// Atualiza os ângulos dos planetas e suas luas
	for _, p := range sim.Planets {
		p.Angle += p.OrbitSpeed
		for _, m := range p.Moons {
			m.Angle += m.OrbitSpeed
		}
	}

	// Verifica colisões e dispara explosão se necessário
	sim.CheckCollisions()

	// Atualiza o tempo da explosão, se ativo
	if sim.ExplosionActive {
		sim.ExplosionTime += dt
		if sim.ExplosionTime >= sim.ExplosionDuration {
			sim.ExplosionActive = false
		}
	}
}

// Desenha as órbitas dos planetas (no plano XZ)
func (sim *Simulation) DrawOrbitPaths() {
	center := rl.NewVector3(0, 0, 0)
	for _, p := range sim.Planets {
		rl.DrawCircle3D(center, float32(p.OrbitRadius), rl.NewVector3(1, 0, 0), 90, rl.LightGray)
	}
}

// Função auxiliar para desenhar uma esfera (usando a função nativa)
func drawSphere(pos rl.Vector3, radius float32, col rl.Color) {
	rl.DrawSphere(pos, radius, col)
}

// Gera um mesh para um anel (para Saturno – no plano XZ)
func generateRingMesh(innerRadius, outerRadius float32, segments int) rl.Mesh {
	vertexCount := segments * 2
	triangleCount := segments * 2

	vertices := make([]float32, vertexCount*3)
	normals := make([]float32, vertexCount*3)
	texcoords := make([]float32, vertexCount*2)
	indices := make([]uint16, triangleCount*3)

	angleStep := 2 * math.Pi / float64(segments)
	vertexIndex, texIndex := 0, 0
	for i := 0; i < segments; i++ {
		angle := float64(i) * angleStep
		cosA, sinA := float32(math.Cos(angle)), float32(math.Sin(angle))
		// Vértice externo
		vertices[vertexIndex] = outerRadius * cosA
		vertices[vertexIndex+1] = 0
		vertices[vertexIndex+2] = outerRadius * sinA
		normals[vertexIndex] = 0
		normals[vertexIndex+1] = 1
		normals[vertexIndex+2] = 0
		texcoords[texIndex] = (cosA + 1) * 0.5
		texcoords[texIndex+1] = (sinA + 1) * 0.5
		vertexIndex += 3
		texIndex += 2

		// Vértice interno
		vertices[vertexIndex] = innerRadius * cosA
		vertices[vertexIndex+1] = 0
		vertices[vertexIndex+2] = innerRadius * sinA
		normals[vertexIndex] = 0
		normals[vertexIndex+1] = 1
		normals[vertexIndex+2] = 0
		texcoords[texIndex] = (cosA + 1) * 0.5
		texcoords[texIndex+1] = (sinA + 1) * 0.5
		vertexIndex += 3
		texIndex += 2
	}

	index := 0
	for i := 0; i < segments; i++ {
		next := (i + 1) % segments
		vi0 := uint16(i * 2)
		vi1 := uint16(i*2 + 1)
		vi2 := uint16(next * 2)
		vi3 := uint16(next*2 + 1)
		indices[index] = vi0
		indices[index+1] = vi2
		indices[index+2] = vi1
		indices[index+3] = vi1
		indices[index+4] = vi2
		indices[index+5] = vi3
		index += 6
	}

	return rl.Mesh{
		VertexCount: int32(vertexCount),
		Vertices:    &vertices[0],
		Normals:     &normals[0],
		Texcoords:   &texcoords[0],
		Indices:     &indices[0],
	}
}

// Desenha a cena 3D usando as funções nativas (esferas, modelo do anel e efeitos)
func (sim *Simulation) Draw3D(ringModel rl.Model) {
	// Desenha o Sol
	drawSphere(rl.NewVector3(0, 0, 0), sim.SunRadius, rl.Yellow)

	// Desenha as órbitas dos planetas
	sim.DrawOrbitPaths()

	// Desenha os planetas e suas luas
	for _, p := range sim.Planets {
		x := float32(p.OrbitRadius * math.Cos(p.Angle))
		z := float32(p.OrbitRadius * math.Sin(p.Angle))
		planetPos := rl.NewVector3(x, 0, z)
		drawSphere(planetPos, p.Radius, p.Color)

		// Se for Saturn, desenha os anéis
		if p.Name == "Saturn" {
			rl.DrawModelEx(ringModel, planetPos, rl.NewVector3(1, 0, 0), 25, rl.NewVector3(p.Radius*3, 1, p.Radius*3), rl.LightGray)
		}
		// Desenha as luas
		for _, m := range p.Moons {
			mx := planetPos.X + float32(m.OrbitRadius*math.Cos(m.Angle))
			mz := planetPos.Z + float32(m.OrbitRadius*math.Sin(m.Angle))
			moonPos := rl.NewVector3(mx, 0, mz)
			drawSphere(moonPos, m.Radius, m.Color)
		}
	}

	// Desenha os asteroides
	for _, a := range sim.Asteroids {
		ax := float32(a.OrbitRadius * math.Cos(a.Angle))
		az := float32(a.OrbitRadius * math.Sin(a.Angle))
		asteroidPos := rl.NewVector3(ax, 0, az)
		drawSphere(asteroidPos, a.Radius, rl.Gray)
	}

	// Desenha o rastro do cometa (meteoro)
	// Primeiro, desenha esferas com alfa decrescente
	for i := 0; i < len(sim.Comet.TailPoints)-1; i++ {
		alpha := uint8(200 * (1 - float32(i)/float32(len(sim.Comet.TailPoints))))
		col := rl.NewColor(255, 255, 255, alpha)
		drawSphere(sim.Comet.TailPoints[i], 2, col)
	}
	// Em seguida, desenha linhas conectando os pontos do rastro para um efeito contínuo
	for i := 0; i < len(sim.Comet.TailPoints)-1; i++ {
		alpha := uint8(200 * (1 - float32(i)/float32(len(sim.Comet.TailPoints))))
		col := rl.NewColor(255, 255, 255, alpha)
		rl.DrawLine3D(sim.Comet.TailPoints[i], sim.Comet.TailPoints[i+1], col)
	}

	// Desenha o cometa (meteoro)
	drawSphere(sim.Comet.Position, 4, rl.White)

	// Desenha as estrelas cintilantes
	for _, star := range sim.Stars {
		brightness := float32(128 + 127*math.Sin(star.Phase))
		if brightness < 0 {
			brightness = 0
		} else if brightness > 255 {
			brightness = 255
		}
		col := rl.NewColor(255, 255, 255, uint8(brightness))
		drawSphere(star.Position, 1, col)
	}

	// Se uma explosão estiver ativa, desenha o efeito de explosão
	if sim.ExplosionActive {
		maxExplosionRadius := float32(30)
		explosionRadius := float32(sim.ExplosionTime/sim.ExplosionDuration) * maxExplosionRadius
		alpha := uint8(255 * (1 - float32(sim.ExplosionTime)/float32(sim.ExplosionDuration)))
		explosionColor := rl.NewColor(255, 200, 0, alpha)
		rl.DrawSphere(sim.ExplosionPosition, explosionRadius, explosionColor)
	}
}

func main() {
	// Configurações da janela
	screenWidth := int32(1280)
	screenHeight := int32(720)
	rl.InitWindow(screenWidth, screenHeight, "Simulação 3D Realista do Sistema Solar - Câmeras, Física & Colisões")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	// Cria a câmera 3D com parâmetros iniciais (modo normal)
	camera := rl.Camera3D{
		Position:   rl.NewVector3(0, 300, 600),
		Target:     rl.NewVector3(0, 0, 0),
		Up:         rl.NewVector3(0, 1, 0),
		Fovy:       45,
		Projection: rl.CameraPerspective,
	}

	// Salvamos o estado normal da câmera para restaurá-lo depois do modo Top View.
	normalCamera := camera

	// Variável que guarda o modo de câmera para os controles normais (Orbital ou Livre)
	currentCameraMode := CameraOrbital

	// Variável que indica se o modo Top View está ativo
	topViewEnabled := false

	// Gera o mesh e o modelo para o anel de Saturno
	ringMesh := generateRingMesh(1.5, 2.0, 100)
	ringModel := rl.LoadModelFromMesh(ringMesh)
	defer rl.UnloadModel(ringModel)

	sim := NewSimulation()

	for !rl.WindowShouldClose() {
		sim.Update()

		// Alterna entre o modo Top View e o normal ao pressionar a tecla P.
		// No modo Top View, a câmera fica fixa, sem reagir ao mouse.
		if rl.IsKeyPressed(rl.KeyP) {
			if !topViewEnabled {
				normalCamera = camera
				camera.Position = rl.NewVector3(0, 800, 0)
				camera.Target = rl.NewVector3(0, 0, 0)
				camera.Up = rl.NewVector3(0, 0, -1)
				camera.Projection = rl.CameraPerspective
				camera.Fovy = 45
				topViewEnabled = true
			} else {
				topViewEnabled = false
				camera = normalCamera
			}
		}

		// Se não estiver no modo Top View, atualiza a câmera com base nas entradas do usuário.
		if !topViewEnabled {
			if rl.IsKeyPressed(rl.KeyOne) {
				currentCameraMode = CameraOrbital
			}
			if rl.IsKeyPressed(rl.KeyTwo) {
				currentCameraMode = CameraFree
			}
			rl.UpdateCamera(&camera, rl.CameraMode(currentCameraMode))
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.BeginMode3D(camera)
		sim.Draw3D(ringModel)
		rl.EndMode3D()

		// Exibe informações na tela
		modeText := ""
		if topViewEnabled {
			modeText = "Top View"
		} else {
			if currentCameraMode == CameraOrbital {
				modeText = "Orbital"
			} else {
				modeText = "Livre"
			}
		}
		rl.DrawText("Simulação 3D Realista do Sistema Solar", 10, 10, 20, rl.White)
		rl.DrawText("Modo da Câmera: "+modeText, 10, 40, 20, rl.White)
		rl.DrawText("Pressione 1: Orbital | 2: Livre (modo normal)", 10, 70, 20, rl.White)
		rl.DrawText("Pressione P: Alternar Top View", 10, 100, 20, rl.White)

		rl.EndDrawing()
	}
}
