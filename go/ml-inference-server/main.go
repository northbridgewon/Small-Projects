package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	gorgonia "gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

const (
	inputSize  = 4  // Example: 4 features
	hiddenSize = 5
	outputSize = 3  // Example: 3 classes
)

// NeuralNetwork represents a simple feed-forward neural network.
type NeuralNetwork struct {
	graph *gorgonia.ExprGraph
	w1, b1, w2, b2 *gorgonia.Node
	output         *gorgonia.Node
	machine        gorgonia.VM
}

// NewNeuralNetwork creates and initializes a simple neural network.
func NewNeuralNetwork() (*NeuralNetwork, error) {
	g := gorgonia.NewGraph()

	// Input layer
	input := gorgonia.NewMatrix(g, tensor.Float64, gorgonia.WithShape(1, inputSize), gorgonia.WithName("input"))

	// Hidden layer weights and biases
	w1 := gorgonia.NewMatrix(g, tensor.Float64, gorgonia.WithShape(inputSize, hiddenSize), gorgonia.WithName("w1"), gorgonia.WithInit(gorgonia.GlorotU(1)))
	b1 := gorgonia.NewMatrix(g, tensor.Float64, gorgonia.WithShape(1, hiddenSize), gorgonia.WithName("b1"), gorgonia.WithInit(gorgonia.Zeroes()))

	// Output layer weights and biases
	w2 := gorgonia.NewMatrix(g, tensor.Float64, gorgonia.WithShape(hiddenSize, outputSize), gorgonia.WithName("w2"), gorgonia.WithInit(gorgonia.GlorotU(1)))
	b2 := gorgonia.NewMatrix(g, tensor.Float64, gorgonia.WithShape(1, outputSize), gorgonia.WithName("b2"), gorgonia.WithInit(gorgonia.Zeroes()))

	// Hidden layer calculation: input * w1 + b1
	hidden1, err := gorgonia.Mul(input, w1)
	if err != nil {
		return nil, fmt.Errorf("failed to multiply input and w1: %w", err)
	}
	hidden1, err = gorgonia.Add(hidden1, b1)
	if err != nil {
		return nil, fmt.Errorf("failed to add b1: %w", err)
	}
	hidden1 = gorgonia.Must(gorgonia.Rectify(hidden1)) // ReLU activation

	// Output layer calculation: hidden1 * w2 + b2
	output, err := gorgonia.Mul(hidden1, w2)
	if err != nil {
		return nil, fmt.Errorf("failed to multiply hidden1 and w2: %w", err)
	}
	output, err = gorgonia.Add(output, b2)
	if err != nil {
		return nil, fmt.Errorf("failed to add b2: %w", err)
	}
	output = gorgonia.Must(gorgonia.SoftMax(output)) // Softmax activation for probabilities

	// Create VM to run the graph
	m := gorgonia.NewTapeMachine(g, gorgonia.BindDualValues(w1, b1, w2, b2))

	return &NeuralNetwork{
		graph:   g,
		w1:      w1,
		b1:      b1,
		w2:      w2,
	b2:      b2,
		output:  output,
		machine: m,
	}, nil
}

// Predict performs inference on the given input data.
func (nn *NeuralNetwork) Predict(inputData []float64) ([]float64, error) {
	if len(inputData) != inputSize {
		return nil, fmt.Errorf("input data size mismatch: expected %d, got %d", inputSize, len(inputData))
	}

	// Set input value
	inputTensor := tensor.New(tensor.WithShape(1, inputSize), tensor.WithBacking(inputData))
	gorgonia.Let(nn.graph.Inputs()[0], inputTensor)

	// Run the computation graph
	if err := nn.machine.RunAll(); err != nil {
		return nil, fmt.Errorf("failed to run machine: %w", err)
	}

	// Get output value
	outputValue := nn.output.Value().Data().([]float64)

	// Reset the machine for the next prediction
	nn.machine.Reset()

	return outputValue, nil
}

// InferenceRequest represents the structure of the incoming JSON request.
type InferenceRequest struct {
	Data []float64 `json:"data"`
}

// InferenceResponse represents the structure of the outgoing JSON response.
type InferenceResponse struct {
	Predictions []float64 `json:"predictions"`
	Error       string    `json:"error,omitempty"`
}

// inferenceHandler handles incoming inference requests.
func inferenceHandler(nn *NeuralNetwork) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
			return
		}

		var req InferenceRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}

		predictions, err := nn.Predict(req.Data)
		resp := InferenceResponse{}
		if err != nil {
			resp.Error = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			resp.Predictions = predictions
			w.WriteHeader(http.StatusOK)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func main() {
	// Initialize the neural network
	nn, err := NewNeuralNetwork()
	if err != nil {
		log.Fatalf("Failed to initialize neural network: %v", err)
	}

	// Set up HTTP server
	http.HandleFunc("/predict", inferenceHandler(nn))

	port := 7000 // Default port
	fmt.Printf("ML Inference Server started on :%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
