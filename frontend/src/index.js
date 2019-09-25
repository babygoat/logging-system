import Express from 'express'
import logger from './logger'

const app = new Express()
const port = '3000'

app.get('/info', function (req, res) {
  logger.info('This is log with severity of info!')
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
  logger.info("success setup")
})
