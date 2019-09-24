import Express from 'express'

const app = new Express()
const port = '3000'

app.get('/info', function (req, res) {
  console.log('Get info')
  res.status(200).send({'status': 'success'})
})

app.get('/errorStackTrace', function (req, res) {
  console.error('Error stacktrace')
  res.status(500).send({'status': 'error'})
})

app.listen(port, (err) => {
  if (err) {
    console.error(err)
  }
  console.info("success setup")
})
