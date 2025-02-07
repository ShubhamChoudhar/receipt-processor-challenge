package main

import (
    "encoding/json"
    "fmt"
    "log"
    "math"
    "net/http"
    "regexp"
    "strconv"
    "strings"
    "time"

    "github.com/google/uuid"
    "github.com/gorilla/mux"
)

type Receipt struct {
    Retailer     string `json:"retailer"`
    PurchaseDate string `json:"purchaseDate"`
    PurchaseTime string `json:"purchaseTime"`
    Total        string `json:"total"`
    Items        []Item `json:"items"`
}


type Item struct {
    ShortDescription string `json:"shortDescription"`
    Price            string `json:"price"`
}

var pointsStore = make(map[string]int)

func handleProcessReceipt(w http.ResponseWriter, r *http.Request) {
    var receipt Receipt
    if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
        http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
        return
    }

    pts, err := calculatePoints(&receipt)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    id := uuid.New().String()

    pointsStore[id] = pts

    w.Header().Set("Content-Type", "application/json")
    response := map[string]string{"id": id}
    json.NewEncoder(w).Encode(response)
}

func handleGetPoints(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    pts, ok := pointsStore[id]
    if !ok {
        http.Error(w, "Receipt not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    response := map[string]int{"points": pts}
    json.NewEncoder(w).Encode(response)
}

func calculatePoints(receipt *Receipt) (int, error) {
    totalPoints := 0
    alnumRe := regexp.MustCompile(`[A-Za-z0-9]`)
    alnumMatches := alnumRe.FindAllString(receipt.Retailer, -1)
    totalPoints += len(alnumMatches)

    purchaseTotal, err := strconv.ParseFloat(receipt.Total, 64)
    if err != nil {
        return 0, fmt.Errorf("invalid total: %v", err)
    }

    if purchaseTotal == float64(int64(purchaseTotal)) {
        totalPoints += 50
    }

    quarter := purchaseTotal / 0.25
    if math.Mod(quarter, 1) == 0 {
        totalPoints += 25
    }

    totalPoints += (len(receipt.Items) / 2) * 5

    for _, item := range receipt.Items {
        trimmedDesc := strings.TrimSpace(item.ShortDescription)
        if len(trimmedDesc) > 0 && len(trimmedDesc)%3 == 0 {
            itemPrice, err := strconv.ParseFloat(item.Price, 64)
            if err != nil {
                return 0, fmt.Errorf("invalid item price: %v", err)
            }
            bonus := math.Ceil(itemPrice * 0.2)
            totalPoints += int(bonus)
        }
    }

    if purchaseTotal > 10.00 {
        totalPoints += 5
    }

    dateLayout := "2006-01-02"
    t, err := time.Parse(dateLayout, receipt.PurchaseDate)
    if err != nil {
        return 0, fmt.Errorf("invalid purchaseDate: %v", err)
    }
    day := t.Day()
    if day%2 == 1 {
        totalPoints += 6
    }

    timeLayout := "15:04"
    tt, err := time.Parse(timeLayout, receipt.PurchaseTime)
    if err != nil {
        return 0, fmt.Errorf("invalid purchaseTime: %v", err)
    }

    hour := tt.Hour()
    if hour == 14 || hour == 15 {
        totalPoints += 10
    }

    return totalPoints, nil
}

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/receipts/process", handleProcessReceipt).Methods("POST")
    r.HandleFunc("/receipts/{id}/points", handleGetPoints).Methods("GET")

    srv := &http.Server{
        Handler: r,
        Addr:    ":8080",
    }

    log.Println("Starting server on :8080...")
    log.Fatal(srv.ListenAndServe())
}