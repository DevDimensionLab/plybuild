package webservice

const generateTemplate = `
<html>
<head>
<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta1/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-giJF6kkoqNQ00vy+HMDP7azOuL0xtbfIcaT9wjKHr8RbDVddVHyTfAAsrekwKmP1" crossorigin="anonymous">
<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta1/dist/js/bootstrap.bundle.min.js" integrity="sha384-ygbV9kiqUc6oa4msXn9868pTtWMgiQaeYH7/t7LECLbyPA2x65Kgf80OJFdroafW" crossorigin="anonymous"></script>
</head>
<body>
<div class="container">
  <div class="row justify-content-md-center">
  <div class="col">
    <h1>Co-pilot - Generate</h1>
	<form class="form-inline" action="/api/generate" method="POST">

	<div class="row">
		<div class="col form-group mb-3">
			<label class="sr-only" for="groupId">GroupId</label>
			<input type="text" class="form-control" id="groupId" name="groupId" value={{.ProjectConfig.GroupId}}>
		</div>
		<div class="col form-group mb-3">
			<label class="sr-only" for="artifactId">ArtifactId</label>
			<input type="text" class="form-control" id="artifactId" name="artifactId" value={{.ProjectConfig.ArtifactId}}>
		</div>
	</div>
	<div class="form-group mb-3">
		<label class="sr-only" for="package">Package</label>
		<input type="text" class="form-control" id="package" name="package" value={{.ProjectConfig.Package}}>
	</div>
	<div class="form-group mb-3">
		<label class="sr-only" for="name">Name</label>
		<input type="text" class="form-control" id="name" name="name" value={{.ProjectConfig.Name}}>
	</div>
	<div class="form-group mb-3">
		<label class="sr-only" for="description">Description</label>
		<input type="text" class="form-control" id="description" name="description" value={{.ProjectConfig.Description}}>
	</div>
	<div class="form-check form-check-inline mb-3">
	  <input class="form-check-input" type="radio" id="language1" value="kotlin" name="language" checked>
	  <label class="form-check-label" for="language1">Kotlin</label>
	</div>
	<div class="form-check form-check-inline mb-3">
	  <input class="form-check-input" type="radio" id="language2" value="java" name="language">
	  <label class="form-check-label" for="language2">Java</label>
	</div>
	<div class="form-group mb-3">
		<label for="templates">Templates</label>
		<select multiple class="form-control" id="templates" name="templates" size=5>
			{{range .CloudConfig.Templates }}
				<option>{{ .Name }}</option>
			{{end}}
		</select>
	</div>
	<div class="form-group mb-3">
		<label for="dependencies">Dependencies</label>
		<select multiple class="form-control" id="dependencies" name="dependencies" size="10">
			{{range .IoResponse.Dependencies.Values }}
				<option disabled>{{ .Name }}</option>
				{{range .Values }}
					<option value="{{ .Id }}">{{ .Name }}</option>
				{{end}}
			{{end}}
		</select>
	</div>
	<button type="submit" class="btn btn-primary btn-block">Submit</button>
	</form>
  </div>
  </div>
</div>
</body>
</html>
`
