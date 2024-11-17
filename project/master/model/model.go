package model

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
)

type Model struct {
	numFeatures    int
	epochs         int
	learningRate   float64
	regularization float64
	R              [][]float64
	P              [][]float64
	Q              [][]float64
}

type ModelConfig struct {
	NumFeatures    int         `json:"numFeatures"`
	Epochs         int         `json:"epochs"`
	LearningRate   float64     `json:"learningRate"`
	Regularization float64     `json:"regularization"`
	R              [][]float64 `json:"R"`
	P              [][]float64 `json:"P"`
	Q              [][]float64 `json:"Q"`
}

func LoadTrainData(filename string) ([][]float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("trainFileError: Error opening file %s: %v", filename, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("trainFileError: Error reading file %s: %v", filename, err)
	}

	maxUserId, maxMovieId := 0, 0
	for _, record := range records[1:] {
		userId, _ := strconv.Atoi(record[0])
		movieId, _ := strconv.Atoi(record[1])
		if userId > maxUserId {
			maxUserId = userId
		}
		if movieId > maxMovieId {
			maxMovieId = movieId
		}
	}

	ratingsMatrix := make([][]float64, maxUserId+1)
	for i := range ratingsMatrix {
		ratingsMatrix[i] = make([]float64, maxMovieId+1)
	}

	for _, record := range records[1:] {
		userId, _ := strconv.Atoi(record[0])
		movieId, _ := strconv.Atoi(record[1])
		rating, _ := strconv.ParseFloat(record[2], 64)
		ratingsMatrix[userId][movieId] = rating
	}

	return ratingsMatrix, nil
}

func NewModel(numFeatures, epochs int, learningRate, regularization float64, R [][]float64, randomState int) Model {
	r := rand.New(rand.NewSource(int64(randomState)))
	numUsers := len(R)
	numItems := len(R[0])
	return Model{
		numFeatures:    numFeatures,
		epochs:         epochs,
		learningRate:   learningRate,
		regularization: regularization,
		R:              R,
		P:              initMatrix(numUsers, numFeatures, r),
		Q:              initMatrix(numItems, numFeatures, r),
	}
}

func LoadModel(modelConfig *ModelConfig) Model {
	return Model{
		numFeatures:    modelConfig.NumFeatures,
		epochs:         modelConfig.Epochs,
		learningRate:   modelConfig.LearningRate,
		regularization: modelConfig.Regularization,
		R:              modelConfig.R,
		P:              modelConfig.P,
		Q:              modelConfig.Q,
	}
}

func initMatrix(rows, cols int, r *rand.Rand) [][]float64 {
	matrix := make([][]float64, rows)
	for i := range matrix {
		matrix[i] = make([]float64, cols)
		for j := range matrix[i] {
			matrix[i][j] = r.Float64()
		}
	}
	return matrix
}

func (model *Model) Train() {
	var wg sync.WaitGroup
	muQ := make([]sync.Mutex, len(model.R[0]))
	for epoch := 0; epoch < model.epochs; epoch++ {
		for userId := range model.R {
			wg.Add(1)
			go func(userId int) {
				defer wg.Done()
				for itemId := range model.R[userId] {
					if model.R[userId][itemId] != 0 {
						muQ[itemId].Lock()
						pred := model.Predict(userId, itemId)
						err := model.R[userId][itemId] - pred
						for k := 0; k < model.numFeatures; k++ {
							userGrad := model.learningRate * (err*model.Q[itemId][k] - model.regularization*model.P[userId][k])
							itemGrad := model.learningRate * (err*model.P[userId][k] - model.regularization*model.Q[itemId][k])
							model.P[userId][k] += userGrad
							model.Q[itemId][k] += itemGrad
						}
						muQ[itemId].Unlock()
					}
				}
			}(userId)
		}
		wg.Wait()
	}
}

func (model *Model) Predict(userId, itemId int) float64 {
	prediction := 0.0
	for k := 0; k < model.numFeatures; k++ {
		prediction += model.P[userId][k] * model.Q[itemId][k]
	}
	return prediction
}

func (model *Model) CalculateRMSE() float64 {
	var squaredErrorSum float64
	var regularizationSum float64
	count := 0

	for u := 0; u < len(model.P); u++ {
		for i := 0; i < len(model.Q); i++ {
			if model.R[u][i] > 0 {
				pred := model.Predict(u, i)
				err := model.R[u][i] - pred
				squaredErrorSum += math.Pow(err, 2)
				for k := 0; k < model.numFeatures; k++ {
					regularizationSum += math.Pow(model.P[u][k], 2) + math.Pow(model.Q[i][k], 2)
				}
				count++
			}
		}
	}

	rmse := math.Sqrt((squaredErrorSum + model.regularization/2*regularizationSum) / float64(count))
	return rmse
}

func (model *Model) ParamsToJson(filename string) error {
	params := ModelConfig{
		NumFeatures:    model.numFeatures,
		Epochs:         model.epochs,
		LearningRate:   model.learningRate,
		Regularization: model.regularization,
		R:              model.R,
		P:              model.P,
		Q:              model.Q,
	}

	jsonData, err := json.MarshalIndent(params, "", "\t")
	if err != nil {
		return fmt.Errorf("error marshalling params to JSON: %v", err)
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON to file: %v", err)
	}

	return nil
}

type ModelGrid struct {
	NumFeatures    []int
	Epochs         []int
	LearningRate   []float64
	Regularization []float64
}

func SearchGrid(grid ModelGrid, R [][]float64) Model {
	var bestModel Model
	bestRMSE := math.Inf(1)
	for _, numFeatures := range grid.NumFeatures {
		for _, epochs := range grid.Epochs {
			for _, learningRate := range grid.LearningRate {
				for _, regularization := range grid.Regularization {
					model := NewModel(numFeatures, epochs, learningRate, regularization, R, 1)
					fmt.Println("--------------------")
					fmt.Println("Model with:")
					fmt.Println("numFeatures:", numFeatures)
					fmt.Println("epochs:", epochs)
					fmt.Println("learningRate:", learningRate)
					fmt.Println("regularization:", regularization)
					model.Train()
					rmse := model.CalculateRMSE()
					fmt.Println("RMSE:", rmse)
					if rmse < bestRMSE {
						bestRMSE = rmse
						bestModel = model
					}
				}
			}
		}
	}
	fmt.Println("Best RMSE:", bestRMSE)
	fmt.Println("numFeatures:", bestModel.numFeatures)
	fmt.Println("epochs:", bestModel.epochs)
	fmt.Println("learningRate:", bestModel.learningRate)
	fmt.Println("regularization:", bestModel.regularization)
	return bestModel
}

func (model *Model) UpdateUserFactors(ratings []float64, userFactors *[]float64, startItemId, endItemId int) ([]float64, int) {
	weightedGrad := make([]float64, model.numFeatures)
	n := endItemId - startItemId
	count := 0
	for epoch := 0; epoch < model.epochs; epoch++ {
		for i := 0; i < n; i++ {
			itemId := startItemId + i
			if ratings[i] != 0 {
				pred := 0.0
				for k := 0; k < model.numFeatures; k++ {
					pred += (*userFactors)[k] * model.Q[itemId][k]
				}
				err := ratings[i] - pred
				for k := 0; k < model.numFeatures; k++ {
					userGrad := model.learningRate * (err*model.Q[itemId][k] - model.regularization*(*userFactors)[k])
					(*userFactors)[k] += userGrad
					weightedGrad[k] += userGrad
				}
				count++
			}
		}
	}

	for k := 0; k < model.numFeatures; k++ {
		weightedGrad[k] *= float64(count) // w * grad
	}
	return weightedGrad, count
}

func (model *Model) PredictUser(userFactors []float64, itemId int) float64 {
	prediction := 0.0
	for k := 0; k < model.numFeatures; k++ {
		prediction += userFactors[k] * model.Q[itemId][k]
	}
	return prediction
}
