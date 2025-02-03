package main

import (
	"math"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Definição de uma cor Silver (já que rl.Silver não está definida)
var Silver = rl.NewColor(192, 192, 192, 255)

// ─────────────────────────────────────────────
// Estruturas da simulação (a lógica permanece semelhante, com mais objetos)
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
	CamAngle  float64 // Para movimentação dinâmica da câmera
}

// NewSimulation cria e inicializa os corpos celestes
func NewSimulation() *Simulation {
	sim := &Simulation{
		SunRadius: 40,
		Planets:   make([]*Planet, 0),
		Time:      0,
		CamAngle:  0,
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
	// Jupiter com múltiplas luas para enriquecer a cena
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
	// Kuiper Belt (asteroides extras além de Plutão)
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

	// Cometa (com rastro)
	sim.Comet = Comet{
		Position:      rl.NewVector3(-50, 0, -50),
		Angle:         math.Pi / 4,
		Speed:         4.0,
		TailPoints:    make([]rl.Vector3, 0),
		TailMaxLength: 20,
	}

	rand.Seed(time.Now().UnixNano())
	return sim
}

func (sim *Simulation) Update() {
	sim.Time += 1.0 / 60.0

	// Atualiza as fases das estrelas para cintilação
	for i := range sim.Stars {
		sim.Stars[i].Phase += sim.Stars[i].Speed
	}

	// Atualiza os ângulos dos asteroides
	for i := range sim.Asteroids {
		sim.Asteroids[i].Angle += sim.Asteroids[i].OrbitSpeed
	}

	// Atualiza posição e rastro do cometa
	sim.Comet.Position.X += float32(sim.Comet.Speed * math.Cos(sim.Comet.Angle))
	sim.Comet.Position.Z += float32(sim.Comet.Speed * math.Sin(sim.Comet.Angle))
	sim.Comet.TailPoints = append([]rl.Vector3{sim.Comet.Position}, sim.Comet.TailPoints...)
	if len(sim.Comet.TailPoints) > sim.Comet.TailMaxLength {
		sim.Comet.TailPoints = sim.Comet.TailPoints[:sim.Comet.TailMaxLength]
	}
	if sim.Comet.Position.X > 800 || sim.Comet.Position.Z > 800 {
		sim.Comet.Position = rl.NewVector3(-50, 0, -50)
		sim.Comet.TailPoints = sim.Comet.TailPoints[:0]
	}

	// Atualiza ângulos dos planetas e de suas luas
	for _, p := range sim.Planets {
		p.Angle += p.OrbitSpeed
		for _, m := range p.Moons {
			m.Angle += m.OrbitSpeed
		}
	}

	// Atualiza o ângulo da câmera para movimento orbital suave
	sim.CamAngle += 0.001
}

// Desenha as órbitas dos planetas (no plano XZ)
func (sim *Simulation) DrawOrbitPaths() {
	center := rl.NewVector3(0, 0, 0)
	for _, p := range sim.Planets {
		rl.DrawCircle3D(center, float32(p.OrbitRadius), rl.NewVector3(1, 0, 0), 90, rl.LightGray)
	}
}

// ─────────────────────────────────────────────
// Shader de Iluminação Phong modificado para luz ponto (o Sol é a fonte)
const vertexShaderSource = `#version 330
in vec3 vertexPosition;
in vec3 vertexNormal;
uniform mat4 mvp;
uniform mat4 model;
out vec3 fragPos;
out vec3 normal;
void main() {
    fragPos = vec3(model * vec4(vertexPosition, 1.0));
    normal = mat3(transpose(inverse(model))) * vertexNormal;
    gl_Position = mvp * vec4(vertexPosition, 1.0);
}`

const fragmentShaderSource = `#version 330
in vec3 fragPos;
in vec3 normal;
uniform vec3 lightPos;      // Posição do Sol
uniform vec3 lightColor;
uniform vec3 ambient;
uniform vec3 viewPos;
uniform float shininess;
uniform vec3 objectColor;
out vec4 finalColor;
void main() {
    // Componente ambiente
    vec3 ambientComponent = ambient * objectColor;
    // Luz difusa
    vec3 norm = normalize(normal);
    vec3 lightDir = normalize(lightPos - fragPos);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = diff * lightColor * objectColor;
    // Componente especular
    vec3 viewDir = normalize(viewPos - fragPos);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), shininess);
    vec3 specular = spec * lightColor;
    // Atenuação (queda de intensidade com a distância)
    float distance = length(lightPos - fragPos);
    float attenuation = 1.0 / (distance * distance * 0.0005 + 1.0);
    vec3 result = (ambientComponent + diffuse + specular) * attenuation;
    finalColor = vec4(result, 1.0);
}`

// ─────────────────────────────────────────────
// Função para desenhar uma esfera com o shader customizado
func drawLitSphere(model rl.Model, shader rl.Shader, pos rl.Vector3, radius float32, col rl.Color) {
	objColor := []float32{
		float32(col.R) / 255.0,
		float32(col.G) / 255.0,
		float32(col.B) / 255.0,
	}
	rl.SetShaderValue(shader, rl.GetShaderLocation(shader, "objectColor"), objColor, rl.ShaderUniformVec3)
	rl.DrawModelEx(model, pos, rl.NewVector3(0, 1, 0), 0, rl.NewVector3(radius, radius, radius), rl.White)
}

// ─────────────────────────────────────────────
// Gera um mesh para um anel (para Saturno – no plano XZ)
func generateRingMesh(innerRadius, outerRadius float32, segments int) rl.Mesh {
	vertexCount := segments * 2
	triangleCount := segments * 2

	vertices := make([]float32, vertexCount*3)  // x,y,z para cada vértice
	normals := make([]float32, vertexCount*3)   // normais
	texcoords := make([]float32, vertexCount*2) // u,v
	indices := make([]uint16, triangleCount*3)

	angleStep := 2 * math.Pi / float64(segments)
	vertexIndex := 0
	texIndex := 0
	for i := 0; i < segments; i++ {
		angle := float64(i) * angleStep
		cosA := float32(math.Cos(angle))
		sinA := float32(math.Sin(angle))
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

	mesh := rl.Mesh{
		VertexCount: int32(vertexCount),
		Vertices:    &vertices[0],
		Normals:     &normals[0],
		Texcoords:   &texcoords[0],
		Indices:     &indices[0],
	}
	return mesh
}

// ─────────────────────────────────────────────
// Desenha a cena 3D usando o shader customizado e os modelos gerados.
// (A funcionalidade de skybox foi removida para evitar erros de compilação.)
func (sim *Simulation) Draw3D(sphereModel, ringModel rl.Model, shader rl.Shader, camera rl.Camera3D) {
	// Atualiza os uniforms do shader
	lightPos := []float32{0.0, 0.0, 0.0}
	lightColor := []float32{1.0, 1.0, 1.0}
	ambient := []float32{0.7, 0.7, 0.7}
	shininess := []float32{32.0}

	rl.SetShaderValue(shader, rl.GetShaderLocation(shader, "lightPos"), lightPos, rl.ShaderUniformVec3)
	rl.SetShaderValue(shader, rl.GetShaderLocation(shader, "lightColor"), lightColor, rl.ShaderUniformVec3)
	rl.SetShaderValue(shader, rl.GetShaderLocation(shader, "ambient"), ambient, rl.ShaderUniformVec3)
	rl.SetShaderValue(shader, rl.GetShaderLocation(shader, "shininess"), shininess, rl.ShaderUniformFloat)
	rl.SetShaderValue(shader, rl.GetShaderLocation(shader, "viewPos"),
		[]float32{camera.Position.X, camera.Position.Y, camera.Position.Z}, rl.ShaderUniformVec3)

	// Desenha o Sol
	drawLitSphere(sphereModel, shader, rl.NewVector3(0, 0, 0), sim.SunRadius, rl.Yellow)

	// Desenha as órbitas dos planetas
	sim.DrawOrbitPaths()

	// Desenha os planetas e suas luas
	for _, p := range sim.Planets {
		x := float32(p.OrbitRadius * math.Cos(p.Angle))
		z := float32(p.OrbitRadius * math.Sin(p.Angle))
		planetPos := rl.NewVector3(x, 0, z)
		drawLitSphere(sphereModel, shader, planetPos, p.Radius, p.Color)
		// Se for Saturn, desenha os anéis
		if p.Name == "Saturn" {
			rl.DrawModelEx(ringModel, planetPos, rl.NewVector3(1, 0, 0), 25, rl.NewVector3(p.Radius*3, 1, p.Radius*3), rl.LightGray)
		}
		for _, m := range p.Moons {
			mx := planetPos.X + float32(m.OrbitRadius*math.Cos(m.Angle))
			mz := planetPos.Z + float32(m.OrbitRadius*math.Sin(m.Angle))
			moonPos := rl.NewVector3(mx, 0, mz)
			drawLitSphere(sphereModel, shader, moonPos, m.Radius, m.Color)
		}
	}
	// Desenha os asteroides
	for _, a := range sim.Asteroids {
		ax := float32(a.OrbitRadius * math.Cos(a.Angle))
		az := float32(a.OrbitRadius * math.Sin(a.Angle))
		asteroidPos := rl.NewVector3(ax, 0, az)
		drawLitSphere(sphereModel, shader, asteroidPos, a.Radius, rl.Gray)
	}
	// Desenha o rastro do cometa
	for i := 0; i < len(sim.Comet.TailPoints)-1; i++ {
		alpha := uint8(200 * (1 - float32(i)/float32(len(sim.Comet.TailPoints))))
		col := rl.NewColor(255, 255, 255, alpha)
		drawLitSphere(sphereModel, shader, sim.Comet.TailPoints[i], 2, col)
	}
	// Desenha o cometa
	drawLitSphere(sphereModel, shader, sim.Comet.Position, 4, rl.White)
	// Desenha as estrelas cintilantes
	for _, star := range sim.Stars {
		brightness := float32(128 + 127*math.Sin(star.Phase))
		if brightness < 0 {
			brightness = 0
		} else if brightness > 255 {
			brightness = 255
		}
		col := rl.NewColor(255, 255, 255, uint8(brightness))
		drawLitSphere(sphereModel, shader, star.Position, 1, col)
	}

	// OBS.: A função de skybox foi removida para evitar erros (rl.DrawSkybox não está disponível nesta versão).
}

func main() {
	// Configurações da janela e MSAA
	screenWidth := int32(1280)
	screenHeight := int32(720)
	rl.SetConfigFlags(32)
	rl.InitWindow(screenWidth, screenHeight, "Simulação 3D Realista do Sistema Solar")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	// Cria a câmera 3D
	camera := rl.Camera3D{
		Position:   rl.NewVector3(0, 300, 600),
		Target:     rl.NewVector3(0, 0, 0),
		Up:         rl.NewVector3(0, 1, 0),
		Fovy:       45,
		Projection: rl.CameraPerspective,
	}

	// Carrega o shader de iluminação
	shader := rl.LoadShaderFromMemory(vertexShaderSource, fragmentShaderSource)
	defer rl.UnloadShader(shader)

	// Modelo de esfera (alta resolução para planetas, luas, etc.)
	sphereMesh := rl.GenMeshSphere(1.0, 32, 32)
	sphereModel := rl.LoadModelFromMesh(sphereMesh)
	sphereModel.Materials.Shader = shader
	defer rl.UnloadModel(sphereModel)

	// Modelo para o anel de Saturno
	ringMesh := generateRingMesh(1.5, 2.0, 100)
	ringModel := rl.LoadModelFromMesh(ringMesh)
	ringModel.Materials.Shader = shader
	defer rl.UnloadModel(ringModel)

	// OBS.: O código para carregar o cubemap (skybox) foi removido,
	// pois as funções rl.LoadTextureCubemap e rl.DrawSkybox não estão definidas na sua versão.

	sim := NewSimulation()

	// Loop principal
	for !rl.WindowShouldClose() {
		// Atualiza a simulação e a câmera (movimento orbital suave)
		sim.Update()
		camRadius := float32(600)
		camera.Position.X = camRadius * float32(math.Cos(sim.CamAngle))
		camera.Position.Z = camRadius * float32(math.Sin(sim.CamAngle))
		camera.Target = rl.NewVector3(0, 0, 0)

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.BeginMode3D(camera)
		// Desenha a cena (sem skybox)
		sim.Draw3D(sphereModel, ringModel, shader, camera)
		rl.EndMode3D()

		rl.DrawText("Simulação 3D Realista do Sistema Solar", 10, 10, 20, rl.White)
		rl.EndDrawing()
	}
}
