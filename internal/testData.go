package testdata

import "ngevent/internal/model"

func GenerateJabodetabekDataset() (users []model.Location, events []model.Location) {
	users = []model.Location{
		{Name: "User-Condet", Lat: -6.3397286, Lon: 106.8331984},
		{Name: "User-PasarMinggu", Lat: -6.3144585, Lon: 106.8268601},
		{Name: "User-PondokLabu", Lat: -6.305911, Lon: 106.7994565},
		{Name: "User-TanjungPriok", Lat: -6.1330885, Lon: 106.8695885},
		{Name: "User-Kemayoran", Lat: -6.1606132, Lon: 106.8422236},
		{Name: "User-PasarRebo", Lat: -6.3313338, Lon: 106.845284},
		{Name: "User-Pluit", Lat: -6.1109418, Lon: 106.7781484},
		{Name: "User-Semanggi", Lat: -6.226084, Lon: 106.8202041},
		{Name: "User-Cakung", Lat: -6.1727261, Lon: 106.9354858},
		{Name: "User-Kebagusan", Lat: -6.3175008, Lon: 106.8269633},
	}

	events = []model.Location{
		{Name: "Event-GrandIndonesia", Lat: -6.1957601, Lon: 106.8214547},
		{Name: "Event-Ancol", Lat: -6.125215, Lon: 106.8362474},
		{Name: "Event-TMI", Lat: -6.3026341, Lon: 106.8952962},
		{Name: "Event-JIS", Lat: -6.1250747, Lon: 106.8609623},
		{Name: "Event-GBK", Lat: -6.2186492, Lon: 106.8036258},
		{Name: "Event-KotaKasablanka", Lat: -6.2232551, Lon: 106.8426972},
		{Name: "Event-TebetEcoPark", Lat: -6.2403737, Lon: 106.8524328},
		{Name: "Event-SummareconMallBekasi", Lat: -6.2260106, Lon: 107.0010614},
		{Name: "Event-BSDTangerang", Lat: -6.300681, Lon: 106.6365721},
		{Name: "Event-MargoCity", Lat: -6.3728858, Lon: 106.8350994},
	}

	return users, events
}
