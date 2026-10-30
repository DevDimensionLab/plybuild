package taskrun

import ()

// Model the pane disappearing after pane get, before the follow's agent get.
// The later inventory/process evidence is synthetic, never real OS coverage.
const al21CloseoutRaceStandin = `
if cmd == ["agent", "get"] and model.get("al21_disappear"):
    model["al21_disappear"] = False
    model["closed_transport"] = True
    path.write_text(json.dumps(model))
`
