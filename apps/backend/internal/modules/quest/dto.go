package quest

// dto.go = bentuk data yang masuk/keluar lewat HTTP (request & response).
// DTO TERPISAH dari entitas domain: entitas boleh berubah tanpa mengubah
// kontrak API, dan sebaliknya. Konversi domain<->DTO ada di handler.
//
// TODO:
//   type CreateQuestRequest struct {
//       Title      string `json:"title"      validate:"required"`
//       Note       string `json:"note"`
//       Category   string `json:"category"   validate:"required"`
//       Difficulty string `json:"difficulty" validate:"required,oneof=easy medium hard"`
//   }
//   type UpdateQuestRequest struct { /* field opsional (pointer) */ }
//   type CompleteQuestRequest struct { /* opsional: progress/value */ }
//
//   type QuestResponse struct {
//       ID string; Title string; Category string; Difficulty string;
//       Points int; Active bool; ...
//   }
//   type QuestLogResponse struct { QuestID string; Date string; Status string; Points int }
//   type TodayQuestsResponse struct {
//       Date  string
//       Items []struct{ Quest QuestResponse; Completed bool }
//   }
//
//   // fungsi mapper: toQuestResponse(Quest) QuestResponse, dst.
