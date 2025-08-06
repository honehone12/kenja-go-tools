package opensearch

import (
	"strings"
)

const PLACE_HOLDER = ""

var __OPENSEARCH_REPLACER = strings.NewReplacer(
	"_", PLACE_HOLDER,
	"+", PLACE_HOLDER,
	"-", PLACE_HOLDER,
	"*", PLACE_HOLDER,
	"%", PLACE_HOLDER,
	"&", PLACE_HOLDER,
	"|", PLACE_HOLDER,
	"=", PLACE_HOLDER,
	"<", PLACE_HOLDER,
	">", PLACE_HOLDER,
	"!", PLACE_HOLDER,
	"?", PLACE_HOLDER,
	"^", PLACE_HOLDER,
	"~", PLACE_HOLDER,
	"{", PLACE_HOLDER,
	"}", PLACE_HOLDER,
	"[", PLACE_HOLDER,
	"]", PLACE_HOLDER,
	"(", PLACE_HOLDER,
	")", PLACE_HOLDER,
	":", PLACE_HOLDER,
	",", PLACE_HOLDER,
	".", PLACE_HOLDER,
	"/", PLACE_HOLDER,
	`\`, PLACE_HOLDER,
	`"`, PLACE_HOLDER,
	"'", PLACE_HOLDER,
	"`", PLACE_HOLDER,
)
