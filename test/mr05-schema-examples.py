"""Validate the public MR05 examples with the required offline interpreter."""
import json
from pathlib import Path
import jsonschema

root = Path(__file__).resolve().parent.parent
schema = json.loads((root / "schemas/workflow-control-v1.schema.json").read_text())
jsonschema.Draft202012Validator.check_schema(schema)
validator = jsonschema.Draft202012Validator(schema, format_checker=jsonschema.FormatChecker())
examples = sorted((root / "schemas/workflow-control-examples").glob("*.json"))
assert examples, "no public examples"
for example in examples:
    value = json.loads(example.read_text())
    errors = list(validator.iter_errors(value))
    assert bool(errors) == (".invalid." in example.name), (example.name, errors)
print("MR05 public examples: valid and invalid cases checked")
