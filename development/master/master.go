package master

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"recommendation-service/master/safecounts"
	"recommendation-service/model"
	"recommendation-service/syncutils"
	"sort"
	"sync"
	"time"
)

type Master struct {
	ip              string
	movieTitles     []string
	movieGenreNames []string
	movieGenreIds   [][]int
	modelConfig     model.ModelConfig
	slaveIps        []string
	slavesInfo      safecounts.SafeCounts
}

type MasterConfig struct {
	SlaveIps        []string          `json:"slaveIps"`
	MovieTitles     []string          `json:"movieTitles"`
	MovieGenreNames []string          `json:"movieGenreNames"`
	MovieGenreIds   [][]int           `json:"movieGenreIds"`
	ModelConfig     model.ModelConfig `json:"modelConfig"`
}

func (master *Master) handleSyncronization() {
	log.Println("INFO: Start synchronization")
	var wg sync.WaitGroup
	for i, ip := range master.slaveIps {
		if !master.slavesInfo.ReadStatustByIndex(i) {
			wg.Add(1)
			go func(slaveId int, ip string) {
				defer wg.Done()
				err := master.handleSlaveSync(i, ip)
				if err != nil {
					log.Println(err)
				}
			}(i, ip)
		}
	}
	wg.Wait()
	for i, status := range master.slavesInfo.ReadStatus() {
		if !status {
			log.Printf("INFO: Slave %d: Not responding.\n", i)
		} else {
			log.Printf("INFO: Slave %d: Syncronized.\n", i)
		}

	}

	for i, ip := range master.slaveIps {
		if !master.slavesInfo.ReadStatustByIndex(i) {
			go func(slaveId int, ip string) {
				for {
					time.Sleep(time.Second * 60)
					err := master.handleSlaveSync(slaveId, ip)
					if err != nil {
						log.Println(err)
						continue
					}
					break
				}
			}(i, ip)
		}
	}
}

func (master *Master) handleSlaveSync(slaveId int, ip string) error {

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, syncutils.SyncronizationPort))
	if err != nil {
		master.slavesInfo.WriteStatusByIndex(false, slaveId)
		return fmt.Errorf("syncError: Slave %d connection error: %v", slaveId, err)
	}
	defer conn.Close()

	timeout := 20 * time.Second
	conn.SetDeadline(time.Now().Add(timeout))

	err = master.sendSyncRequest(&conn)
	if err != nil {
		master.slavesInfo.WriteStatusByIndex(false, slaveId)
		return fmt.Errorf("syncError: Slave %d sync request error: %v", slaveId, err)
	}
	var response syncutils.SlaveSyncResponse
	master.receiveSyncResponse(&conn, &response)
	if response.Status == 0 {
		master.slavesInfo.WriteStatusByIndex(true, slaveId)
	}
	log.Println("INFO: Slave", slaveId, "synchronized")
	return nil
}

func (master *Master) sendSyncRequest(conn *net.Conn) error {
	request := syncutils.MasterSyncRequest{
		MasterIp:      master.ip,
		MovieGenreIds: master.movieGenreIds,
		ModelConfig:   master.modelConfig,
	}
	request.ModelConfig.R = nil
	request.ModelConfig.P = nil

	err := syncutils.SendObjectAsJsonMessage(&request, conn)
	if err != nil {
		return fmt.Errorf("syncRequestErr: Error sending request object as json: %v", err)
	}
	return nil
}

func (master *Master) receiveSyncResponse(conn *net.Conn, response *syncutils.SlaveSyncResponse) error {
	err := syncutils.ReceiveJsonMessageAsObject(response, conn)
	if err != nil {
		return fmt.Errorf("syncResponseError: Error reading data: %v", err)
	}
	return nil
}

func (master *Master) loadConfig(filename string) error {
	var config MasterConfig
	err := syncutils.LoadJsonFile(filename, &config)
	if err != nil {
		return fmt.Errorf("loadConfig: Error loading config file: %v", err)
	}
	master.slaveIps = config.SlaveIps
	master.movieTitles = config.MovieTitles
	master.movieGenreNames = config.MovieGenreNames
	master.movieGenreIds = config.MovieGenreIds
	master.modelConfig = config.ModelConfig
	log.Println("INFO: Config loaded")
	return nil
}

func (master *Master) Init() error {
	master.ip = syncutils.GetOwnIp()
	err := master.loadConfig("config/master.json")
	if err != nil {
		return fmt.Errorf("initError: Error loading config: %v", err)
	}
	Banner()
	numSlaves := len(master.slaveIps)
	master.slavesInfo.Counts = make([]int, numSlaves)
	master.slavesInfo.Status = make([]bool, numSlaves)

	return nil
}

func (master *Master) Run() error {
	log.Println("INFO: Running")
	defer log.Println("INFO: Stopped")

	master.handleSyncronization()
	master.handleService()

	return nil
}

const handleServicePrefix = "handleService"

func (master *Master) handleService() {
	http.HandleFunc("/recommendations", master.serviceRecommendation)
	http.HandleFunc("/movies/titles", master.moviesTitlesHandler)
	http.HandleFunc("/genres", master.genresHandler)
	http.HandleFunc("/genres/movies", master.getMoviesByGenresHandler)
	http.HandleFunc("/movies/genres", master.MoviesGenresHandler)

	serviceAdress := syncutils.JoinAddress(master.ip, syncutils.ServicePort)

	log.Printf("INFO: %s: Service running on %s", handleServicePrefix, serviceAdress)
	defer log.Printf("INFO: %s: Service stopped", handleServicePrefix)

	if err := http.ListenAndServe(serviceAdress, enableCORS(http.DefaultServeMux)); err != nil {
		log.Printf("ERROR: %s: Server initialization error: %v\n", handleServicePrefix, err)
	}
}

func (master *Master) serviceRecommendation(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		log.Println("POST/", request.URL.Path)
		master.handleRecommendation(&response, request)
	default:
		http.Error(response, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

const handleRecommendationPrefix = "handleRec"

func (master *Master) handleRecommendation(apiResponse *http.ResponseWriter, apiRequest *http.Request) {

	log.Printf("INFO: %s: Handling recommendation\n", handleRecommendationPrefix)
	defer log.Printf("INFO: %s: Recommendation handled\n", handleRecommendationPrefix)

	var request ClientRecToSend
	err := receiveRecommendationRequest(apiResponse, apiRequest, &request)

	if err != nil {
		log.Printf("ERROR: %s: %v\n", handleRecommendationPrefix, err)
		return
	}

	clientRecRequest := syncutils.ClientRecRequest{
		UserId:   request.UserId,
		Quantity: request.Quantity,
		GenreIds: request.GenreIds,
	}
	moviesTitle := MoviesTitles{Title: master.movieTitles}

	clientRecRequest.Ratings = MappRatingsClient(request.MoviesRatings, &moviesTitle)

	var response syncutils.MasterRecResponse
	err = master.processRecommendationRequest(apiResponse, &response, &clientRecRequest)
	if err != nil {
		log.Printf("ERROR: %s: %v\n", handleRecommendationPrefix, err)
		return
	}
	err = respondRecommendationRequest(apiResponse, &response)
	if err != nil {
		log.Printf("ERROR: %s: %v\n", handleRecommendationPrefix, err)
		return
	}
}

const receiveRecommendationRequestPrefix = "receiveRecRequest"

func receiveRecommendationRequest(apiResponse *http.ResponseWriter, apiRequest *http.Request, request *ClientRecToSend) error {
	err := json.NewDecoder(apiRequest.Body).Decode(request)
	if err != nil {
		http.Error(*apiResponse, "Invalid request payload", http.StatusBadRequest)
		return fmt.Errorf("%s: Error decoding request: %v", receiveRecommendationRequestPrefix, err)
	}
	return nil
}

const processRecommendationRequestPrefix = "processRecRequest"

func (master *Master) processRecommendationRequest(apiResponse *http.ResponseWriter, response *syncutils.MasterRecResponse, request *syncutils.ClientRecRequest) error {
	var predictions []syncutils.Prediction
	var sum float64
	var max float64
	var min float64
	var count int

	if len(request.Ratings) != len(master.modelConfig.Q) {
		http.Error(*apiResponse, "Invalid request payload", http.StatusBadRequest)
		return fmt.Errorf("%s: Incorrect ratings quantity", processRecommendationRequestPrefix)
	}

	err := master.handleModelRecommendation(&predictions, &sum, &max, &min, &count, request)
	if err != nil {
		http.Error(*apiResponse, "Internal server error", http.StatusInternalServerError)
		return fmt.Errorf("%s: %v", processRecommendationRequestPrefix, err)
	}

	(*response).UserId = request.UserId
	if count == 0 {
		(*response).Recommendations = []syncutils.Recommendation{}
		return nil
	}

	(*response).Recommendations = make([]syncutils.Recommendation, len(predictions))
	mean := sum / float64(count)

	for i, prediction := range predictions {
		(*response).Recommendations[i].Id = prediction.MovieId
		(*response).Recommendations[i].Title = master.movieTitles[prediction.MovieId]
		(*response).Recommendations[i].Rating = prediction.Rating
		(*response).Recommendations[i].Genres = []string{}
		movieGenreIds := master.movieGenreIds[prediction.MovieId]
		(*response).Recommendations[i].Genres = make([]string, len(movieGenreIds))
		for j, genreId := range movieGenreIds {
			(*response).Recommendations[i].Genres[j] = master.movieGenreNames[genreId]
		}
		(*response).Recommendations[i].Comment = getComment(prediction.Rating, max, min, mean)
	}

	return nil
}

const respondRecommendationRequestPrefix = "respondRecRequest"

func respondRecommendationRequest(apiResponse *http.ResponseWriter, response *syncutils.MasterRecResponse) error {
	bytes, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		http.Error(*apiResponse, "Error marshalling response", http.StatusInternalServerError)
		return fmt.Errorf("%s: Error marshalling response: %v", respondRecommendationRequestPrefix, err)
	}
	_, err = io.Writer.Write(*apiResponse, bytes)
	if err != nil {
		http.Error(*apiResponse, "Error writing response", http.StatusInternalServerError)
		return fmt.Errorf("%s: Error writing response: %v", respondRecommendationRequestPrefix, err)
	}
	(*apiResponse).Header().Set("Content-Type", "application/json")
	(*apiResponse).WriteHeader(http.StatusOK)

	return nil
}

const handleModelRecommendationPrefix = "handleModelRec"

func (master *Master) handleModelRecommendation(predictions *[]syncutils.Prediction, sum, max, min *float64, count *int, request *syncutils.ClientRecRequest) error {
	log.Printf("INFO: %s: Handling model recommendation", handleModelRecommendationPrefix)
	defer log.Printf("INFO: %s: Model recommendation handled", handleModelRecommendationPrefix)

	beginStatus := master.slavesInfo.ReadStatus()
	log.Printf("INFO: %s: Slaves status: %v", handleModelRecommendationPrefix, beginStatus)

	nBatches := 0
	for _, active := range beginStatus {
		if active {
			nBatches++
		}
	}
	if nBatches == 0 {
		return fmt.Errorf("RecRequestErr: No active slaves")
	}
	log.Printf("INFO: %s: Created (%d) batches:", handleModelRecommendationPrefix, nBatches)

	cond := sync.NewCond(&sync.Mutex{})
	partialUserFactorsCh := make(chan *syncutils.SlavePartialUserFactors, nBatches)
	partialRecommendationCh := make(chan *syncutils.SlaveRecResponse, nBatches)
	masterUserFactors := syncutils.MasterUserFactors{
		UserId:      request.UserId,
		UserFactors: initializeUserFactors(master.modelConfig.NumFeatures),
	}

	batches := master.createBatches(nBatches, request.UserId, request.Ratings, request.Quantity, request.GenreIds, masterUserFactors.UserFactors)

	activeSlaveIds := master.slavesInfo.GetActiveIdsByStatus(true)

	go func() {
		for batchId, batch := range batches {
			go func(batchId int) {
				var slaveId int
				if batchId < len(activeSlaveIds) {
					slaveId = activeSlaveIds[batchId]
				} else {
					slaveId = master.slavesInfo.GetMinCountIdByStatus(true)
				}
				master.handleRecommendationRequestBatch(cond, partialUserFactorsCh, partialRecommendationCh, batchId, slaveId, &batch, &masterUserFactors)
			}(batchId)
		}
	}()

	userFactorsGrads := make([]float64, master.modelConfig.NumFeatures)
	weightCount := 0

	for i := 0; i < nBatches; i++ {
		partialUserFactors := <-partialUserFactorsCh
		for j := range partialUserFactors.WeightedGrad {
			userFactorsGrads[j] += partialUserFactors.WeightedGrad[j]
		}
		weightCount += partialUserFactors.Count
	}

	if weightCount != 0 {
		for i := range userFactorsGrads {
			userFactorsGrads[i] /= float64(weightCount)
		}
	}

	masterUserFactors.UserId = request.UserId
	masterUserFactors.UserFactors = userFactorsGrads

	log.Printf("INFO: %s: User factors updated", handleModelRecommendationPrefix)
	cond.L.Lock()
	cond.Broadcast()
	cond.L.Unlock()

	for i := 0; i < nBatches; i++ {
		partialRecommendation := <-partialRecommendationCh
		if partialRecommendation.Count > 0 {
			*predictions = append(*predictions, partialRecommendation.Predictions...)
			*sum += partialRecommendation.Sum
			*count += partialRecommendation.Count
			if partialRecommendation.Max > *max {
				*max = partialRecommendation.Max
			}
			if partialRecommendation.Min < *min {
				*min = partialRecommendation.Min
			}
			if i > 0 {
				sort.Slice(*predictions, func(i, j int) bool {
					return (*predictions)[i].Rating > (*predictions)[j].Rating
				})
				if len(*predictions) > request.Quantity {
					*predictions = (*predictions)[:request.Quantity]
				}
			}
		}
	}

	return nil
}

func (master *Master) handleRecommendationRequestBatch(cond *sync.Cond, partialUserFactorsCh chan *syncutils.SlavePartialUserFactors, partialRecommendationCh chan *syncutils.SlaveRecResponse, batchId, slaveId int, batch *syncutils.MasterRecRequest, masterUserFactors *syncutils.MasterUserFactors) {
	var err error
	var conn net.Conn
	log.Printf("INFO: RequestBatch: Handling batch (%d).\n", batchId)
	defer log.Printf("INFO: RequestBatch: Batch (%d) handled.\n", batchId)
	for {
		log.Printf("INFO: Trying to connect batch to slaveId (%d)\n", slaveId)
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", master.slaveIps[slaveId], syncutils.RecommendationPort))
		if err != nil {
			master.slavesInfo.WriteStatusByIndex(false, slaveId)
			log.Printf("ERROR: RequestBatchErr: Error connecting to slave node %d for batch %d: %v", slaveId, batchId, err)
			log.Printf("Slaves status: %v", master.slavesInfo.ReadStatus())
			slaveId = master.slavesInfo.GetMinCountIdByStatus(true)
			continue
		}
		timeout := 20 * time.Second
		conn.SetDeadline(time.Now().Add(timeout))

		err = master.handlePartialRecommendation(&conn, cond, partialUserFactorsCh, partialRecommendationCh, slaveId, batchId, batch, masterUserFactors)
		if err != nil {
			log.Println("ERROR: RequestBatchErr: ", err)
			master.slavesInfo.WriteStatusByIndex(false, slaveId)
			slaveId = master.slavesInfo.GetMinCountIdByStatus(true)
			continue
		}
		break
	}
}

func (master *Master) handlePartialRecommendation(conn *net.Conn, cond *sync.Cond, partialUserFactorsCh chan *syncutils.SlavePartialUserFactors, partialRecommendationCh chan *syncutils.SlaveRecResponse, slaveId, batchId int, batch *syncutils.MasterRecRequest, masterUserFactors *syncutils.MasterUserFactors) error {
	// sendRequest
	err := syncutils.SendObjectAsJsonMessage(batch, conn)
	if err != nil {
		return fmt.Errorf("partialRecommendErr: Error sending batch to slave node (%d) for batch (%d): %v", slaveId, batchId, err)
	}
	// ReceivePartialUserFactors
	var partialUserFactors syncutils.SlavePartialUserFactors
	err = syncutils.ReceiveJsonMessageAsObject(&partialUserFactors, conn)
	if err != nil {
		return fmt.Errorf("partialRecommendErr: Error receiving partial user factors from slave node (%d): %v", slaveId, err)
	}

	partialUserFactorsCh <- &partialUserFactors

	cond.L.Lock()
	cond.Wait()
	cond.L.Unlock()

	// SendUserFactors
	err = syncutils.SendObjectAsJsonMessage(masterUserFactors, conn)
	if err != nil {
		return fmt.Errorf("partialRecommendErr: Error sending user factors to slave node (%d): %v", slaveId, err)
	}

	// ReceivePartialRecommendation
	var response syncutils.SlaveRecResponse
	err = syncutils.ReceiveJsonMessageAsObject(&response, conn)
	if err != nil {
		return fmt.Errorf("partialRecommendErr: Error receiving response from slave node (%d): %v", slaveId, err)
	}
	partialRecommendationCh <- &response
	return nil
}

func (master *Master) createBatches(nBatches, userId int, ratings []float64, quantity int, genreIds []int, userFactors []float64) []syncutils.MasterRecRequest {
	batches := make([]syncutils.MasterRecRequest, nBatches)
	var rangeSize int = len(master.movieTitles) / nBatches
	var startMovieId int = 0
	for i := 0; i < nBatches; i++ {
		var endMovieId int = startMovieId + rangeSize
		if i == nBatches-1 {
			endMovieId = len(master.movieTitles)
		}
		batches[i] = syncutils.MasterRecRequest{
			UserId:       userId,
			UserRatings:  ratings[startMovieId:endMovieId],
			StartMovieId: startMovieId,
			EndMovieId:   endMovieId,
			Quantity:     quantity,
			GenreIds:     genreIds,
			UserFactors:  userFactors,
		}
		startMovieId = endMovieId
	}
	return batches
}

func initializeUserFactors(numFeatures int) []float64 {
	userFactors := make([]float64, numFeatures)

	for i := 0; i < numFeatures; i++ {
		userFactors[i] = rand.Float64()*0.02 - 0.01 // Valores en el rango [-0.01, 0.01]
	}

	return userFactors
}

func FedAvg(gradients [][]float64, weights []float64) []float64 {
	numFeatures := len(gradients[0]) // Asumimos que todos los vectores tienen la misma longitud

	avgGrad := make([]float64, numFeatures)

	for i := 0; i < numFeatures; i++ {
		for j, grad := range gradients {
			avgGrad[i] += grad[i] * weights[j] // Pondera cada gradiente
		}
	}

	sumWeights := 0.0
	for _, w := range weights {
		sumWeights += w
	}

	for i := range avgGrad {
		avgGrad[i] /= sumWeights
	}

	return avgGrad // Retorna el gradiente promedio para actualizar el modelo global
}
