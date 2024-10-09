package main

import (
	"image/color"
	"math/rand"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

const (
	width  = 60 // Width of the game area in segments
	height = 40 // Height of the game area in segments
	size   = 20 // Size of each segment (square) in pixels
)

type Point struct {
	x, y int
}

type Game struct {
	head      Point
	direction Point // Direction in which the snake moves
	food      []Point
	rocks     []Point
	snake     []Point       // List to track snake body segments
	score     int           // Score counter for food eaten
	scoreText *canvas.Text  // Text object to display the score
	gameOver  bool          // Flag to indicate if the game is over
	speed     time.Duration // Speed of the game
}

func newGame() *Game {
	scoreText := canvas.NewText("Score: 0", color.White)
	scoreText.TextSize = 18
	scoreText.Move(fyne.NewPos(float32(width*size-100), 10)) // Position the score text
	return &Game{
		head:      Point{x: width / 2, y: height / 2}, // Starting position of the snake head
		direction: Point{x: 0, y: 0},                  // Initial direction is set to zero
		food:      generateFood(10),                   // Start with 10 food pieces
		rocks:     generateRocks(10),                  // Start with 10 rocks
		snake:     []Point{{5, 5}},                    // Start with the snake's head
		score:     0,                                  // Initial score
		scoreText: scoreText,                          // Initialize the score text
		speed:     time.Millisecond * 80,              // Initial speed
	}
}

// Generate random food positions
func generateFood(count int) []Point {
	food := make([]Point, 0, count)
	occupied := make(map[Point]bool)

	for len(food) < count {
		newFood := Point{x: rand.Intn(width), y: rand.Intn(height)}
		if !occupied[newFood] && !occupied[Point{x: 5, y: 5}] { // Ensure it doesn't overlap with the snake's head
			food = append(food, newFood)
			occupied[newFood] = true
		}
	}
	return food
}

// Generate random rock positions
func generateRocks(count int) []Point {
	rocks := make([]Point, 0, count)
	occupied := make(map[Point]bool)

	for len(rocks) < count {
		newRock := Point{x: rand.Intn(width), y: rand.Intn(height)}
		if !occupied[newRock] && !occupied[Point{x: 5, y: 5}] { // Ensure it doesn't overlap with the snake's head
			rocks = append(rocks, newRock)
			occupied[newRock] = true
		}
	}
	return rocks
}

func (g *Game) CreateHead() *canvas.Rectangle {
	headSquare := canvas.NewRectangle(&color.RGBA{R: 255, G: 255, B: 255, A: 255})
	headSquare.Resize(fyne.NewSize(size, size))
	headSquare.Move(fyne.NewPos(float32(g.head.x*size), float32(g.head.y*size)))
	return headSquare
}

func (g *Game) CreateFood() []*canvas.Circle {
	var foodSquares []*canvas.Circle
	for _, f := range g.food {
		foodSquare := canvas.NewCircle(&color.RGBA{R: 0, G: 255, B: 0, A: 255}) // Red food
		foodSquare.Resize(fyne.NewSize(size, size))
		foodSquare.Move(fyne.NewPos(float32(f.x*size), float32(f.y*size)))
		foodSquares = append(foodSquares, foodSquare)
	}
	return foodSquares
}

func (g *Game) CreateRocks() []*canvas.Rectangle {
	var rockSquares []*canvas.Rectangle
	for _, f := range g.rocks {
		rockSquare := canvas.NewRectangle(&color.RGBA{R: 255, G: 0, B: 0, A: 255})
		rockSquare.Resize(fyne.NewSize(size, size))
		rockSquare.Move(fyne.NewPos(float32(f.x*size), float32(f.y*size)))
		rockSquares = append(rockSquares, rockSquare)
	}
	return rockSquares
}

func (g *Game) Update() {
	if g.gameOver {
		return // Don't update if the game is over
	}

	// Update the position of the snake head based on the direction
	newX := g.head.x + g.direction.x
	newY := g.head.y + g.direction.y

	// Check for boundaries and update head position if within limits
	if newX >= 0 && newX < width && newY >= 0 && newY < height {
		// Move the snake's body
		if len(g.snake) > 0 {
			// Move each segment to the position of the previous segment
			for i := len(g.snake) - 1; i > 0; i-- {
				g.snake[i] = g.snake[i-1]
			}
			// Move the first segment to the previous position of the head
			g.snake[0] = g.head
		}

		// Check for food collision
		for i, f := range g.food {
			if newX == f.x && newY == f.y {
				// Remove the food and grow the snake
				g.food = append(g.food[:i], g.food[i+1:]...)               // Remove food from the list
				g.snake = append(g.snake, Point{x: g.head.x, y: g.head.y}) // Add a new segment at the end
				g.food = append(g.food, generateFood(1)[0])                // Add a new piece of food
				g.rocks = append(g.rocks, generateRocks(1)[0])             // Add a new rock
				g.score++                                                  // Increment the score
				g.scoreText.Text = "Score: " + strconv.Itoa(g.score)       // Update the score text
				g.scoreText.Refresh()                                      // Refresh the score text
				if time.Millisecond*20 < g.speed-2*time.Millisecond {
					g.speed -= 2 * time.Millisecond
				} else {
					g.speed = time.Millisecond * 20
				}
				break
			}
		}

		// Check for rock collision
		for _, f := range g.rocks {
			if newX == f.x && newY == f.y {
				g.gameOver = true
				break
			}
		}

		// Update the head's position
		g.head.x = newX
		g.head.y = newY

		// Check for self-collision
		for _, segment := range g.snake[1:] { // Start from 1 to avoid checking the head
			if g.head.x == segment.x && g.head.y == segment.y {
				// Stop the game if it collides with itself
				g.gameOver = true
				break
			}
		}

	} else {
		g.direction = Point{0, 0} // Stop moving if hitting the wall
	}
}

func (g *Game) DrawSnake() []*canvas.Circle {
	var snakeSquares []*canvas.Circle
	for _, segment := range g.snake {
		snakeSquare := canvas.NewCircle(&color.RGBA{R: 255, G: 255, B: 255, A: 255}) // White snake body
		snakeSquare.Resize(fyne.NewSize(size, size))
		snakeSquare.Move(fyne.NewPos(float32(segment.x*size), float32(segment.y*size)))
		snakeSquares = append(snakeSquares, snakeSquare)
	}
	return snakeSquares
}

func (g *Game) GameOverDisplay(w fyne.Window) *fyne.Container {
	// Create game over text
	gameOverText := canvas.NewText("Game Over!", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	gameOverText.TextSize = 24 // Increase text size for visibility
	gameOverText.Refresh()     // Refresh to apply the size change

	// Create score text
	scoreText := canvas.NewText("Score: "+strconv.Itoa(g.score), color.RGBA{R: 255, G: 255, B: 255, A: 255})
	scoreText.TextSize = 18 // Set text size
	scoreText.Refresh()     // Refresh to apply the size change

	// Create instruction text
	instructionText := canvas.NewText("Press Spacebar to play again", color.RGBA{R: 255, G: 255, B: 255, A: 255})
	instructionText.TextSize = 18 // Set text size
	instructionText.Refresh()     // Refresh to apply the size change

	// Center the texts manually using MinSize
	width, height := w.Canvas().Size().Width, w.Canvas().Size().Height
	gameOverText.Move(fyne.NewPos((width-gameOverText.MinSize().Width)/2, height/3))
	scoreText.Move(fyne.NewPos((width-scoreText.MinSize().Width)/2, height/2))
	instructionText.Move(fyne.NewPos((width-instructionText.MinSize().Width)/2, height*2/3))

	// Create a container to hold the text
	return container.NewWithoutLayout(gameOverText, scoreText, instructionText)
}

func (g *Game) Reset() {
	// Reset the game state
	g.head = Point{x: width / 2, y: height / 2}
	g.direction = Point{x: 0, y: 0}
	g.snake = []Point{{5, 5}}
	g.food = generateFood(10)
	g.rocks = generateRocks(10)
	g.score = 0
	g.speed = 80 * time.Millisecond // Reset speed
	g.gameOver = false
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	game := newGame()

	a := app.New()
	w := a.NewWindow("Snake Game")
	w.Resize(fyne.NewSize(width*size, height*size)) // Resize window to fit the game area

	// Create a black background container
	background := canvas.NewRectangle(&color.RGBA{R: 0, G: 0, B: 0, A: 255}) // Black background
	background.Resize(fyne.NewSize(width*size, height*size))

	// Create a container for the game area
	gameContainer := container.NewWithoutLayout(background)

	// Create the snake head
	head := game.CreateHead()
	gameContainer.Add(head)

	// Create food items
	foodItems := game.CreateFood()
	for _, food := range foodItems {
		gameContainer.Add(food)
	}

	// Create rocks
	rockItems := game.CreateRocks()
	for _, rock := range rockItems {
		gameContainer.Add(rock)
	}

	// Add the score text to the container
	gameContainer.Add(game.scoreText)

	// Set the content of the window to the game container
	w.SetContent(gameContainer)

	// Handle key events for arrow keys
	w.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		if game.gameOver {
			if ke.Name == fyne.KeySpace { // Check if the spacebar is pressed
				game.Reset()
				gameContainer.Objects = append(gameContainer.Objects[:0], background) // Clear previous drawings
				gameContainer.Add(game.CreateHead())
				foodItems = game.CreateFood() // Generate new food
				for _, food := range foodItems {
					gameContainer.Add(food) // Add food to the container
				}
				for _, rock := range rockItems {
					gameContainer.Add(rock)
				}
				gameContainer.Add(game.scoreText) // Add the score text to the container
			}
			return // Ignore other key presses if the game is over
		}
		switch ke.Name {
		case fyne.KeyUp:
			if game.direction.y == 0 { // Prevent moving in the opposite direction
				game.direction = Point{x: 0, y: -1}
			}
		case fyne.KeyDown:
			if game.direction.y == 0 {
				game.direction = Point{x: 0, y: 1}
			}
		case fyne.KeyLeft:
			if game.direction.x == 0 {
				game.direction = Point{x: -1, y: 0}
			}
		case fyne.KeyRight:
			if game.direction.x == 0 {
				game.direction = Point{x: 1, y: 0}
			}
		}
	})

	// Game loop to update the snake's position
	go func() {
		for {
			time.Sleep(game.speed) // Control the game update rate
			game.Update()          // Update the game state

			// Clear the previous drawings
			gameContainer.Objects = []fyne.CanvasObject{background} // Keep only the background

			// Draw the snake
			snakeSegments := game.DrawSnake()
			for _, segment := range snakeSegments {
				gameContainer.Add(segment)
			}

			// Draw the food
			foodItems = game.CreateFood()
			for _, food := range foodItems {
				gameContainer.Add(food)
			}

			rockItems = game.CreateRocks()
			for _, rock := range rockItems {
				gameContainer.Add(rock)
			}

			// Add the score text
			gameContainer.Add(game.scoreText) // Add the score text to the container

			// Check for game over and display message if needed
			if game.gameOver {
				gameContainer.Add(game.GameOverDisplay(w)) // Display the game over message
			}

			// Refresh the container to redraw everything
			gameContainer.Refresh()
		}
	}()

	w.ShowAndRun()
}
