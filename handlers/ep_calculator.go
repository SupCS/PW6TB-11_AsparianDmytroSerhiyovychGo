package handlers

import (
	"html/template"
	"math"
	"net/http"
	"strconv"
)

// Структура для обладнання
type EquipmentData struct {
	Name                     string
	EfficiencyRating         float64
	LoadPowerFactor          float64
	LoadVoltage              float64
	Quantity                 int
	NominalPower             float64
	UsageFactor              float64
	ReactivePowerCoefficient float64
	MultipliedPower          float64
	Current                  float64
}

// Структура для передачі даних в шаблон
type EpCalculatorResult struct {
	EquipmentList           []EquipmentData
	KvGroup                 float64
	EffEpAmount             float64
	TotalDepartmentUtilCoef float64
	EffEpDepartmentAmount   float64
	RozrahActNav            float64
	RozrahReactNav          float64
	FullPower               float64
	RozrahGroupStrumShr1    float64
	RozrahActNavShin        float64
	RozrahReactNavShin      float64
	FullPowerShin           float64
	RozrahGroupStrumShin    float64
}

func EpCalculatorHandler(w http.ResponseWriter, r *http.Request) {
	// Жорстко заданий список обладнання
	equipmentList := []EquipmentData{
		{"Шліфувальний верстат", 0.92, 0.9, 0.38, 4, 20, 0.15, 1.33, 0, 0},
		{"Свердлильний верстат", 0.92, 0.9, 0.38, 2, 14, 0.12, 1.0, 0, 0},
		{"Фугувальний верстат", 0.92, 0.9, 0.38, 4, 42, 0.15, 1.33, 0, 0},
		{"Циркулярна пила", 0.92, 0.9, 0.38, 1, 36, 0.3, 1.52, 0, 0},
		{"Прес", 0.92, 0.9, 0.38, 1, 20, 0.5, 0.75, 0, 0},
		{"Полірувальний верстат", 0.92, 0.9, 0.38, 1, 40, 0.2, 1.0, 0, 0},
		{"Фрезерний верстат", 0.92, 0.9, 0.38, 2, 32, 0.2, 1.0, 0, 0},
		{"Вентилятор", 0.92, 0.9, 0.38, 1, 20, 0.65, 0.75, 0, 0},
	}

	if r.Method == http.MethodPost {
		var sumNPnKvProduct, sumNPnProduct, sumNPnPnProduct float64

		for i := range equipmentList {
			quantity, _ := strconv.Atoi(r.FormValue("quantity" + strconv.Itoa(i)))
			if quantity > 0 {
				equipmentList[i].Quantity = quantity
				equipmentList[i].MultipliedPower = float64(quantity) * equipmentList[i].NominalPower
				equipmentList[i].Current = equipmentList[i].MultipliedPower / (math.Sqrt(3) * equipmentList[i].LoadVoltage * equipmentList[i].LoadPowerFactor * equipmentList[i].EfficiencyRating)
				sumNPnKvProduct += equipmentList[i].MultipliedPower * equipmentList[i].UsageFactor
				sumNPnProduct += equipmentList[i].MultipliedPower
				sumNPnPnProduct += float64(quantity) * equipmentList[i].NominalPower * equipmentList[i].NominalPower
			}
		}

		KvGroup := sumNPnKvProduct / sumNPnProduct
		EffEpAmount := math.Ceil((sumNPnProduct * sumNPnProduct) / sumNPnPnProduct)

		Kr := 1.25
		PH := 27.0
		TanPhi := 1.63
		Un := 0.28

		Pp := Kr * sumNPnKvProduct
		Qp := KvGroup * PH * TanPhi
		Sp := math.Sqrt((Pp * Pp) + (Qp * Qp))
		Ip := Pp / Un

		KvDepartment := 752.0 / 2330.0
		NE := 2330.0 * 2330.0 / 96399.0

		Kr2 := 0.7
		PpShin := Kr2 * 752.0
		QpShin := Kr2 * 657.0
		SpShin := math.Sqrt((PpShin * PpShin) + (QpShin * QpShin))
		IpShin := PpShin / 0.38

		result := EpCalculatorResult{
			EquipmentList:           equipmentList,
			KvGroup:                 KvGroup,
			EffEpAmount:             EffEpAmount,
			TotalDepartmentUtilCoef: KvDepartment,
			EffEpDepartmentAmount:   NE,
			RozrahActNav:            Pp,
			RozrahReactNav:          Qp,
			FullPower:               Sp,
			RozrahGroupStrumShr1:    Ip,
			RozrahActNavShin:        PpShin,
			RozrahReactNavShin:      QpShin,
			FullPowerShin:           SpShin,
			RozrahGroupStrumShin:    IpShin,
		}

		tmpl, _ := template.ParseFiles("templates/ep_calculator.html")
		tmpl.Execute(w, result)
		return
	}

	tmpl, _ := template.ParseFiles("templates/ep_calculator.html")
	tmpl.Execute(w, nil)
}
