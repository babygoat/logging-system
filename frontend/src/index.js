import Express from 'express'
import logger from './logger'

const app = new Express()
const port = '3000'

const errorFunction = () => {
  throw Error("Unknown error")
}

app.get('/info', function (req, res) {
  logger.info('This is log with severity of info!')
  res.status(200).send({'status': 'success'})
})

app.get('/errorStackTrace', function (req, res) {
  try {
    errorFunction()
  } catch (err) {
    logger.error(err.stack)
  }
  res.status(500).send({'status': 'error'})
})

app.listen(port, (err) => {
  if (err) {
    logger.error(err)
  }
  logger.info("success setup")
})
