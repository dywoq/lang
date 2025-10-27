package symbol

// Symbol represents the symbol of code,
// containing the information.
type Symbol struct {
	Exported   bool   `json:"exported"`
	Const      bool   `json:"const"`
	Consteval  bool   `json:"consteval"`
	Copied     bool   `json:"copied"`
	CopiedFrom string `json:"copied_from"`
}
