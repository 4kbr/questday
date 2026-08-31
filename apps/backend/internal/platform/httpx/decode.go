package httpx

// TODO (baca request):
//   // Decode membaca JSON body ke dst, menolak field asing (DisallowUnknownFields),
//   // membatasi ukuran body, dan mengembalikan error yang ramah untuk JSON rusak.
//   func Decode(w http.ResponseWriter, r *http.Request, dst any) error
//
//   // DecodeAndValidate: Decode lalu jalankan validator (platform/validator).
//   func DecodeAndValidate(w, r, dst) error
