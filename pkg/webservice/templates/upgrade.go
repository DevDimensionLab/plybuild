package templates

const Upgrade = header + upgrade + footer

const upgrade = `
<div class="container">
  <div class="row justify-content-md-center">
    <div class="col">
    <h1>Co-pilot - Upgrade</h1>
    <form class="form-inline" action="/api/upgrade" method="POST">
      <button type="submit" class="btn btn-primary btn-block">Submit</button>
    </form>
    </div>
  </div>
</div>
`
