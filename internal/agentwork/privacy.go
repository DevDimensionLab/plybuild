package agentwork

import "regexp"

var privateTextPattern = regexp.MustCompile(`(?i)(/Users/|/private/|/home/|/tmp/|~[/\\]|xox[baprs]-|xapp-|token=|password=|secret=|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`)

func privateText(s string) bool { return privateTextPattern.MatchString(s) }
